package waves

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Sine struct {
	ys         []float64
	sampleRate float64
}

func (s *Sine) Vec() []rl.Vector2 {
	if len(s.ys) == 0 || s.sampleRate == 0 {
		return []rl.Vector2{}
	}

	res := make([]rl.Vector2, len(s.ys))
	dt := float32(1 / s.sampleRate)
	for i := range s.ys {
		res[i] = rl.NewVector2(dt*float32(i), float32(s.ys[i]))
	}

	return res
}

func Add(os ...*Sine) *Sine {
	if len(os) == 0 {
		return nil
	}

	nOut := 0
	sampleRate := os[0].sampleRate
	for oi := range os {
		n := len(os[oi].ys)
		nOut = max(nOut, n)

		if os[oi].sampleRate != sampleRate {
			panic("sample rate not matching")
		}
	}

	res := make([]float64, nOut)
	for oi := range os {
		for i := range os[oi].ys {
			res[i] += os[oi].ys[i]

		}
	}

	return &Sine{
		ys:         res,
		sampleRate: sampleRate,
	}

}
func GetSine(freq, amplitude, offset, phase, sampleRate float64, n int) *Sine {
	dt := float64(1 / sampleRate)

	res := make([]float64, n)
	t := float64(0)
	for i := range n {
		t += dt
		y := math.Sin(freq*t*2*math.Pi+phase)*amplitude + offset
		res[i] = y
	}

	return &Sine{
		ys:         res,
		sampleRate: sampleRate,
	}
}
