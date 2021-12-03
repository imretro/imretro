package imretro

import (
	"errors"
	"image"
	"image/color"
	"io"
)

// ImretroSignature is the "magic string" used for identifying an imretro file.
const ImretroSignature = "IMRETRO"

// BitsPerPixelIndex is the position of the two bits for the bits-per-pixel
// mode (7 is left-most).
const bitsPerPixelIndex byte = 6

// DecodeError is an error signifying that something unexpected happened when
// decoding the imretro reader.
type DecodeError string

// DecodeConfig returns the color model and dimensions of an imretro image
// without decoding the entire image.
func DecodeConfig(r io.Reader) (image.Config, error) {
	var buff []byte
	var err error

	buff = make([]byte, len(ImretroSignature)+1)
	mode, err := checkHeader(r, buff)
	if err != nil {
		return image.Config{}, err
	}

	bitsPerPixel := mode & (0b11 << bitsPerPixelIndex)
	hasPalette := mode&WithPalette != 0

	buff = make([]byte, 4)
	_, err = io.ReadFull(r, buff)
	if err != nil {
		return image.Config{}, err
	}

	width := (uint16(buff[0]) << 8) | uint16(buff[1])
	height := (uint16(buff[2]) << 8) | uint16(buff[3])

	var model color.Model
	switch bitsPerPixel {
	case OneBit:
		model, err = decode1bit(r, hasPalette)
	default:
		err = errors.New("Not implemented")
	}

	return image.Config{model, int(width), int(height)}, err
}

func decode1bit(r io.Reader, hasPalette bool) (color.Model, error) {
	if !hasPalette {
		return Default1BitColorModel, nil
	}

	buff := make([]byte, 8)
	if _, err := io.ReadFull(r, buff); err != nil {
		return nil, err
	}
	model := NewOneBitColorModel(ColorFromBytes(buff[:4]), ColorFromBytes(buff[4:]))

	return model, nil
}

// CheckHeader confirms the reader is an imretro image by checking the "magic bytes",
// and returns the "mode".
func checkHeader(r io.Reader, buff []byte) (mode byte, err error) {
	_, err = io.ReadFull(r, buff)
	if err != nil {
		return
	}

	for i, b := range buff[:len(buff)-1] {
		if b != ImretroSignature[i] {
			return mode, DecodeError("unexpected signature byte")
		}
	}
	return buff[len(buff)-1], nil
}

// Error reports that the format could not be decoded as imretro.
func (e DecodeError) Error() string {
	return string(e)
}
