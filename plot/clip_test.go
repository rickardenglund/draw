package plot

import (
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func TestClipSegment(t *testing.T) {
	r := rl.NewRectangle(0, 0, 10, 10)
	v := rl.NewVector2
	for _, c := range []struct {
		name   string
		a, b   rl.Vector2
		ok     bool
		ca, cb rl.Vector2
	}{
		{"inside", v(1, 1), v(9, 9), true, v(1, 1), v(9, 9)},
		{"one end out", v(5, 5), v(5, 20), true, v(5, 5), v(5, 10)},
		{"both out, crossing", v(-5, 5), v(15, 5), true, v(0, 5), v(10, 5)},
		{"both out, diagonal", v(-5, -5), v(15, 15), true, v(0, 0), v(10, 10)},
		{"both out, missing", v(-5, 4), v(4, -5), false, v(0, 0), v(0, 0)},
		{"outside, parallel", v(-1, 0), v(-1, 10), false, v(0, 0), v(0, 0)},
		{"point", v(3, 3), v(3, 3), true, v(3, 3), v(3, 3)},
	} {
		a, b, ok := clipSegment(c.a, c.b, r)
		if ok != c.ok || (ok && (a != c.ca || b != c.cb)) {
			t.Errorf("%s: got %v %v %v, want %v %v %v", c.name, a, b, ok, c.ca, c.cb, c.ok)
		}
	}
}
