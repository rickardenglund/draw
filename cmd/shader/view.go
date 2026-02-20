package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/draw"
)

type view struct {
	shader rl.Shader
}

func newView() *view {
	shader := rl.LoadShader("cmd/shader/base.vert", "cmd/shader/base.frag")
	pointsLoc := rl.GetShaderLocation(shader, "ps")

	const nPoints = 400
	points := make([]float32, nPoints*2)
	for i := 0; i < nPoints*2; i += 2 {
		points[i] = (float32(i) / float32(nPoints)) - 1

		t := 10 * 2 * math.Pi * float64(points[i])
		points[i+1] =
			.5*float32(math.Sin(t+float64(rl.GetTime()))) +
				.2*float32(math.Sin(t/5))*.5
	}

	rl.SetShaderValueV(shader, pointsLoc, points, rl.ShaderUniformVec2, nPoints)

	return &view{
		shader: shader,
	}
}

func (v view) Draw(target rl.Rectangle) {

	rl.BeginShaderMode(v.shader)
	defer rl.EndShaderMode()

	rl.DrawRectangle(int32(target.X),
		int32(target.Y),
		int32(target.Width),
		int32(target.Height),
		rl.Green)
}

func (v view) GetSize(target rl.Rectangle) rl.Vector2 {
	return rl.NewVector2(target.X, target.Y)
}

var _ draw.Drawable = new(view)
