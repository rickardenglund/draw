package main

import (
	"fmt"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/draw"
	"github.com/rickardenglund/draw/plot"
	"github.com/rickardenglund/draw/shapes"
	"github.com/rickardenglund/draw/views"
	"github.com/rickardenglund/draw/widget"
)

type myp struct {
	X float32
	Y float32
}

func (p myp) GetXString() string {
	return fmt.Sprintf("(%.2f)", p.X)
}
func (p myp) GetYString() string {
	return fmt.Sprintf("%.2f", p.Y)
}

func main() {
	ps := []myp{
		{X: -5, Y: 5},
		{-3, 8},
		{-1, 5},
		{0, 0},
		{3, 3},
		{5, 0},
	}
	n := len(ps)
	p := plot.NewPlot(ps)
	btn := widget.NewButton("do", func() {
		newV := float32(time.Now().Nanosecond() % 30)
		p.ShiftIn(myp{X: float32(n), Y: newV}, 10)
		n++
	})

	v := views.NewColumnView(
		shapes.NewCool(20, rl.Red),
		btn,
		p,
	)

	draw.NewWindow(rl.NewVector2(900, 550), func() {}, v)
}
