package scrollphathd

import (
	"fmt"
	"testing"

	"periph.io/x/periph/conn/i2c/i2ctest"
)

func TestDriver_New(t *testing.T) {
	// Validate that the driver is able to run through setup without any issues
	testCases := []struct {
		rotation Rotation
	}{
		{rotation: Rotation0},
		{rotation: Rotation90},
		{rotation: Rotation180},
		{rotation: Rotation270},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("rotation %d", tc.rotation), func(t *testing.T) {
			_, err := NewDriver(&i2ctest.Record{}, WithRotation(tc.rotation))
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

// TODO: Validate low level write behavior
