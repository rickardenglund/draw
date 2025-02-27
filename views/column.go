package views

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/draw"
)

type ColumnView struct {
	objs []draw.Drawable
}

func (c *ColumnView) Draw(target rl.Rectangle) {
	pos := rl.NewVector2(target.X, target.Y)
	size := rl.NewVector2(target.Width, target.Height)

	dy := size.Y / float32(len(c.objs))
	partSize := rl.NewVector2(size.X, dy)
	offset := rl.NewVector2(0, partSize.Y)

	for i := range c.objs {
		p := rl.Vector2Add(pos, rl.Vector2Scale(offset, float32(i)))

		margin := float32(0)
		pInBorder := rl.Vector2Add(p, rl.NewVector2(margin, margin))
		sizeInBorder := rl.Vector2Subtract(partSize, rl.NewVector2(margin*2, margin*2))

		tar := rl.NewRectangle(pInBorder.X, pInBorder.Y, sizeInBorder.X, sizeInBorder.Y)
		c.objs[i].Draw(tar)
	}
}

func NewColumnView(objs ...draw.Drawable) *ColumnView {
	return &ColumnView{objs: objs}
}
