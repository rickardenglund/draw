package text

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/draw"
	"github.com/rickardenglund/draw/theme"
)

type Blinkable struct {
	s       string
	blinked float64
}

func (b *Blinkable) Draw(target rl.Rectangle) {
	p := rl.NewVector2(target.X, target.Y)
	clr := theme.Charcoal
	delta := rl.GetTime() - b.blinked
	dur := .2
	if delta < dur {
		v := 1 - delta/dur

		clr = rl.ColorBrightness(clr, float32(v*.5))
	}

	theme.DrawTextPad(b.s, p, clr)
}

func (b *Blinkable) GetSize(target rl.Rectangle) rl.Vector2 {
	return theme.MeaureTextPad(b.s)
}

func (b *Blinkable) Blink() {
	b.blinked = rl.GetTime()
}

var _ draw.Drawable = new(Blinkable)

func NewBlinkablef(s string, args ...any) *Blinkable {
	b := &Blinkable{
		s:       fmt.Sprintf(s, args...),
		blinked: -1000000,
	}

	return b
}
