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
		RectFill(5, 5, 5, 15),
		RectFill(25, 0, 30, 15),
		Rect(1, 1, 20, 40),
		Rect(40, 20, 20, 20),
		Rect(40, 20, 15, 15),
		Rect(40, 20, 10, 10),
		Rect(40, 20, 5, 5),
	)
}
