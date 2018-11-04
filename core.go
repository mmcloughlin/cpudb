package cpudb

// CPU represents CPU properties.
type CPU struct {
	Type     string
	Alias    string
	Platform string
	Stepping string
	Leaves   map[uint32][]Leaf
}

// Leaf contains values from one CPUID leaf.
type Leaf struct {
	EAX uint32
	EBX uint32
	ECX uint32
	EDX uint32
}

// CPUID simulates a CPUID call on this CPU. Boolean return value is false if the leaf does not exist.
func (c *CPU) CPUID(eax, ecx uint32) (Leaf, bool) {
	sub, ok := c.Leaves[eax]
	if !ok || int(ecx) >= len(sub) {
		return Leaf{}, false
	}
	return sub[ecx], true
}
