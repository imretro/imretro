package imretro

import (
	"bytes"
	"io"
	"testing"
)

// TestPassCheckSignature tests that a reader starting with "IMRETRO" bytes will
// pass.
func TestPassCheckSignature(t *testing.T) {
	r := Make1Bit(t)
	if err := checkSignature(r); err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}
}

// TestFailCheckSignature tests that a reader with unexpected magic bytes will
// fail.
func TestFailCheckSignature(t *testing.T) {
	partialSignature := "IMRET"
	jpgSignature := "\xFF\xD8\xFF\xE0\x00\x10\x4A\x46\x49\x46\x00\x01"

	partialr := bytes.NewBufferString(partialSignature)
	if err, want := checkSignature(partialr), io.ErrUnexpectedEOF; err != want {
		t.Errorf(`err = %v, want %v`, err, want)
	}

	jpgr := bytes.NewBufferString(jpgSignature)
	if err, want := checkSignature(jpgr), DecodeError("unexpected signature byte"); err != want {
		t.Fatalf(`err = %v, want %v`, err, want)
	}
}

// Make1Bit makes a 1-bit imretro reader.
func Make1Bit(t *testing.T) *bytes.Buffer {
	t.Helper()
	return bytes.NewBuffer([]byte{
		// signature/magic bytes
		'I', 'M', 'R', 'E', 'T', 'R', 'O',
	})
}
