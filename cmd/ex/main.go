package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"draw/draw"
	"draw/plot"
	"draw/shapes"
	"draw/sound"
	"draw/theme"
	"draw/views"
	"draw/widget"
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
					widget.NewButton(d.update),
					widget.NewButton(f),
					myC,
				),
				sound.NewPlayer(d.getTWF),
			),
		),
		views.NewRowView(
			plot.NewPlot(d.getTWF),
		),
	)

	draw.NewWindow(rl.NewVector2(800, 400), v)
}
