package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/draw"
	"github.com/rickardenglund/draw/shapes"
	"github.com/rickardenglund/draw/views"
	"github.com/rickardenglund/draw/views/text"
	"github.com/rickardenglund/draw/widget"
)

func main() {
	// Top row: three growable boxes that share the extra space 1:2:1,
	// wrapped in space-between justification with a gap.
	grow := views.NewFlexView(views.FlexRow).
		SetJustifyContent(views.JustifyStart).
		//SetAlignItems(views.AlignStretch).
		SetGap(10)
	grow.AddItem(views.Item(widget.NewBorder(shapes.NewSquare(rl.NewVector2(40, 40), rl.Red))).Grow(1))
	grow.AddItem(views.Item(widget.NewBorder(shapes.NewSquare(rl.NewVector2(10, 40), rl.Green))).Grow(2))
	grow.AddItem(views.Item(shapes.NewSquare(rl.NewVector2(40, 40), rl.Blue)).Grow(1))

	// Middle row: fixed-size circles with different align-self values,
	// centered on the main axis.
	align := views.NewFlexView(views.FlexRow).
		SetJustifyContent(views.JustifyCenter).
		SetAlignItems(views.AlignCenter).
		SetGap(20)
	align.AddItem(views.Item(shapes.NewCircle(15, rl.Purple)).AlignSelf(views.AlignStart))
	align.AddItem(views.Item(shapes.NewCircle(15, rl.Orange)).AlignSelf(views.AlignCenter))
	align.AddItem(views.Item(text.NewTextf("kanin")).AlignSelf(views.AlignStart))
	align.AddItem(views.Item(shapes.NewCircle(15, rl.SkyBlue)).AlignSelf(views.AlignEnd))

	// Bottom row: several fixed-width labels that wrap onto new lines
	// once they no longer fit the available width.
	wrap := views.NewFlexView(views.FlexRow).
		SetWrap(views.Wrap).
		SetJustifyContent(views.JustifySpaceEvenly).
		SetGap(8)
	for _, s := range []string{"one", "two", "three", "four", "five", "six", "seven", "eight"} {
		wrap.Add(text.NewTextf(s))
	}
	for _, s := range []string{"one", "two", "three", "four", "five", "six", "seven", "eight"} {
		wrap.Add(text.NewTextf(s))
	}

	c := views.NewColumnView(grow, align, wrap, widget.NewButton("hek", func() {
		wrap.Add(widget.NewBorder(text.NewTextf("new")))

	}))

	draw.NewWindow(rl.NewVector2(500, 400), func() {}, c)
}
