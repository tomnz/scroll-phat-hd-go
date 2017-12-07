# Scroll pHAT HD Go library

Provides a Go implementation for interfacing with Pimoroni's [Scroll pHAT HD](https://shop.pimoroni.com/products/scroll-phat-hd). The top-level library provides a lot of the same functionality as the reference [Python library](http://docs.pimoroni.com/scrollphathd/), including:

* An internal display buffer that automatically expands as needed.
* Scrolling.
* Flipping.
* Tiling.

Coming soon:

* Text rendering.
* Graph rendering.

## Overview

The library depends on the [periph.io](https://periph.io) framework for low level device communication. There are two primary ways that the library allows you to interact with the device:

Display wraps the Driver (or any other struct providing appropriate functionality), and extends it with the higher level capabilities described above, such as an auto-expanding internal buffer, scrolling, flipping, etc. In most cases, the Display offers a safer and more fully-featured way to interact with the device.

Driver abstracts the low level I2C hardware device, and handles all communication. This does include some basic drawing functionality such as SetPixel, SetBrightness, and support for rotation. It's possible to use the Driver directly in your projects - this can be particularly useful in performance-critical situations where you want to incur minimum overhead in memory usage and copying.

Please refer to the [godocs](https://godoc.org/github.com/tomnz/scroll-phat-hd-go) for full API reference.
