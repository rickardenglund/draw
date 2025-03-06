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
	return rl.NewVector2(target.Width, target.Height)
}

func (c *MultiColumnView) Draw(target rl.Rectangle) {
	if len(c.objs) == 0 {
		return
	}

	c.handleKeys(target)

	panelWidth := float32(10)

	markerRadius := float32(5)
	panelPartSize := target.Height / float32(len(c.objs))
	var mp1Y, mp2Y float32
	for i, o := range c.objs {
		y := target.Y + panelPartSize*float32(i) + panelPartSize/2
		p := rl.NewVector2(target.X+markerRadius, y)
		labelSize := o.Label.GetSize(target)
		labelRect := rl.NewRectangle(p.X, p.Y-labelSize.Y/2, labelSize.X, labelSize.Y)
		o.Label.Draw(labelRect)

		panelWidth = max(panelWidth, labelSize.X+2*markerRadius)
		if i == c.active {
			mp1Y = target.Y + panelPartSize*c.markerY.Get() + .5*panelPartSize - .5*labelSize.Y
			mp2Y = target.Y + panelPartSize*c.markerY.Get() + .5*panelPartSize + .5*labelSize.Y

		}
	}

	mpX := target.X + panelWidth
	mp0 := rl.NewVector2(mpX, target.Y)
	mp1 := rl.NewVector2(mpX, mp1Y)
	mp1Left := rl.NewVector2(target.X, mp1Y)
	mp2Left := rl.NewVector2(target.X, mp2Y)
	mp2 := rl.NewVector2(mpX, mp2Y)
	mp3 := rl.NewVector2(target.X+panelWidth, target.Y+target.Height)

	rl.DrawLineEx(mp0, mp1, 2, theme.Charcoal)
	rl.DrawLineEx(mp1, mp1Left, 2, theme.Charcoal)
	rl.DrawLineEx(mp1Left, mp2Left, 2, theme.Charcoal)
	rl.DrawLineEx(mp2Left, mp2, 2, theme.Charcoal)
	rl.DrawLineEx(mp2, mp3, 2, theme.Charcoal)

	mainTarget := rl.NewRectangle(target.X+panelWidth, target.Y, target.Width-panelWidth, target.Height)
	c.objs[c.active].Full.Draw(mainTarget)
}

func (c *MultiColumnView) handleKeys(target rl.Rectangle) {
	if rl.IsKeyPressed(rl.KeyJ) || rl.IsKeyPressed(rl.KeyDown) {
		c.active = (c.active + 1) % len(c.objs)
		c.markerY.Set(float32(c.active))
	}

	if rl.IsKeyPressed(rl.KeyK) || rl.IsKeyPressed(rl.KeyUp) {
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

func (c *MultiColumnView) Clear() {
	c.objs = []MultiItem{}
	c.active = 0
	c.markerY.Set(0)
}

func (c *MultiColumnView) GetActive() int {
	return c.active
}

func (c *MultiColumnView) SetActive(a int) {
	if a >= len(c.objs) {
		c.active = len(c.objs) - 1
		c.markerY.SetForce(float32(c.active))
		return
	}

	c.active = max(a, 0)
	c.markerY.SetForce(float32(c.active))
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
