package device

import "fmt"

// Option allows specifying behavior for the device.
type Option func(*options)

// WithGamma allows overriding the gamma curve for the device. Must include 256
// level mappings.
func WithGamma(gamma []byte) Option {
	return func(options *options) {
		if len(gamma) != 256 {
			panic("Must pass 256 gamma levels")
		}
		options.gamma = gamma
	}
}

// Rotation specifies a rotation amount.
type Rotation uint16

const (
	// Rotation0 specifies no display rotation.
	Rotation0 Rotation = 0
	// Rotation90 specifies 90 degree display rotation.
	Rotation90 Rotation = 90
	// Rotation180 specifies 180 degree display rotation.
	Rotation180 Rotation = 180
	// Rotation270 specifies 270 degree display rotation.
	Rotation270 Rotation = 270
)

// WithRotation applies rotation to the internal buffer before pushing pixels to the device.
// Note that this can alter the final width/height of the device. If you need to dynamically check
// these values, use the Width and Height functions.
func WithRotation(rotation Rotation) Option {
	return func(options *options) {
		if rotation != Rotation0 && rotation != Rotation90 && rotation != Rotation180 && rotation != Rotation270 {
			panic(fmt.Sprintf("received invalid rotation %d - must be a right angle", rotation))
		}
		options.rotation = rotation
	}
}

type options struct {
	gamma    []byte
	rotation Rotation
}

var defaultOptions = options{
	gamma:    defaultGamma,
	rotation: Rotation0,
}
