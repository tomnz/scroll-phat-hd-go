package scrollphathd

import "fmt"

// DriverOption allows specifying behavior for the driver.
type DriverOption func(*driverOptions)

// WithGamma allows overriding the gamma curve for the driver. Must include 256
// level mappings.
func WithGamma(gamma []byte) DriverOption {
	return func(options *driverOptions) {
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
func WithRotation(rotation Rotation) DriverOption {
	return func(options *driverOptions) {
		if rotation != Rotation0 && rotation != Rotation90 && rotation != Rotation180 && rotation != Rotation270 {
			panic(fmt.Sprintf("received invalid rotation %d - must be a right angle", rotation))
		}
		options.rotation = rotation
	}
}

type driverOptions struct {
	gamma    []byte
	rotation Rotation
}

var defaultDriverOptions = driverOptions{
	gamma:    defaultGamma,
	rotation: Rotation0,
}
