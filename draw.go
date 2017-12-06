package scrollphathd

// TODO: Implement text and graphing methods for parity with Python lib

// SetPixel sets the given coordinate to the given value.
// Results must be explicitly pushed to the device with Show.
func (d *Display) SetPixel(x, y int, val byte) {
	d.growBuffer(x, y)
	d.buffer[y][x] = val
}

// Fill fills the given rectable with the given value.
// Results must be explicitly pushed to the device with Show.
func (d *Display) Fill(x, y, width, height int, val byte) {
	d.growBuffer(x+width, y+height)
	for ix := 0; ix < width; ix++ {
		for iy := 0; iy < height; iy++ {
			d.buffer[y+iy][x+ix] = val
		}
	}
}

// ClearRect clears the given rectangle.
// Results must be explicitly pushed to the device with Show.
func (d *Display) ClearRect(x, y, width, height int) {
	d.Fill(x, y, width, height, 0)
}
