package cpudb

import "testing"

func TestHaveSomeCPUs(t *testing.T) {
	n := len(CPUs)
	t.Logf("len(CPUs)=%d", n)
	if n == 0 {
		t.Fatal("expected some CPUs")
	}
}
