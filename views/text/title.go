package text

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/draw"
	"github.com/rickardenglund/draw/theme"
)

var _ draw.Drawable = &Title{}

type Title struct {
	title draw.Drawable
	child draw.Drawable
}

func (t Title) GetSize(target rl.Rectangle) rl.Vector2 {
	return rl.NewVector2(target.X, target.Y)
}

func (t Title) Draw(target rl.Rectangle) {
	s := t.title.GetSize(target)
	t.title.Draw(rl.NewRectangle(target.X, target.Y, s.X, s.Y))
	childTarget := rl.NewRectangle(target.X, target.Y+s.Y, target.Width, target.Height-s.Y)
	t.child.Draw(childTarget)

	rl.DrawLine(int32(childTarget.X), int32(childTarget.Y), int32(childTarget.X+childTarget.Width), int32(childTarget.Y), theme.Charcoal)
}

func NewTitle(s draw.Drawable, c draw.Drawable) *Title {
	return &Title{title: s, child: c}
}
