package imretro

import "testing"

// TestIndexByte tests that a byte can be easily indexed.
func TestIndexByte(t *testing.T) {
	var b byte = 0b0001_0000

	for i := 0; i < 8; i++ {
		var want byte = 0
		if i == 3 {
			want = 1
		}

		t.Logf(`Checking index %d`, i)
		if bit := IndexByte(b, i); bit != want {
			t.Errorf(`b[%d] = %b, want %b`, i, bit, want)
		}
	}
}

// TestByteOutOfBounds tests that a byte cannot receive invalid indices.
func TestByteOutOfBounds(t *testing.T) {
	t.Log(`Testing an OK index`)
	IndexByte(0, 4)
	ExpectOutOfBoundsByte(t, -1)
	ExpectOutOfBoundsByte(t, 8)
}

// ExpectOutOfBoundsByte expects an out-of-bounds error when indexing into a
// byte.
func ExpectOutOfBoundsByte(t *testing.T, index int) {
	t.Helper()
	defer func() { recover() }()
	IndexByte(0, index)
	t.Fatalf(`Failed to get out-of-bounds error with index %d`, index)
}
