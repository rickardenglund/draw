package draw

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Drawable interface {
	Draw(target rl.Rectangle)
	GetSize(target rl.Rectangle) rl.Vector2
}
