package imretro

import (
	"errors"
	"image"
	"image/color"
	"io"

	"github.com/spenserblack/go-byteutils"
)

// ImretroSignature is the "magic string" used for identifying an imretro file.
const ImretroSignature = "IMRETRO"

// BitsPerPixelIndex is the position of the two bits for the bits-per-pixel
// mode (7 is left-most).
const bitsPerPixelIndex byte = 6

// DecodeError is an error signifying that something unexpected happened when
// decoding the imretro reader.
type DecodeError string

// Decode decodes an image in the imretro format.
func Decode(r io.Reader) (ImretroImage, error) {
	config, err := DecodeConfig(r)
	if err != nil {
		return nil, err
	}
	mode, err := ModelBitMode(config.ColorModel)
	if err != nil {
		return nil, err
	}
	pixels, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	switch mode {
	case OneBit:
		return &image1Bit{imretroImage{config, pixels}}, nil
	case TwoBit:
		return &image2Bit{imretroImage{config, pixels}}, nil
	case EightBit:
		return &image8Bit{imretroImage{config, pixels}}, nil
	}
	return nil, errors.New("Not implemented")
}

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
	hasPalette := byteutils.BitAsBool(byteutils.GetL(mode, PaletteIndex))

	buff = make([]byte, 4)
	_, err = io.ReadFull(r, buff)
	if err != nil {
		return image.Config{}, err
	}

	width := byteutils.ToUint16(buff[0:2], byteutils.LittleEndian)
	height := byteutils.ToUint16(buff[2:4], byteutils.LittleEndian)

	var model color.Model
	switch bitsPerPixel {
	case OneBit:
		model, err = decode1bitModel(r, hasPalette)
	case TwoBit:
		model, err = decode2bitModel(r, hasPalette)
	case EightBit:
		model, err = decode8bitModel(r, hasPalette)
	default:
		err = errors.New("Not implemented")
	}

	return image.Config{model, int(width), int(height)}, err
}

func decode1bitModel(r io.Reader, hasPalette bool) (color.Model, error) {
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

func decode2bitModel(r io.Reader, hasPalette bool) (color.Model, error) {
	if !hasPalette {
		return Default2BitColorModel, nil
	}

	buff := make([]byte, 16)
	if _, err := io.ReadFull(r, buff); err != nil {
		return nil, err
	}
	model := NewTwoBitColorModel(
		ColorFromBytes(buff[:4]),
		ColorFromBytes(buff[4:8]),
		ColorFromBytes(buff[8:12]),
		ColorFromBytes(buff[12:]),
	)

	return model, nil
}

func decode8bitModel(r io.Reader, hasPalette bool) (color.Model, error) {
	if !hasPalette {
		return Default8BitColorModel, nil
	}

	colors := make(Palette, 0, 256)
	buff := make([]byte, 4)

	for i := 0; i < cap(colors); i++ {
		if _, err := io.ReadFull(r, buff); err != nil {
			return nil, err
		}
		colors = append(colors, ColorFromBytes(buff))
	}

	return NewEightBitColorModel(colors), nil
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

func init() {
	image.RegisterFormat("imretro", ImretroSignature, globalDecode, DecodeConfig)
}

// GlobalDecode returns an image.Image instead of an ImretroImage so that it
// can be registered as a format.
func globalDecode(r io.Reader) (image.Image, error) {
	i, err := Decode(r)
	return i.(image.Image), err
}
