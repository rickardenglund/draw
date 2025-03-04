package views

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/animated"
	"github.com/rickardenglund/draw/theme"
)

type MultiRowView struct {
	objs    []MultiItem
	active  int
	markerX *animated.Animated
}

func (r MultiRowView) GetSize(target rl.Rectangle) rl.Vector2 {
	return rl.NewVector2(target.Width, target.Height)
}

func (r *MultiRowView) Draw(target rl.Rectangle) {
	if len(r.objs) == 0 {
		return
	}

	r.handleKeys(target)

	panelHeight := float32(0)

	panelPartSize := target.Width / float32(len(r.objs))

	for i, o := range r.objs {
		labelSize := o.Label.GetSize(target)
		x := target.X + panelPartSize*float32(i) + panelPartSize/2 - labelSize.X/2
		p := rl.NewVector2(x, target.Y)
		labelTarget := rl.NewRectangle(p.X, p.Y, labelSize.X, labelSize.Y)
		o.Label.Draw(labelTarget)
		panelHeight = max(panelHeight, labelSize.Y)
	}

	mainTarget := rl.NewRectangle(target.X, target.Y+panelHeight, target.Width, target.Height-panelHeight)
	r.objs[r.active].Full.Draw(mainTarget)

	mpHeight := rl.NewVector2(0, -panelHeight)
	pleft := rl.NewVector2(mainTarget.X, mainTarget.Y)
	pRight := rl.NewVector2(mainTarget.X+mainTarget.Width, mainTarget.Y)
	mp0 := rl.NewVector2(target.X+r.markerX.Get()*panelPartSize, mainTarget.Y)
	mp0U := rl.Vector2Add(mp0, mpHeight)
	mp1 := rl.NewVector2(target.X+(r.markerX.Get()+1)*panelPartSize, mainTarget.Y)
	mp1U := rl.Vector2Add(mp1, mpHeight)
	rl.DrawLineEx(pleft, mp0, 2, theme.Charcoal)
	rl.DrawLineEx(mp0, mp0U, 2, theme.Charcoal)
	rl.DrawLineEx(mp0U, mp1U, 2, theme.Charcoal)
	rl.DrawLineEx(mp1U, mp1, 2, theme.Charcoal)
	rl.DrawLineEx(mp1, pRight, 2, theme.Charcoal)
}

func (r *MultiRowView) handleKeys(target rl.Rectangle) {
	if rl.IsKeyPressed(rl.KeyL) {
		r.active = (r.active + 1) % len(r.objs)
		r.markerX.Set(float32(r.active))
	}

	if rl.IsKeyPressed(rl.KeyH) {
		r.active -= 1
		if r.active < 0 {
			r.active = len(r.objs) - 1
		}
		r.markerX.Set(float32(r.active))
	}
}

func (r *MultiRowView) Add(i MultiItem) {
	r.objs = append(r.objs, i)
}

func NewMultiRowView(objs ...MultiItem) *MultiRowView {
	return &MultiRowView{
		objs:    objs,
		active:  0,
		markerX: animated.NewAnimated(0, .2),
	}
}
