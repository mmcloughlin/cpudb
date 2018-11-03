package cpuidb

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

type Property struct {
	Key   string
	Value string
}

type Section struct {
	Name       string
	Properties []Property
}

func NewSection(name string) *Section {
	return &Section{
		Name: name,
	}
}

func (s *Section) AddProperty(k, v string) {
	s.Properties = append(s.Properties, Property{
		Key:   k,
		Value: v,
	})
}

func (s *Section) Property(key string) string {
	for _, p := range s.Properties {
		if p.Key == key {
			return p.Value
		}
	}
	return ""
}

type Config struct {
	Sections []*Section
}

func (c *Config) LookupSection(name string) (*Section, bool) {
	for _, s := range c.Sections {
		if s.Name == name {
			return s, true
		}
	}
	return nil, false
}

var headingRegexp = regexp.MustCompile(`^------\[ (.+) \]------$`)

// ParseConfig parses the overall sections and key-value pairs.
func ParseConfig(r io.Reader) (*Config, error) {
	cfg := &Config{}
	var s *Section

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		// Match a heading.
		m := headingRegexp.FindSubmatch(scanner.Bytes())
		if m != nil {
			if s != nil {
				cfg.Sections = append(cfg.Sections, s)
			}
			s = NewSection(string(m[1]))
			continue
		}

		// Match a key-value pair.
		if s == nil {
			continue
		}

		parts := strings.SplitN(scanner.Text(), ":", 2)
		if len(parts) != 2 {
			continue
		}

		s.AddProperty(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if s != nil {
		cfg.Sections = append(cfg.Sections, s)
	}

	return cfg, nil
}

func BuildCPUIDLeaves(s *Section) (map[uint32][]Leaf, error) {
	leaves := make(map[uint32][]Leaf)

	for _, p := range s.Properties {
		var eax uint32
		n, err := fmt.Sscanf(p.Key, "CPUID %X", &eax)
		if n != 1 || err != nil {
			continue
		}

		var l Leaf
		n, err = fmt.Sscanf(p.Value, "%X-%X-%X-%X", &l.EAX, &l.EBX, &l.ECX, &l.EDX)
		if n != 4 || err != nil {
			return nil, err
		}

		leaves[eax] = append(leaves[eax], l)
	}

	return leaves, nil
}

// ParseCPU parses CPU data.
func ParseCPU(r io.Reader) (*CPU, error) {
	// Parse config sections.
	cfg, err := ParseConfig(r)
	if err != nil {
		return nil, err
	}

	// Fetch CPU Info.
	info, found := cfg.LookupSection("CPU Info")
	if !found {
		return nil, errors.New("missing CPU Info section")
	}

	// Fetch CPUID Info.
	cpu0, found := cfg.LookupSection("Logical CPU #0")
	if !found {
		return nil, errors.New("missing CPUID for CPU #0")
	}

	leaves, err := BuildCPUIDLeaves(cpu0)
	if err != nil {
		return nil, err
	}

	// Construct CPU.
	cpu := &CPU{
		Type:     info.Property("CPU Type"),
		Alias:    info.Property("CPU Alias"),
		Platform: info.Property("CPU Platform"),
		Stepping: info.Property("CPU Stepping"),
		Leaves:   leaves,
	}

	return cpu, nil
}

// ParseCPUFile parses CPU data from a given file. This is just a convenience around ParseCPU.
func ParseCPUFile(filename string) (*CPU, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ParseCPU(f)
}
