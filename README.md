# fontsheet

Command-line utility for generating template sheets for bitmap fonts.

Example output file:

![Example font sheet](/images/example_x3.gif)

## Installation

* Install Go.
* Clone this repository.
* Run `go build` inside the repository.
* Copy `fontsheet` to a directory in your `PATH`, such as `/usr/local/bin`.

## Usage

```
Usage: fontsheet <output-file> [flags]

Arguments:
  <output-file>    GIF file to output

Flags:
  -h, --help                         Show context-sensitive help.
  -f, --first=32                     First character
  -l, --last=127                     Last character
  -c, --columns=16                   Maximum number of columns, set to zero for single row
  -W, --width=10                     Width of each character cell
  -H, --height=16                    Height of each character cell
  -b, --baseline=4                   Distance of baseline from bottom of character cell
  -s, --strip                        Strip top and left border
  -n, --no-text                      No placeholder characters
  -B, --background-color="ffffff"    Background color
  -G, --grid-color="cccccc"          Color of grid and template characters
  -F, --font-color="000000"          Color reserved for font
```

## Notes on output file

The output file has three colors: 0 (the background color, white by default), 1 (reserved for font, black by default) and 2 (grid and placeholder characters, light gray by default).

It consists of tiles of size `--width` times `--height` pixels, with grid lines drawn in the rightmost pixel column and the bottom pixel row. Unless `--strip` is given, grid lines are added to the left and top of the image, offsetting the tiles one pixel to the right and bottom.

If `--width` and `--height` allows, and `--no-text` is not given, placeholder characters are inserted using the grid color for characters in the ASCII range. The idea is to manually erase each placeholder and draw the real character in its place using the font color.
