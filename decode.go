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

// Image1Bit is the underlying type for a 1-bit mode image.
type image1Bit struct {
	config image.Config
	pixels []byte
}

// Image2Bit is the underlying type for a 2-bit mode image.
type image2Bit struct {
	config image.Config
	pixels []byte
}

// Decode decodes an image in the imretro format.
func Decode(r io.Reader) (image.Image, error) {
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
		return &image1Bit{config, pixels}, nil
	case TwoBit:
		return &image2Bit{config, pixels}, nil
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

// ColorModel returns the Image's color model.
func (i *image1Bit) ColorModel() color.Model {
	return i.config.ColorModel
}

// Bounds returns the boundaries of the image.
func (i *image1Bit) Bounds() image.Rectangle {
	return image.Rect(0, 0, i.config.Width, i.config.Height)
}

// At returns the color at the given pixel.
func (i *image1Bit) At(x, y int) color.Color {
	if !image.Pt(x, y).In(i.Bounds()) {
		return NoColor
	}
	index := (y * i.config.Width) + x
	byteIndex := index / 8
	bitIndex := byte(index % 8)

	if byteIndex > len(i.pixels) {
		return NoColor
	}

	b := i.pixels[byteIndex]
	bit := byteutils.GetL(b, bitIndex)

	model := i.ColorModel().(OneBitColorModel)
	return model.colors[int(bit)]
}

// ColorModel returns the Image's color model.
func (i *image2Bit) ColorModel() color.Model {
	return i.config.ColorModel
}

// Bounds returns the boundaries of the image.
func (i *image2Bit) Bounds() image.Rectangle {
	return image.Rect(0, 0, i.config.Width, i.config.Height)
}

// At returns the color at the given pixel.
func (i *image2Bit) At(x, y int) color.Color {
	if !image.Pt(x, y).In(i.Bounds()) {
		return NoColor
	}
	index := (y * i.config.Width) + x
	byteIndex := index / 4
	bitIndex := byte(index%4) * 2

	bits := byteutils.SliceL(i.pixels[byteIndex], bitIndex, bitIndex+2)

	model := i.ColorModel().(TwoBitColorModel)
	return model.colors[int(bits)]
}

func init() {
	image.RegisterFormat("imretro", ImretroSignature, Decode, DecodeConfig)
}
