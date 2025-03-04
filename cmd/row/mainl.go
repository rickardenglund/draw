package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/draw"
	"github.com/rickardenglund/draw/shapes"
	"github.com/rickardenglund/draw/views"
	"github.com/rickardenglund/draw/views/text"
)

func main() {
	c := views.NewColumnView(
		views.NewRowView(
			text.NewTextf("hej"),
			views.NewMultiRowView(
				views.MultiItem{text.NewTextf("r1"), shapes.NewCircle(15, rl.Red)},
				views.MultiItem{text.NewTextf("r2"), shapes.NewCircle(15, rl.Green)},
				views.MultiItem{text.NewTextf("r4"), shapes.NewCircle(15, rl.Blue)},
			),
		),

		views.NewRowView(
			text.NewTextf("hej"),
			views.NewMultiColumnView(
				views.MultiItem{text.NewTextf("r1"), shapes.NewCircle(15, rl.Red)},
				views.MultiItem{text.NewTextf("r2"), shapes.NewCircle(15, rl.Green)},
				views.MultiItem{text.NewTextf("r4"), shapes.NewCircle(15, rl.Blue)},
			),
		),
	)

	draw.NewWindow(rl.NewVector2(400, 400), func() {}, c)
}
