package shapes

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/animated"
	"github.com/rickardenglund/draw/theme"
)

type Circle struct {
	radius    *animated.Animated
	clr       rl.Color
	prevRad   float32
	animStart float64
}

func NewCircle(radius float32, clr rl.Color) *Circle {
	return &Circle{
		radius:    animated.NewAnimated(radius, 1),
		clr:       clr,
		prevRad:   0,
		animStart: rl.GetTime(),
	}
}

func (c *Circle) Set(newR float32) {
	c.radius.Set(newR)
}

func (c *Circle) Draw(target rl.Rectangle) {
	pos := rl.NewVector2(target.X, target.Y)
	size := rl.NewVector2(target.Width, target.Height)
	mp := rl.GetMousePosition()
	p := rl.Vector2Add(pos, rl.Vector2Scale(size, .5))

	clr := c.clr
	r := c.radius.Get()
	if rl.CheckCollisionPointCircle(mp, p, r) {
		clr = rl.Color{
			R: 255 - c.clr.R,
			G: 255 - c.clr.G,
			B: 255 - c.clr.B,
			A: c.clr.A,
		}
	}

	rl.DrawCircleV(p, r, clr)
	rl.DrawCircleLinesV(p, r, theme.Charcoal)
}
