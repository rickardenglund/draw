package plot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func labels(ts []tick) []string {
	var ls []string
	for _, t := range ts {
		ls = append(ls, t.label)
	}
	return ls
}

func TestTicksAreRoundAndDistinct(t *testing.T) {
	require.Equal(t, []string{"0.0", "0.2", "0.4", "0.6", "0.8", "1.0"}, labels(ticks(-.005, 1.005, 6)))
	require.Equal(t, []string{"0", "100", "200", "300", "400", "500"}, labels(ticks(-2.5, 514, 6)))
	require.Equal(t, []string{"-1", "0", "1"}, labels(ticks(-1.3, 1.3, 3)))
	require.Equal(t, []string{"-0.05", "0.00", "0.05"}, labels(ticks(-.06, .06, 3)))
}

func TestNiceStep(t *testing.T) {
	require.InDelta(t, 0.2, niceStep(1, 6), 1e-12)
	require.InDelta(t, 5, niceStep(20, 4), 1e-12)
	require.InDelta(t, 10, niceStep(95, 10), 1e-12)
	require.Equal(t, 1.0, niceStep(0, 5))
}
