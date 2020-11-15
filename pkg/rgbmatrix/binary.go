package rgbmatrix

import "fmt"

// Binary creates a binary image of a uint16
// TODO: Shift from green to red to purple as we get closer to 500
func Binary(n uint16) [25]Rgb {
	digits := fmt.Sprintf("%25b\n", n)

	image := Clear
	for i, b := range digits {
		if b == '1' {
			image[i] = Rgb{0xff, 0xff, 0x00}
		}
	}

	return image
}
