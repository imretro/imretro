# imretro

An image format for retro-style images

I made this format so that you can make your own image with a hex editor for a retro-ish
artstyle. In the future I may need to create an editor to allow making larger images easily.

:warning: This image format is *not* for color accuracy, storage efficiency, or honestly anything
useful. I went through part of an SNES programming tutorial and was inspired by the idea of
manually writing a sprite's data using a hex editor. If there is any useful feature, it is
that this image format can easily be manually written by anyone knowledgable of binary numbers.

## File

### Contents

#### Header

Each file should start with a header to provide some information about the image.

##### Signature

The first 7 bytes are `"IMRETRO"`. This is the file signature.

##### Mode Byte

The next byte is used for flags for enabled/disabled features used by the decoder.

The first two bits declare the bits-per-pixel of the image. `00` for 1-bit, `01` for
2-bit, and `10` for 8-bit.

Next is a single bit for palette usage: `0` for no palette, `1` to declare that
the file contains a palette. When no palette is present in the file, this means that
the file decoder should choose its own default palette.

The following 2 bits are unused. They are reserved for potential future features.

The sixth and seventh bits are a flag for how many channels each color in the
in-file palette will have. `00` for grayscale (1 channel), `01` for RGB (3
channels), and `10` for RGBA (4 channels). The decoder should ignore this flag
if the in-file palette flag is not set.

The eighth bit is a flag for color accuracy for the [in-file palette][palette].
`0` for 2 bits per color channel, `1` for 8 bits (a byte) per color
channel. The decoder should ignore this flag if the in-file palette flag is not
set.

#### Dimensions

Following that are 24 bits for the dimensions: 12 for width and 12 for height. The reason for this limited
range of dimensions is to be faithful to the retro-ish goal of this format. In fact, the maximum dimensions
are fairly large, but with the goal of supporting at least 480p.

This results in 11 bytes for the header.

#### Palette (Optional)

The palette will declare the possible colors in the image. The number of colors in your
palette depend on the number of bits you chose to use in your header. In 1-Bit mode, you
will declare 2 colors, in 2-Bit mode, 4 colors, and in 8-bit mode you will declare 256
colors.

##### Note on Byte Count

The bits used to declare the palette depends on the color channel and color accuracy flags
in the modes.

For example, in 8-bit mode with 4 color channels and the color accuracy flag
set (8 bits per channel), the palette would be 8192 bits, or 1024 bytes.

Conversely, in 1-bit mode with 1 color channel and the color accuracy flag
unset (2 bits per channel), only 4 bits would need to be written. Note that this
is less than a full byte. The last byte of the palette should be 0-filled (e.g.
`1011` are the colors in `10110000`).

#### Pixels

After the header comes the actual declaration of the pixels. The number of bits used in each
pixel will depend on the "mode" you chose. Only 1 bit for each pixel in 1-Bit mode, 2 bits in
2-Bit mode, etc. Each value for each pixel maps to a color in the palette.

##### Bit Order

When multiple pixels are stored in a single byte, the bits are used from *left to right*.
For example, if a byte in 2-bit mode contains the pixels values `1 2 3` and the remaining 2 bits are
unused (this can happen when the number of pixels is not a multiple of 8), then the byte would be
`01101100`. This style is used so that unused bits can simply be left as `0`, and so that the bits
can be read from left to right.

### Modes

#### 1-Bit Mode

This mode has only two colors: off and on.

##### Default

The default palette is for "off" to be black and "on" to be white, like you
might expect from a Pong console.

![1-Bit Palette](./assets/1-bit-palette.png "1-Bit Palette")

#### 2-Bit Mode

This mode has four colors: off, light, strong, and full.

##### Default

![2-Bit Palette](./assets/2-bit-palette.png "2-Bit Palette")

#### 8-Bit Mode

256 colors.

##### Default

For the nth color, where n is in \[0, 256\), the RGBA values are `n & 3`, `n >> 2 & 3`,
`n >> 4 & 3`, and `n >> 6 & 3`. In other words, if `n` is represented as a byte, then `r` is the
smallest 2 bits, `g` is the next 2 bits, etc.
The first 64 colors are completely transparent, but technically have different RGB values.

![8-Bit Pallete](./assets/8-bit-palette.png "8-Bit Palette")

## Implementations

- [Go](https://github.com/imretro/go)
- [TypeScript (WIP)](https://github.com/imretro/ts)

[palette]: #palette-optional
