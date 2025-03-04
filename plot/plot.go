package plot

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/floating"
	"github.com/rickardenglund/draw/theme"
)

type Plot struct {
	getPs func() []rl.Vector2
}

func (p Plot) GetSize(target rl.Rectangle) rl.Vector2 {
	return rl.NewVector2(target.Width, target.Height)
}

func (p *Plot) Draw(targetWidget rl.Rectangle) {
	ps := p.getPs()
	if len(ps) < 2 {
		return
	}

	ls := minmax(ps)
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

	s := newScale(ps, targetData)
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

func NewPlot(getPs func() []rl.Vector2) *Plot {
	return &Plot{getPs: getPs}
}

func inter(v, fromMin, fromMax, toMin, toMax float32) float32 {
	f := (v - fromMin) / (fromMax - fromMin)
	return f*(toMax-toMin) + toMin
}
