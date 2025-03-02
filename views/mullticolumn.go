package views

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/animated"
	"github.com/rickardenglund/draw/draw"
	"github.com/rickardenglund/draw/theme"
)

type MultiColumnView struct {
	objs    []MultiItem
	active  int
	markerY *animated.Animated
}

func (c MultiColumnView) GetSize(target rl.Rectangle) rl.Vector2 {
	return rl.NewVector2(target.X, target.Y)
}

func (c *MultiColumnView) Draw(target rl.Rectangle) {
	if len(c.objs) == 0 {
		return
	}

	c.handleKeys(target)

	panelWidth := float32(10)

	markerRadius := float32(5)
	panelPartSize := target.Height / float32(len(c.objs))
	mawrkerPos := rl.NewVector2(
		target.X+markerRadius,
		target.Y+c.markerY.Get()*panelPartSize+panelPartSize/2,
	)

	for i, o := range c.objs {
		y := target.Y + panelPartSize*float32(i) + panelPartSize/2
		p := rl.NewVector2(target.X+markerRadius, y)
		labelSize := o.Label.GetSize(target)
		labelRect := rl.NewRectangle(p.X, p.Y-labelSize.Y/2, labelSize.X, labelSize.Y)
		o.Label.Draw(labelRect)

		panelWidth = max(panelWidth, labelSize.X+2*markerRadius)
	}

	rl.DrawCircleV(mawrkerPos, markerRadius, theme.Charcoal)

	mainTarget := rl.NewRectangle(target.X+panelWidth, target.Y, target.Width-panelWidth, target.Height)
	c.objs[c.active].Full.Draw(mainTarget)
}

func (c *MultiColumnView) handleKeys(target rl.Rectangle) {
	if rl.IsKeyPressed(rl.KeyJ) {
		c.active = (c.active + 1) % len(c.objs)
		c.markerY.Set(float32(c.active))
	}

	if rl.IsKeyPressed(rl.KeyK) {
		c.active -= 1
		if c.active < 0 {
			c.active = len(c.objs) - 1
		}
		c.markerY.Set(float32(c.active))
	}

}

func (c *MultiColumnView) Add(i MultiItem) {
	c.objs = append(c.objs, i)
}

type MultiItem struct {
	Label draw.Drawable
	Full  draw.Drawable
}

func NewMultiColumnView(objs ...MultiItem) *MultiColumnView {
	return &MultiColumnView{
		objs:    objs,
		active:  0,
		markerY: animated.NewAnimated(0, .2),
	}
}
