package widget

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/draw"
	"github.com/rickardenglund/draw/input"
	"github.com/rickardenglund/draw/theme"
)

// Danger is a button color for destructive actions.
var Danger = rl.NewColor(160, 40, 40, 255)

type Button struct {
	action  func()
	title   string
	color   rl.Color
	outline bool
}

const buttonPadding = float32(5)

func (b Button) GetSize(target rl.Rectangle) rl.Vector2 {
	ms := theme.MeaureText(b.title)
	return rl.NewVector2(ms.X+2*buttonPadding, ms.Y+2*buttonPadding)
}

func (b *Button) Draw(target rl.Rectangle) {
	ms := theme.MeaureText(b.title)

	mp := input.MousePosition()
	padding := buttonPadding

	h := ms.Y + 2*padding
	w := ms.X + 2*padding
	btn := rl.NewRectangle(target.X+target.Width/2-w/2, target.Y+target.Height/2-h/2, w, h)

	hover := rl.CheckCollisionPointRec(mp, btn)
	clr := b.color
	if hover {
		clr = rl.ColorBrightness(clr, .2)
	}

	if hover && rl.IsMouseButtonDown(rl.MouseButtonLeft) {
		clr = rl.ColorBrightness(clr, .5)
	}

	textPos := rl.NewVector2(btn.X+padding, btn.Y+padding)
	if b.outline {
		if hover {
			rl.DrawRectangleRec(btn, rl.ColorAlpha(b.color, .15))
		}
		rl.DrawRectangleLinesEx(btn, padding/2, b.color)
		theme.DrawText(b.title, textPos, b.color)
	} else {
		rl.DrawRectangleRec(btn, clr)
		rl.DrawRectangleLinesEx(btn, padding/2, rl.ColorBrightness(b.color, -.5))
		theme.DrawText(b.title, textPos, theme.Nude)
	}

	if hover && rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
		b.action()
	}
}

func NewButton(title string, action func()) draw.Drawable {
	return &Button{action: action, title: title, color: theme.Charcoal}
}

// NewOutlineButton creates an unfilled button with border and text in clr, e.g. Danger.
func NewOutlineButton(title string, clr rl.Color, action func()) draw.Drawable {
	return &Button{action: action, title: title, color: clr, outline: true}
}
