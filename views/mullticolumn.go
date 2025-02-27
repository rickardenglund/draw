package views

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/draw"
	"github.com/rickardenglund/draw/theme"
)

type MultiColumnView struct {
	objs   []MultiItem
	active int
}

func (c *MultiColumnView) Init() {
	for _, o := range c.objs {
		o.Full.Init()
	}
}

func (c *MultiColumnView) Draw(target rl.Rectangle) {
	c.handleKeys()
	panelWidth := float32(0)

	panelPartSize := target.Height / float32(len(c.objs))
	for i, o := range c.objs {
		y := panelPartSize*float32(i) + panelPartSize/2
		markerRadius := float32(5)
		markerPos := rl.NewVector2(target.X+markerRadius, target.Y+y)
		fontsize := float32(20)
		spacing := float32(0.5)

		label := o.Label
		ms := rl.MeasureTextEx(theme.Font, label, fontsize, spacing)
		textpos := rl.NewVector2(markerPos.X+markerRadius, markerPos.Y-ms.Y/2)
		rl.DrawTextEx(theme.Font, label, textpos, fontsize, spacing, theme.Charcoal)

		panelWidth = max(panelWidth, ms.X+2*markerRadius)

		if i == c.active {
			rl.DrawCircleV(markerPos, markerRadius, theme.Charcoal)
		}
	}

	mainTarget := rl.NewRectangle(target.X+panelWidth, target.Y, target.Width-panelWidth, target.Height)
	//panelTarget := rl.NewRectangle(target.X, target.Y, panelWidth, target.Height)
	c.objs[c.active].Full.Draw(mainTarget)
}

func (c *MultiColumnView) handleKeys() {
	if rl.IsKeyPressed(rl.KeyJ) {
		c.active = (c.active + 1) % len(c.objs)
	}
	if rl.IsKeyPressed(rl.KeyK) {
		c.active -= 1
		if c.active < 0 {
			c.active = len(c.objs) - 1
		}
	}

}

type MultiItem struct {
	Label string
	Full  draw.Drawable
}

func NewMultiColumnView(objs ...MultiItem) *MultiColumnView {
	return &MultiColumnView{
		objs: objs, active: 0,
	}
}
