package device

import (
	"fmt"
	"time"

	"periph.io/x/periph/conn"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/devices"
)

// New returns a new Scroll pHAT HD hardware device.
func New(bus i2c.Bus, opts ...Option) (*ScrollPhatHD, error) {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}

	width, height := devWidth, devHeight
	if options.rotation == Rotation90 || options.rotation == Rotation270 {
		width, height = devHeight, devWidth
	}

	d := &ScrollPhatHD{
		options:    options,
		i2c:        &i2c.Dev{Bus: bus, Addr: addr},
		brightness: 255,
		width:      width,
		height:     height,
	}
	d.setup()
	return d, nil
}

// ScrollPhatHD is a handle to a scrollphathd device.
type ScrollPhatHD struct {
	options options
	// Device handle for I2C bus
	i2c conn.Conn
	// Hardware frame currently in use
	frame         byte
	buffer        [][]byte
	brightness    byte
	width, height int
}

// Width returns the width of the device in pixels.
func (s *ScrollPhatHD) Width() int {
	return s.width
}

// Height returns the height of the device in pixels.
func (s *ScrollPhatHD) Height() int {
	return s.height
}

// SetPixel sets the pixel at the given coordinate to the given value.
func (s *ScrollPhatHD) SetPixel(x, y int, val byte) error {
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
func (s *ScrollPhatHD) SetPixels(pixels [][]byte) error {
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

// SetPixelsUnsafe allows setting all of the pixels at once by swapping out the internal buffer.
// This does NOT copy any of the data. This is exposed for performance reasons, but caution should
// be exercised! If the buffer is later updated externally, the contents of the internal buffer
// will also change!
// The dimensions of the incoming buffer are also not checked.
// When the final values are written to the device via Show, the internal buffer is copied, so
// this may increase safety some.
// Note that the array should be indexed in row, col order.
func (s *ScrollPhatHD) SetPixelsUnsafe(pixels [][]byte) {
	s.buffer = pixels
}

// SetBrightness sets the brightness of the device. This is applied to all pixels on Show.
// 0 is off, 255 is maximum brightness.
func (s *ScrollPhatHD) SetBrightness(brightness byte) {
	s.brightness = brightness
}

// Clear turns off all pixels on the device.
func (s *ScrollPhatHD) Clear() error {
	for _, row := range s.buffer {
		for x := range row {
			row[x] = 0
		}
	}

	return s.Show()
}

// Show renders the contents of the internal buffer to the device. Brightness is applied.
func (s *ScrollPhatHD) Show() error {
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
func (s *ScrollPhatHD) scaleVal(val byte) byte {
	return byte(uint16(val) * uint16(s.brightness) / 255)
}

// pixelAddr maps an x, y coordinate to the physical LED index that should be updated, after rotating
// the coordinates.
func (s *ScrollPhatHD) pixelAddr(x, y int) int {
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
func (s *ScrollPhatHD) setup() error {
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
func (s *ScrollPhatHD) reset() error {
	if err := s.writeRegister(regShutdown, 0); err != nil {
		return err
	}
	time.Sleep(time.Millisecond * 10)
	return s.writeRegister(regShutdown, 1)
}

// writeRegister writes the corresponding value into the given register in the configuration bank
// on the device.
func (s *ScrollPhatHD) writeRegister(register, value byte) error {
	if err := s.bank(configBank); err != nil {
		return err
	}
	return s.write(register, value)
}

// bank switches the active bank on the device. The device uses multiple banks to multiplex the
// amount of data that it needs to access.
func (s *ScrollPhatHD) bank(bank byte) error {
	return s.write(bankAddr, bank)
}

func (s *ScrollPhatHD) write(cmd byte, value ...byte) error {
	msg := []byte{cmd}
	msg = append(msg, value...)
	return s.i2c.Tx(msg, nil)
}

// Halt implements devices.Device.
func (s *ScrollPhatHD) Halt() error {
	return s.writeRegister(regShutdown, 0)
}

// Ensure the device actually implements the periph.io interface.
var _ devices.Device = &ScrollPhatHD{}
