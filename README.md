# Scroll pHAT HD Go library

[![build](https://travis-ci.org/tomnz/scroll-phat-hd-go.svg?branch=master)](https://travis-ci.org/tomnz/scroll-phat-hd-go)
[![godocs](https://godoc.org/github.com/tomnz/scroll-phat-hd-go?status.svg)](https://godoc.org/github.com/tomnz/scroll-phat-hd-go)

Provides a Go implementation for interfacing with Pimoroni's [Scroll pHAT HD](https://shop.pimoroni.com/products/scroll-phat-hd). The top-level library provides a lot of the same functionality as the reference [Python library](http://docs.pimoroni.com/scrollphathd/), including:

* An internal display buffer that automatically expands as needed.
* Scrolling.
* Flipping.
* Tiling.

Coming soon:

* Text rendering.
* Graph rendering.

## Overview

There are two primary ways that the library allows you to interact with the device:

`Display` wraps the `Driver` (or any other struct providing appropriate functionality), and extends it with the higher level capabilities described above, such as an auto-expanding internal buffer, scrolling, flipping, etc. In most cases, the `Display` offers a safer and more fully-featured way to interact with the device.

`Driver` abstracts the low-level I2C hardware device, and handles all communication. This does include some basic drawing functionality such as `SetPixel`, `SetBrightness`, and support for rotation. It's possible to use the `Driver` directly in your projects - this can be particularly useful in performance-critical situations where you want to incur minimum overhead in memory usage and copying.

## Installation

First, clone the project into your GOPATH:

```bash
go get github.com/tomnz/scroll-phat-hd-go
```

The library depends on the [periph.io](https://periph.io) framework for low level device communication. You can install this manually with `go get`, or (preferred) use `dep`:

```bash
go get -u github.com/golang/dep/cmd/dep
cd $GOPATH/src/github.com/tomnz/scroll-phat-hd-go
dep ensure
```

## Usage

First, initialize a periph.io I2C bus, and instantiate the display with it:

```go
package main

import (
    "github.com/tomnz/scroll-phat-hd-go"
    "periph.io/x/periph/conn/i2c/i2creg"
    "periph.io/x/periph/host"
)

func main() {
    // TODO: Handle errors
    _, _ := host.Init()
    bus, _ := i2creg.Open("1")
    display, _ := scrollphathd.New(bus)
}
```

Now, you can use `display` to interact with the hardware. For example:

```go
display.SetBrightness(127)
display.Fill(0, 0, 5, 5, 255)
display.Show()
```

Please refer to the [godocs](https://godoc.org/github.com/tomnz/scroll-phat-hd-go) for full API reference.

## Contributing

Contributions welcome! Please refer to the [contributing guide](https://github.com/tomnz/scroll-phat-hd-go/blob/master/CONTRIBUTING.md).
