package sound

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/theme"
)

type Player struct {
	sound rl.Sound
	getPs func() []rl.Vector2
}

func (p *Player) loadAudio() {
	ps := p.getPs()
	minY, maxY := float32(0), float32(0)
	for _, s := range ps {
		if s.Y < minY {
			minY = s.Y
		}
		if s.Y > maxY {
			maxY = s.Y
		}

	}
	peak := (maxY - minY) / 2.0
	wantPeak := float32(50)

	bs := make([]byte, len(ps))
	for i := range ps {
		v := byte((ps[i].Y / peak) * wantPeak)
		bs[i] = v + 128
	}

	if len(bs) == 0 {
		return
	}
	if rl.IsSoundValid(p.sound) {
		rl.UnloadSound(p.sound)
	}

	w := rl.NewWave(uint32(len(bs)), 44100, 8, 1, bs)
	s := rl.LoadSoundFromWave(w)

	p.sound = s
}

func (p *Player) Draw(target rl.Rectangle) {
	if rl.IsAudioDeviceReady() {
		clr := theme.Salmon
		playing := rl.IsSoundPlaying(p.sound)
		if playing {
			clr = theme.Charcoal
		}

		btnRad := float32(20)
		btnPos := rl.Vector2Add(
			rl.NewVector2(target.X, target.Y),
			rl.NewVector2(target.Width/2, target.Height/2),
		)
		if playing {
			t := rl.GetTime() * 2 * math.Pi
			f := t * 2
			r := float32(math.Cos(t*.125)) * 5
			x := float32(math.Cos(f)) * r
			y := float32(math.Sin(f)) * r
			btnPos = rl.Vector2Add(btnPos, rl.NewVector2(x, y))
		}
		mp := rl.GetMousePosition()
		if rl.Vector2Distance(mp, btnPos) < btnRad && !playing {
			clr = rl.ColorBrightness(clr, .2)
		}
		rl.DrawCircleV(btnPos, btnRad, clr)

		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) &&
			rl.CheckCollisionPointRec(mp, target) {
			p.Play()
		}
	}
}

func (p *Player) Play() {
	if !rl.IsSoundPlaying(p.sound) {
		p.loadAudio()
		rl.PlaySound(p.sound)
	}
}

func NewPlayer(getps func() []rl.Vector2) *Player {
	return &Player{getPs: getps}

}
