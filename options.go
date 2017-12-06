package scrollphathd

// Option allows specifying behavior for the display.
type Option func(*options)

// WithTiling specifies whether the buffer should tile when scrolling (default true).
func WithTiling(tile bool) Option {
	return func(options *options) {
		options.tile = tile
	}
}

type options struct {
	tile bool
}

var defaultOptions = options{
	tile: true,
}
