package photomosaic

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
)

const (
	subimageWidthPx  = 30
	subimageHeightPx = 30
)

type SubbableImage interface {
	image.Image
	SubImage(r image.Rectangle) image.Image
	Set(x, y int, c color.Color)
}

type ImageMatrix struct {
	img SubbableImage
}

func NewImageMatrix(i SubbableImage) (*ImageMatrix, error) {
	b := i.Bounds()
	w := b.Max.X - b.Min.X
	h := b.Max.Y - b.Min.Y

	if w%subimageWidthPx != 0 {
		err := errors.New(fmt.Sprintf("width must be divisible by %d", subimageWidthPx))
		return nil, err
	}

	if h%subimageHeightPx != 0 {
		err := errors.New(fmt.Sprintf("height must be divisible by %d", subimageHeightPx))
		return nil, err
	}

	return &ImageMatrix{img: i}, nil
}

func (p *ImageMatrix) genSubImageRectangles() []image.Rectangle {
	var out []image.Rectangle
	b := p.img.Bounds()

	for y := b.Min.Y; y < b.Max.Y; y += subimageHeightPx {
		for x := b.Min.X; x < b.Max.X; x += subimageWidthPx {
			out = append(out, image.Rect(x, y, x+subimageWidthPx, y+subimageHeightPx))
		}
	}
	return out
}

func (p *ImageMatrix) SubImages() []image.Image {
	listRec := p.genSubImageRectangles()
	out := make([]image.Image, len(listRec))

	for i := range listRec {
		out[i] = p.img.SubImage(listRec[i])
	}
	return out
}

func (p *ImageMatrix) generateImageMatrix(in []image.Image) {
	listRec := p.genSubImageRectangles()

	for i := range listRec {
		draw.Draw(p.Img, listRec[i], in[i], in[i].Bounds().Min, draw.Src)
	}
}
