package photomosaic

import (
	"image"
	"image/color"
	"image/draw"
	"strings"
	"testing"
)

func TestNewImageMatrixAcceptsImageOfSubimageSize(t *testing.T) {
	input := image.NewNRGBA(image.Rect(100, 100, 400, 400))
	im, err := NewImageMatrix(input)
	if err != nil {
		t.Errorf("Error while creating the image: %s", err)
	}
	if im == nil {
		t.Errorf("Image was nil")
	}
	if input != im.Img {
		t.Errorf("Input image was not stored in the returned struct")
	}
}

func TestNewImageMatrixRejectsIncompatibleWidthImage(t *testing.T) {
	input := image.NewNRGBA(image.Rect(100, 100, 123, 250))
	im, err := NewImageMatrix(input)
	if err == nil {
		t.Errorf("NewImageMatrix accepted an image of incompatible width")
	}
	if !strings.Contains(err.Error(), "width") {
		t.Errorf("Returned error message did not mention the incompatible width: %s", err)
	}
	if im != nil {
		t.Errorf("NewImageMatrix returned both an image and an error for an image of incompatible width")
	}
}

func TestNewImageMatrixRejectsIncompatibleHeightImage(t *testing.T) {
	input := image.NewNRGBA(image.Rect(100, 100, 400, 123))
	im, err := NewImageMatrix(input)
	if err == nil {
		t.Errorf("NewImageMatrix accepted an image of incompatible height")
	}
	if !strings.Contains(err.Error(), "height") {
		t.Errorf("Returned error message did not mention the incompatible height: %s", err)
	}
	if im != nil {
		t.Errorf("NewImageMatrix returned both an image and an error for an image of incompatible height")
	}
}

func TestgenSubImagesRectanglesReturnsCorrectNumberOfRectangles(t *testing.T) {
	input := image.NewNRGBA(image.Rect(100, 100, 300, 700))
	im, err := NewImageMatrix(input)
	if err != nil {
		t.Fatalf("Error creating ImageMatrix: %s", err)
	}
	numRec := len(im.genSubImageRectangles())
	if numRec != 4 {
		t.Errorf("genSubImageRectangles() doesn't have correct number of rectangles: %d", numRec)
	}
}

func TestgenSubImagesRectanglesReturnsCorrectCoordinates(t *testing.T) {
	input := image.NewNRGBA(image.Rect(0, 0, 400, 450))
	coord := []image.Point{{0, 0}, {200, 0}, {0, 150}, {200, 150}, {0, 300}, {200, 300}}
	im, _ := NewImageMatrix(input)
	listRec := im.genSubImageRectangles()
	for i := range coord {
		if coord[i] != listRec[i].Min {
			t.Errorf("Rectangle coordinate (%d, %d )doesn't match with (%d, %d)", listRec[i].Min.X, listRec[i].Min.Y, coord[i].X, coord[i].Y)
		}
	}
}

func TestSubImagesReturnsExpectedImages(t *testing.T) {
	input := image.NewNRGBA(image.Rect(0, 0, 600, 600))

	blue := color.NRGBA{0, 0, 255, 255}
	green := color.NRGBA{0, 255, 0, 255}
	red := color.NRGBA{255, 0, 0, 255}
	purple := color.NRGBA{200, 0, 200, 255}

	draw.Draw(input, image.Rect(0, 0, 300, 300), &image.Uniform{blue}, image.ZP, draw.Src)
	draw.Draw(input, image.Rect(300, 0, 600, 300), &image.Uniform{green}, image.ZP, draw.Src)
	draw.Draw(input, image.Rect(0, 300, 300, 600), &image.Uniform{red}, image.ZP, draw.Src)
	draw.Draw(input, image.Rect(300, 300, 600, 600), &image.Uniform{purple}, image.ZP, draw.Src)

	im, err := NewImageMatrix(input)
	if err != nil {
		t.Fatalf("Error creating ImageMatrix: %s", err)
	}

	imgRec := im.SubImages()

	for i := 0; i < 4; i++ {
		b := imgRec[i].Bounds()
		for y := b.Min.Y; y < b.Max.Y; y++ {
			for x := b.Min.X; x < b.Max.X; x++ {
				if imgRec[i].At(x, y) != input.At(x, y) {
					t.Errorf("Images are not the same at (%d, %d) point", x, y)
				}
			}
		}

	}
}

func TestgenerateImageMatrix(t *testing.T) {
	input := image.NewNRGBA(image.Rect(0, 0, 600, 600))

	blue := color.NRGBA{0, 0, 255, 255}
	green := color.NRGBA{0, 255, 0, 255}
	red := color.NRGBA{255, 0, 0, 255}
	purple := color.NRGBA{200, 0, 200, 255}

	draw.Draw(input, image.Rect(0, 0, 300, 300), &image.Uniform{blue}, image.ZP, draw.Src)
	draw.Draw(input, image.Rect(300, 0, 600, 300), &image.Uniform{green}, image.ZP, draw.Src)
	draw.Draw(input, image.Rect(0, 300, 300, 600), &image.Uniform{red}, image.ZP, draw.Src)
	draw.Draw(input, image.Rect(300, 300, 600, 600), &image.Uniform{purple}, image.ZP, draw.Src)

	colors := [4]color.NRGBA{blue, green, red, purple}
	colorimgs := make([]image.Image, 4)
	for i := range colors {
		colorimg := image.NewNRGBA(image.Rect(0, 0, 600, 600))
		draw.Draw(colorimg, image.Rect(0, 0, 300, 300), &image.Uniform{colors[i]}, image.ZP, draw.Src)
		colorimgs[i] = colorimg
	}

	newimg := image.NewNRGBA(image.Rect(0, 0, 600, 600))
	im, err := NewImageMatrix(newimg)
	if err != nil {
		t.Fatalf("Error creating ImageMatrix: %s", err)
	}

	im.generateImageMatrix(colorimgs)

	b := input.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if input.At(x, y) != im.img.At(x, y) {
				t.Errorf("Images are not the same at (%d, %d) point", x, y)
			}
		}
	}

}
