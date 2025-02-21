package shapes

import (
	"testing"

	"github.com/gen2brain/raylib-go/easings"
)

func TestName(t *testing.T) {
	v := easings.LinearIn(1, 10, 10, 2)
	println(v)

}
