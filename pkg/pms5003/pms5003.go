package pms5003

import (
	"encoding/binary"
	"io"
	"time"

	"github.com/jacobsa/go-serial/serial"
	"gobot.io/x/gobot/platforms/raspi"
)

// TODO: We should probably port.Close() on exit
var port io.ReadCloser

// Setup enables and resets the PMS5003
func Setup(pi *raspi.Adaptor) (err error) {
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

	if err := reset(pi); err != nil {
		return err
	}

	return nil
}

func reset(pi *raspi.Adaptor) (err error) {
	// reset LOW
	if err := pi.DigitalWrite(pinReset, 0); err != nil {
		return err
	}

	// TODO: flush the serial
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

	// we can throw this out and add it to the frame start (but probably improve framestart)
	var length uint16
	if err := binary.Read(port, binary.BigEndian, &length); err != nil {
		return frame, err
	}

	// read the frame
	if err := binary.Read(port, binary.BigEndian, &frame); err != nil {
		return frame, err
	}

	// reserved
	reserved := make([]byte, 2)
	if _, err := port.Read(reserved); err != nil {
		return frame, err
	}

	// checksum? (or just don't worry about it?)
	var checksum uint16
	if err := binary.Read(port, binary.BigEndian, &checksum); err != nil {
		return frame, err
	}

	return frame, nil
}

func waitForFrame() error {
	b := make([]byte, 2)

	for {
		if _, err := port.Read(b); err != nil {
			return err
		}

		if b[0] == startOfFrame[0] && b[1] == startOfFrame[1] {
			break
		}
	}

	return nil
}
