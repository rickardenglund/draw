// Package input wraps raylib input so it can be overridden, e.g. to render
// previews with a simulated mouse position.
package input

import (
	"fmt"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var mouseOverride *rl.Vector2

// The DRAW_MOUSE env var ("x,y") sets an initial mouse override.
func init() {
	v := os.Getenv("DRAW_MOUSE")
	if v == "" {
		return
	}

	var x, y float32
	if _, err := fmt.Sscanf(v, "%f,%f", &x, &y); err != nil {
		fmt.Printf("ignoring bad DRAW_MOUSE %q: %s\n", v, err)
		return
	}
	SetMousePosition(rl.NewVector2(x, y))
}

// MousePosition returns the overridden mouse position if set, otherwise the
// real one.
func MousePosition() rl.Vector2 {
	if mouseOverride != nil {
		return *mouseOverride
	}

	return rl.GetMousePosition()
}

// SetMousePosition overrides the position returned by MousePosition.
func SetMousePosition(p rl.Vector2) {
	mouseOverride = &p
}

// ClearMousePosition removes the mouse position override.
func ClearMousePosition() {
	mouseOverride = nil
}
