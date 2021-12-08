package imretro

import "testing"

// TestModelBitMode checks that the correct bit mode is interpreted from the
// color model.
func TestModelBitMode(t *testing.T) {
	if mode, err := ModelBitMode(Default1BitColorModel); err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	} else if mode != OneBit {
		t.Errorf(
			`mode = %v (%08b), want %v (%08b)`,
			mode, mode,
			OneBit, OneBit,
		)
	}
}