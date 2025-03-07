package plot

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/floating"
	"github.com/rickardenglund/draw/theme"
)

type Plot struct {
	cur          []rl.Vector2
	prev         []rl.Vector2
	updated      float64
	curLimits    limits
	prevLimits   limits
	animDuration float64
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
	drawAxisY(targetYAxis, s)
	drawAxisX(targetXAxis, s)

	screenPs := make([]rl.Vector2, len(ps))
	for i := range ps {
		sp := s.transform(ps[i])
		screenPs[i] = sp
	}

	origo := s.transform(rl.Vector2{})
	o1 := rl.NewVector2(screenPs[0].X, origo.Y)
	o2 := rl.NewVector2(o1.X+targetData.Width, o1.Y)
	if rl.CheckCollisionPointRec(o1, targetData) {
		rl.DrawLineEx(o1, o2, 2, theme.Salmon)
	}
	dt := screenPs[1].X - screenPs[0].X
	thickness := float32(2)
	if dt < thickness && dt > 1 {
		thickness = dt
	}
	for i := range screenPs {
		rl.DrawCircleV(screenPs[i], thickness/2, theme.Charcoal)
	}
	rl.DrawLineStrip(screenPs, theme.Charcoal)

	mp := rl.GetMousePosition()
	ci, d := findClosest(screenPs, mp)
	if d < 30 && rl.CheckCollisionPointRec(mp, targetWidget) {
		p := ps[ci]
		fmtString := fmt.Sprintf("Y: %s\nX: %s", getFmt(0, p.Y), getFmt(0, p.X))
		msg := fmt.Sprintf(fmtString, p.Y, p.X)
		floating.TextArea(mp, msg)
		rl.DrawLineEx(mp, screenPs[ci], 1, theme.Charcoal)
		rl.DrawCircleV(mp, 3, theme.Charcoal)
		rl.DrawCircleV(screenPs[ci], 3, theme.Charcoal)
	}
}

func (p *Plot) getPs() []rl.Vector2 {
	var ps []rl.Vector2
	now := rl.GetTime()
	animTime := now - p.updated
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
	p.prevLimits = p.getLimits()
	p.curLimits = minmax(ps)
	p.prev = p.getPs()
	p.cur = ps
	p.updated = rl.GetTime()
}

func (p Plot) getLimits() limits {
	var ls limits
	now := rl.GetTime()
	animTime := now - p.updated
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
	ls := minmax(ps)
	return &Plot{
		cur:          ps,
		prev:         ps,
		updated:      rl.GetTime(),
		curLimits:    ls,
		prevLimits:   ls,
		animDuration: .2,
	}
}

func inter(v, fromMin, fromMax, toMin, toMax float32) float32 {
	f := (v - fromMin) / (fromMax - fromMin)
	return f*(toMax-toMin) + toMin
}
