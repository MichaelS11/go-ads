// +build !windows

package ads

import (
	"encoding/binary"
	"time"

	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
)

// NewADS to create a new ADS struct.
// sensorType is not use at this time, it is for possible future use
func NewADS(busName string, address uint16, sensorType string) (*ADS, error) {
	busCloser, err := i2creg.Open(busName)
	if err != nil {
		return nil, err
	}

	ads := &ADS{
		busCloser:      &busCloser,
		dev:            &i2c.Dev{Bus: busCloser, Addr: address},
		configDataRate: uint16(ConfigDataRate128),
		config:         make([]byte, 2),
		write:          make([]byte, 3),
		read:           make([]byte, 3),
	}

	binary.BigEndian.PutUint16(ads.config, configDefault|ads.configGain|ads.configDataRate)
	ads.write[0] = registerPointerConfig
	ads.write[1] = ads.config[0]
	ads.write[2] = ads.config[1]

	return ads, nil
}

// Read reads the ads chip
func (ads *ADS) Read() (uint16, error) {
	// send config
	ads.write[0] = registerPointerConfig
	err := ads.dev.Tx(ads.write, nil)
	if err != nil {
		return 0, nil
	}

	// wait for conversion
	time.Sleep(8 * time.Millisecond)

	// send register pointer config
	ads.write[0] = registerPointerConversion
	err = ads.dev.Tx(ads.write, ads.read)
	if err != nil {
		return 0, nil
	}

	return binary.BigEndian.Uint16(ads.read), nil
}
