package scrollphathd

// DisplayOption allows specifying behavior for the display.
type DisplayOption func(*displayOptions)

// WithTiling specifies whether the buffer should tile when scrolling (default true).
func WithTiling(tile bool) DisplayOption {
	return func(options *displayOptions) {
		options.tile = tile
	}
}

type displayOptions struct {
	tile bool
}

var defaultDisplayOptions = displayOptions{
	tile: true,
}
