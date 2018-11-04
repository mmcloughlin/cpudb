// Package parse provides parsing for InstLatx64 CPUID dumps.
package parse

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/mmcloughlin/cpudb"
)

// Property is a key-value pair.
type Property struct {
	Key   string
	Value string
}

// Section is a section from the input, with a heading and a collection of properties.
type Section struct {
	Name       string
	Properties []Property // Slice preferred over a map to preserve order.
}

// NewSection constructs a section with the given name.
func NewSection(name string) *Section {
	return &Section{
		Name: name,
	}
}

// AddProperty appends a Property to the section.
func (s *Section) AddProperty(k, v string) {
	s.Properties = append(s.Properties, Property{
		Key:   k,
		Value: v,
	})
}

// Property returns the first Properties entry with the given key, or the empty string if not found.
func (s *Section) Property(key string) string {
	for _, p := range s.Properties {
		if p.Key == key {
			return p.Value
		}
	}
	return ""
}

// Config represents all data found in a CPUID dump file.
type Config struct {
	Sections []*Section
}

// LookupSection returns the first Section with the given name. Returns nil if not found.
func (c *Config) LookupSection(name string) *Section {
	for _, s := range c.Sections {
		if s.Name == name {
			return s
		}
	}
	return nil
}

var headingRegexp = regexp.MustCompile(`^------\[ (.+) \]------$`)

// ConfigSections parses all sections and key-value pairs.
func ConfigSections(r io.Reader) (*Config, error) {
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

// BuildCPUIDLeaves parses properties starting with "CPUID" into a structured format.
func BuildCPUIDLeaves(s *Section) (map[uint32][]cpudb.Leaf, error) {
	leaves := make(map[uint32][]cpudb.Leaf)

	for _, p := range s.Properties {
		var eax uint32
		n, err := fmt.Sscanf(p.Key, "CPUID %X", &eax)
		if n != 1 || err != nil {
			continue
		}

		var l cpudb.Leaf
		n, err = fmt.Sscanf(p.Value, "%X-%X-%X-%X", &l.EAX, &l.EBX, &l.ECX, &l.EDX)
		if n != 4 || err != nil {
			return nil, err
		}

		leaves[eax] = append(leaves[eax], l)
	}

	return leaves, nil
}

// CPU parses CPU data.
func CPU(r io.Reader) (*cpudb.CPU, error) {
	// Parse config sections.
	cfg, err := ConfigSections(r)
	if err != nil {
		return nil, err
	}

	// Fetch CPU Info.
	info := cfg.LookupSection("CPU Info")
	if info == nil {
		return nil, errors.New("missing CPU Info section")
	}

	// Fetch CPUID Info.
	cpu0 := cfg.LookupSection("Logical CPU #0")
	if cpu0 == nil {
		return nil, errors.New("missing CPUID for CPU #0")
	}

	leaves, err := BuildCPUIDLeaves(cpu0)
	if err != nil {
		return nil, err
	}

	// Construct CPU.
	cpu := &cpudb.CPU{
		Type:     info.Property("CPU Type"),
		Alias:    info.Property("CPU Alias"),
		Platform: info.Property("CPU Platform"),
		Stepping: info.Property("CPU Stepping"),
		Leaves:   leaves,
	}

	return cpu, nil
}

// CPUFile parses CPU data from a given file. This is just a convenience around ParseCPU.
func CPUFile(filename string) (*cpudb.CPU, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return CPU(f)
}
