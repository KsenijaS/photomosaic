package photomosaic

import (
	"errors"
	"fmt"
	"image"
)

const (
	subimageWidthPx  = 200
	subimageHeightPx = 150
)

type ImageMatrix struct {
	img *image.NRGBA
}

func NewImageMatrix(i *image.NRGBA) (*ImageMatrix, error) {
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
