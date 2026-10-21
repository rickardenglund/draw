package plot

import (
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/stretchr/testify/require"
)

func Test_limits_zoomedAt(t *testing.T) {
	l := limits{minX: 0, maxX: 100, minY: 0, maxY: 100}

	z := l.zoomedAt(.1, rl.NewVector2(50, 50))
	require.Equal(t, limits{minX: -5, maxX: 105, minY: -5, maxY: 105}, z)

	// Zoomed at a corner, that corner stays put.
	z = l.zoomedAt(.1, rl.NewVector2(0, 100))
	require.Equal(t, limits{minX: 0, maxX: 110, minY: -10, maxY: 100}, z)

	// Zooming in (negative f) at the centre.
	z = l.zoomedAt(-.1, rl.NewVector2(50, 50))
	require.Equal(t, limits{minX: 5, maxX: 95, minY: 5, maxY: 95}, z)
}

func Test_limits_CenterZoomed(t *testing.T) {
	l := limits{minX: 0, maxX: 100, minY: 0, maxY: 100}
	z := l.CenterZoomed(.1)
	require.Equal(t, limits{minX: -5, maxX: 105, minY: -5, maxY: 105}, z)
}
