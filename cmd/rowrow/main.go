package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/data"
	"github.com/rickardenglund/draw/draw"
	"github.com/rickardenglund/draw/views"
	"github.com/rickardenglund/draw/views/text"
	"github.com/rickardenglund/draw/widget"
)

func main() {
	d := data.NewData()
	d.Update()

	b :=
		views.NewRowView(
			views.NewColumnView(
				widget.NewBorder(
					views.NewRowView(
						widget.NewButton("btn", func() {}),
						widget.NewButton("btn", func() {}),
						widget.NewButton("btn", func() {}),
					),
				),
				widget.NewBorder(
					text.NewTextf("colasldfjasldfjalsdfjalsdfjalsdfjalsdfjalsdfjalsdfjalsdfjalsdfj alsdfjalsdfjalsdfjalsdfjalsdfjalsdfjalsdfjalsdfjalsdfjalsdfj"),
				),
			),
			text.NewTextf("row"),
		)

	draw.NewWindow(rl.NewVector2(900, 550), func() {}, b)
}
