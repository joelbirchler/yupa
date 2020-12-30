package rgbmatrix

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Fill creates an image fill based on the color r with random variation
func Fill(r Rgb, variance int) [25]Rgb {
	image := Clear

	for i := range image {
		image[i] = Rgb{R: vary(r.R, variance), G: vary(r.G, variance), B: vary(r.B, variance)}
	}

	return image
}

func vary(u uint8, variance int) uint8 {
	return limitUint8(int(u) - (variance / 2) + rand.Intn(variance))
}

// limitUint8 doesn't wrap like uint8(u)
func limitUint8(u int) uint8 {
	if u < 0 {
		return 0
	} else if u > 255 {
		return 255
	}

	return uint8(u)
}
