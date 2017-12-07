package scrollphathd

import (
	"periph.io/x/periph/conn/i2c"
)

// New instantiates a new Scroll pHAT HD display, the supplied options. This method requires
// an I2C bus to be supplied, which will be used to connect to the actual hardware device.
// For example:
//
//	import (
//		"github.com/tomnz/scroll-phat-hd-go"
//		"periph.io/x/periph/conn/i2c/i2creg"
//		"periph.io/x/periph/host"
//	)
//	_, _ := host.Init()
//	bus, _ := i2creg.Open("1")
//	display, _ := scrollphathd.New(bus)
//
func New(bus i2c.Bus, opts ...DisplayOption) (*Display, error) {
	device, err := NewDriver(bus)
	if err != nil {
		return nil, err
	}
	disp := NewWithDevice(device, opts...)
	return disp, nil
}

// NewWithDevice instantiates a new Scroll pHAT HD display, using the supplied device and
// options. For using the standard I2C hardware device, you likely want to just use New
// instead.
// This constructor is useful for passing a non-standard device implementation, such as a
// mock or terminal emulator.
// You can also override some of the options for the standard driver by declaring it first,
// then passing it to this constructor.
func NewWithDevice(device Device, opts ...DisplayOption) *Display {
	options := defaultDisplayOptions
	for _, opt := range opts {
		opt(&options)
	}

	outBuf := make([][]byte, device.Height())
	for y := range outBuf {
		outBuf[y] = make([]byte, device.Width())
	}

	d := &Display{
		options: options,
		device:  device,
		outBuf:  outBuf,
	}
	d.resetBuffer()
	return d
}

// Display is the primary struct for interacting with the Scroll pHAT HD device.
type Display struct {
	options displayOptions
	device  Device
	buffer  [][]byte
	width, height,
	scrollX, scrollY int
	flipX, flipY bool

	// We maintain the output buffer for the device ourselves, to reduce the amount of
	// memory allocation and copying that goes on
	outBuf [][]byte

	// TODO: Make this goroutine-safe? Would involve wrapping any buffer operations with a mutex.
}

// Device is an abstraction that defines the capabilities that the display requires from
// its actual device (hardware or otherwise).
type Device interface {
	SetBuffer(buffer [][]byte)
	SetBrightness(brightness byte)
	Show() error
	Width() int
	Height() int
}

// SetBrightness configures the display's brightness.
// 0 is off, 255 is maximum brightness.
func (d *Display) SetBrightness(brightness byte) {
	d.device.SetBrightness(brightness)
}

// SetFlip configures flipping for the display.
func (d *Display) SetFlip(flipX, flipY bool) {
	d.flipX = flipX
	d.flipY = flipY
}

// ScrollTo configures the top left coordinate to use from the buffer for display.
func (d *Display) ScrollTo(scrollX, scrollY int) {
	d.scrollX = scrollX
	d.scrollY = scrollY
}

// Scroll scrolls the buffer relative to its current position.
func (d *Display) Scroll(deltaX, deltaY int) {
	d.scrollX += deltaX
	d.scrollY += deltaY
}

// Show renders the current state of the display to the device. Scrolling and flipping are applied,
// and the relevant subset of the display is sent to the device for actual rendering.
func (d *Display) Show() {
	for y, row := range d.outBuf {
		for x := range row {
			row[x] = d.getSourcePixel(x, y)
		}
	}
	d.device.SetBuffer(d.outBuf)
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
	if x < 0 || x >= d.width {
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

	if y < 0 || y >= d.height {
		return 0
	}

	return d.buffer[y][x]
}

// Clear clears the entire display.
func (d *Display) Clear() {
	d.resetBuffer()
	d.Show()
}
