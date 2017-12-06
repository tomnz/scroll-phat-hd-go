# Scroll PHAT HD Go library

Provides a Go implementation for interfacing with Pimoroni's [Scroll pHAT HD product](https://shop.pimoroni.com/products/scroll-phat-hd).

The top-level library provides a lot of the same functionality as the reference [Python library](http://docs.pimoroni.com/scrollphathd/), including:

* Scrolling, flipping and tiling, based on an internal display buffer.
* Text rendering.
* Graph rendering.

This library separates the higher-level rendering functionality from the lower-level hardware driver in the `device` package. Projects may use the driver directly if they do not require the rendering helpers.

Please refer to the [godocs](https://godoc.org/github.com/tomnz/scroll-phat-hd-go) for full API reference.