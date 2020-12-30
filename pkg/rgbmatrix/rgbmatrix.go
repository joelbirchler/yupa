package rgbmatrix

import (
	"log"

	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

type Rgb struct {
	R, G, B uint8
}

// RgbMatrix is a driver for a 5x5 LED RGB matrix for the IS31FL3731.
var (
	conn  i2c.I2cOperations
	buff  [25]Rgb
	frame uint8
)

// Setup connects to the RgbMatrix device and resets
func Open(pi *raspi.Adaptor) (err error) {
	log.Println("opening RGB Matrix")

	// connect
	conn, err = pi.GetConnection(deviceAddress, pi.GetDefaultBus())
	if err != nil {
		return
	}

	// config writes
	writes := []struct {
		reg   uint8
		bytes []byte
	}{
		{bankRegister, []byte{configBank}}, // switch to config bank
		{shutdownRegister, []byte{0}},      // reset
		{shutdownRegister, []byte{1}},
		{modeRegister, []byte{0}},      // picture mode
		{audioSyncRegister, []byte{0}}, // disable audio sync feature
	}

	for _, w := range writes {
		if err = write(w.reg, w.bytes); err != nil {
			return
		}
	}

	// Write enable pattern to frame 1 bank
	// The enable pattern enables and disables the LED addresses that we use for the 5x5 RGB matrix.
	// The IS31FL3731 can drive other LED layouts with more LEDs. So these instructions instruct each frame
	// which addresses are on or off. We could do this for more than 2 frames. I think there are up to 8?
	for f := 0; f < 8; f++ {
		if err = write(bankRegister, []byte{uint8(f)}); err != nil {
			return
		}
		if err = write(0x00, enablePattern); err != nil {
			return
		}
	}

	return
}

func Close() {
	log.Println("closing RGB Matrix")

	Set(Clear)
	Render()
	conn.Close()
}

func write(reg uint8, b []byte) error {
	if err := conn.WriteBlockData(reg, b); err != nil {
		return err
	}

	return nil
}

func Set(image [25]Rgb) {
	buff = image
}

func Render() (err error) {
	f := nextFrame()

	// write to the next frame
	if err = write(bankRegister, []byte{f}); err != nil {
		return
	}

	addresses := make([]uint8, 135)
	for i, b := range buff {
		a := rgbAddresses[i]
		addresses[a.R] = b.R
		addresses[a.G] = b.G
		addresses[a.B] = b.B
	}

	for i := 0; i < len(addresses); i += 32 {
		end := min(i+32, len(addresses))

		if err = write(uint8(offset+i), addresses[i:end]); err != nil {
			return
		}
	}

	// advance the frame
	if err = write(bankRegister, []byte{configBank}); err != nil {
		return
	}

	if err = write(frameRegister, []byte{f}); err != nil {
		return
	}

	frame = f

	return
}

func nextFrame() uint8 {
	if frame == 0 {
		return 1
	}

	return 0
}

func min(x, y int) int {
	if x < y {
		return x
	}

	return y
}
