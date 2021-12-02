package imretro

import (
	"errors"
	"image"
	"io"
)

// ImretroSignature is the "magic string" used for identifying an imretro file.
const ImretroSignature = "IMRETRO"

// DecodeError is an error signifying that something unexpected happened when
// decoding the imretro reader.
type DecodeError string

// DecodeConfig returns the color model and dimensions of an imretro image
// without decoding the entire image.
func DecodeConfig(r io.Reader) (image.Config, error) {
	return image.Config{}, errors.New("Not implemented")
}

// CheckSignature confirms the reader is an imretro image by checking the "magic string".
func checkSignature(r io.Reader) error {
	buff := make([]byte, len(ImretroSignature))
	_, err := io.ReadFull(r, buff)
	if err != nil {
		return err
	}

	for i, b := range buff {
		if b != ImretroSignature[i] {
			return DecodeError("unexpected signature byte")
		}
	}
	return nil
}

// Error reports that the format could not be decoded as imretro.
func (e DecodeError) Error() string {
	return string(e)
}
