package plot

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/floating"
	"github.com/rickardenglund/draw/input"
	"github.com/rickardenglund/draw/theme"
)

type point interface {
	~struct {
		X float32
		Y float32
	}
}

type Labeler interface {
	GetXString() string
	GetYString() string
}

type Plot[P point] struct {
	cur       []P
	prev      []P
	psUpdated float64

	limitsUpdated float64
	curLimits     limits
	prevLimits    limits

	selection    selection
	animDuration float64

	mover mover

	markers *markers
}

func (p Plot[point]) GetSize(target rl.Rectangle) rl.Vector2 {
	return rl.NewVector2(target.Width, target.Height)
}

func (p *Plot[point]) Draw(targetWidget rl.Rectangle) {
	if len(p.cur) < 2 {
		return
	}

	ps := p.getPs()

	ls := p.getLimits()
	offset := p.mover.GetOffset()
	ls = ls.translated(offset)

	YTickSize := theme.MeaureTextPad(fmt.Sprintf(getFmt(ls.maxY-ls.minY), ls.maxY))
	XTickSize := theme.MeaureTextPad(fmt.Sprintf(getFmt(ls.maxX), ls.maxX))
	targetPlot := theme.Pad(targetWidget)

	yAxisWidth := YTickSize.X
	xAxisHeight := float32(XTickSize.Y)

	targetYAxis := targetPlot
	targetYAxis.Width = yAxisWidth
	targetYAxis.Height -= xAxisHeight

	targetXAxis := rl.NewRectangle(
		targetYAxis.X+yAxisWidth,
		targetPlot.Y+targetPlot.Height-xAxisHeight,
		targetPlot.Width-yAxisWidth,
		xAxisHeight,
	)

	targetData := rl.NewRectangle(
		targetPlot.X+yAxisWidth,
		targetPlot.Y,
		targetPlot.Width-yAxisWidth,
		targetPlot.Height-xAxisHeight,
	)

	s := newScale(ls, targetData)
	p.handleKeys(targetData, s)

	drawAxisY(targetYAxis, s)
	drawAxisX(targetXAxis, s)

	screenPs := make([]rl.Vector2, len(ps))
	for i := range ps {
		sp := s.transform(rl.Vector2(ps[i]))
		screenPs[i] = sp
	}

	o1 := s.transform(rl.NewVector2(ls.minX, 0))
	o2 := s.transform(rl.NewVector2(ls.maxX, 0))
	if rl.CheckCollisionPointRec(o1, targetData) {
		rl.DrawLineEx(o1, o2, 2, theme.Salmon)
	}
	dt := screenPs[1].X - screenPs[0].X
	thickness := float32(2)
	if dt < thickness && dt > 1 {
		thickness = dt
	}
	for i := range screenPs {
		sp := screenPs[i]
		clr := theme.Charcoal
		th := thickness / 2
		if !rl.CheckCollisionPointRec(sp, targetData) {
			sp = clamp(sp, targetData)
			//screenPs[i] = sp
			clr = rl.Red
			th *= 4
		}
		rl.DrawCircleV(sp, th, clr)
	}

	//rl.DrawLineStrip(screenPs, theme.Charcoal)
	//for i := 1; i < len(screenPs); i++ {
	//	if rl.CheckCollisionPointRec(screenPs[i], targetData) &&
	//		rl.CheckCollisionPointRec(screenPs[i-1], targetData) {
	//		rl.DrawLineV(screenPs[i-1], screenPs[i], theme.Charcoal)
	//	}
	//}
	for i := 1; i < len(screenPs); i++ {
		lineInSquare(screenPs[i-1], screenPs[i], targetData)
	}

	if p.selection.IsActive() {
		rl.DrawRectangleLinesEx(p.selection.GetRect(), 2, rl.ColorAlpha(rl.White, .5))
	}

	mp := input.MousePosition()
	ci, d := findClosest(screenPs, mp)
	if !p.selection.IsActive() &&
		d < 30 &&
		rl.CheckCollisionPointRec(mp, targetWidget) &&
		rl.IsKeyDown(rl.KeySpace) {
		p2 := rl.Vector2(ps[ci])
		lab, ok := any(ps[ci]).(Labeler)
		msg := "msg"
		if ok {
			msg = fmt.Sprintf(
				"Y: %s\nX: %s",
				lab.GetYString(),
				lab.GetXString(),
			)
		} else {
			fmtString := fmt.Sprintf(
				"Y: %s\nX: %s",
				getFmt(ls.DeltaY()),
				getFmt(ls.DeltaX()),
			)
			msg = fmt.Sprintf(fmtString, p2.Y, p2.X)
		}

		rl.DrawLineEx(mp, screenPs[ci], 1, theme.Charcoal)
		rl.DrawCircleV(mp, 3, theme.Charcoal)
		rl.DrawCircleV(screenPs[ci], 3, theme.Charcoal)
		floating.TextArea(targetWidget, mp, msg)
	}

	p.markers.Draw(s)

	p.mover.Draw(targetPlot)
}

func lineInSquare(a, b rl.Vector2, target rl.Rectangle) {
	aIn := rl.CheckCollisionPointRec(a, target)
	bIn := rl.CheckCollisionPointRec(b, target)
	if aIn && bIn {
		rl.DrawLineV(a, b, theme.Charcoal)
		return
	}
	if !aIn && !bIn {
		//mid := rl.Vector2Add(a,
		//	rl.Vector2Scale(rl.Vector2Subtract(b, a), 0.5),
		//)

		//rl.DrawCircleV(mid, 3, rl.Blue)
		//lineInSquare(a, mid, target)
		//lineInSquare(b, mid, target)
		return
		//b = mid
	}

	if !aIn && bIn {
		a, b = b, a
	}

	col := rl.Vector2{}
	if colLineRect(a, b, target, &col) {
		rl.DrawLineV(a, col, theme.Charcoal)
	}
}

func colLineRect(a rl.Vector2, b rl.Vector2, target rl.Rectangle, res *rl.Vector2) bool {
	if rl.CheckCollisionLines( // top
		a,
		b,
		rl.NewVector2(target.X, target.Y),
		rl.NewVector2(target.X+target.Width, target.Y),
		res) {
		return true
	}

	if rl.CheckCollisionLines( // left
		a,
		b,
		rl.NewVector2(target.X, target.Y),
		rl.NewVector2(target.X, target.Y+target.Height),
		res) {
		return true
	}

	if rl.CheckCollisionLines( // right
		a,
		b,
		rl.NewVector2(target.X+target.Width, target.Y),
		rl.NewVector2(target.X+target.Width, target.Y+target.Height),
		res) {
		return true
	}
	if rl.CheckCollisionLines( // bottom
		a,
		b,
		rl.NewVector2(target.X, target.Y+target.Height),
		rl.NewVector2(target.X+target.Width, target.Y+target.Height),
		res) {
		return true
	}

	return false
}

func clamp(p rl.Vector2, target rl.Rectangle) rl.Vector2 {
	x := max(target.X, p.X)
	x = min(x, target.X+target.Width)
	y := max(target.Y, p.Y)
	y = min(y, target.Y+target.Height)

	return rl.NewVector2(x, y)
}

func (p *Plot[point]) getPs() []point {
	var ps []point
	now := rl.GetTime()
	animTime := now - p.psUpdated
	if animTime > p.animDuration {
		return p.cur
	}

	n := max(len(p.cur), len(p.prev))
	cur := make([]point, n)
	copy(cur, p.cur)
	cur = cur[:len(p.cur)]

	prev := make([]point, n)
	copy(prev, p.prev)
	prev = prev[:len(p.prev)]

	if len(cur) > len(prev) {
		prev = extend(prev, len(cur))
	}
	if len(cur) < len(prev) {
		cur = extend(cur, len(prev))
	}

	ps = make([]point, len(cur))
	f := animTime / p.animDuration
	for i := range cur {
		diff := rl.Vector2Subtract(rl.Vector2(cur[i]), rl.Vector2(prev[i]))
		scaled := rl.Vector2Scale(diff, float32(f))
		ps[i] = point(rl.Vector2Add(rl.Vector2(prev[i]), scaled))
	}
	return ps
}

func extend[T any](short []T, n int) []T {
	rest := make([]T, n-len(short))
	if len(short) == 0 {
		return rest
	}

	lastVal := short[len(short)-1]
	for i := range rest {
		rest[i] = lastVal
	}

	return append(short, rest...)
}

func (p *Plot[point]) Set(ps []point) {
	p.prev = p.getPs()
	p.cur = ps
	p.psUpdated = rl.GetTime()
	p.SetLimits(minmax(ps).CenterZoomed(.01))
}

func (p *Plot[point]) ShiftIn(in point, maxSize int) {
	newPrev := p.cur
	if len(newPrev) >= maxSize {
		newPrev = newPrev[1:]
	}

	p.prev = newPrev
	p.cur = append(p.prev, in)

	p.psUpdated = rl.GetTime()
	p.SetLimits(minmax(p.cur).CenterZoomed(.01))
}

func (p *Plot[point]) SetLimits(ls limits) {
	p.prevLimits = p.getLimits()

	if ls.minY == ls.maxY {
		ls.minY -= 1
		ls.maxY += 10
	}
	if ls.minX == ls.maxX {
		ls.minX -= 1
		ls.maxX += 10
	}
	p.curLimits = ls
	p.limitsUpdated = rl.GetTime()
}

func (p Plot[point]) getLimits() limits {
	var ls limits
	now := rl.GetTime()
	animTime := now - p.limitsUpdated
	if animTime < p.animDuration {
		f := float32(animTime / p.animDuration)
		ls = limits{
			minX: rl.Lerp(p.prevLimits.minX, p.curLimits.minX, f),
			maxX: rl.Lerp(p.prevLimits.maxX, p.curLimits.maxX, f),
			minY: rl.Lerp(p.prevLimits.minY, p.curLimits.minY, f),
			maxY: rl.Lerp(p.prevLimits.maxY, p.curLimits.maxY, f),
		}
	} else {
		ls = p.curLimits
	}

	return ls
}

func findClosest(ps []rl.Vector2, mp rl.Vector2) (int, float32) {
	minDist := rl.Vector2Distance(ps[0], mp)
	closest := 0
	for i := range ps {
		if d := rl.Vector2Distance(mp, ps[i]); d < minDist {
			minDist = d
			closest = i
		}

	}

	return closest, minDist
}

func NewPlot[P point](ps []P) *Plot[P] {
	ls := minmax(ps).CenterZoomed(.01)
	return &Plot[P]{
		cur:           ps,
		prev:          ps,
		psUpdated:     rl.GetTime(),
		limitsUpdated: rl.GetTime(),
		curLimits:     ls,
		prevLimits:    ls,
		animDuration:  .2,
		markers:       newMarkers(),
	}
}

func inter(v, fromMin, fromMax, toMin, toMax float32) float32 {
	f := (v - fromMin) / (fromMax - fromMin)
	return f*(toMax-toMin) + toMin
}
