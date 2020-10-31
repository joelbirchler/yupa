package rgbmatrix

import (
	"fmt"
	"log"
	"time"

	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

// TODOx: Just get the first steps of this config to work
// TODO: Struct for storing the connection, current frame, and buffer/matrix
// TODO: Get the rendering working
// TODO: Clean / organize
// TODO: conn.Close() and other tear down reseting

type Rgb struct {
	r, g, b uint8
}

// RgbMatrix is a driver for a 5x5 LED RGB matrix for the IS31FL3731.
type RgbMatrix struct {
	pi    *raspi.Adaptor
	conn  i2c.I2cOperations
	buff  [25]Rgb
	frame uint8
}

// New connects to the RgbMatrix device and returns an RgbMatrix
func New(pi *raspi.Adaptor) (*RgbMatrix, error) {
	return &RgbMatrix{pi: pi}, nil
}

func (m *RgbMatrix) write(reg uint8, b []byte) error {
	log.Printf("writing to reg %d: %v", reg, b)

	if err := m.conn.WriteBlockData(reg, b); err != nil {
		return err
	}

	return nil
}

func (m *RgbMatrix) Connect() error {
	// Connect part
	bus := m.pi.GetDefaultBus()
	conn, err := m.pi.GetConnection(deviceAddress, bus)
	if err != nil {
		return fmt.Errorf("could not connect to RgbMatrix device: %w", err)
	}
	m.conn = conn

	// Switch display driver memory bank
	// self.i2c.write_i2c_block_data(self.address, _BANK_ADDRESS, [_CONFIG_BANK])
	if err := m.write(bankAddress, []byte{configBank}); err != nil {
		return fmt.Errorf("error writing to configBank: %w", err)
	}

	// Reset (not sure why we do this)
	// self.i2c.write_i2c_block_data(self.address, _SHUTDOWN_REGISTER, [False])
	if err := m.write(shutdownRegister, []byte{0}); err != nil {
		return fmt.Errorf("error writing to shutdownRegister: %w", err)
	}

	// time.sleep(0.00001)
	time.Sleep(time.Microsecond)

	// self.i2c.write_i2c_block_data(self.address, _SHUTDOWN_REGISTER, [True])
	if err := m.write(shutdownRegister, []byte{1}); err != nil {
		return fmt.Errorf("error writing to shutdownRegister: %w", err)
	}

	// Show?? (Do we really need to do this here? Maybe so that we can turn them off first?)
	m.Render() // FIXME: handle err

	// Note that show flips the bank to the frame... gross side-effects... so.. back to the config bank
	// self.i2c.write_i2c_block_data(self.address, _BANK_ADDRESS, [_CONFIG_BANK])
	if err := m.write(bankAddress, []byte{configBank}); err != nil {
		return fmt.Errorf("error writing to configBank: %w", err)
	}

	// Picture mode
	// self.i2c.write_i2c_block_data(self.address, _MODE_REGISTER, [_PICTURE_MODE])
	// TODO: Not sure if we need the const pictureMode
	if err := m.write(modeRegister, []byte{0}); err != nil {
		return fmt.Errorf("error writing to modeRegister: %w", err)
	}

	// Disable audio sync
	// self.i2c.write_i2c_block_data(self.address, _AUDIOSYNC_REGISTER, [0])
	if err := m.write(audioSyncRegister, []byte{0}); err != nil {
		return fmt.Errorf("error writing to audioSyncRegister: %w", err)
	}

	// Write enable pattern to frame 1 bank
	// The enable pattern enables and disables the LED addresses that we use for the 5x5 RGB matrix.
	// The IS31FL3731 can drive other LED layouts with more LEDs. So these instructions instruct each frame
	// which addresses are on or off. We could do this for more than 2 frames. I think there are up to 8?
	for f := 0; f < 8; f++ {
		if err := m.write(bankAddress, []byte{uint8(f)}); err != nil {
			return err
		}
		if err := m.write(0x00, enablePattern); err != nil { // TODO: isn't 0x00 the modeRegister?? Pimoroni!?
			return err
		}
	}

	// TODO: we could clean this up quite a bit by looping over a list of reg,b pairs or maybe a helper function?

	return nil
}

func (m *RgbMatrix) SetPixel(x, y, r, g, b uint8) {
	// set the pixel in the buffer
	m.buff[(y*5)+x] = Rgb{r: r, g: g, b: b}
	// TODO: probably add an out of range error
}

func (m *RgbMatrix) Render() error {
	// translate buffer to pixel addresses (note that python version already has a weird buffer)
	//
	// Switch the bank to the next frame (we alt between 0 and 1)
	// self.i2c.write_i2c_block_data(self.address, _BANK_ADDRESS, [whatever_next_frame_is])
	f := m.nextFrame()

	if err := m.write(bankAddress, []byte{f}); err != nil {
		return err
	}

	// write in blocks of 32 for some reason.
	//
	// offset = _COLOR_OFFSET
	// for chunk in self._chunk(output, 32):
	// 	 self.i2c.write_i2c_block_data(self.address, offset, chunk)
	//   offset += 32
	// FIXME: This is a terrible mess
	bufa := m.bufAddresses()
	for i := 0; i < len(bufa); i += 32 {
		if err := m.write(uint8(offset+i), bufa[i:i+32]); err != nil {
			return err
		}
	}

	// advance the frame counter (alt between 0 and 1)
	// whateverFrame = the new frame yo
	// self.i2c.write_i2c_block_data(self.address, _BANK_ADDRESS, [_CONFIG_BANK])
	if err := m.write(bankAddress, []byte{configBank}); err != nil {
		return err
	}
	// self.i2c.write_i2c_block_data(self.address, _FRAME_REGISTER, [whateverFrame])
	if err := m.write(frameRegister, []byte{f}); err != nil {
		return err
	}

	m.frame = f

	return nil
}

func (m *RgbMatrix) bufAddresses() []uint8 {
	addresses := make([]uint8, 160) // that's kind of arbitrary

	for i, b := range m.buff {
		a := rgbAddresses[i]
		addresses[a.r] = b.r
		addresses[a.g] = b.g
		addresses[a.b] = b.b
	}

	return addresses
}

func (m *RgbMatrix) nextFrame() uint8 {
	if m.frame == 0 {
		return 1
	}

	return 0
}
