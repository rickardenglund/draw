package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/draw"
	"github.com/rickardenglund/draw/groups"
	"github.com/rickardenglund/draw/shapes"
	"github.com/rickardenglund/draw/theme"
	"github.com/rickardenglund/draw/views"
	"github.com/rickardenglund/draw/views/text"
)

func main() {
	m := map[groups.ID]groups.Node{
		"c5_Lile": {Parent: "qul"},

		"r2": {Parent: ""},
		"c4": {Parent: "r2"},

		"c":  {Parent: ""},
		"c1": {Parent: "c"},
		"c2": {Parent: "c1"},
		"c3": {Parent: "c2"},
	}
	for k, v := range m {
		d := views.NewRowView(
			text.NewTextf("%s", k),
			shapes.NewSquare(rl.NewVector2(20, 20), theme.Muave),
		)

		v.Draw = d //shapes.NewCircle(20, theme.Salmon)
		m[k] = v
	}
	g := groups.NewGroup(m)
	v := views.NewColumnView(
		text.NewTextf("hej"),
		views.NewRowView(
			g,
		),
	)

	draw.NewWindow(rl.NewVector2(900, 550), func() {}, v)
}
