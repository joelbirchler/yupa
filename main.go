package main

import (
	"fmt"

	"github.com/joelbirchler/yupa/pkg/rgbmatrix"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	fmt.Println("yo")

	pi := raspi.NewAdaptor()
	m, err := rgbmatrix.New(pi)
	fmt.Printf("%v ; %v", m, err)

	err = m.Connect()
	fmt.Printf("%v", err)

	m.SetPixel(0, 0, 255, 0, 0)
	m.SetPixel(4, 4, 0, 255, 0)
	m.SetPixel(0, 4, 0, 0, 255)
	m.SetPixel(4, 0, 255, 0, 255)

	m.SetPixel(1, 1, 125, 0, 0)
	m.SetPixel(2, 2, 125, 125, 125)
	m.SetPixel(3, 3, 0, 125, 0)
	m.SetPixel(1, 3, 0, 125, 125)
	m.SetPixel(3, 1, 125, 0, 125)
	m.Render()
}
