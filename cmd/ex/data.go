package main

import (
	"math"
	"math/cmplx"
	"math/rand"

	"github.com/gen2brain/raylib-go/easings"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/mjibson/go-dsp/fft"
)

type data struct {
	tw      []rl.Vector2
	sp      []rl.Vector2
	prevTW  []rl.Vector2
	updated float64
	prevSP  []rl.Vector2
}

func newData() *data {
	return &data{
		tw:      []rl.Vector2{},
		sp:      []rl.Vector2{},
		updated: rl.GetTime(),
	}
}

func (d *data) update() {
	wave := getSine()
	ps := make([]float64, len(wave))

	spectra := make([]rl.Vector2, len(wave)/2)
	for i := range ps {
		ps[i] = float64(wave[i].Y)
	}

	sampleRate := float32(44100)
	binSize := sampleRate / float32(len(ps))

	polars := fft.FFTReal(ps)
	for i := range spectra {
		m, _ := cmplx.Polar(polars[i])
		f := float32(i) * binSize
		spectra[i] = rl.NewVector2(f, float32(m))
	}

	d.prevTW = d.tw
	d.updated = rl.GetTime()
	d.tw = wave
	d.prevSP = d.sp
	d.sp = spectra
}

func anim(old, cur []rl.Vector2, updated, dur float32) []rl.Vector2 {
	now := float32(rl.GetTime())
	if now > updated+dur || len(old) != len(cur) {
		return cur
	}

	mod := make([]rl.Vector2, len(cur))
	for i := range cur {
		y := easings.QuadIn(float32(now-updated), old[i].Y, cur[i].Y-old[i].Y, dur)
		mod[i] = rl.NewVector2(cur[i].X, y)
	}

	return mod
}

func (d *data) getSP() []rl.Vector2 {
	return anim(d.prevSP, d.sp, float32(d.updated), .2)
}

func (d *data) getTWF() []rl.Vector2 {
	return anim(d.prevTW, d.tw, float32(d.updated), .2)
}

func getSine() []rl.Vector2 {
	noiseLevel := float64(1)
	sr := 44100.0
	dt := 1.0 / sr

	bf := rand.Float64()*5000 + 200
	rpm := float64(1200)

	fs := []func(t float64) float64{}
	nOvertones := 4
	for io := range nOvertones {
		f := bf * float64(io+1)
		nSideBand := 4
		for i := range nSideBand {
			fs = append(fs, func(t float64) float64 {
				return math.Sin(t*math.Pi*2*(f-float64(i)*rpm)) * 5 * (1 / float64(i+1)) * (1.0 / float64(io+1))
			})
			fs = append(fs, func(t float64) float64 {
				return math.Sin(t*math.Pi*2*(f+float64(i)*rpm)) * 5 * (1 / float64(i+1)) * (1.0 / float64(io+1))
			})
		}
	}

	n := 44100 / 1 //64
	ps := make([]rl.Vector2, n)
	for i := range n {
		t := float64(i) * dt
		v := float64(0)
		for _, f := range fs {
			v += f(t)
		}

		v += rand.Float64() * noiseLevel

		ps[i] = rl.NewVector2(float32(t), float32(v))

	}
	return ps
}
