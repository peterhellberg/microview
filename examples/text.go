// +build example
//
// Do not build by default.

package main

import (
	"flag"
	"log"
	"strings"

	. "github.com/peterhellberg/microview"
)

func main() {
	name := flag.String("name", "/dev/cu.usbserial-DA00SSM3", "name of the serial port")
	text := flag.String("text", "Hello From Go!", "text to display on the MicroView")

	flag.Parse()

	mv, err := OpenMicroView(*name)
	if err != nil {
		log.Fatal(err)
	}

	for i, s := range strings.Split(*text, " ") {
		mv.DrawString(0, uint8(10*i), s)
	}
}
