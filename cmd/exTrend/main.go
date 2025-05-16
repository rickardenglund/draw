package main

import (
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/draw"
	"github.com/rickardenglund/draw/plot"
	"github.com/rickardenglund/draw/views"
	"github.com/rickardenglund/draw/widget"
)

func main() {
	ps := []rl.Vector2{
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
		p.ShiftIn(rl.Vector2{X: float32(n), Y: newV}, 10)
		n++
	})

	v := views.NewColumnView(
		btn,
		p,

	)

	draw.NewWindow(rl.NewVector2(900, 550), func() {}, v)
}
