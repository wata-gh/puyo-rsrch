package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"strings"
)

type canvas struct {
	Width     int
	Height    int
	Image     *image.NRGBA
	PuyoImage image.Image
}

func newCanvas(width int, height int) *canvas {
	fpuyo, err := os.Open("images/puyos.png")
	if err != nil {
		panic(err)
	}
	defer fpuyo.Close()

	canvas := new(canvas)
	canvas.PuyoImage, _, err = image.Decode(fpuyo)
	if err != nil {
		panic(err)
	}
	canvas.Width = width
	canvas.Height = height
	canvas.Image = image.NewNRGBA(image.Rectangle{image.Pt(0, 0), image.Pt(32*width, 32*height)})
	return canvas
}

func (c *canvas) placePuyo(puyo rune, x int, y int) {
	ix := 0
	iy := 0
	switch puyo {
	case 'r':
		ix = 0
	case 'g':
		ix = 32
	case 'b':
		ix = 64
	case 'y':
		ix = 96
	case 'p':
		ix = 128
	case '.':
		return
	}

	point := image.Pt((x+1)*32, (y+1)*32)
	draw.Draw(c.Image, image.Rectangle{image.Pt(x*32, y*32), point}, c.PuyoImage, image.Pt(ix, iy), draw.Over)
}

func (c *canvas) export(name string) {
	outfile, _ := os.Create(name)
	defer outfile.Close()
	png.Encode(outfile, c.Image)
}

func main() {
	dir := flag.String("dir", "", "output directory")
	height := flag.Int("height", 20, "height of canvas")
	flag.Parse()

	if *dir != "" {
		_, err := os.Stat(*dir)
		if os.IsNotExist(err) {
			err := os.Mkdir(*dir, 0755)
			if err != nil {
				panic(err)
			}
		}
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		haipuyo := strings.TrimRight(scanner.Text(), "\n")

		c := newCanvas(4, *height)
		puyoCnt := map[rune]int{
			'r': 0,
			'g': 0,
			'b': 0,
			'y': 0,
			'p': 0,
		}
		row := 0
		puyoRow := map[rune]int{}
		for i, puyo := range haipuyo {
			_, ok := puyoRow[puyo]
			if ok == false {
				puyoRow[puyo] = row
				row++
			}
			c.placePuyo(puyo, puyoRow[puyo], (*height-1)-puyoCnt[puyo])
			puyoCnt[puyo]++
			if (i+1)%2 == 0 {
				c.export(fmt.Sprintf("%s/%d.png", *dir, (i+1)/2))
			}
		}
	}
}
