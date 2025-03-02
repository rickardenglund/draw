package views

import (
	"fmt"
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func TestName(t *testing.T) {
	size := rl.Vector2{100, 200}
	dir := rl.NewVector2(0, 1)

	delta := rl.Vector2DotProduct(dir, size)
	fmt.Printf("d: %.0f\n", delta)

}
