package pms5003

const (
	pinReset  = "13" // physical pin 13, gpio/bcm 27
	pinEnable = "15" // physical pin 15, gpio/bcm 22
)

var startOfFrame = []byte{0x42, 0x4d, 0x00, 0x1c} // the second two bytes are actually the frame length (which is always 28)
var headerSum uint16 = 0x42 + 0x4d + 0x1c
