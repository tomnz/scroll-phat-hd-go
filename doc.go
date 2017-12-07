/*
Package scrollphathd provides a helper library for interacting with a Pimoroni Scroll pHAT HD device:

https://shop.pimoroni.com/products/scroll-phat-hd

The library depends on the periph.io framework for low level device communication. There are two primary ways
that the library allows you to interact with the device:

Display wraps the Driver (or any other struct providing appropriate functionality), and extends it with higher
level capabilities, such as an auto-expanding internal buffer, scrolling, flipping, etc. In most cases, the
Display offers a safer and more fully-featured way to interact with the device.

Driver abstracts the low level I2C hardware device, and handles all communication. This does include some basic
drawing functionality such as SetPixel, SetBrightness, and supports rotation. It's possible to use the Driver
directly in your projects. This can be particularly useful in performance-critical situations where you want
to incur minimum overhead.
*/
package scrollphathd
