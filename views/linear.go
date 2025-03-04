package views

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/draw"
)

type LinearView struct {
	objs []draw.Drawable
	dir  rl.Vector2
}

func (v LinearView) GetSize(target rl.Rectangle) rl.Vector2 {
	return rl.NewVector2(target.Width, target.Height)
}

func (v LinearView) Draw(target rl.Rectangle) {
	childSizeUni := rl.Vector2Scale(
		v.GetSize(target),
		1/float32(len(v.objs)),
	)
	ortho := rl.NewVector2(v.dir.Y, v.dir.X)
	orhoSize := rl.Vector2Multiply(ortho, v.GetSize(target))

	step := rl.Vector2Multiply(v.dir, childSizeUni)
	p := rl.NewVector2(target.X, target.Y)
	for i := range v.objs {
		//childSize := r.objs[i].GetSize(target) // should not be not target
		childSize := rl.NewVector2(
			max(orhoSize.X, childSizeUni.X),
			max(orhoSize.Y, childSizeUni.Y),
		)
		childTarget := rl.NewRectangle(p.X, p.Y, childSize.X, childSize.Y)

		v.objs[i].Draw(childTarget)

		p = rl.Vector2Add(p, step)
	}
}

func NewRowView(objs ...draw.Drawable) draw.Drawable {
	return LinearView{objs: objs, dir: rl.NewVector2(1, 0)}
}

func NewColumnView(objs ...draw.Drawable) *LinearView {
	return &LinearView{objs: objs, dir: rl.NewVector2(0, 1)}
}
