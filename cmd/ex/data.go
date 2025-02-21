package main

import (
	"math"
	"math/cmplx"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/mjibson/go-dsp/fft"
)

type data struct {
	tw []rl.Vector2
	sp []rl.Vector2
}

func newData() *data {
	return &data{
		tw: []rl.Vector2{},
		sp: []rl.Vector2{},
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

	d.tw = wave
	d.sp = spectra
}

func (d *data) getSP() []rl.Vector2 {
	return d.sp
}

func (d *data) getTWF() []rl.Vector2 {
	return d.tw
}

func getSine() []rl.Vector2 {
	sr := 44100.0
	dt := 1.0 / sr

	bf := rand.Float64()*5000 + 200

	fs := []func(t float64) float64{}
	nf := 4
	for i := range nf {
		f := bf * float64(i+1)
		fs = append(fs, func(t float64) float64 {
			return math.Sin(t*math.Pi*2*f) * 5 * (1 / float64(i+1))
		})
	}

	n := 44100 / 64
	ps := make([]rl.Vector2, n)
	for i := range n {
		t := float64(i) * dt
		//v := math.Sin(t*math.Pi*2*f) * 5
		//v += math.Sin(t*math.Pi*2*f*2) * 10
		v := float64(0)
		for _, f := range fs {
			v += f(t)
		}
		ps[i] = rl.NewVector2(float32(t), float32(v))

	}
	return ps
}
