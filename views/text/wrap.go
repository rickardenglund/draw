package text

import (
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/draw"
	"github.com/rickardenglund/draw/theme"
)

func wrap(s string, size rl.Vector2) string {
	ms := theme.MeaureTextPad(s)
	if ms.X < size.X {
		return s
	}

	sb := strings.Builder{}
	words := strings.Split(s, " ")
	first := true
	for _, w := range words {
		ms := theme.MeaureTextPad(sb.String() + " " + w)
		if ms.Y > size.Y {
			break
		}

		addingLine := ms.X > size.X
		if addingLine {
			sb.WriteString("\n")
		}
		if !addingLine && !first {
			sb.WriteString(" ")
		}
		sb.WriteString(w)
		first = false
	}

	return sb.String()
}

var _ draw.Drawable = &Wrapped{}

type Wrapped struct {
	s string
}

func (w Wrapped) Draw(target rl.Rectangle) {
	s := wrap(w.s, rl.NewVector2(target.Width, target.Height))
	theme.DrawTextPad(s, rl.NewVector2(target.X, target.Y), theme.Charcoal)
}

func (w Wrapped) GetSize(target rl.Rectangle) rl.Vector2 {
	return rl.NewVector2(target.Width, target.Height)
}

func NewWrap(s string) *Wrapped {
	return &Wrapped{s: s}
}
