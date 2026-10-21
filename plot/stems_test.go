package plot

import (
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func TestXOnly(t *testing.T) {
	p := Plot[rl.Vector2]{curLimits: limits{minX: 0, maxX: 10, minY: 0, maxY: 5}}
	zoomed := limits{minX: 2, maxX: 4, minY: 1, maxY: 2}

	if got := p.xOnly(zoomed); got != zoomed {
		t.Errorf("line plot: got %+v, want %+v", got, zoomed)
	}

	p.stems = true
	want := limits{minX: 2, maxX: 4, minY: 0, maxY: 5}
	if got := p.xOnly(zoomed); got != want {
		t.Errorf("stem plot: got %+v, want %+v", got, want)
	}
}
