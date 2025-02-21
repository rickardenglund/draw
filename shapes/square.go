package shapes

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Square struct {
	size rl.Vector2
	clr  rl.Color
}

func NewSquare(size rl.Vector2, clr rl.Color) Square {
	return Square{size: size, clr: clr}
}

func (s Square) Draw(pos, size rl.Vector2) {
	p := rl.Vector2Add(pos, rl.Vector2Scale(size, .5))
	p = rl.Vector2Subtract(p, rl.Vector2Scale(s.size, .5))

	rl.DrawRectangleV(p, s.size, s.clr)
}
