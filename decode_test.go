package imretro

import (
	"bytes"
	"io"
	"testing"
)

// TestPassCheckSignature tests that a reader starting with "IMRETRO" bytes will
// pass.
func TestPassCheckSignature(t *testing.T) {
	buff := make([]byte, 8)
	r := Make1Bit(t)
	mode, err := checkSignature(r, buff)
	if err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}
	if pixelMode := mode & (0b1100_0000); pixelMode != EightBit {
		t.Errorf(
			`pixelMode = %d (%08b), want %d (%08b)`,
			pixelMode, pixelMode,
			EightBit, EightBit,
		)
	}
	if hasPalette := mode & (0b0010_0000); hasPalette != 0x20 {
		t.Error("mode does not signify in-file palette")
	}
}

// TestFailCheckSignature tests that a reader with unexpected magic bytes will
// fail.
func TestFailCheckSignature(t *testing.T) {
	buff := make([]byte, 8)
	partialSignature := "IMRET"
	jpgSignature := "\xFF\xD8\xFF\xE0\x00\x10\x4A\x46\x49\x46\x00\x01"

	partialr := bytes.NewBufferString(partialSignature)
	if _, err := checkSignature(partialr, buff); err != io.ErrUnexpectedEOF {
		t.Errorf(`err = %v, want %v`, err, io.ErrUnexpectedEOF)
	}

	jpgr := bytes.NewBufferString(jpgSignature)
	if _, err := checkSignature(jpgr, buff); err != DecodeError("unexpected signature byte") {
		t.Fatalf(`err = %v, want %v`, err, DecodeError("unexpected signature byte"))
	}
}

// Make1Bit makes a 1-bit imretro reader.
func Make1Bit(t *testing.T) *bytes.Buffer {
	t.Helper()
	return bytes.NewBuffer([]byte{
		// signature/magic bytes
		'I', 'M', 'R', 'E', 'T', 'R', 'O',
		// Mode byte (8-bit, in-file palette)
		0b1010_0000,
	})
}
