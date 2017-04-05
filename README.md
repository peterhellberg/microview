# microview

[![Build Status](https://travis-ci.org/peterhellberg/microview.svg?branch=master)](https://travis-ci.org/peterhellberg/microview)
[![Go Report Card](https://goreportcard.com/badge/github.com/peterhellberg/microview)](https://goreportcard.com/report/github.com/peterhellberg/microview)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/peterhellberg/microview)
[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg?style=flat)](https://github.com/peterhellberg/microview/blob/master/LICENSE)


Go library used to remote control a [MicroView](http://microview.io/).

## Requirements

- [MicroView - OLED Arduino Module](https://www.sparkfun.com/products/12923)
- [USB Programmer](https://www.sparkfun.com/products/12924)
- [Arduino IDE](https://www.arduino.cc/en/Main/Software)

## Install

### Go package

    go get -u github.com/peterhellberg/microview

### MicroView Arduino Library

**Note:** This package requires the use of a newer version (`v1.07b` or later) of the
[MicroView Arduino Library](https://github.com/geekammo/MicroView-Arduino-Library)
than what is currently available in the Arduino Library Manager.

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

![MicroView With Programmer](http://microview.io/images/MicroViewWithProgrammer.png)
