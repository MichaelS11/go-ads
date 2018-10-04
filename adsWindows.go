// +build windows

package ads

// NewADS to create a new ADS struct.
// sensorType is not use at this time, it is for possible future use
func NewADS(busName string, address uint16, sensorType string) (*ADS, error) {
	return &ADS{configDataRate: uint16(ConfigDataRate128)}, nil
}

// Read reads the ads chip
func (ads *ADS) Read() (uint16, error) {
	return 0, nil
}
