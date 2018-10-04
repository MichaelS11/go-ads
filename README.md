# Go ADS1113, ADS1114, and ADS1115 interface

Golang ADS1113, ADS1114, and ADS1115 interface using periph.io driver

[![GoDoc Reference](https://godoc.org/github.com/MichaelS11/go-ads?status.svg)](http://godoc.org/github.com/MichaelS11/go-ads)
[![Go Report Card](https://goreportcard.com/badge/github.com/MichaelS11/go-ads)](https://goreportcard.com/report/github.com/MichaelS11/go-ads)

Should be easily extendable for ADS1111, ADS1112, and ADS1113, devices (ADS111x) with very little work. If you are intrested, let me know and we can work togeter to get it to work.


## Please note

Please make sure to setup your ADS1113, ADS1114, and ADS1115 correctly. Do a search on the internet to find guide. Here is an example of a guide:

https://learn.adafruit.com/adafruit-4-channel-adc-breakouts/

The examples below are from using a Raspberry Pi 3 with I2C1 for the bus and 0x48 for the address. Your setup may be different, if so, your bus and address would need to change in each example.

Tested on Raspberry Pi 3 with ADS1115. Please open an issue if there are any issues.


## Get

go get github.com/MichaelS11/go-ads


## ReadRetry example

```go
package main

import (
	"fmt"

	"github.com/MichaelS11/go-ads"
)

func main() {
	// call HostInit once
	err := ads.HostInit()
	if err != nil {
		fmt.Println(err)
	}

	// create new ads with wanted busName and address. 
	var ads1 *ads.ADS
	ads1, err = ads.NewADS("I2C1", 0x48, "")
	if err != nil {
		fmt.Println(err)
	}

	// example changing config gain (2/3 is default, so only an example)
	ads1.SetConfigGain(ads.ConfigGain2_3)

	// read retry from ads chip
	var result uint16
	result, err = ads1.ReadRetry(5)
	if err != nil {
		ads1.Close()
		fmt.Println(err)
	}

	// close ads bus
	err = ads1.Close()
	if err != nil {
		fmt.Println(err)
	}

	// print results
	fmt.Println("result:", result)
	volts := (float64(result) / 32767.0 * 5.0)
	fmt.Println("volts:", volts)
	psi := 360.0*volts - 25.0
	fmt.Println("psi:", psi)
}
```


## ReadBackground example

```go
package main

import (
	"fmt"
	"time"

	"github.com/MichaelS11/go-ads"
)

func main() {
	// call HostInit once
	err := ads.HostInit()
	if err != nil {
		fmt.Println(err)
	}

	// create new ads with wanted busName and address
	var ads1 *ads.ADS
	ads1, err = ads.NewADS("I2C1", 0x48, "")
	if err != nil {
		fmt.Println(err)
	}

	// changing config gain
	ads1.SetConfigGain(ads.ConfigGain1)

	stop := make(chan struct{})
	stopped := make(chan struct{})
	var result uint16

	// get ads chip reading every 5 seconds in background
	go ads1.ReadBackground(&result, 5*time.Second, stop, stopped)

	// should have at least read the ads chip twice after 15 seconds
	time.Sleep(15 * time.Second)

	// to stop ReadBackground after done with reading, close the stop channel
	close(stop)

	// can check stopped channel to know when ReadBackground has stopped
	<-stopped

	// close ads bus
	err = ads1.Close()
	if err != nil {
		fmt.Println(err)
	}

	// print results
	fmt.Println("result:", result)
	volts := (float64(result) / 32767.0 * 5.0)
	fmt.Println("volts:", volts)
	psi := 360.0*volts - 25.0
	fmt.Println("psi:", psi)
}
```


## Read example

```go
	// call HostInit once
	err := ads.HostInit()
	if err != nil {
		fmt.Println(err)
	}

	// create new ads with wanted busName and address. 
	var ads1 *ads.ADS
	ads1, err = ads.NewADS("I2C1", 0x48, "")
	if err != nil {
		fmt.Println(err)
	}

	// example changing config gain (2/3 is default, so only an example)
	ads1.SetConfigGain(ads.ConfigGain2_3)

	// read from ads chip
	var result uint16
	result, err = ads1.Read()
	if err != nil {
		ads1.Close()
		fmt.Println(err)
	}

	// close ads bus
	err = ads1.Close()
	if err != nil {
		fmt.Println(err)
	}

	// print results
	fmt.Println("result:", result)
	volts := (float64(result) / 32767.0 * 5.0)
	fmt.Println("volts:", volts)
	psi := 360.0*volts - 25.0
	fmt.Println("psi:", psi)
}
```
