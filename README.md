# microview

<img src="http://microview.io/images/MicroViewWithProgrammer.png" width="256" align="right">

[![Build Status](https://travis-ci.org/peterhellberg/microview.svg?branch=master)](https://travis-ci.org/peterhellberg/microview)
[![Go Report Card](https://goreportcard.com/badge/github.com/peterhellberg/microview?branch=master)](https://goreportcard.com/report/github.com/peterhellberg/microview)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/peterhellberg/microview)
[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg?style=flat)](https://github.com/peterhellberg/microview/blob/master/LICENSE)


Go library used to remote control a [MicroView](http://microview.io/)

## Requirements

- [MicroView - OLED Arduino Module](https://www.sparkfun.com/products/12923)
- [USB Programmer](https://www.sparkfun.com/products/12924)
- [Arduino IDE](https://www.arduino.cc/en/Main/Software)

## Install

### Go package

    go get -u github.com/peterhellberg/microview

### MicroView Arduino Library

**Note:** This package requires the use of a newer version of the
[MicroView Arduino Library](https://github.com/geekammo/MicroView-Arduino-Library)
(`v1.07b` or later) than what is currently available in the Arduino Library Manager.

So just follow these steps instead:

1. Change directory to Arduino's main directory (`~/Documents/Arduino/` in my case)
2. cd libraries
3. mkdir MicroView
4. cd MicroView
5. git clone https://github.com/geekammo/MicroView-Arduino-Library.git .

### Arduino sketch

Now you need to upload the following sketch to your MicroView using the Arduino IDE

```arduino
#include <MicroView.h>

void setup() {
  uView.begin();
  uView.clear(PAGE);

  Serial.begin(115200);
  Serial.print("MicroView");
}

void loop() {
  uView.checkComm();
}
```

## Example

```go
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
```
