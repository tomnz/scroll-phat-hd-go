package scrollphathd

import (
	"fmt"
	"time"

	"periph.io/x/periph/conn"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/devices"
)

// NewDriver returns a new Scroll pHAT HD hardware driver. This implements the Device
// interface, and can be used by Display. Connects to the device on the given I2C bus
// at its standard address.
func NewDriver(bus i2c.Bus, opts ...DriverOption) (*Driver, error) {
	return NewDriverWithConn(&i2c.Dev{Bus: bus, Addr: addr}, opts...)
}

// NewDriverWithConn returns a new Scroll pHAT HD hardware driver, using the given periph.io
// conn.Conn object. Typically NewDriver should be used instead, but this may be useful for
// testing using mocks, or a custom I2C connection with a different address than the default.
func NewDriverWithConn(periphConn conn.Conn, opts ...DriverOption) (*Driver, error) {
	options := defaultDriverOptions
	for _, opt := range opts {
		opt(&options)
	}

	width, height := devWidth, devHeight
	if options.rotation == Rotation90 || options.rotation == Rotation270 {
		width, height = devHeight, devWidth
	}

	d := &Driver{
		options:    options,
		i2c:        periphConn,
		brightness: 255,
		width:      width,
		height:     height,
	}
	if err := d.setup(); err != nil {
		return nil, err
	}
	return d, nil
}

// Driver handles low level communication with a Scroll pHAT HD hardware device. It can
// be used directly if you do not have need for some of the higher-level features of the
// Display object.
type Driver struct {
	options driverOptions
	// Device handle for I2C bus
	i2c conn.Conn
	// Hardware frame currently in use
	frame         byte
	buffer        [][]byte
	brightness    byte
	width, height int
}

// Width returns the width of the device in pixels.
func (s *Driver) Width() int {
	return s.width
}

// Height returns the height of the device in pixels.
func (s *Driver) Height() int {
	return s.height
}

// SetPixel sets the pixel at the given coordinate to the given value.
func (s *Driver) SetPixel(x, y int, val byte) error {
	if x < 0 || x > s.width-1 {
		return fmt.Errorf("received invalid x coordinate %d", x)
	}
	if y < 0 || y > s.height-1 {
		return fmt.Errorf("received invalid y coordinate %d", y)
	}
	s.buffer[x][y] = val
	return nil
}

// SetPixels copies all of the given pixels at once to the internal buffer.
// Dimensions of the incoming buffer are checked to ensure they match the width and height of
// the device.
// Note that the array should be indexed in row, col order.
func (s *Driver) SetPixels(pixels [][]byte) error {
	if len(pixels) != s.height {
		return fmt.Errorf("received invalid buffer of height %d", len(pixels))
	}

	for y, row := range pixels {
		if len(row) != s.width {
			return fmt.Errorf("received invalid buffer with row %d of width %d", y, len(row))
		}
		for x, val := range row {
			s.buffer[y][x] = val
		}
	}
	return nil
}

// SetBuffer allows setting all of the pixels at once by swapping out the internal buffer.
// This does NOT copy any of the data. This is exposed for performance reasons, but caution should
// be exercised! If the buffer is later updated externally, the contents of the internal buffer
// will also change!
// The dimensions of the incoming buffer are also not checked.
// When the final values are written to the device via Show, the internal buffer is copied, so
// this may increase safety some.
// Note that the array should be indexed in row, col order.
func (s *Driver) SetBuffer(buffer [][]byte) {
	s.buffer = buffer
}

// SetBrightness sets the brightness of the device. This is applied to all pixels on Show.
// 0 is off, 255 is maximum brightness.
func (s *Driver) SetBrightness(brightness byte) {
	s.brightness = brightness
}

// Clear turns off all pixels on the device.
func (s *Driver) Clear() error {
	for _, row := range s.buffer {
		for x := range row {
			row[x] = 0
		}
	}

	return s.Show()
}

// Show renders the contents of the internal buffer to the device. Brightness is applied.
func (s *Driver) Show() error {
	output := make([]byte, 144)
	for y, row := range s.buffer {
		for x, val := range row {
			output[s.pixelAddr(x, y)] = s.options.gamma[s.scaleVal(val)]
		}
	}

	nextFrame := (s.frame + 1) % 2
	if err := s.bank(nextFrame); err != nil {
		return err
	}

	// Write the pixel data in chunks
	offset := byte(0)
	for len(output) > chunkSize {
		if err := s.write(offsetColor+offset, output[:chunkSize]...); err != nil {
			return err
		}
		output = output[chunkSize:]
		offset += chunkSize
	}
	// Write rest
	if err := s.write(offsetColor+offset, output...); err != nil {
		return err
	}
	// Switch the active frame to the new frame
	if err := s.writeRegister(regFrame, nextFrame); err != nil {
		return err
	}
	s.frame = nextFrame
	return nil
}

// scaleVal applies brightness to the given value.
func (s *Driver) scaleVal(val byte) byte {
	return byte(uint16(val) * uint16(s.brightness) / 255)
}

// pixelAddr maps an x, y coordinate to the physical LED index that should be updated, after rotating
// the coordinates.
func (s *Driver) pixelAddr(x, y int) int {
	switch s.options.rotation {
	case Rotation0:
	case Rotation90:
		x, y = devWidth-1-y, x
	case Rotation180:
		x, y = devWidth-1-x, devHeight-1-y
	case Rotation270:
		x, y = y, devHeight-1-x
	default:
		panic("unknown rotation")
	}

	y = s.height - y - 1
	if x > 8 {
		x -= 8
		y = -y - 2
	} else {
		x = 8 - x
	}
	return x*16 + y
}

// setup performs initial setup of the I2C hardware device by sending initialization messages.
func (s *Driver) setup() error {
	if err := s.reset(); err != nil {
		return err
	}

	if err := s.writeRegister(regFrame, 0); err != nil {
		return err
	}
	if err := s.writeRegister(regMode, modePicture); err != nil {
		return err
	}
	if err := s.writeRegister(regAudioSync, 0); err != nil {
		return err
	}

	// Need to "turn on" all of the LEDs with an enable bit in the frames
	// that we are going to use
	enableRows := make([]byte, 17)
	for i := range enableRows {
		enableRows[i] = 255
	}

	if err := s.bank(1); err != nil {
		return err
	}
	if err := s.write(offsetEnable, enableRows...); err != nil {
		return err
	}
	if err := s.bank(0); err != nil {
		return err
	}
	if err := s.write(offsetEnable, enableRows...); err != nil {
		return err
	}

	s.buffer = make([][]byte, s.height)
	for y := range s.buffer {
		s.buffer[y] = make([]byte, s.width)
	}

	return s.Clear()
}

// reset reboots the hardware device.
func (s *Driver) reset() error {
	if err := s.writeRegister(regShutdown, 0); err != nil {
		return err
	}
	time.Sleep(time.Millisecond * 10)
	return s.writeRegister(regShutdown, 1)
}

// writeRegister writes the corresponding value into the given register in the configuration bank
// on the device.
func (s *Driver) writeRegister(register, value byte) error {
	if err := s.bank(configBank); err != nil {
		return err
	}
	return s.write(register, value)
}

// bank switches the active bank on the device. The device uses multiple banks to multiplex the
// amount of data that it needs to access.
func (s *Driver) bank(bank byte) error {
	return s.write(bankAddr, bank)
}

func (s *Driver) write(cmd byte, value ...byte) error {
	msg := []byte{cmd}
	msg = append(msg, value...)
	return s.i2c.Tx(msg, nil)
}

// Halt implements devices.Device.
func (s *Driver) Halt() error {
	return s.writeRegister(regShutdown, 0)
}

// Ensure the device actually implements the periph.io interface.
var _ devices.Device = &Driver{}
