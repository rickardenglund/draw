package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/draw"
	"github.com/rickardenglund/draw/views"
	"github.com/rickardenglund/draw/views/text"
	"github.com/rickardenglund/draw/widget"
)

func main() {
	v := views.NewColumnView()
	btn := widget.NewClickable(text.NewTextf("more"), func() {
		v.Add(newView())
	})
	v.Add(btn)
	init := func() {
		v.Add(newView())
	}
	draw.NewWindow(rl.NewVector2(500, 500), init, v)
}
func apa() {
	rl.InitWindow(800, 500, "shader")
	defer rl.CloseWindow()

	shader := rl.LoadShader("cmd/shader/base.vert", "cmd/shader/base.frag")
	//texture := rl.LoadTexture("/Users/rickard/repo/raylib-go/examples/shaders/texture_waves/space.png")

	secondsLoc := rl.GetShaderLocation(shader, "seconds")
	centerLoc := rl.GetShaderLocation(shader, "center")
	pointsLoc := rl.GetShaderLocation(shader, "ps")

	center := []float32{.5, .5}
	const nPoints = 400
	points := make([]float32, nPoints*2)

	screenSize := []float32{float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())}

	rl.SetShaderValue(shader, rl.GetShaderLocation(shader, "size"), screenSize, rl.ShaderUniformVec2)

	seconds := []float32{0}

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		mp := rl.GetMousePosition()
		px := mp.X // float32(texture.Width)
		py := mp.Y // float32(texture.Height)
		center[0] = px

		center[1] = py

		seconds[0] += rl.GetFrameTime()
		for i := 0; i < nPoints*2; i += 2 {
			points[i] = (float32(i) / float32(nPoints)) - 1

			t := 10 * 2 * math.Pi * float64(points[i])
			points[i+1] =
				.5*float32(math.Sin(t+float64(seconds[0]*2))) +
					.2*float32(math.Sin(t/5))*.5
		}

		rl.SetShaderValue(shader, secondsLoc, seconds, rl.ShaderUniformFloat)
		rl.SetShaderValue(shader, centerLoc, center, rl.ShaderUniformVec2)
		rl.SetShaderValueV(shader, pointsLoc, points, rl.ShaderUniformVec2, nPoints)

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		rl.BeginShaderMode(shader)
		//rl.DrawTexture(texture, textureX, textureY, rl.Blue)
		rl.DrawCircle(250, 200, 50, rl.Yellow)
		rl.DrawRectangle(200, 200, 200, 200, rl.Green)
		rl.EndShaderMode()

		rl.DrawFPS(10, 10)
		rl.EndDrawing()
	}

}
