// Package ads1115 allows interfacing with Texas Instruments ADS1115 16-bit
// analog-to-digitial coverter (ADC) with an I2C interfcace.

// TODO - add documentation
// TODO - test neg numbners

package ads1115

import (
	"math"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/kidoman/embd"
)

const (
	voltageReg = 0x00
	configReg  = 0x01
)

const (
	Single0 = 4
	Single1 = 5
	Single2 = 6
	Single3 = 7
	Diff01  = 0
	Diff03  = 1
	Diff13  = 2
	Diff23  = 3
)

const (
	Range6V       = 0
	Range4V       = 1
	Range2V       = 2
	Range1V       = 3
	RangeHalfV    = 4
	RangeQuarterV = 5
)

// ADS1115 represents an ADS1115 voltage sensor
type ADS1115 struct {
	Bus embd.I2CBus

	address byte
	ch      byte
	gain    byte
	config  uint16

	initialized bool
	mu          sync.RWMutex
}

func New(bus embd.I2CBus, addr byte, rng byte) *ADS1115 {
	ads := ADS1115{
		Bus:     bus,
		address: addr,
		ch:      Single0,
		gain:    rng,
		config:  0,
	}
	return &ads
}

func (d *ADS1115) setup() error {
	d.mu.RLock()
	if d.initialized {
		d.mu.RUnlock()
		return nil
	}
	d.mu.RUnlock()

	d.mu.Lock()
	defer d.mu.Unlock()

	start := uint16(0) & 0x1
	mux := uint16(d.ch) & 0x7
	pga := uint16(d.gain) & 0x7
	mode := uint16(0) & 0x1 // continuous
	dr := uint16(0) & 0x7   // 8 SPS - if you change this, must also changing timing
	//dr := uint16(4) & 0x7   // 128 SPS - if you change this, must also changing timing
	compMode := uint16(0) & 0x1
	compPol := uint16(0) & 0x1
	compLat := uint16(0) & 0x1
	que := uint16(3) & 0x3 // disable

	compConfig := (compMode << 4) | (compPol << 3) | (compLat << 2) | (que << 0)
	d.config = (start << 15) | (mux << 12) | (pga << 9) | (mode << 8) | (dr << 5) | compConfig

	err := d.Bus.WriteWordToReg(d.address, configReg, d.config)
	if err != nil {
		return err
	}

	glog.V(1).Infof("ads1115: initaliized")

	d.initialized = true

	return nil
}

func (d *ADS1115) setChannel(ch byte) error {
	if d.ch == ch {
		return nil
	}

	mux := uint16(ch) & 0x7

	d.config = (d.config & 0x8FFF) | (mux << 12)
	d.ch = ch
	err := d.Bus.WriteWordToReg(d.address, configReg, d.config)
	if err != nil {
		return err
	}

	return nil
}

func (d *ADS1115) Voltage(ch byte) (float64, error) {
	err := d.setup()
	if err != nil {
		return math.NaN(), err
	}

	if ch != d.ch {
		err = d.setChannel(ch)
		if err != nil {
			return math.NaN(), err
		}
		time.Sleep( 375 * time.Millisecond) // wait for next conversion
	}

	v, err := d.Bus.ReadWordFromReg(d.address, voltageReg)
	if err != nil {
		return math.NaN(), err
	}

	fs := 0.0 // full scale voltage
	switch {
	case d.gain == Range6V:
		fs = 6.144
	case d.gain == Range4V:
		fs = 4.096
	case d.gain == Range2V:
		fs = 2.048
	case d.gain == Range1V:
		fs = 1.024
	case d.gain == RangeHalfV:
		fs = 0.512
	case d.gain == RangeQuarterV:
		fs = 0.256
	}

	sv := int16( v ) 
	voltage := float64(sv) * fs / 32768.0

	return voltage, nil
}

// Close
func (d *ADS1115) Close() error {
	// put in power down mode
	config := uint16(0x0103) // single shot mode
	err := d.Bus.WriteWordToReg(d.address, configReg, config)
	if err != nil {
		return err
	}
	return nil
}
