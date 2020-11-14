package pms5003

const (
	pinReset  = "13" // physical pin 13, gpio/bcm 27
	pinEnable = "15" // physical pin 15, gpio/bcm 22
)

var startOfFrame = [2]byte{0x42, 0x4d}
