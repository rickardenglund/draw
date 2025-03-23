package plot

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/theme"
)

type mover struct {
	start *rl.Vector2
	moved rl.Vector2
}

func (m *mover) Draw(target rl.Rectangle) {
	mp := rl.GetMousePosition()
	if m.start != nil {
		rl.DrawLineV(*m.start, mp, rl.White)
	}

	if m.start != nil {
		s := fmt.Sprintf("%v", *m.start)
		p := rl.NewVector2(target.X, target.Y)
		theme.DrawTextPad(s, p, theme.Charcoal)
	}
}

func (m *mover) handleKeys(active bool, s scale) {
	mp := rl.GetMousePosition()
	if active {

		if rl.IsKeyDown(rl.KeyD) {
			if m.start == nil {
				m.start = &mp
				m.moved = rl.NewVector2(0, 0)
			} else {
				delta := rl.Vector2Subtract(s.transformR(*m.start), s.transformR(mp))

				m.moved = delta
			}
		}
		if rl.IsKeyReleased(rl.KeyD) {
			m.start = nil
		}
	} else {
		m.start = nil
		m.moved = rl.NewVector2(0, 0)
	}

}

func (m *mover) GetOffset() rl.Vector2 {
	return m.moved
}
