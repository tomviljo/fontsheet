package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"log"
	"os"
	"strconv"

	"github.com/alecthomas/kong"
)

//go:embed images/font.gif
var font []byte

var fontGif *gif.GIF

const (
	fontW        = 6
	fontH        = 9
	fontBaseline = 1
)

func init() {
	var err error
	fontGif, err = gif.DecodeAll(bytes.NewReader(font))
	if err != nil {
		panic(err)
	}
}

func parseRGB(text string) (color.NRGBA, error) {
	if len(text) != 6 {
		return color.NRGBA{}, fmt.Errorf("Not an RRGGBB value: %q", text)
	}
	value24, err := strconv.ParseUint(text, 16, 24)
	if err != nil {
		return color.NRGBA{}, err
	}
	return color.NRGBA{
		R: uint8(value24 >> 16),
		G: uint8(value24 >> 8),
		B: uint8(value24),
		A: uint8(0xff),
	}, nil
}

func drawHorizontal(img *image.Paletted, x0, y0, dx int, index uint8) {
	for i := range dx {
		img.SetColorIndex(x0+i, y0, index)
	}
}

func drawVertical(img *image.Paletted, x0, y0, dy int, index uint8) {
	for i := range dy {
		img.SetColorIndex(x0, y0+i, index)
	}
}

func drawChar(img *image.Paletted, x0, y0, char int, index uint8) {
	if char < 32 || char > 127 {
		return
	}
	fontX := fontW * ((char - 32) % 16)
	fontY := fontH * ((char - 32) / 16)
	for i := range fontH {
		for j := range fontW {
			if fontGif.Image[0].ColorIndexAt(fontX+j, fontY+i) == 1 {
				//img.Pix[img.PixOffset(x0+j, y0+i)] ^= index
				img.SetColorIndex(x0+j, y0+i, img.ColorIndexAt(x0+j, y0+i)^index)
			}
		}
	}
}

func command() error {
	var cli struct {
		First           int    `short:"f" default:"32" help:"First character"`
		Last            int    `short:"l" default:"127" help:"Last character"`
		Columns         int    `short:"c" default:"16" help:"Maximum number of columns, set to zero for single row"`
		Width           int    `short:"W" default:"10" help:"Width of each character cell"`
		Height          int    `short:"H" default:"16" help:"Height of each character cell"`
		Baseline        int    `short:"b" default:"4" help:"Distance of baseline from bottom of character cell"`
		Strip           bool   `short:"s" help:"Strip top and left border"`
		NoText          bool   `short:"n" help:"No placeholder characters"`
		BackgroundColor string `short:"B" default:"ffffff" help:"Background color"`
		GridColor       string `short:"G" default:"cccccc" help:"Color of grid and template characters"`
		FontColor       string `short:"F" default:"000000" help:"Color reserved for font"`
		OutputFile      string `arg:"" help:"GIF file to output"`
	}
	_ = kong.Parse(&cli)

	chars := (cli.Last - cli.First) + 1
	var columns, rows int
	if cli.Columns == 0 {
		columns = chars
		rows = 1
	} else {
		columns = cli.Columns
		rows = (chars + cli.Columns - 1) / cli.Columns
	}
	totalWidth := columns * cli.Width
	totalHeight := rows * cli.Height
	offsetX := 0
	offsetY := 0
	if !cli.Strip {
		totalWidth++
		totalHeight++
		offsetX++
		offsetY++
	}

	backgroundColor, err := parseRGB(cli.BackgroundColor)
	if err != nil {
		return err
	}
	gridColor, err := parseRGB(cli.GridColor)
	if err != nil {
		return err
	}
	fontColor, err := parseRGB(cli.FontColor)
	if err != nil {
		return err
	}
	palette := color.Palette{backgroundColor, fontColor, gridColor}

	img := image.NewPaletted(image.Rect(0, 0, totalWidth, totalHeight), palette)

	gridIndex := uint8(2)
	if !cli.Strip {
		drawHorizontal(img, 0, 0, totalWidth, gridIndex)
		drawVertical(img, 0, 0, totalHeight, gridIndex)
	}
	for i := range columns {
		rightX := offsetX + (i+1)*cli.Width - 1
		drawVertical(img, rightX, 0, totalHeight, gridIndex)
	}
	for i := range rows {
		bottomY := offsetY + (i+1)*cli.Height - 1
		drawHorizontal(img, 0, bottomY, totalWidth, gridIndex)
		if cli.Baseline > 0 {
			drawHorizontal(img, 0, bottomY-cli.Baseline, totalWidth, gridIndex)
		}
	}

	if !cli.NoText && cli.Width >= fontW && (cli.Height-cli.Baseline) >= fontH {
		charOffsetX := (cli.Width - fontW) / 2
		charOffsetY := (cli.Height - cli.Baseline) - fontH
		for i := range cli.Last - cli.First + 1 {
			column := i % columns
			row := i / columns
			x := offsetX + cli.Width*column + charOffsetX
			y := offsetY + cli.Height*row + charOffsetY
			drawChar(img, x, y, cli.First+i, gridIndex)
		}
	}

	file, err := os.Create(cli.OutputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	gifFile := &gif.GIF{
		Image: []*image.Paletted{img},
		Delay: []int{0},
		Config: image.Config{
			ColorModel: palette,
			Width:      totalWidth,
			Height:     totalHeight,
		},
	}
	return gif.EncodeAll(file, gifFile)
}

func main() {
	err := command()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
