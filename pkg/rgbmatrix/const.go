package rgbmatrix

const (
	deviceAddress = 0x74

	modeRegister      = 0x00
	frameRegister     = 0x01
	audioSyncRegister = 0x06
	shutdownRegister  = 0x0a
	bankRegister      = 0xfd

	configBank = 0x0b

	pictureMode = 0x00
	offset      = 0x24
)

// enable pattern defines which LED addresses that we use for the 5x5 RGB matrix.
var enablePattern = []byte{
	0b00000000, 0b10111111,
	0b00111110, 0b00111110,
	0b00111111, 0b10111110,
	0b00000111, 0b10000110,
	0b00110000, 0b00110000,
	0b00111111, 0b10111110,
	0b00111111, 0b10111110,
	0b01111111, 0b11111110,
	0b01111111, 0b00000000,
}

var rgbAddresses = []Rgb{
	{118, 69, 85}, {117, 68, 101}, {116, 84, 100}, {115, 83, 99}, {114, 82, 98},
	{132, 19, 35}, {133, 20, 36}, {134, 21, 37}, {112, 80, 96}, {113, 81, 97},
	{131, 18, 34}, {130, 17, 50}, {129, 33, 49}, {128, 32, 48}, {127, 47, 63},
	{125, 28, 44}, {124, 27, 43}, {123, 26, 42}, {122, 25, 58}, {121, 41, 57},
	{126, 29, 45}, {15, 95, 111}, {8, 89, 105}, {9, 90, 106}, {10, 91, 107},
}
