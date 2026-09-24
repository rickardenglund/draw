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
	AlignStretch AlignItems = iota
	AlignStart
	AlignEnd
	AlignCenter
	// AlignAuto (only valid as an item's AlignSelf) means "inherit the view's AlignItems".
	AlignAuto
)

// FlexItem wraps a Drawable with CSS flexbox-like sizing hints.
type FlexItem struct {
	Drawable draw.Drawable
	Grow     float32    // flex-grow
	Shrink   float32    // flex-shrink
	Basis    float32    // flex-basis on the main axis, -1 means "auto" (use GetSize)
	AlignSelf AlignItems // align-self, defaults to AlignAuto (inherit view's AlignItems)
}

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

// Add appends items with default grow/shrink/basis (0, 1, auto).
func (f *FlexView) Add(objs ...draw.Drawable) *FlexView {
	for _, o := range objs {
		f.items = append(f.items, FlexItem{Drawable: o, Grow: 0, Shrink: 1, Basis: -1, AlignSelf: AlignAuto})
	}
	return f
}

// AddItem appends a fully configured FlexItem.
func (f *FlexView) AddItem(item FlexItem) *FlexView {
	f.items = append(f.items, item)
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
	if item.Basis >= 0 {
		return item.Basis
	}
	s := item.Drawable.GetSize(target)
	if f.isRow() {
		return s.X
	}
	return s.Y
}

func (f FlexView) itemCross(item FlexItem, target rl.Rectangle) float32 {
	s := item.Drawable.GetSize(target)
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

	order := make([]int, len(f.items))
	for i := range order {
		order[i] = i
	}
	if f.isReverse() {
		for l, r := 0, len(order)-1; l < r; l, r = l+1, r-1 {
			order[l], order[r] = order[r], order[l]
		}
	}

	// Split into lines according to wrap.
	var lines []flexLine
	cur := flexLine{}
	curMain := float32(0)
	for _, idx := range order {
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

	crossStart := float32(0)
	if f.isRow() {
		crossStart = target.Y
	} else {
		crossStart = target.X
	}

	for _, line := range lines {
		lineCross := crossAvail
		if len(lines) > 1 {
			lineCross = 0
			for _, idx := range line.indexes {
				lineCross = max(lineCross, f.itemCross(f.items[idx], target))
			}
		}

		f.drawLine(target, line, mainAvail, crossStart, lineCross)

		crossStart += lineCross
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
		growSum += item.Grow
		shrinkWeightSum += item.Shrink * line.basis[i]
	}

	finalMain := make([]float32, n)
	for i := range line.indexes {
		finalMain[i] = line.basis[i]
	}

	remaining := freeSpace
	if freeSpace > 0 && growSum > 0 {
		for i, idx := range line.indexes {
			item := f.items[idx]
			finalMain[i] += freeSpace * (item.Grow / growSum)
		}
		remaining = 0
	} else if freeSpace < 0 && shrinkWeightSum > 0 {
		for i, idx := range line.indexes {
			item := f.items[idx]
			w := item.Shrink * line.basis[i]
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

	mainPos := startOffset
	if f.isRow() {
		mainPos += target.X
	} else {
		mainPos += target.Y
	}

	for i, idx := range line.indexes {
		item := f.items[idx]
		itemMainSize := finalMain[i]
		itemCrossSize := f.itemCross(item, target)

		align := item.AlignSelf
		if align == AlignAuto {
			align = f.align
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

		var childTarget rl.Rectangle
		if f.isRow() {
			childTarget = rl.NewRectangle(mainPos, crossPos, itemMainSize, crossSizeFinal)
		} else {
			childTarget = rl.NewRectangle(crossPos, mainPos, crossSizeFinal, itemMainSize)
		}

		item.Drawable.Draw(childTarget)

		mainPos += itemMainSize + between
	}
}
