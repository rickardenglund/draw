package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/data"
	"github.com/rickardenglund/draw/draw"
	"github.com/rickardenglund/draw/plot"
	"github.com/rickardenglund/draw/shapes"
	"github.com/rickardenglund/draw/sound"
	"github.com/rickardenglund/draw/theme"
	"github.com/rickardenglund/draw/views"
	"github.com/rickardenglund/draw/widget"
)

func main() {
	d := data.NewData()
	d.Update()

	myC := shapes.NewCircle(20, theme.Muave)
	r := float32(10)
	f := func() {
		r += 20
		if r > 80 {
			r = 10
		}

		myC.Set(r)
	}

	player := sound.NewPlayer(d.GetTWF)

	update := func() {
		d.Update()
		player.Play()
	}
	v := views.NewColumnView(
		views.NewRowView(
			plot.NewPlot(d.GetSP),
			views.NewColumnView(
				views.NewRowView(
					widget.NewButton("Update Plot", update),
					widget.NewButton("New Size", f),
					myC,
				),
				player,
			),
		),
		views.NewRowView(
			plot.NewPlot(d.GetTWF),
		),
	)

	draw.NewWindow(rl.NewVector2(900, 550), v)
}
