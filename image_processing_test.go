package photomosaic

import (
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"math"
	"os"
	"testing"
)

func almostEqual(x float64, y float64) bool {
	return math.Abs(x-y) < 0.0001
}

func TestNewImageSignatureReturnsCorrectNumbers(t *testing.T) {
	input := image.NewNRGBA(image.Rect(0, 0, 600, 150))

	blue := color.NRGBA{0, 0, 255, 255}
	green := color.NRGBA{0, 255, 0, 255}
	red := color.NRGBA{255, 0, 0, 255}

	draw.Draw(input, image.Rect(0, 0, 200, 150), &image.Uniform{blue}, image.ZP, draw.Src)
	draw.Draw(input, image.Rect(200, 0, 400, 150), &image.Uniform{green}, image.ZP, draw.Src)
	draw.Draw(input, image.Rect(400, 0, 600, 150), &image.Uniform{red}, image.ZP, draw.Src)

	sign := NewImageSignature(input)

	if !almostEqual(sign.r, 1.0/3) || !almostEqual(sign.g, 1.0/3) || !almostEqual(sign.b, 1.0/3) {
		t.Errorf("The concentration of red, blue or green color is not correct r %f, g %f, b %f", sign.r, sign.g, sign.b)
	}
}

func TestNewImageSignatureReturnsCorrectNumbersOnlyBlueAndRed(t *testing.T) {
	input := image.NewNRGBA(image.Rect(0, 0, 600, 150))

	blue := color.NRGBA{0, 0, 255, 255}
	red := color.NRGBA{255, 0, 0, 255}
	black := color.NRGBA{0, 0, 0, 255}

	draw.Draw(input, image.Rect(0, 0, 200, 150), &image.Uniform{blue}, image.ZP, draw.Src)
	draw.Draw(input, image.Rect(200, 0, 400, 150), &image.Uniform{black}, image.ZP, draw.Src)
	draw.Draw(input, image.Rect(400, 0, 600, 150), &image.Uniform{red}, image.ZP, draw.Src)

	sign := NewImageSignature(input)

	if !almostEqual(sign.r, 1.0/2) || !almostEqual(sign.g, 0.0) || !almostEqual(sign.b, 1.0/2) {
		t.Errorf("The concentration of red, blue or green color is not correct r %f, g %f, b %f", sign.r, sign.g, sign.b)
	}
}

func TestImageColorDistanceReturnsCorrectDistance(t *testing.T) {
	input := image.NewNRGBA(image.Rect(0, 0, 600, 150))

	blue := color.NRGBA{0, 0, 255, 255}
	red := color.NRGBA{255, 0, 0, 255}
	green := color.NRGBA{0, 255, 0, 255}

	draw.Draw(input, image.Rect(0, 0, 200, 150), &image.Uniform{blue}, image.ZP, draw.Src)
	draw.Draw(input, image.Rect(200, 0, 400, 150), &image.Uniform{green}, image.ZP, draw.Src)
	draw.Draw(input, image.Rect(400, 0, 600, 150), &image.Uniform{red}, image.ZP, draw.Src)

	sign := NewImageSignature(input)
	d := sign.distance(&sign)

	if !almostEqual(d, 0.0) {
		t.Errorf("Color distance between images is not correct %f", d)
	}
}

func TestNewIndexedImagesReturnCorrectValue(t *testing.T) {
	index, err := NewIndexedImages("dir_test")
	if err != nil {
		t.Fatalf("Cannot create IndexedImages: %s", err)
	}

	if len(index.images) != 3 || len(index.signatures) != 3 {
		t.Errorf("Length is not correct")
	}

	files, err := ioutil.ReadDir("dir_test")
	if err != nil {
		t.Fatalf("Cannot open directory")
	}

	isign := make([]ImageSignature, 3)
	present := make([]int, 3)

	for i := range files {
		file, err := os.Open("dir_test/" + files[i].Name())
		if err != nil {
			t.Fatalf("Cannot open the file: %s", err)
		}
		defer file.Close()

		idx, _, err := image.Decode(file)
		if err != nil {
			t.Fatalf("Cannot decode image: %s", err)
		}

		isign[i] = NewImageSignature(idx)
	}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if index.signatures[j] == isign[i] {
				if present[i] == 1 {
					t.Errorf("Duplicate image %s", isign[i].debugString())
				}
				present[i] = 1
			}

		}
	}

	for i := range present {
		if present[i] == 0 {
			t.Errorf("Image not found: %s", isign[i].debugString())
		}
	}

}

func TestFindClosestReturnsCorrectImage(t *testing.T) {
	index, err := NewIndexedImages("dir_test")
	if err != nil {
		t.Fatalf("Cannot create IndexedImages: %s", err)
	}

	files, err := ioutil.ReadDir("dir_test")
	if err != nil {
		t.Fatalf("Cannot open directory")
	}

	file, err := os.Open("dir_test/" + files[0].Name())
	if err != nil {
		t.Fatalf("Cannot open the file: %s", err)
	}
	defer file.Close()

	idx, _, err := image.Decode(file)
	if err != nil {
		t.Fatalf("Cannot decode image: %s", err)
	}

	img := index.FindClosest(idx)

	sign1 := NewImageSignature(idx)
	sign2 := NewImageSignature(img)

	d := sign1.distance(&sign2)

	if !almostEqual(d, 0.0) {
		t.Errorf("Correct image is not found, image signature is %s instead of %s", sign2.debugString(), sign1.debugString())
	}
}
