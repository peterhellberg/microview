// +build example
//
// Do not build by default.

package main

import (
	"flag"
	"log"
	"time"

	. "github.com/peterhellberg/microview"
)

func main() {
	name := flag.String("name", "/dev/cu.usbserial-DA00SSM3", "name of the serial port")

	flag.Parse()

	mv, err := OpenMicroView(*name, Delay(90*time.Millisecond))
	if err != nil {
		log.Fatal(err)
	}

	mv.Run(
		CircleFill(5, 5, 5),
		Circle(1, 1, 20),
		Circle(40, 20, 20),
		Circle(40, 20, 15),
		Circle(40, 20, 10),
		Circle(40, 20, 5),
	)
}
