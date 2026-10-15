package plot

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/input"
)

type mover struct {
	start *rl.Vector2

	offset    rl.Vector2
	tmpOffset rl.Vector2
}

func (m *mover) Draw(target rl.Rectangle) {
	mp := input.MousePosition()
	if m.start != nil {
		rl.DrawLineV(*m.start, mp, rl.ColorAlpha(rl.White, .5))
	}
}

func (m *mover) handleKeys(active bool, s scale) {
	mp := input.MousePosition()
	if active {

		if rl.IsKeyDown(rl.KeyD) {
			if m.start == nil {
				m.start = &mp
				//m.offset = rl.NewVector2(0, 0)
				m.tmpOffset = rl.Vector2Zero()
			} else {
				delta := rl.Vector2Subtract(s.transformR(*m.start), s.transformR(mp))

				m.tmpOffset = delta
			}
		}
		if rl.IsKeyReleased(rl.KeyD) {
			m.start = nil
			m.offset = rl.Vector2Add(m.tmpOffset, m.offset)
			m.tmpOffset = rl.Vector2Zero()
		}
	} else {
		m.start = nil
		m.tmpOffset = rl.Vector2Zero()
	}

}

func (m *mover) GetOffset() rl.Vector2 {
	return rl.Vector2Add(m.offset, m.tmpOffset)
}

func (m *mover) Reset() {
	m.offset = rl.NewVector2(0, 0)
}
