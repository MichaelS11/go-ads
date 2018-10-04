package ads

import (
	"periph.io/x/periph/conn/i2c"
)

// https://cdn-shop.adafruit.com/datasheets/ads1115.pdf

type (
	// ConfigGain config gain amplifier
	ConfigGain uint16
	// ConfigDataRate config data rate
	ConfigDataRate uint16
)

const (
	// Bit [15] OS: Operational status/single-shot conversion start
	configOperationNone   uint16 = 0x0000 // 0 : No effect
	configOperationSingle uint16 = 0x8000 // 1 : Begin a single conversion (when in power-down mode)

	// Bits [14:12] MUX[2:0]: Input multiplexer configuration (ADS1115 only)
	configInputMultiplexerDifferential01 uint16 = 0x0000 // 000 : AINP = AIN0 and AINN = AIN1 (default)
	configInputMultiplexerDifferential03 uint16 = 0x1000 // 001 : AINP = AIN0 and AINN = AIN3
	configInputMultiplexerDifferential13 uint16 = 0x2000 // 010 : AINP = AIN1 and AINN = AIN3
	configInputMultiplexerDifferential23 uint16 = 0x3000 // 011 : AINP = AIN2 and AINN = AIN3
	configInputMultiplexerSingle0        uint16 = 0x4000 // 100 : AINP = AIN0 and AINN = GND
	configInputMultiplexerSingle1        uint16 = 0x5000 // 101 : AINP = AIN1 and AINN = GND
	configInputMultiplexerSingle2        uint16 = 0x6000 // 110 : AINP = AIN2 and AINN = GND
	configInputMultiplexerSingle3        uint16 = 0x7000 // 111 : AINP = AIN3 and AINN = GND

	// Bits [11:9] PGA[2:0]: Programmable gain amplifier configuration (ADS1114 and ADS1115 only)

	// ConfigGain2_3 gain amplifier 2/3 is +/-6.144V range (default)
	ConfigGain2_3 ConfigGain = 0x0000 // 000
	// ConfigGain1 gain amplifier 1 is +/-4.096V range
	ConfigGain1 ConfigGain = 0x0200 // 001
	// ConfigGain2 gain amplifier 2 is +/-2.048V range
	ConfigGain2 ConfigGain = 0x0400 // 010
	// ConfigGain4 gain amplifier 4 is +/-1.024V range
	ConfigGain4 ConfigGain = 0x0600 // 011
	// ConfigGain8 gain amplifier 8 is +/-0.512V range
	ConfigGain8 ConfigGain = 0x0800 // 100
	// ConfigGain16 gain amplifier 16 +/-0.256V range
	ConfigGain16 ConfigGain = 0x0A00 // 101

	// Bit [8] MODE: Device operating mode
	configOperatingModeContinuous uint16 = 0x0000 // 0 : Continuous conversion mode
	configOperatingModeSingle     uint16 = 0x0100 // 1 : Power-down single-shot mode (default)

	// Bits [7:5] DR[2:0]: Data rate

	// ConfigDataRate8 data rate of 8 samples per second
	ConfigDataRate8 ConfigDataRate = 0x0000 // 000
	// ConfigDataRate16 data rate of 16 samples per second
	ConfigDataRate16 ConfigDataRate = 0x0020 // 001
	// ConfigDataRate32 data rate of 32 samples per second
	ConfigDataRate32 ConfigDataRate = 0x0040 // 010
	// ConfigDataRate64 data rate of 64 samples per second
	ConfigDataRate64 ConfigDataRate = 0x0060 // 011
	// ConfigDataRate128 data rate of 128 samples per second (default)
	ConfigDataRate128 ConfigDataRate = 0x0080 // 100
	// ConfigDataRate250 data rate of 250 samples per second
	ConfigDataRate250 ConfigDataRate = 0x00A0 // 101
	// ConfigDataRate475 data rate of 475 samples per second
	ConfigDataRate475 ConfigDataRate = 0x00C0 // 110
	// ConfigDataRate860 data rate of 860 samples per second
	ConfigDataRate860 ConfigDataRate = 0x00E0 // 111

	// Bit [4] COMP_MODE: Comparator mode (ADS1114 and ADS1115 only)
	configComparatorModeTraditional uint16 = 0x0000 // 0 : Traditional comparator with hysteresis (default)
	configComparatorModeWindow      uint16 = 0x0010 // 1 : Window comparator

	// Bit [3] COMP_POL: Comparator polarity (ADS1114 and ADS1115 only)
	configComparatorPolarityLow  uint16 = 0x0000 // 0 : Active low (default)
	configComparatorPolarityHigh uint16 = 0x0008 // 1 : Active high

	// Bit [2] COMP_LAT: Latching comparator (ADS1114 and ADS1115 only)
	configLatchingComparatorOff uint16 = 0x0000 // 0 : Non-latching comparator (default)
	configLatchingComparatorOn  uint16 = 0x0004 // 1 : Latching comparator

	// Bits [1:0] COMP_QUE: Comparator queue and disable (ADS1114 and ADS1115 only)
	configComparatorQueue1   uint16 = 0x0000 // 00 : Assert after one conversion
	configComparatorQueue2   uint16 = 0x0001 // 01 : Assert after two conversions
	configComparatorQueue4   uint16 = 0x0002 // 10 : Assert after four conversions
	configComparatorQueueOff uint16 = 0x0003 // 11 : Disable comparator (default)

	// register pointer
	registerPointerConversion byte = 0x0000
	registerPointerConfig     byte = 0x0001

	// configDefault is the default config
	configDefault = configOperationSingle | configInputMultiplexerSingle0 | // missing ConfigGain
		configOperatingModeSingle | // missing ConfigDataRate
		configComparatorModeTraditional | configComparatorPolarityLow | configLatchingComparatorOff | configComparatorQueueOff
)

// ADS struct to interface with the Analog-to-Digital Converter.
// Call NewADS to create a new one.
type ADS struct {
	busCloser      *i2c.BusCloser
	dev            *i2c.Dev
	configGain     uint16
	configDataRate uint16
	config         []byte
	write          []byte
	read           []byte
}
