/*
Package scrollphathd provides a helper library for interacting with a Pimoroni Scroll pHAT HD device:

https://shop.pimoroni.com/products/scroll-phat-hd

The top-level library attempts to provide similar functionality to the reference Python library, which includes
scrolling, flipping, a dynamically sized display buffer, text rendering, etc.

Users may also take a dependency instead on the lower level device package, which exposes a barebones hardware
driver, eschewing a lot of the more complex features of the top-level Display object. This may be useful in
performance-critical applications.
*/
package scrollphathd
