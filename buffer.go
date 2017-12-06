package scrollphathd

// resetBuffer will clear and recreate the buffer with the device's width and height.
func (d *Display) resetBuffer() {
	d.width = d.device.Width()
	d.height = d.device.Height()
	d.buffer = make([][]byte, d.width)
	for x := range d.buffer {
		d.buffer[x] = make([]byte, d.height)
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

	if newY+1 > d.height {
		for x, col := range d.buffer {
			// Expand each column as needed
			newCol := make([]byte, newY+1)
			copy(newCol, col)
			d.buffer[x] = newCol
		}
	}
	for x := d.width; x <= newX; x++ {
		// Insert new columns as needed
		d.buffer = append(d.buffer, make([]byte, newY+1))
	}
	d.width = newX + 1
	d.height = newY + 1
}
