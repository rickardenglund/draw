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

func (p *Plot) Draw(targetWidget rl.Rectangle) {
	margin := float32(10)
	targetPlot := rl.NewRectangle(
		targetWidget.X+margin,
		targetWidget.Y+margin,
		targetWidget.Width-2*margin,
		targetWidget.Height-2*margin,
	)

	axisWidth := float32(20)
	axisHeight := float32(20)
	targetYAxis := rl.NewRectangle(
		targetPlot.X,
		targetPlot.Y,
		axisWidth,
		targetPlot.Height,
	)
	targetXAxis := rl.NewRectangle(
		targetYAxis.X+axisWidth,
		targetPlot.Y+targetPlot.Height-axisHeight,
		targetPlot.Width-axisWidth,
		axisHeight,

	)

	targetData := rl.NewRectangle(
		targetPlot.X+axisWidth,
		targetPlot.Y,
		targetPlot.Width-axisWidth,
		targetPlot.Height-axisHeight,
	)

	ps := p.getPs()
	if len(ps) < 2 {
		return
	}

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
	rl.DrawLineEx(o1, o2, 2, theme.Salmon)
	dt := screenPs[1].X - screenPs[0].X
	thickness := float32(2)
	if dt < thickness && dt > 1 {
		thickness = dt
	}
	for i := range screenPs {
		//z := rl.NewVector2(screenPs[i].X, origo.Y)
		//rl.DrawLineEx(z, screenPs[i], thickness, theme.Charcoal)
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
