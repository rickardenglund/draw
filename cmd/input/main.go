package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/draw"
	"github.com/rickardenglund/draw/views"
	"github.com/rickardenglund/draw/views/text"
)

func main() {
	c := views.NewColumnView(
		text.NewTextf("hej"),
	)
	c.Add(
		text.NewTextInput("lbl2", func(v string) {
			c.Add(text.NewTextf(v))
		}))
	v := text.NewTitle(text.NewTextf("my title"), c)
	draw.NewWindow(rl.NewVector2(400, 400), func() {}, v)
}
