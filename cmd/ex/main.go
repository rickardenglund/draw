package main

import (
	"time"

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
	twf := plot.NewPlot(d.GetTWF())
	sp := plot.NewPlot(d.GetSP())
	updatePlotsF := func() {
		d.Update()
		twf.Set(d.GetTWF())
		sp.Set(d.GetSP())
	}
	currentShape := 0
	myShapes := [][]rl.Vector2{
		{
			{20, 10},
			{11, 11},
			{10, 20},
			{20, 20},
			{20, 10},
		},
		{
			{15, 15},
			{10, 10},
			{15, 20},
			{20, 10},
		},
		{
			{15, 15},
			{10, 10},
			{15, 25},
			{20, 10},
			{30, 10},
			{30, 5},
		},
		{
			{5, 15},
			{10, 15},
		},
	}
	plt := plot.NewPlot(myShapes[0])
	nextShapeF := func() {
		currentShape = (currentShape + 1) % len(myShapes)
		plt.Set(myShapes[currentShape])

	}
	blinkables := []draw.Drawable{
		text.NewBlinkablef("hejsan: %d", 1),
		text.NewBlinkablef("hejsan: %d", 2),
		text.NewBlinkablef("hejsan: %d", 3),
		text.NewBlinkablef("hejsan: %d", 4),
	}

	updateLabel := widget.NewClickable(
		text.NewTextf("updated"), updatePlotsF,
	)

	v := views.NewColumnView(
		views.NewRowView(
			views.NewMultiRowView(
				views.MultiItem{Label: text.NewTextf("Plt"), Full: plt},
				views.MultiItem{Label: text.NewTextf("Muave"), Full: shapes.NewSquare(rl.NewVector2(50, 30), theme.Muave)},
				views.MultiItem{Label: text.NewTextf("Salmon"), Full: shapes.NewSquare(rl.NewVector2(20, 50), theme.Salmon)},
				views.MultiItem{Label: text.NewTextf("Green"), Full: shapes.NewSquare(rl.NewVector2(20, 50), rl.Green)},
			),
			sp,
			views.NewMultiColumnView(
				views.MultiItem{Label: updateLabel, Full: views.NewColumnView(
					widget.NewButton("Update Plot", updatePlotsF),
					widget.NewButton("Next Shape", nextShapeF),
					widget.NewButton("blink", blink(blinkables)),

				)},
				views.MultiItem{Label: text.NewTextf("size"), Full: widget.NewButton("New Size", f)},
				views.MultiItem{Label: text.NewTextf("myC"), Full: myC},
			),
			views.NewColumnView(blinkables...),
		),
		views.NewRowView(
			twf,
		),
	)

	draw.NewWindow(rl.NewVector2(900, 550), func() {}, v)
}

func blink(bs []draw.Drawable) func() {
	return func() {
		go func() {
			for range 3 {
				for i := range bs {
					b, ok := bs[i].(*text.Blinkable)
					if ok {
						b.Blink()
					}
					time.Sleep(200 * time.Millisecond)
				}
			}
		}()

	}
}
