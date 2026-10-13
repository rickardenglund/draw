package views

import (
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type box struct {
	size rl.Vector2
	got  rl.Rectangle
}

func (b *box) Draw(t rl.Rectangle)             { b.got = t }
func (b *box) GetSize(rl.Rectangle) rl.Vector2 { return b.size }
func newBox(w, h float32) *box                 { return &box{size: rl.NewVector2(w, h)} }
func rect(x, y, w, h float32) rl.Rectangle     { return rl.NewRectangle(x, y, w, h) }

func check(t *testing.T, name string, b *box, want rl.Rectangle) {
	t.Helper()
	if b.got != want {
		t.Errorf("%s: got %+v, want %+v", name, b.got, want)
	}
}

func TestItemLiteralDefaults(t *testing.T) {
	a, b := newBox(10, 10), newBox(10, 10)
	f := NewFlexView(FlexRow).SetAlignItems(AlignStart)
	f.AddItem(Item(a).Grow(1), Item(b))
	f.Draw(rect(0, 0, 100, 50))
	check(t, "grow", a, rect(0, 0, 90, 10))
	check(t, "auto basis + inherited align", b, rect(90, 0, 10, 10))
}

func TestRowReverse(t *testing.T) {
	a, b := newBox(10, 10), newBox(20, 10)
	f := NewFlexView(FlexRowReverse).Add(a, b)
	f.Draw(rect(0, 0, 100, 10))
	check(t, "first at right", a, rect(90, 0, 10, 10))
	check(t, "second left of first", b, rect(70, 0, 20, 10))
}

func TestColumnReverseJustifyEnd(t *testing.T) {
	a := newBox(10, 10)
	f := NewFlexView(FlexColumnReverse).SetJustifyContent(JustifyEnd).Add(a)
	f.Draw(rect(0, 0, 10, 100))
	check(t, "end is top", a, rect(0, 0, 10, 10))
}

func TestWrapReverseKeepsLineOrder(t *testing.T) {
	a, b, c := newBox(40, 10), newBox(40, 10), newBox(40, 10)
	f := NewFlexView(FlexRowReverse).SetWrap(Wrap).SetAlignItems(AlignStart).Add(a, b, c)
	f.Draw(rect(0, 0, 100, 20))
	check(t, "a", a, rect(60, 0, 40, 10))
	check(t, "b", b, rect(20, 0, 40, 10))
	check(t, "c on second line", c, rect(60, 10, 40, 10))
}

func TestWrapDistributesCrossSpace(t *testing.T) {
	a, b := newBox(60, 10), newBox(60, 10)
	f := NewFlexView(FlexRow).SetWrap(Wrap).Add(a, b)
	f.Draw(rect(0, 0, 100, 100))
	check(t, "a", a, rect(0, 0, 60, 50))
	check(t, "b", b, rect(0, 50, 60, 50))

	c := newBox(10, 10)
	single := NewFlexView(FlexRow).SetWrap(Wrap).Add(c)
	single.Draw(rect(0, 0, 100, 100))
	check(t, "single wrapped line fills", c, rect(0, 0, 10, 100))
}
