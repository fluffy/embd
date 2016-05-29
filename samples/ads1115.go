// +build ignore

package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/kidoman/embd"
	"github.com/kidoman/embd/sensor/ads1115"

	_ "github.com/kidoman/embd/host/all"
)

func main() {
	flag.Parse()

	if err := embd.InitI2C(); err != nil {
		panic(err)
	}
	defer embd.CloseI2C()

	bus := embd.NewI2CBus(1)

	ads := ads1115.New(bus, 0x48, ads1115.RangeQuarterV)
	defer ads.Close()

	for {

		if false {
			v, err := ads.Voltage(ads1115.Diff23)
			if err != nil {
				panic(err)
			}
			fmt.Printf("CH23 diff Voltage=%v scale to %v V \n", v, -v*100332.0/332.0)
		}
		
		if true {
			v, err := ads.Voltage(ads1115.Single1)
			if err != nil {
			panic(err)
			}
			fmt.Printf("CH01 diff Voltage=%v scale to %v A \n", v, -v*1000.0)
			//fmt.Printf("%v A \n", -v*1000.0)
		}
		
		time.Sleep(20 * time.Millisecond)
	}
}
