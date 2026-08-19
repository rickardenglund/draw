package text

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/draw"
	"github.com/rickardenglund/draw/theme"
)

var _ draw.Drawable = &TextInput{}

type TextInput struct {
	label    string
	v        string
	callback func(string)
}

func (t *TextInput) GetSize(target rl.Rectangle) rl.Vector2 {
	ms := theme.MeaureTextPad(t.label)
	ms.Y += theme.FontSize * 2
	return ms
}

func (t *TextInput) Draw(target rl.Rectangle) {
	p := rl.NewVector2(target.X, target.Y)
	theme.DrawTextPad(t.label, p, theme.Charcoal)

	if pastePressed() {
		t.v = rl.GetClipboardText()
	}
	ms := theme.MeaureTextPad(t.label)
	tgt := rl.NewRectangle(
		target.X,
		target.Y+theme.FontSize*1.5,
		target.Width,
		min(target.Height-ms.Y, theme.FontSize*2),
	)

	done := gui.TextBox(tgt, &t.v, int(theme.FontSize), true)
	if done {
		t.callback(t.v)
	}
}

func pastePressed() bool {
	ctrl := rl.IsKeyDown(rl.KeyRightSuper) ||
		rl.IsKeyDown(rl.KeyLeftSuper) ||
		rl.IsKeyDown(rl.KeyLeftControl) ||
		rl.IsKeyDown(rl.KeyRightControl)
	return ctrl && rl.IsKeyPressed(rl.KeyV)
}

func NewTextInput(s string, callback func(string)) *TextInput {
	return &TextInput{
		label:    s,
		callback: callback,
	}
}
