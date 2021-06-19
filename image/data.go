package image

import "github.com/RH12503/Triangula/color"

type Data interface {
	RGBAt(x, y int) color.RGB
	Size() (int, int)
}
