package text

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/draw"
	"github.com/rickardenglund/draw/theme"
)

var _ draw.Drawable = &Text{}

type Text struct {
	s string
}

func (t Text) Draw(target rl.Rectangle) {
	p := rl.NewVector2(target.X, target.Y)
	theme.DrawTextPad(t.s, p, theme.Charcoal)
}

func NewTextf(s string, args ...any) *Text {
	return &Text{s: fmt.Sprintf(s, args...)}
}
