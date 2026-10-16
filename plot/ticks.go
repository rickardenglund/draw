package plot

import (
	"math"
	"strconv"

	"github.com/rickardenglund/draw/theme"
)

type tick struct {
	v     float32
	label string
}

// niceStep is the smallest round step (1, 2 or 5 times a power of ten)
// giving at most maxCount steps over span.
func niceStep(span float64, maxCount int) float64 {
	maxCount = max(1, maxCount)
	if span <= 0 || math.IsInf(span, 0) || math.IsNaN(span) {
		return 1
	}
	raw := span / float64(maxCount)
	mag := math.Pow(10, math.Floor(math.Log10(raw)))
	for _, m := range []float64{1, 2, 5} {
		if m*mag >= raw*(1-1e-9) {
			return m * mag
		}
	}
	return 10 * mag
}

// ticks are the multiples of a round step within [lo, hi], at most about
// maxCount of them, labelled with just the decimals the step needs.
func ticks(lo, hi float32, maxCount int) []tick {
	step := niceStep(float64(hi-lo), maxCount)
	dec := max(0, int(-math.Floor(math.Log10(step)+1e-9)))
	var ts []tick
	for i := math.Ceil(float64(lo)/step - 1e-6); i*step <= float64(hi)+step*1e-6; i++ {
		v := i * step
		if i == 0 {
			v = 0 // never "-0.0"
		}
		ts = append(ts, tick{float32(v), strconv.FormatFloat(v, 'f', dec, 64)})
		if len(ts) > 1000 {
			break
		}
	}
	return ts
}

// xTicks are as many ticks over [lo, hi] as fit side by side in width.
func xTicks(lo, hi, width float32) []tick {
	gap := theme.MeaureTextPad("0").X
	for n := max(1, int(width/(2*gap))); n > 1; n-- {
		ts := ticks(lo, hi, n)
		if float32(len(ts))*(widest(ts)+gap) <= width {
			return ts
		}
	}
	return ticks(lo, hi, 1)
}

// yTicks are as many ticks over [lo, hi] as fit, with some air, in height.
func yTicks(lo, hi, height float32) []tick {
	h := theme.MeaureTextPad("0").Y
	ts := ticks(lo, hi, max(1, int(height/(1.2*h))))
	if len(ts) < 2 && height >= 2*h {
		// A lone tick, often just 0, says nothing about the scale.
		ts = ticks(lo, hi, 3)
	}
	return ts
}

func widest(ts []tick) float32 {
	w := float32(0)
	for _, t := range ts {
		w = max(w, theme.MeaureTextPad(t.label).X)
	}
	return w
}
