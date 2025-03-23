package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/data"
	"github.com/rickardenglund/draw/draw"
	"github.com/rickardenglund/draw/plot"
	"github.com/rickardenglund/draw/shapes"
	"github.com/rickardenglund/draw/views"
	"github.com/rickardenglund/draw/views/text"
	"github.com/rickardenglund/draw/waves"
)

func main() {
	sampleRate := float64(100)
	d := data.NewData()
	w := waves.Add(
		waves.GetSine(13, 3, 0, 0, sampleRate, 100),
		waves.GetSine(19, 6, 0, 0, sampleRate, 100),
	)

	d.Set(w.Vec(), float32(sampleRate))

	pWf := plot.NewPlot(d.GetTWF())
	pSp := plot.NewPlot(d.GetSP())

	pSp.AddMarker(21)

	c := views.NewColumnView(
		pWf, pSp,
	)
	draw.NewWindow(rl.NewVector2(400, 400), func() {}, c)
}

//Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
func main2() {
	c := views.NewColumnView(
		views.NewRowView(
			text.NewWrap("Hej jag heter Rickard vad heter du? Vilken är din favoritglass? Jag gillara jorggubbsglass! Men jag tycker också att det är fantastiskt med banan i glass"),
			views.NewMultiRowView(
				views.MultiItem{text.NewTextf("r1"), shapes.NewCircle(15, rl.Red)},
				views.MultiItem{text.NewTextf("Kanin"), shapes.NewCircle(15, rl.Green)},
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
