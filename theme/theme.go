package theme

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Dreams
//var (
//	Muave = rl.NewColor(224, 176, 255, 255)
//	Lavender = rl.NewColor(230, 230, 250, 255)
//	Black = rl.Black
//	Orchid = rl.NewColor(218, 112, 214, 255)
//)

func init() {
	Font = rl.GetFontDefault()
}

var (
	Muave    = rl.NewColor(224, 176, 255, 255)
	Salmon   = rl.NewColor(250, 128, 114, 255)
	Charcoal = rl.NewColor(54, 69, 79, 255)
	Nude     = rl.NewColor(227, 188, 154, 255)
	Font     = rl.GetFontDefault()
)
