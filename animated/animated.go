package animated

import (
	"github.com/gen2brain/raylib-go/easings"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Animated struct {
	cur  float32
	prev float32
	set  float64
	dur  float32
}

func NewAnimated(cur, dur float32) *Animated {
	return &Animated{
		cur:  cur,
		prev: cur,
		set:  rl.GetTime(),
		dur:  dur,
	}
}

func (a *Animated) Set(v float32) {
	a.prev = a.Get()
	a.cur = v
	a.set = rl.GetTime()
}

func (a *Animated) Get() float32 {
	t := float32(rl.GetTime() - a.set)
	if t > a.dur {
		return a.cur
	}
	return easings.LinearOut(t, a.prev, a.cur-a.prev, a.dur)
}

func (a *Animated) SetForce(f float32) {
	a.Set(f)
	a.prev = f
}
