# imretro

An image format for retro-style images

I made this format so that you can make your own image with a hex editor for a retro-ish
artstyle. In the future I may need to create an editor to allow making larger images.

## File

### Endianness

The file is little-endian, meaning that `0A 00` is 10 (`0x000A`), not 2560 (`0x0A00`).
This is used so that any unused bit can simple be a trailing `0` instead of having to
pad `0`s in front of the last bit in the final byte.

### Contents

#### Header

Each file should start with a header to provide some information about the image.
The first 7 *bytes* are `"IMRETRO"` (for a total of 56 bits). This is the file signature.
The next two bits map to the "modes" that declare the bits-per-pixel. Following that are
32 bits for the dimensions: 16 for width and 16 for height. The reason for this limited
range of dimensions is to be faithful to the retro-ish goal of this format.
Next is the palette. First, a single bit: `0` for no palette, `1` to declare that a palette
will be present in the file. When no palette is present in the file, this means that the
file decoder should choose its own default palette.

The palette will declare the possible colors in the image. The number of colors in your
palette depend on the number of bits you chose to use in your header. In 1-Bit mode, you
will declare 2 colors, in 2-Bit mode, 4 colors, etc. Each color in the palette will be 4 bytes:
RGB and an alpha value. So, in 8-Bit mode, with 256 possible colors, the palette will be 1024
bytes.

#### Pixels

After the header comes the actual declaration of the pixels. The number of bits used in each
pixel will depend on the "mode" you chose. Only 1 bit for each pixel in 1-Bit mode, 2 bits in
2-Bit mode, etc. Each value for each pixel maps to a color in the palette.

### Modes

#### 1-Bit Mode

This mode has only two colors: "off" and "on".
The default palette is for "off" to be black and "on" to be white, like you
might expect from a Pong console.

#### 2-Bit Mode

Black, white, and two shades of gray.

#### 8-Bit Mode

256 colors. Inspired by the NES, the default palette only has
54 colors. The remaining unused colors all mapping to alpha.
The palette can be viewed [here][NES palette].

[NES Palette]: https://en.wikipedia.org/wiki/List_of_video_game_console_palettes#NES
