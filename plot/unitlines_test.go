package plot

import "testing"

func TestUnitLines(t *testing.T) {
	measureUnit = func(s string) float32 { return float32(10 * len(s)) }
	one, word := measureUnit("mm/s RMS"), measureUnit("mm/s")

	lines, ok := unitLines("mm/s RMS", one)
	if !ok || len(lines) != 1 {
		t.Fatalf("fits on one line: %q %v", lines, ok)
	}
	lines, ok = unitLines("mm/s RMS", word)
	if !ok || len(lines) != 2 || lines[0] != "mm/s" || lines[1] != "RMS" {
		t.Fatalf("split at the space: %q %v", lines, ok)
	}
	if _, ok := unitLines("mm/s RMS", word-1); ok {
		t.Fatal("a word that doesn't fit falls back to the tick")
	}
	if lines, ok := unitLines("", 10); !ok || len(lines) != 0 {
		t.Fatalf("no unit: %q %v", lines, ok)
	}
}

func TestActiveMarkGlides(t *testing.T) {
	var a activeMark
	a.moveTo(10, 0)
	if a.x(0) != 10 {
		t.Fatalf("starts where first set: %v", a.x(0))
	}
	a.moveTo(20, 1)
	if got := a.x(1 + activeMarkDuration/2); got != 15 {
		t.Fatalf("halfway: %v", got)
	}
	a.moveTo(20, 1+activeMarkDuration/2) // same target: keeps going
	if got := a.x(1 + activeMarkDuration); got != 20 {
		t.Fatalf("arrives: %v", got)
	}
	a.moveTo(0, 1+activeMarkDuration/2+activeMarkDuration) // from where it is
	if got := a.x(1 + activeMarkDuration/2 + activeMarkDuration); got != 20 {
		t.Fatalf("restarts from current: %v", got)
	}
	if a.prev != 20 || a.prevSize != 1 {
		t.Fatalf("old ring shrinks from full size: %v %v", a.prev, a.prevSize)
	}
	if got := a.progress(1 + activeMarkDuration/2 + 2*activeMarkDuration); got != 1 {
		t.Fatalf("new ring fully grown: %v", got)
	}
}
