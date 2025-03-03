package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/data"
	"github.com/rickardenglund/draw/draw"
	"github.com/rickardenglund/draw/plot"
	"github.com/rickardenglund/draw/shapes"
	"github.com/rickardenglund/draw/theme"
	"github.com/rickardenglund/draw/views"
	"github.com/rickardenglund/draw/views/text"
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
	b := text.NewBlinkablef("hejsan: %d", 5)
	v := views.NewColumnView(
		views.NewRowView(
			views.NewMultiRowView(
				views.MultiItem{text.NewTextf("Muave"), shapes.NewSquare(rl.NewVector2(50, 30), theme.Muave)},
				views.MultiItem{text.NewTextf("Salmon"), shapes.NewSquare(rl.NewVector2(20, 50), theme.Salmon)},
				views.MultiItem{text.NewTextf("Green"), shapes.NewSquare(rl.NewVector2(20, 50), rl.Green)},
			),
			plot.NewPlot(d.GetSP),
			views.NewMultiColumnView(
				views.MultiItem{Label: text.NewTextf("update"), Full: widget.NewButton("Update Plot", d.Update)},
				views.MultiItem{Label: text.NewTextf("size"), Full: widget.NewButton("New Size", f)},
				views.MultiItem{Label: text.NewTextf("myC"), Full: myC},
			),
			views.NewColumnView(
				widget.NewButton("blink", b.Blink),
				b,
				text.NewTextf("hoppsan"),
				text.NewTextf("Kalle kanin\när 🐇"),
			),
		),
		views.NewRowView(
			plot.NewPlot(d.GetTWF),
		),
	)

	draw.NewWindow(rl.NewVector2(900, 550), func() {}, v)
}
