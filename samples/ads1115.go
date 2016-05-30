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

	for n := 0; n < 10; n++ {
		channels := []byte{
			ads1115.Diff01,
			ads1115.Diff23,
		}

		fmt.Printf("-------------------\n")

		for _, ch := range channels {

			for i := 0; i < 1; i++ {

				v, err := ads.Voltage(ch)
				if err != nil {
					panic(err)
				}
				fmt.Printf("CH=%v Voltage=%v scale to %5.2f V scale to %5.2f A\n", ch, v, -v*99370.0/331.42, v*1000)
				//time.Sleep(125 * time.Millisecond)

			}

		}
		time.Sleep(1000 * time.Millisecond)
	}
}
