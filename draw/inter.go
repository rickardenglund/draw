package draw

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Drawable interface {
	Init()
	Draw(target rl.Rectangle)
}
