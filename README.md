# Scroll pHAT HD Go library

Provides a Go implementation for interfacing with Pimoroni's [Scroll pHAT HD](https://shop.pimoroni.com/products/scroll-phat-hd). The top-level library provides a lot of the same functionality as the reference [Python library](http://docs.pimoroni.com/scrollphathd/), including:

* An internal display buffer that automatically expands as needed.
* Scrolling.
* Flipping.
* Tiling.

Coming soon:

* Text rendering.
* Graph rendering.

This library separates the higher-level rendering functionality from the lower-level hardware driver in the `device` package. Projects may use the driver directly if they do not require the rendering helpers.

The device driver itself depends on the excellent [periph.io](https://periph.io) framework for communication with the I2C bus - make sure you've made it available in your Go path with `dep ensure`, or `go get -u periph.io/x/periph`.

Please refer to the [godocs](https://godoc.org/github.com/tomnz/scroll-phat-hd-go) for full API reference.