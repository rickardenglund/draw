package draw

import (
	"fmt"
	"os"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/rickardenglund/draw/theme"
)

type fontPref struct {
	path    string
	spacing float32 // extra glyph spacing for fonts that render too tight
}

// fontPreferences returns candidate font files, most preferred first.
func fontPreferences() []fontPref {
	return []fontPref{
		{"/Users/rickardenglund/Library/Fonts/AurulentSansMNerdFontMono-Regular.otf", 0},
		{"/System/Library/Fonts/Supplemental/Arial.ttf", 0},
		{"/Library/Fonts/Arial.ttf", 0},
		{"/usr/share/fonts/truetype/dejavu/DejaVuSansMono.ttf", 0},
		{"/usr/share/fonts/TTF/DejaVuSansMono.ttf", 0},
		{"C:\\Windows\\Fonts\\consola.ttf", 0},
		{"C:\\Windows\\Fonts\\arial.ttf", 0},
	}
}

// ScreenshotPath, when set, makes NewWindow render ScreenshotFrames frames,
// save the last one as an image at ScreenshotPath and return.
// Defaults to the DRAW_SCREENSHOT and DRAW_SCREENSHOT_FRAMES env vars.
var (
	ScreenshotPath   = os.Getenv("DRAW_SCREENSHOT")
	ScreenshotFrames = envInt("DRAW_SCREENSHOT_FRAMES", 30)
)

func envInt(key string, def int) int {
	if v, err := strconv.Atoi(os.Getenv(key)); err == nil {
		return v
	}

	return def
}

func NewWindow(size rl.Vector2, init func(), objs ...Drawable) {
	rl.SetTraceLogLevel(rl.LogWarning)

	rl.SetConfigFlags(rl.FlagWindowHighdpi)
	rl.SetConfigFlags(rl.FlagMsaa4xHint)

	rl.InitWindow(int32(size.X), int32(size.Y), "Super Window")
	theme.Font = rl.GetFontDefault()

	for _, fp := range fontPreferences() {
		if _, err := os.Stat(fp.path); err != nil {
			continue
		}
		mf := rl.LoadFontEx(fp.path, 24, nil)
		if rl.IsFontValid(mf) {
			theme.Font = mf
			theme.Spacing = fp.spacing
			defer rl.UnloadFont(mf)
			fmt.Printf("selected font: %s\n", fp.path)
			break
		}
	}

	defer rl.CloseWindow()

	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	init()
	for frame := 1; !rl.WindowShouldClose(); frame++ {
		rl.BeginDrawing()
		rl.ClearBackground(theme.Nude)
		for i := range objs {
			border := rl.NewVector2(1, 1)
			ms := rl.Vector2Subtract(size, rl.Vector2Scale(border, 2))
			p := border
			tar := rl.NewRectangle(p.X, p.Y, ms.X, ms.Y)
			objs[i].Draw(tar)
		}

		if ScreenshotPath != "" && frame >= ScreenshotFrames {
			saveScreenshot(ScreenshotPath)
			rl.EndDrawing()
			return
		}
		rl.EndDrawing()
	}
}

func saveScreenshot(path string) {
	rl.DrawRenderBatchActive()
	img := rl.LoadImageFromScreen()
	defer rl.UnloadImage(img)

	if !rl.ExportImage(*img, path) {
		fmt.Printf("failed to save screenshot: %s\n", path)
		return
	}
	fmt.Printf("saved screenshot: %s\n", path)
}
