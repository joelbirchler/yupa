package main

import (
	"fmt"
	"log"

	"github.com/joelbirchler/yupa/pkg/rgbmatrix"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	fmt.Println("yo")

	pi := raspi.NewAdaptor()
	if err := rgbmatrix.Setup(pi); err != nil {
		log.Fatalf("unable to setup rgbmatrix: %v", err)
	}

	rgbmatrix.SetPixel(1, 0, 255, 60, 0)
	rgbmatrix.SetPixel(2, 0, 255, 60, 0)
	rgbmatrix.SetPixel(3, 0, 255, 60, 0)

	rgbmatrix.SetPixel(0, 1, 225, 20, 0)
	rgbmatrix.SetPixel(1, 1, 225, 20, 0)
	rgbmatrix.SetPixel(2, 1, 225, 20, 0)
	rgbmatrix.SetPixel(3, 1, 225, 20, 0)
	rgbmatrix.SetPixel(4, 1, 225, 20, 0)

	rgbmatrix.SetPixel(0, 2, 185, 10, 0)
	rgbmatrix.SetPixel(2, 2, 185, 10, 0)
	rgbmatrix.SetPixel(4, 2, 185, 10, 0)

	rgbmatrix.SetPixel(0, 3, 135, 0, 0)
	rgbmatrix.SetPixel(1, 3, 135, 0, 0)
	rgbmatrix.SetPixel(2, 3, 135, 0, 0)
	rgbmatrix.SetPixel(3, 3, 135, 0, 0)
	rgbmatrix.SetPixel(4, 3, 135, 0, 0)

	rgbmatrix.SetPixel(0, 4, 105, 0, 10)
	rgbmatrix.SetPixel(2, 4, 105, 0, 10)
	rgbmatrix.SetPixel(4, 4, 105, 0, 10)

	if err := rgbmatrix.Render(); err != nil {
		log.Printf("rendering error: %v", err)
	}
}
