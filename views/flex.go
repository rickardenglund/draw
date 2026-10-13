package views

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/draw"
)

// FlexDirection mirrors the CSS flex-direction property.
type FlexDirection int

const (
	FlexRow FlexDirection = iota
	FlexRowReverse
	FlexColumn
	FlexColumnReverse
)

// FlexWrap mirrors the CSS flex-wrap property.
type FlexWrap int

const (
	NoWrap FlexWrap = iota
	Wrap
)

// JustifyContent mirrors the CSS justify-content property (main axis alignment).
type JustifyContent int

const (
	JustifyStart JustifyContent = iota
	JustifyEnd
	JustifyCenter
	JustifySpaceBetween
	JustifySpaceAround
	JustifySpaceEvenly
)

// AlignItems mirrors the CSS align-items / align-self properties (cross axis alignment).
type AlignItems int

const (
	// AlignAuto is the zero value. As an item's align-self it inherits the
	// view's AlignItems; as the view's AlignItems it behaves like AlignStretch.
	AlignAuto AlignItems = iota
	AlignStretch
	AlignStart
	AlignEnd
	AlignCenter
)

// FlexItem wraps a Drawable with CSS flexbox-like sizing hints.
// Create one with Item so that it starts from the CSS defaults.
type FlexItem struct {
	drawable  draw.Drawable
	grow      float32
	shrink    float32
	basis     float32
	autoBasis bool
	alignSelf AlignItems
}

// Item returns a FlexItem with CSS defaults: grow 0, shrink 1, basis auto, align-self auto.
func Item(d draw.Drawable) FlexItem {
	return FlexItem{drawable: d, shrink: 1, autoBasis: true}
}

// Grow sets flex-grow.
func (i FlexItem) Grow(g float32) FlexItem { i.grow = g; return i }

// Shrink sets flex-shrink.
func (i FlexItem) Shrink(s float32) FlexItem { i.shrink = s; return i }

// Basis sets a fixed flex-basis on the main axis.
func (i FlexItem) Basis(b float32) FlexItem { i.basis = b; i.autoBasis = false; return i }

// BasisAuto makes the flex-basis come from the drawable's GetSize.
func (i FlexItem) BasisAuto() FlexItem { i.autoBasis = true; return i }

// AlignSelf sets align-self.
func (i FlexItem) AlignSelf(a AlignItems) FlexItem { i.alignSelf = a; return i }

// FlexView lays out its items similar to a CSS flexbox container.
type FlexView struct {
	items     []FlexItem
	direction FlexDirection
	wrap      FlexWrap
	justify   JustifyContent
	align     AlignItems
	gap       float32
}

func NewFlexView(direction FlexDirection) *FlexView {
	return &FlexView{
		direction: direction,
		wrap:      NoWrap,
		justify:   JustifyStart,
		align:     AlignStretch,
	}
}

func (f *FlexView) SetWrap(w FlexWrap) *FlexView                 { f.wrap = w; return f }
func (f *FlexView) SetJustifyContent(j JustifyContent) *FlexView { f.justify = j; return f }
func (f *FlexView) SetAlignItems(a AlignItems) *FlexView         { f.align = a; return f }
func (f *FlexView) SetGap(gap float32) *FlexView                 { f.gap = gap; return f }

// Add appends drawables with default item settings (see Item).
func (f *FlexView) Add(objs ...draw.Drawable) *FlexView {
	for _, o := range objs {
		f.items = append(f.items, Item(o))
	}
	return f
}

// AddItem appends configured FlexItems.
func (f *FlexView) AddItem(items ...FlexItem) *FlexView {
	f.items = append(f.items, items...)
	return f
}

func (f *FlexView) Clear() {
	f.items = nil
}

func (f FlexView) isRow() bool {
	return f.direction == FlexRow || f.direction == FlexRowReverse
}

func (f FlexView) isReverse() bool {
	return f.direction == FlexRowReverse || f.direction == FlexColumnReverse
}

func (f FlexView) GetSize(target rl.Rectangle) rl.Vector2 {
	return rl.NewVector2(target.Width, target.Height)
}

func (f FlexView) mainSize(target rl.Rectangle) float32 {
	if f.isRow() {
		return target.Width
	}
	return target.Height
}

func (f FlexView) crossSize(target rl.Rectangle) float32 {
	if f.isRow() {
		return target.Height
	}
	return target.Width
}

func (f FlexView) itemMain(item FlexItem, target rl.Rectangle) float32 {
	if !item.autoBasis {
		return item.basis
	}
	s := item.drawable.GetSize(target)
	if f.isRow() {
		return s.X
	}
	return s.Y
}

func (f FlexView) itemCross(item FlexItem, target rl.Rectangle) float32 {
	s := item.drawable.GetSize(target)
	if f.isRow() {
		return s.Y
	}
	return s.X
}

type flexLine struct {
	indexes []int
	basis   []float32
}

func (f *FlexView) Draw(target rl.Rectangle) {
	if len(f.items) == 0 {
		return
	}

	mainAvail := f.mainSize(target)
	crossAvail := f.crossSize(target)

	// Split into lines in source order; reverse directions are handled when
	// positioning items, so the first line still holds the first items.
	var lines []flexLine
	cur := flexLine{}
	curMain := float32(0)
	for idx := range f.items {
		basis := f.itemMain(f.items[idx], target)
		add := basis
		if len(cur.indexes) > 0 {
			add += f.gap
		}
		if f.wrap == Wrap && len(cur.indexes) > 0 && curMain+add > mainAvail {
			lines = append(lines, cur)
			cur = flexLine{}
			curMain = 0
			add = basis
		}
		cur.indexes = append(cur.indexes, idx)
		cur.basis = append(cur.basis, basis)
		curMain += add
	}
	if len(cur.indexes) > 0 {
		lines = append(lines, cur)
	}

	// A nowrap container is single-line and its line fills the cross axis.
	// With wrap, each line is as tall as its largest item and leftover cross
	// space is shared equally between lines (align-content: stretch).
	lineCross := make([]float32, len(lines))
	if f.wrap == NoWrap {
		lineCross[0] = crossAvail
	} else {
		total := float32(0)
		for i, line := range lines {
			for _, idx := range line.indexes {
				lineCross[i] = max(lineCross[i], f.itemCross(f.items[idx], target))
			}
			total += lineCross[i]
		}
		if extra := crossAvail - total; extra > 0 {
			for i := range lineCross {
				lineCross[i] += extra / float32(len(lines))
			}
		}
	}

	crossStart := target.X
	if f.isRow() {
		crossStart = target.Y
	}
	for i, line := range lines {
		f.drawLine(target, line, mainAvail, crossStart, lineCross[i])
		crossStart += lineCross[i]
	}
}

func (f *FlexView) drawLine(target rl.Rectangle, line flexLine, mainAvail float32, crossStart float32, lineCross float32) {
	n := len(line.indexes)

	sumBasis := float32(0)
	for _, b := range line.basis {
		sumBasis += b
	}
	gapTotal := f.gap * float32(n-1)
	freeSpace := mainAvail - sumBasis - gapTotal

	growSum := float32(0)
	shrinkWeightSum := float32(0)
	for i, idx := range line.indexes {
		item := f.items[idx]
		growSum += item.grow
		shrinkWeightSum += item.shrink * line.basis[i]
	}

	finalMain := make([]float32, n)
	for i := range line.indexes {
		finalMain[i] = line.basis[i]
	}

	remaining := freeSpace
	if freeSpace > 0 && growSum > 0 {
		for i, idx := range line.indexes {
			item := f.items[idx]
			finalMain[i] += freeSpace * (item.grow / growSum)
		}
		remaining = 0
	} else if freeSpace < 0 && shrinkWeightSum > 0 {
		for i, idx := range line.indexes {
			item := f.items[idx]
			w := item.shrink * line.basis[i]
			finalMain[i] += freeSpace * (w / shrinkWeightSum)
			if finalMain[i] < 0 {
				finalMain[i] = 0
			}
		}
		remaining = 0
	}

	// Determine spacing/offset from justify-content using any leftover space.
	startOffset := float32(0)
	between := f.gap
	if remaining > 0 {
		switch f.justify {
		case JustifyStart:
			// nothing to do
		case JustifyEnd:
			startOffset = remaining
		case JustifyCenter:
			startOffset = remaining / 2
		case JustifySpaceBetween:
			if n > 1 {
				between = f.gap + remaining/float32(n-1)
			}
		case JustifySpaceAround:
			pad := remaining / float32(n)
			startOffset = pad / 2
			between = f.gap + pad
		case JustifySpaceEvenly:
			pad := remaining / float32(n+1)
			startOffset = pad
			between = f.gap + pad
		}
	}

	mainOrigin := target.Y
	if f.isRow() {
		mainOrigin = target.X
	}
	mainPos := startOffset

	for i, idx := range line.indexes {
		item := f.items[idx]
		itemMainSize := finalMain[i]
		itemCrossSize := f.itemCross(item, target)

		align := item.alignSelf
		if align == AlignAuto {
			align = f.align
		}
		if align == AlignAuto {
			align = AlignStretch
		}

		crossPos := crossStart
		crossSizeFinal := itemCrossSize
		switch align {
		case AlignStretch:
			crossSizeFinal = lineCross
			crossPos = crossStart
		case AlignStart:
			crossPos = crossStart
		case AlignEnd:
			crossPos = crossStart + lineCross - itemCrossSize
		case AlignCenter:
			crossPos = crossStart + (lineCross-itemCrossSize)/2
		}

		// Reverse directions mirror the layout so main-start is the right/bottom edge.
		p := mainOrigin + mainPos
		if f.isReverse() {
			p = mainOrigin + mainAvail - mainPos - itemMainSize
		}

		var childTarget rl.Rectangle
		if f.isRow() {
			childTarget = rl.NewRectangle(p, crossPos, itemMainSize, crossSizeFinal)
		} else {
			childTarget = rl.NewRectangle(crossPos, p, crossSizeFinal, itemMainSize)
		}

		item.drawable.Draw(childTarget)

		mainPos += itemMainSize + between
	}
}
