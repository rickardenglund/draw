package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(600, 400, "Andriana")
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		drawHeart(100, 100)
		drawHeart(300, 100)
		rl.EndDrawing()
	}
}

func drawHeart(x, y int32) {
	rl.DrawLine(x, y, x+100, y+100, rl.Red)
	rl.DrawLine(x+100, y+100, x+200, y, rl.Red)
	rl.DrawLine(x+200, y, x+150, y-25, rl.Red)
	rl.DrawLine(x+150, y-25, x+100, y, rl.Red)
	rl.DrawLine(x+100, y, x+50, y-25, rl.Red)
	rl.DrawLine(x+50, y-25, x, y, rl.Red)
}
