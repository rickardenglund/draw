package plot

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/theme"
)

type marker struct {
	v         float32
	harmonics int
}
type markers struct {
	xs []marker
}

func newMarkers() *markers {
	return &markers{
		xs: []marker{},
	}
}

func (m *markers) addMarker(x marker) {
	m.xs = append(m.xs, x)
}

func (m *markers) Draw(s scale) {
	for _, v := range m.xs {
		for h := range v.harmonics {
			c := v.v * float32(h+1)
			if c > s.l.minX && c < s.l.maxX {
				markerDown := s.transform(rl.NewVector2(c, s.l.minY))
				markerUp := s.transform(rl.NewVector2(c, s.l.maxY))
				rl.DrawLineEx(markerDown, markerUp, 1, rl.ColorAlpha(theme.Charcoal, .7))
			}
		}
	}
}
