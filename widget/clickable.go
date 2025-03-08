package widget

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/animated"
	"github.com/rickardenglund/draw/draw"
	"github.com/rickardenglund/draw/theme"
)

type ClickAble struct {
	o       draw.Drawable
	onClick func()

	f *animated.Animated
}

func (c *ClickAble) Draw(target rl.Rectangle) {
	mp := rl.GetMousePosition()
	if rl.CheckCollisionPointRec(mp, target) {
		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			c.onClick()
			c.f.SetForce(.2)
			c.f.Set(0)
		}
	}

	rl.DrawRectangleRec(
		target,
		rl.ColorBrightness(theme.Nude, c.f.Get()),
	)
	c.o.Draw(target)
}

func (c *ClickAble) GetSize(target rl.Rectangle) rl.Vector2 {
	return c.o.GetSize(target)
}

var _ draw.Drawable = &ClickAble{}

func NewClickable(o draw.Drawable, f func()) *ClickAble {
	return &ClickAble{
		o:       o,
		onClick: f,
		f:       animated.NewAnimated(0, .2),
	}

}
