#!/usr/bin/env python3
"""
Prints the default 8-bit palette
"""
from PIL import Image, ImageDraw

colors = [
    {'r': n & 0b11, 'g': (n >> 2) & 0b11, 'b': (n >> 4) & 0b11, 'a': n >> 6}
    for n in range(256)
]
for color in colors:
    for channel in color:
        for n in range(3):
            color[channel] |= color[channel] << 2

img = Image.new('RGBA', (256, 256))
draw = ImageDraw.Draw(img)

for index, color in enumerate(colors):
    x = (index % 16) * 16
    y = (index // 16) * 16
    rgba = (color['r'], color['g'], color['b'], color['a'])
    draw.rectangle([(x, y), (x + 16, y + 16)], fill=rgba)

img.save('assets/8-bit-palette.png')
