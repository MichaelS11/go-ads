package ads

import (
	"encoding/binary"
	"time"

	"periph.io/x/periph/host"
)

// HostInit calls periph.io host.Init(). This needs to be done before ADS can be used.
func HostInit() error {
	_, err := host.Init()
	return err
}

// Close closes bus
func (ads *ADS) Close() error {
	if ads.busCloser == nil {
		return nil
	}

	busCloser := *ads.busCloser
	ads.busCloser = nil

	return busCloser.Close()
}

// SetConfigInputMultiplexer sets input multiplexer
func (ads *ADS) SetConfigInputMultiplexer(configInputMultiplexer ConfigInputMultiplexer) {
	ads.configInputMultiplexer = uint16(configInputMultiplexer)
	binary.BigEndian.PutUint16(ads.config, configDefault|ads.configInputMultiplexer|ads.configGain|ads.configDataRate)
	ads.write[1] = ads.config[0]
	ads.write[2] = ads.config[1]
}

// SetConfigGain sets gain
func (ads *ADS) SetConfigGain(configGain ConfigGain) {
	ads.configGain = uint16(configGain)
	binary.BigEndian.PutUint16(ads.config, configDefault|ads.configInputMultiplexer|ads.configGain|ads.configDataRate)
	ads.write[1] = ads.config[0]
	ads.write[2] = ads.config[1]
}

// SetConfigDataRate sets data rate
func (ads *ADS) SetConfigDataRate(configDataRate ConfigDataRate) {
	ads.configDataRate = uint16(configDataRate)
	binary.BigEndian.PutUint16(ads.config, configDefault|ads.configInputMultiplexer|ads.configGain|ads.configDataRate)
	ads.write[1] = ads.config[0]
	ads.write[2] = ads.config[1]
}

// ReadRetry will call Read until there is no errors or the maxRetries is hit.
// Suggest maxRetries to be set around 5.
func (ads *ADS) ReadRetry(maxRetries int) (result uint16, err error) {
	for i := 0; i < maxRetries; i++ {
		result, err = ads.Read()
		if err == nil {
			return
		}
	}
	return
}

// ReadBackground it meant to be run in the background, run as a Goroutine.
// sleepDuration is how long it will try to sleep between reads.
// If there is ongoing read errors there will be no notice except that the result will not be updated.
// Will continue to read ads chip until stop is closed.
// After it has been stopped, the stopped chan will be closed.
// Will panic if result or stop are nil.
func (ads *ADS) ReadBackground(result *uint16, sleepDuration time.Duration, stop chan struct{}, stopped chan struct{}) {
	var resultTemp uint16
	var err error
	var startTime time.Time

Loop:
	for {
		startTime = time.Now()
		resultTemp, err = ads.Read()
		if err == nil {
			// no read error, save result
			*result = resultTemp
			// wait for sleepDuration or stop
			select {
			case <-time.After(sleepDuration - time.Since(startTime)):
			case <-stop:
				break Loop
			}
		} else {
			// read error, just check for stop
			select {
			case <-stop:
				break Loop
			default:
			}
		}
	}

	close(stopped)
}
