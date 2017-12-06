package scrollphathd

// New instantiates a new Scroll pHAT HD display, using the supplied device and options.
// Typical usage will be to provide a hardware device. For example (include error handling!)
//
//	import (
//		"github.com/tomnz/scroll-phat-hd-go/device"
//		"periph.io/x/periph/conn/i2c/i2creg"
//		"periph.io/x/periph/host"
//	)
//	_, _ := host.Init()
//	bus, _ := i2creg.Open("1")
//	device, _ := device.New(bus)
//  display := scrollphathd.New(device)
//
// Because the device is an interface, it will also accept mocks, or alternative output
// implementations.
func New(device Device, opts ...Option) *Display {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}

	outBuf := make([][]byte, device.Width())
	for x := range outBuf {
		outBuf[x] = make([]byte, device.Height())
	}

	d := &Display{
		options: options,
		device:  device,
		outBuf:  outBuf,
	}
	d.resetBuffer()
	return d
}

// Display is the core struct that manages the display state.
type Display struct {
	options options
	device  Device
	buffer  [][]byte
	width, height,
	scrollX, scrollY int
	flipX, flipY bool

	// We maintain the output buffer for the device ourselves, to reduce the amount of
	// memory allocation that goes on
	outBuf [][]byte

	// TODO: Make this goroutine-safe? Would involve wrapping any buffer operations with a mutex.
}

// Device is an abstraction that defines the capabilities that the display requires from
// its actual device (hardware or otherwise).
type Device interface {
	SetPixel(x, y int, val byte) error
	SetPixels(pixels [][]byte) error
	SetPixelsUnsafe(pixels [][]byte)
	SetBrightness(brightness byte)
	Clear() error
	Show() error
	Width() int
	Height() int
}

// TODO: Implement text and graphing methods for parity with Python lib

// SetPixel sets the given coordinate to the given value.
// Results must be explicitly pushed to the device with Show.
func (d *Display) SetPixel(x, y int, val byte) {
	d.growBuffer(x, y)
	d.buffer[x][y] = val
}

// Fill fills the given rectable with the given value.
// Results must be explicitly pushed to the device with Show.
func (d *Display) Fill(x, y, width, height int, val byte) {
	d.growBuffer(x+width, y+height)
	for ix := 0; ix < width; ix++ {
		for iy := 0; iy < height; iy++ {
			d.buffer[x+ix][y+iy] = val
		}
	}
}

// ClearRect clears the given rectangle.
// Results must be explicitly pushed to the device with Show.
func (d *Display) ClearRect(x, y, width, height int) {
	d.Fill(x, y, width, height, 0)
}

// Show renders the current state of the display to the device. Scrolling and flipping are applied,
// and the relevant subset of the display is sent to the device for actual rendering.
func (d *Display) Show() {
	for x, col := range d.outBuf {
		for y := range col {
			col[y] = d.getSourcePixel(x, y)
		}
	}
	d.device.SetPixelsUnsafe(d.outBuf)
	d.device.Show()
}

func (d *Display) getSourcePixel(devX, devY int) byte {
	x := devX
	x += d.scrollX
	if d.options.tile {
		x %= d.width
	}
	if d.flipX {
		x = d.width - x - 1
	}

	// Fail early if x is nonsense
	if x >= d.width {
		return 0
	}

	y := devY
	y += d.scrollY
	if d.options.tile {
		y %= d.height
	}
	if d.flipY {
		y = d.height - y - 1
	}

	if y >= d.height {
		return 0
	}

	return d.buffer[x][y]
}

// Clear clears the entire display.
func (d *Display) Clear() {
	d.resetBuffer()
	d.device.Clear()
}
