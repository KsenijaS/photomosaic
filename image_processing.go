package photomosaic

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"math"
	"os"
	"strings"
)

type ImageSignature struct {
	r, g, b float64
}

type IndexedImages struct {
	images     []image.Image
	signatures []ImageSignature
}

func NewImageSignature(img image.Image) ImageSignature {
	var sumr, sumg, sumb uint64
	var c color.Color

	bounds := img.Bounds()
	numpix := (bounds.Max.Y - bounds.Min.Y) * (bounds.Max.X - bounds.Min.X)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c = img.At(x, y)
			r, g, b, _ := c.RGBA()
			sumr += uint64(r)
			sumg += uint64(g)
			sumb += uint64(b)
		}
	}

	maxpix := 65535
	perr := float64(sumr) / float64(numpix)
	perg := float64(sumg) / float64(numpix)
	perb := float64(sumb) / float64(numpix)

	return ImageSignature{float64(perr) / float64(maxpix), float64(perg) / float64(maxpix), float64(perb) / float64(maxpix)}
}

func (p *ImageSignature) distance(other *ImageSignature) float64 {
	return math.Abs(p.r-other.r) + math.Abs(p.g-other.g) + math.Abs(p.b-other.b)
}

func (p *ImageSignature) debugString() string {
	return fmt.Sprintf("%f:%f:%f", p.r, p.g, p.b)
}

func init() {
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
}

func NewIndexedImages(dirpath string) (IndexedImages, error) {
	var indexImages IndexedImages
	var errstrings []string

	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		return IndexedImages{nil, nil}, err
	}

	count := len(files)

	indexImages.images = make([]image.Image, count)
	indexImages.signatures = make([]ImageSignature, count)
	errs := make([]error, count)

	semaphore := make(chan struct{}, 50)
	collector := make(chan struct {
		img image.Image
		sgn ImageSignature
		err error
	})
	for i := range files {
		go indexImage(dirpath+"/"+files[i].Name(), collector, semaphore)
	}
	for i := 0; i < count; i++ {
		tuple := <-collector
		indexImages.images[i] = tuple.img
		indexImages.signatures[i] = tuple.sgn
		errs[i] = tuple.err
	}

	for i := range errs {
		if errs[i] != nil {
			errstrings = append(errstrings, errs[i].Error())
		}
	}

	if errstrings != nil {
		reterr := fmt.Errorf(strings.Join(errstrings, "\n"))
		return IndexedImages{nil, nil}, reterr
	}

	return indexImages, nil
}

func readImage(filename string, semaphore chan struct{}) (image.Image, error) {
	semaphore <- struct{}{}
	defer func() { <-semaphore }()

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func indexImage(filename string, out chan<- struct {
	img image.Image
	sgn ImageSignature
	err error
}, semaphore chan struct{}) {
	img, err := readImage(filename, semaphore)
	if err != nil {
		out <- struct {
			img image.Image
			sgn ImageSignature
			err error
		}{nil, ImageSignature{0.0, 0.0, 0.0}, err}
		return
	}

	signature := NewImageSignature(img)

	out <- struct {
		img image.Image
		sgn ImageSignature
		err error
	}{img, signature, nil}
}

func (p *IndexedImages) FindClosest(img image.Image) image.Image {
	var mini int
	var d float64

	minIndex := math.MaxFloat64
	imageSignature := NewImageSignature(img)

	for i := range p.images {
		d = p.signatures[i].distance(&imageSignature)
		if minIndex > d {
			minIndex = d
			mini = i
		}
	}

	return p.images[mini]
}
