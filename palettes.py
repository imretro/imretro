#!/usr/bin/env python3
"""
Prints the default 8-bit palette
"""
from PIL import Image, ImageDraw

def draw_color(x, y, draw, color):
    """
    Draws the color at the given coordinates.
    """
    _x, _y = x * 16, y * 16
    rgb = (color['r'], color['g'], color['b'])
    fill = (*rgb, color['a']) if 'a' in color else rgb
    draw.rectangle([(_x, _y), (_x + 16, _y + 16)], fill=fill)

colors1bit = [{c: 0 for c in 'rgb'}, {c: 0xFF for c in 'rgb'}]
colors2bit = [{c: v for c in 'rgb'} for v in (0, 0x55, 0xAA, 0xFF)]
colors8bit = [
    {'r': n & 0b11, 'g': (n >> 2) & 0b11, 'b': (n >> 4) & 0b11, 'a': n >> 6}
    for n in range(256)
]

for color in colors8bit:
    for channel in color:
        for n in range(3):
            color[channel] |= color[channel] << 2

palette1bit = Image.new('RGB', (32, 16))
draw1bit = ImageDraw.Draw(palette1bit)
palette2bit = Image.new('RGB', (32, 32))
draw2bit = ImageDraw.Draw(palette2bit)
palette8bit = Image.new('RGBA', (256, 256))
draw8bit = ImageDraw.Draw(palette8bit)

for index, color in enumerate(colors1bit):
    x = index
    y = 0
    draw_color(x, y, draw1bit, color)

for index, color in enumerate(colors2bit):
    x = index % 2
    y = index // 2
    draw_color(x, y, draw2bit, color)

for index, color in enumerate(colors8bit):
    x = (index % 16)
    y = (index // 16)
    draw_color(x, y, draw8bit, color)

for img, bits in [(palette1bit, 1), (palette2bit, 2), (palette8bit, 8)]:
    img.save(f'assets/{bits}-bit-palette.png')
