package draw

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/theme"
)

func NewWindow(size rl.Vector2, objs ...Drawable) {
	rl.SetTraceLogLevel(rl.LogWarning)
	rl.InitWindow(int32(size.X), int32(size.Y), "Super Window")
	theme.Font = rl.GetFontDefault()
	theme.Font = rl.LoadFontEx("/Users/rickard/Library/Fonts/AurulentSansMNerdFontMono-Regular.otf", 24, nil)
	defer rl.UnloadFont(theme.Font)

	defer rl.CloseWindow()

	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(theme.Nude)
		for i := range objs {
			border := rl.NewVector2(1, 1)
			ms := rl.Vector2Subtract(size, rl.Vector2Scale(border, 2))
			p := border
			tar := rl.NewRectangle(p.X, p.Y, ms.X, ms.Y)
			objs[i].Draw(tar)
		}
		rl.EndDrawing()
	}
}
