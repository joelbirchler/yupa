package pms5003

// Frame is data measured by the sensor. The standard and environment values
// are PM concentrations (μ g/m3). Environment values are the measure of the current
// concentrations in the atmospheric environment. Standard values are adjusted to
// the U.S. Standard Atmosphere at 0 km sea level.
// Reference:
//   * Datasheet: https://www.aqmd.gov/docs/default-source/aq-spec/resources-page/plantower-pms5003-manual_v2-3.pdf
//   * Discussion: https://publiclab.org/questions/samr/04-07-2019/how-to-interpret-pms5003-sensor-values
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
