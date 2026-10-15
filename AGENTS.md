# Agent notes

## Screenshots and simulated input

Any program using `draw.NewWindow` can save a screenshot and exit, so you can
inspect rendering with the image viewer:

```sh
DRAW_SCREENSHOT=/tmp/x.png go run ./cmd/flex
DRAW_MOUSE=250,349 DRAW_SCREENSHOT=/tmp/x.png go run ./cmd/flex   # with hover
```

- `DRAW_SCREENSHOT` / `DRAW_SCREENSHOT_FRAMES` (default 30) map to
  `draw.ScreenshotPath` / `draw.ScreenshotFrames`.
- `DRAW_MOUSE=x,y` (or `input.SetMousePosition`) overrides the mouse position.
- Always read the mouse through `input.MousePosition()`, never
  `rl.GetMousePosition()`, so the override works.
- Screenshots are HiDPI (2x on macOS): mouse coordinates are PNG pixel
  coordinates divided by 2.
- A window briefly opens on screen; this needs a desktop session (not `go test`).

## Testing

- `go test ./plot` has a failing `Test_limits_Zoomed` and vet warnings in
  `plot/plot_test.go` that predate this note.
