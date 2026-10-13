package widget

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/draw"
	"github.com/rickardenglund/draw/theme"
)

const (
	borderLineInset  = 5
	borderChildInset = 10
)

type Border struct {
	color rl.Color
	child draw.Drawable
}

var _ draw.Drawable = &Border{}

func (b Border) GetSize(target rl.Rectangle) rl.Vector2 {
	cs := b.child.GetSize(shrink(target, borderChildInset))
	return rl.NewVector2(cs.X+2*borderChildInset, cs.Y+2*borderChildInset)
}

func (b Border) Draw(target rl.Rectangle) {
	rl.DrawRectangleLinesEx(shrink(target, borderLineInset), 2, b.color)
	b.child.Draw(shrink(target, borderChildInset))
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
