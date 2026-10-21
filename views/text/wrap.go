package text

import (
	"strings"
	"unicode/utf8"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/draw"
	"github.com/rickardenglund/draw/theme"
)

func wrap(s string, size rl.Vector2) string {
	return wrapWith(s, size, theme.MeaureTextPad)
}

// wrapWith breaks s into lines no wider than size.X, at spaces where it
// can; a word wider than a line is split, after punctuation if possible.
// Lines that don't fit in size.Y are left out, but the first one is kept.
func wrapWith(s string, size rl.Vector2, measure func(string) rl.Vector2) string {
	if measure(s).X <= size.X {
		return s
	}
	fits := func(line string) bool { return measure(line).X <= size.X }

	var lines []string
	for _, para := range strings.Split(s, "\n") {
		before := len(lines)
		line := ""
		for _, w := range strings.Fields(para) {
			if line != "" {
				if fits(line + " " + w) {
					line += " " + w
					continue
				}
				lines = append(lines, line)
			}
			// Split a word too wide for a line of its own.
			for !fits(w) {
				cut := longestFit(w, fits)
				lines = append(lines, w[:cut])
				w = w[cut:]
			}
			line = w
		}
		if line != "" || len(lines) == before {
			lines = append(lines, line)
		}
	}

	n := 1
	for n < len(lines) && measure(strings.Join(lines[:n+1], "\n")).Y <= size.Y {
		n++
	}
	return strings.Join(lines[:n], "\n")
}

// longestFit is where to cut w so that w[:cut] fits: after the last
// punctuation in the second half of what fits, otherwise as late as
// possible, but always after at least one character.
func longestFit(w string, fits func(string) bool) int {
	cut := 0
	for i := range w {
		if i > 0 && !fits(w[:i]) {
			break
		}
		cut = i
	}
	if cut == 0 {
		_, size := utf8.DecodeRuneInString(w)
		return size
	}
	if i := strings.LastIndexAny(w[:cut], "./:-_"); i+1 >= cut/2 {
		return i + 1
	}
	return cut
}

var _ draw.Drawable = &Wrapped{}

type Wrapped struct {
	s string
}

func (w Wrapped) Draw(target rl.Rectangle) {
	s := wrap(w.s, rl.NewVector2(target.Width, target.Height))
	NewTextf(s).Draw(target)
}

func (w Wrapped) GetSize(target rl.Rectangle) rl.Vector2 {
	s := wrap(w.s, rl.NewVector2(target.Width, target.Height))
	return NewTextf(s).GetSize(target)
}

func NewWrap(s string) *Wrapped {
	return &Wrapped{s: s}
}
