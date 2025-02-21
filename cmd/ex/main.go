package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/draw"
	"github.com/rickardenglund/draw/plot"
	"github.com/rickardenglund/draw/shapes"
	"github.com/rickardenglund/draw/sound"
	"github.com/rickardenglund/draw/theme"
	"github.com/rickardenglund/draw/views"
	"github.com/rickardenglund/draw/widget"
)

func main() {
	d := newData()
	d.update()

	myC := shapes.NewCircle(20, theme.Muave)
	r := float32(10)
	f := func() {
		r += 20
		if r > 80 {
			r = 10
		}

		myC.Set(r)
	}

	v := views.NewColumnView(
		views.NewRowView(
			plot.NewPlot(d.getSP),
			views.NewColumnView(
				views.NewRowView(
					widget.NewButton("Update Plot", d.update),
					widget.NewButton("New Size", f),
					myC,
				),
				sound.NewPlayer(d.getTWF),
			),
		),
		views.NewRowView(
			plot.NewPlot(d.getTWF),
		),
	)

	draw.NewWindow(rl.NewVector2(900, 550), v)
}
