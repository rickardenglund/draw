package plot

import (
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/stretchr/testify/require"
)

func Test_transform(t *testing.T) {
	tt := []struct {
		target rl.Rectangle
		ps     []rl.Vector2
		p      rl.Vector2
		want   rl.Vector2
	}{
		{
			target: rl.NewRectangle(100, 100, 100, 100),
			ps:     []rl.Vector2{{0, 0}, {100, 100}},
			p:      rl.Vector2{X: 50},
			want:   rl.Vector2{X: 150, Y: 200},
		},
		{
			target: rl.NewRectangle(100, 200, 100, 200),
			ps:     []rl.Vector2{{0, 0}, {100, 100}},
			p:      rl.Vector2{X: 50, Y: 100},
			want:   rl.Vector2{X: 150, Y: 200},
		},
		{
			target: rl.NewRectangle(100, 200, 200, 200),
			ps:     []rl.Vector2{{0, 0}, {100, 100}},
			p:      rl.Vector2{X: 100, Y: 100},
			want:   rl.Vector2{X: 300, Y: 200},
		},
	}
	for _, tt := range tt {
		t.Run("", func(t *testing.T) {
			s := newScale(tt.ps, tt.target)
			sp := s.transform(tt.p)
			require.Equal(t, tt.want, sp)
		})
	}
}

func Test_minmax(t *testing.T) {
	tests := []struct {
		ps       []rl.Vector2
		wantMinX float32
		wantMaxX float32
		wantMinY float32
		wantMaxY float32
	}{
		{
			ps:       []rl.Vector2{rl.NewVector2(0, 0)},
			wantMinX: 0,
			wantMaxX: 0,
			wantMinY: 0,
			wantMaxY: 0,
		},
		{
			ps: []rl.Vector2{
				rl.NewVector2(0, 0),
				rl.NewVector2(10, 10),
			},
			wantMinX: 0,
			wantMaxX: 10,
			wantMinY: 0,
			wantMaxY: 10,
		},
		{
			ps: []rl.Vector2{
				rl.NewVector2(5, 5),
				rl.NewVector2(10, 10),
			},
			wantMinX: 5,
			wantMaxX: 10,
			wantMinY: 5,
			wantMaxY: 10,
		}, {
			ps: []rl.Vector2{
				rl.NewVector2(-15, -5),
				rl.NewVector2(-1, -10),
			},
			wantMinX: -15,
			wantMaxX: -1,
			wantMinY: -10,
			wantMaxY: -5,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			minx, maxX, minY, maxY := minmax(tt.ps)
			require.Equal(t, tt.wantMinX, minx)
			require.Equal(t, tt.wantMaxX, maxX)
			require.Equal(t, tt.wantMinY, minY)
			require.Equal(t, tt.wantMaxY, maxY)
		})
	}
}
