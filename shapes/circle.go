package shapes

import (
	ea "github.com/gen2brain/raylib-go/easings"
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/theme"
)

type Circle struct {
	radius    float32
	clr       rl.Color
	prevRad   float32
	animStart float64
}

func (c *Circle) Init() {
}

func NewCircle(radius float32, clr rl.Color) *Circle {
	return &Circle{
		radius:    radius,
		clr:       clr,
		prevRad:   0,
		animStart: rl.GetTime(),
	}
}

func (c *Circle) Set(newR float32) {
	c.prevRad = c.radius
	c.radius = newR
	c.animStart = rl.GetTime()
}

func (c *Circle) Draw(target rl.Rectangle) {
	pos := rl.NewVector2(target.X, target.Y)
	size := rl.NewVector2(target.Width, target.Height)
	mp := rl.GetMousePosition()
	p := rl.Vector2Add(pos, rl.Vector2Scale(size, .5))

	clr := c.clr
	if rl.CheckCollisionPointCircle(mp, p, c.radius) {
		clr = rl.Color{
			R: 255 - c.clr.R,
			G: 255 - c.clr.G,
			B: 255 - c.clr.B,
			A: c.clr.A,
		}
	}

	now := float32(rl.GetTime() - c.animStart)
	dur := float32(.2)
	r := ea.QuadIn(now, c.prevRad, c.radius-c.prevRad, dur)
	if now > dur {
		r = c.radius
	}

	rl.DrawCircleV(p, r, clr)
	rl.DrawCircleLinesV(p, r, theme.Charcoal)
}
