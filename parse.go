package cpuidb

import (
	"bufio"
	"errors"
	"io"
	"os"
	"regexp"
	"strings"
)

type Section struct {
	Name       string
	Properties map[string]string
}

func NewSection(name string) *Section {
	return &Section{
		Name:       name,
		Properties: make(map[string]string),
	}
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

		s.Properties[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if s != nil {
		cfg.Sections = append(cfg.Sections, s)
	}

	return cfg, nil
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

	cpu := &CPU{
		Type:  info.Properties["CPU Type"],
		Alias: info.Properties["CPU Alias"],
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
