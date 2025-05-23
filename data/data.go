package data

import (
	"math"
	"math/cmplx"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/mjibson/go-dsp/fft"
)

type Data struct {
	tw []rl.Vector2
	sp []rl.Vector2
}

func NewData() *Data {
	return &Data{
		tw: []rl.Vector2{},
		sp: []rl.Vector2{},
	}
}

func (d *Data) Update() {
	wave := getSine()
	d.Set(wave, 44100)
}

func (d *Data) Set(wave []rl.Vector2, sampleRate float32) {
	ps := make([]float64, len(wave))

	for i := range ps {
		ps[i] = float64(wave[i].Y)
	}

	binSize := sampleRate / float32(len(ps))

	polars := fft.FFTReal(ps)

	spectra := make([]rl.Vector2, len(wave)/2)
	for i := range spectra {
		m, _ := cmplx.Polar(polars[i])
		f := float32(i) * binSize
		m /= float64(len(spectra))
		spectra[i] = rl.NewVector2(f, float32(m))
	}
	spectra[0].Y *= .5

	d.tw = wave
	d.sp = spectra
}

func (d *Data) GetSP() []rl.Vector2 {
	return d.sp
}

func (d *Data) GetTWF() []rl.Vector2 {
	return d.tw
}

func getSine() []rl.Vector2 {
	dc := rand.Float64() * 100
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

	n := 44100 / 64
	ps := make([]rl.Vector2, n)
	for i := range n {
		t := float64(i) * dt
		v := float64(dc)
		for _, f := range fs {
			v += f(t)
		}

		v += rand.Float64() * noiseLevel

		ps[i] = rl.NewVector2(float32(t), float32(v))

	}
	return ps
}
