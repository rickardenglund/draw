package widget

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/draw"
	"github.com/rickardenglund/draw/theme"
)

type Border struct {
	color rl.Color
	child draw.Drawable
}

var _ draw.Drawable = &Border{}

func (b Border) GetSize(target rl.Rectangle) rl.Vector2 {
	cs := b.child.GetSize(target)
	return rl.NewVector2(cs.X, cs.Y)
}

func (b Border) Draw(target rl.Rectangle) {
	rl.DrawRectangleLinesEx(shrink(target, 5), 2, b.color)
	b.child.Draw(shrink(target, 10))
}

func shrink(rect rl.Rectangle, v float32) rl.Rectangle {
	return rl.NewRectangle(rect.X+v, rect.Y+v, rect.Width-2*v, rect.Height-2*v)
}

func NewBorder(child draw.Drawable) *Border {
	return &Border{
		child: child,
		color: theme.Salmon,
	}
}
