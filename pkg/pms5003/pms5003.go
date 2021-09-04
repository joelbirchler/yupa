package pms5003

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/jacobsa/go-serial/serial"
	"gobot.io/x/gobot/platforms/raspi"
)

var (
	port io.ReadCloser
)

// Open resets and connects to the PMS5003
func Open(pi *raspi.Adaptor) (err error) {
	log.Println("opening PMS5003")

	// Setup GPIO enable and reset pins (both start high)
	if err := pi.DigitalWrite(pinEnable, 1); err != nil {
		return err
	}

	if err := pi.DigitalWrite(pinReset, 1); err != nil {
		return err
	}

	// Connect serial
	options := serial.OpenOptions{
		PortName:        "/dev/ttyAMA0",
		BaudRate:        9600,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 1,
	}

	port, err = serial.Open(options)
	if err != nil {
		return err
	}

	if err := Reset(pi); err != nil {
		return err
	}

	return nil
}

// Close closes the connection with the PMS5003 serial port
func Close() {
	log.Println("closing PMS5003")
	port.Close()
}

// Reset attempts to reset the PMS5003
func Reset(pi *raspi.Adaptor) (err error) {
	log.Println("resetting PMS5003")

	// reset LOW
	if err := pi.DigitalWrite(pinReset, 0); err != nil {
		return err
	}

	time.Sleep(1 * time.Millisecond)

	// reset HIGH
	if err := pi.DigitalWrite(pinReset, 1); err != nil {
		return err
	}

	return nil
}

// ReadFrame reads a Frame of data from the sensor
func ReadFrame() (Frame, error) {
	var frame Frame

	if err := waitForFrame(); err != nil {
		return frame, err
	}

	// read the frame
	if err := binary.Read(port, binary.BigEndian, &frame); err != nil {
		return frame, err
	}

	// reserved (we'll need these to check the checksum)
	reserved := make([]byte, 2)
	if _, err := port.Read(reserved); err != nil {
		return frame, err
	}

	// check the checksum
	var checksum uint16
	if err := binary.Read(port, binary.BigEndian, &checksum); err != nil {
		return frame, err
	}

	sum := headerSum + frame.sum() + uint16(reserved[0]+reserved[1])
	if sum != checksum {
		return frame, fmt.Errorf("checksum error (got %d, expected %d)", sum, checksum)
	}

	return frame, nil
}

func waitForFrame() error {
	i := 0

	for i < len(startOfFrame) {
		if readByte() == startOfFrame[i] {
			i++
		} else {
			i = 0
		}
	}

	return nil
}

func readByte() byte {
	b := make([]byte, 1)
	port.Read(b) // ignoring errors
	return b[0]
}
