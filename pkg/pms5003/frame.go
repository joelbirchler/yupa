package pms5003

// Frame is the data measured by the sensor. The standard and environment values
// are PM concentrations (μ g/m3). Environment values are the measure of the
// concentrations in the atmospheric environment. Standard values are adjusted to
// the U.S. Standard Atmosphere at 0 km sea level.
//
// Reference:
//   * Datasheet: https://www.aqmd.gov/docs/default-source/aq-spec/resources-page/plantower-pms5003-manual_v2-3.pdf
//   * Discussion: https://publiclab.org/questions/samr/04-07-2019/how-to-interpret-pms5003-sensor-values
//   * AQI: https://en.wikipedia.org/wiki/Air_quality_index
type Frame struct {
	Standard10     uint16
	Standard25     uint16
	Standard100    uint16
	Environment10  uint16
	Environment25  uint16
	Environment100 uint16
	Count3um       uint16
	Count5um       uint16
	Count10um      uint16
	Count25um      uint16
	Count50um      uint16
	Count100um     uint16
}

func (f Frame) sum() uint16 {
	return byteSum(f.Standard10) +
		byteSum(f.Standard25) +
		byteSum(f.Standard100) +
		byteSum(f.Environment10) +
		byteSum(f.Environment25) +
		byteSum(f.Environment100) +
		byteSum(f.Count3um) +
		byteSum(f.Count5um) +
		byteSum(f.Count10um) +
		byteSum(f.Count25um) +
		byteSum(f.Count50um) +
		byteSum(f.Count100um)
}

func byteSum(u uint16) uint16 {
	return (u >> 8) + (u & 0xff)
}
