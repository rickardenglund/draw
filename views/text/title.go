package text

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/draw"
	"github.com/rickardenglund/draw/theme"
)

var _ draw.Drawable = &TopBar{}

type TopBar struct {
	topbar draw.Drawable
	main   draw.Drawable
}

func (t TopBar) GetSize(target rl.Rectangle) rl.Vector2 {
	return rl.NewVector2(target.X, target.Y)
}

func (t TopBar) Draw(target rl.Rectangle) {
	s := t.topbar.GetSize(target)
	t.topbar.Draw(rl.NewRectangle(target.X, target.Y, s.X, s.Y))
	childTarget := rl.NewRectangle(target.X, target.Y+s.Y, target.Width, target.Height-s.Y)
	t.main.Draw(childTarget)

	rl.DrawLine(int32(childTarget.X), int32(childTarget.Y), int32(childTarget.X+childTarget.Width), int32(childTarget.Y), theme.Charcoal)
}

func NewTopBar(s draw.Drawable, c draw.Drawable) *TopBar {
	return &TopBar{topbar: s, main: c}
}
