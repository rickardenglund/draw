package plot

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/floating"
	"github.com/rickardenglund/draw/theme"
)

type Plot struct {
	cur           []rl.Vector2
	prev          []rl.Vector2
	psUpdated     float64
	limitsUpdated float64
	curLimits     limits
	prevLimits    limits
	animDuration  float64
	selection     selection
	markers       *markers
}

func (p Plot) GetSize(target rl.Rectangle) rl.Vector2 {
	return rl.NewVector2(target.Width, target.Height)
}

func (p *Plot) Draw(targetWidget rl.Rectangle) {
	if len(p.cur) < 2 {
		return
	}

	ps := p.getPs()

	ls := p.getLimits()

	YTickSize := theme.MeaureTextPad(fmt.Sprintf(getFmt(ls.minY, ls.maxY), ls.maxY))
	XTickSize := theme.MeaureTextPad(fmt.Sprintf(getFmt(ls.minX, ls.maxX), ls.maxX))
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
		sp := s.transform(ps[i])
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
			screenPs[i] = sp
			clr = rl.Red
			th *= 4
		}
		rl.DrawCircleV(sp, th, clr)
	}
	rl.DrawLineStrip(screenPs, theme.Charcoal)

	if p.selection.IsActive() {
		rl.DrawRectangleLinesEx(p.selection.GetRect(), 2, rl.ColorAlpha(rl.White, .5))
	}

	mp := rl.GetMousePosition()
	ci, d := findClosest(screenPs, mp)
	if !p.selection.IsActive() &&
		d < 30 &&
		rl.CheckCollisionPointRec(mp, targetWidget) &&
		rl.IsKeyDown(rl.KeySpace) {
		p := ps[ci]
		fmtString := fmt.Sprintf("Y: %s\nX: %s", getFmt(0, p.Y), getFmt(0, p.X))
		msg := fmt.Sprintf(fmtString, p.Y, p.X)
		floating.TextArea(mp, msg)
		rl.DrawLineEx(mp, screenPs[ci], 1, theme.Charcoal)
		rl.DrawCircleV(mp, 3, theme.Charcoal)
		rl.DrawCircleV(screenPs[ci], 3, theme.Charcoal)
	}

	p.markers.Draw(s)
}

func clamp(p rl.Vector2, target rl.Rectangle) rl.Vector2 {
	x := max(target.X, p.X)
	x = min(x, target.X+target.Width)
	y := max(target.Y, p.Y)
	y = min(y, target.Y+target.Height)

	return rl.NewVector2(x, y)
}

func (p *Plot) getPs() []rl.Vector2 {
	var ps []rl.Vector2
	now := rl.GetTime()
	animTime := now - p.psUpdated
	if animTime > p.animDuration {
		return p.cur
	}

	n := max(len(p.cur), len(p.prev))
	cur := make([]rl.Vector2, n)
	copy(cur, p.cur)
	cur = cur[:len(p.cur)]

	prev := make([]rl.Vector2, n)
	copy(prev, p.prev)
	prev = prev[:len(p.prev)]

	if len(cur) > len(prev) {
		prev = extend(prev, len(cur))
	}
	if len(cur) < len(prev) {
		cur = extend(cur, len(prev))
	}

	ps = make([]rl.Vector2, len(cur))
	f := animTime / p.animDuration
	for i := range cur {
		diff := rl.Vector2Subtract(cur[i], prev[i])
		scaled := rl.Vector2Scale(diff, float32(f))
		ps[i] = rl.Vector2Add(prev[i], scaled)
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

func (p *Plot) Set(ps []rl.Vector2) {
	p.prev = p.getPs()
	p.cur = ps
	p.psUpdated = rl.GetTime()
	p.SetLimits(minmax(ps).CenterZoomed(.01))
}

func (p *Plot) SetLimits(ls limits) {
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

func (p Plot) getLimits() limits {
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

func NewPlot(ps []rl.Vector2) *Plot {
	ls := minmax(ps).CenterZoomed(.01)
	return &Plot{
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
