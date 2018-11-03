package cpuidb

// CPU represents CPU properties.
type CPU struct {
	Type     string
	Alias    string
	Platform string
	Stepping string
	Leaves   map[uint32][]Leaf
}

type Leaf struct {
	EAX uint32
	EBX uint32
	ECX uint32
	EDX uint32
}
