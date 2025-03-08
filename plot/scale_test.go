package plot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_limits_Zoomed(t *testing.T) {
	l := limits{
		minX: 0,
		maxX: 100,
		minY: 0,
		maxY: 100,
	}
	z := l.Zoomed(.1)
	require.Equal(t, float32(105), z.maxX)
	require.Equal(t, float32(-5), z.minX)
	require.Equal(t, float32(105), z.maxY)
	require.Equal(t, float32(-5), z.minY)
}
