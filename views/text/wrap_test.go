package text

import (
	"strings"
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/stretchr/testify/require"
)

// monoMeasure is a fake font: every character is 10 wide and every line 20
// high.
func monoMeasure(s string) rl.Vector2 {
	lines := strings.Split(s, "\n")
	w := 0
	for _, l := range lines {
		w = max(w, len(l))
	}
	return rl.NewVector2(float32(10*w), float32(20*len(lines)))
}

func Test_wrap(t *testing.T) {
	tt := []struct {
		name, in, want string
		size           rl.Vector2
	}{
		{"fits", "my fine string", "my fine string", rl.NewVector2(200, 40)},
		{"wraps at spaces", "my fine string", "my fine\nstring", rl.NewVector2(120, 40)},
		{"cut at height", "my fine string", "my\nfine", rl.NewVector2(60, 40)},
		{"keeps first line", "my fine string", "my", rl.NewVector2(60, 10)},
		{
			"long word split after punctuation",
			"to host.plant-north.example.com:502 now",
			"to\nhost.plant-\nnorth.\nexample.com:\n502 now",
			rl.NewVector2(120, 200),
		},
		{
			"long word without punctuation",
			"abcdefghijklmnopqrstuvwxyz",
			"abcdefghij\nklmnopqrst\nuvwxyz",
			rl.NewVector2(100, 200),
		},
		{"words after a long one share lines", "aaaaaaaaaaaaaaa b c d", "aaaaaaaaaa\naaaaa b c\nd", rl.NewVector2(100, 200)},
		{"keeps newlines", "one two\nthree four five", "one two\nthree four\nfive", rl.NewVector2(100, 200)},
		{"narrower than a character", "ab", "a\nb", rl.NewVector2(5, 200)},
	}

	for _, tt := range tt {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, wrapWith(tt.in, tt.size, monoMeasure))
		})
	}
}
