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

func (s Square) Draw(target rl.Rectangle) {
	targetPos := rl.NewVector2(target.X, target.Y)
	targetSize := rl.NewVector2(target.Width, target.Height)
	center := rl.Vector2Add(targetPos, rl.Vector2Scale(targetSize, .5))

	topLeft := rl.Vector2Add(center, rl.Vector2Scale(s.size, -.5))

	rl.DrawRectangleV(topLeft, s.size, s.clr)
}
