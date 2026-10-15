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

	// top is the first visible row; scroll animates towards it.
	top    int
	scroll *animated.Animated
	snap   bool
}

// scrollMargin is how many rows to keep between the selection and the list
// edge before scrolling.
const scrollMargin = 2

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
	// Rows get their natural height, stacked from the top. When they don't
	// all fit, the list scrolls to keep the selection away from the edges.
	rowPad := float32(8)
	rowHeight := float32(0)
	for _, o := range c.objs {
		rowHeight = max(rowHeight, o.Label.GetSize(target).Y+2*rowPad)
	}
	c.updateScroll(max(1, int(target.Height/rowHeight)))
	scroll := c.scroll.Get()

	bottom := target.Y + target.Height
	for i, o := range c.objs {
		rowY := target.Y + rowHeight*(float32(i)-scroll)
		labelSize := o.Label.GetSize(target)
		panelWidth = max(panelWidth, labelSize.X+2*markerRadius)
		if rowY < target.Y-0.5 || rowY+rowHeight > bottom+0.5 {
			continue
		}
		labelRect := rl.NewRectangle(target.X+markerRadius, rowY+(rowHeight-labelSize.Y)/2, labelSize.X, labelSize.Y)
		o.Label.Draw(labelRect)
	}

	activeSize := c.objs[c.active].Label.GetSize(target)
	markerMid := target.Y + rowHeight*(c.markerY.Get()-scroll) + rowHeight/2
	mp1Y := rl.Clamp(markerMid-activeSize.Y/2, target.Y, bottom)
	mp2Y := rl.Clamp(markerMid+activeSize.Y/2, target.Y, bottom)

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

// updateScroll moves the first visible row so the active row stays at least
// scrollMargin rows from the edges.
func (c *MultiColumnView) updateScroll(visible int) {
	margin := min(scrollMargin, (visible-1)/2)
	top := c.top
	if c.active < top+margin {
		top = c.active - margin
	}
	if c.active > top+visible-1-margin {
		top = c.active - visible + 1 + margin
	}
	top = max(0, min(top, len(c.objs)-visible))

	if c.snap {
		c.snap = false
		c.top = top
		c.scroll.SetForce(float32(top))
		return
	}
	if top != c.top {
		c.top = top
		c.scroll.Set(float32(top))
	}
}

func (c *MultiColumnView) handleKeys(target rl.Rectangle) {
	// Shift+J/K is left to the active item, e.g. to move within its own list.
	if rl.IsKeyDown(rl.KeyLeftShift) || rl.IsKeyDown(rl.KeyRightShift) {
		return
	}
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
	c.top = 0
	c.scroll.SetForce(0)
}

func (c *MultiColumnView) GetActive() int {
	return c.active
}

func (c *MultiColumnView) SetActive(a int) {
	c.snap = true
	if a >= len(c.objs) {
		c.active = len(c.objs) - 1
		c.markerY.SetForce(float32(c.active))
		return
	}

	c.active = max(a, 0)
	c.markerY.SetForce(float32(c.active))
}

// Top returns the first visible row.
func (c *MultiColumnView) Top() int {
	return c.top
}

// SetTop sets the first visible row without animating, e.g. to keep the
// scroll position when the view is rebuilt. It is clamped on the next Draw.
func (c *MultiColumnView) SetTop(top int) {
	c.top = max(top, 0)
	c.scroll.SetForce(float32(c.top))
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
		scroll:  animated.NewAnimated(0, .2),
	}
}
