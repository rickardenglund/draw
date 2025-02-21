package views

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/draw"
	"github.com/rickardenglund/draw/theme"
)

type RowView struct {
	objs []draw.Drawable
}

func (r RowView) Draw(target rl.Rectangle) {
	pos := rl.NewVector2(target.X, target.Y)
	size := rl.NewVector2(target.Width, target.Height)
	dx := size.X / float32(len(r.objs))
	partSize := rl.NewVector2(dx, size.Y)
	offset := rl.NewVector2(partSize.X, 0)

	for i := range r.objs {
		p := rl.Vector2Add(pos, rl.Vector2Scale(offset, float32(i)))

		rect := rl.NewRectangle(p.X, p.Y, partSize.X, partSize.Y)
		rl.DrawRectangleLinesEx(rect, 2, rl.ColorAlpha(theme.Charcoal, .5))

		margin := float32(0)
		pInBorder := rl.Vector2Add(p, rl.NewVector2(margin, margin))
		sizeInBorder := rl.Vector2Subtract(partSize, rl.NewVector2(margin*2, margin*2))

		tar := rl.NewRectangle(pInBorder.X, pInBorder.Y, sizeInBorder.X, sizeInBorder.Y)
		r.objs[i].Draw(tar)
	}
}

func (r RowView) Init() {
	for _, o := range r.objs {
		o.Init()
	}
}

func NewRowView(objs ...draw.Drawable) draw.Drawable {
	return RowView{objs: objs}
}
