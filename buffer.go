package scrollphathd

// resetBuffer will clear and recreate the buffer with the device's width and height.
func (d *Display) resetBuffer() {
	d.width = d.device.Width()
	d.height = d.device.Height()
	d.buffer = make([][]byte, d.height)
	for y := range d.buffer {
		d.buffer[y] = make([]byte, d.width)
	}
}

// growBuffer will optionally grow the internal buffer as necessary to be able to capture
// the given x, y coordinate.
func (d *Display) growBuffer(newX, newY int) {
	if newX < 0 || newY < 0 {
		panic("coordinates must be 0 or greater")
	}
	if newX < d.width && newY < d.height {
		// Coords already within buffer
		return
	}

	newWidth, newHeight := d.width, d.height
	if newX >= newWidth {
		newWidth = newX + 1
	}
	if newY >= newHeight {
		newHeight = newY + 1
	}

	if newWidth > d.width {
		for y, row := range d.buffer {
			// Expand each row as needed
			newRow := make([]byte, newWidth)
			copy(newRow, row)
			d.buffer[y] = newRow
		}
	}
	for y := d.height - 1; y < newHeight; y++ {
		// Append new rows as needed, of width x coord
		d.buffer = append(d.buffer, make([]byte, newWidth))
	}
	d.width = newWidth
	d.height = newHeight
}
