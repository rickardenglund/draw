package text

import (
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/stretchr/testify/require"

	"github.com/rickardenglund/draw/theme"
)

func Test_wrap(t *testing.T) {
	theme.Font = rl.LoadFontEx("/Users/rickard/Library/Fonts/AurulentSansMNerdFontMono-Regular.otf", 24, nil)
	tt := []struct {
		in, want string
		size     rl.Vector2
	}{
		{
			in:   "my fine string",
			size: rl.NewVector2(120, 40),
			want: "my fine\nstring",
		},
		{
			in:   "my fine string",
			size: rl.NewVector2(70, 40),
			want: "my\nfine",
		},
	}

	for _, tt := range tt {
		t.Run(tt.in, func(t *testing.T) {
			got := wrap(tt.in, tt.size)
			require.Equal(t, tt.want, got)

		})
	}
}
