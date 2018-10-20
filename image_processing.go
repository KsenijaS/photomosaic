package photomosaic

import (
    "os"
    "io/ioutil"
    "image"
    "image/color"
    "image/jpeg"
    "math"
    "fmt"
    "strings"
)


type ImageSignature struct {
    r, g, b float64
}

type IndexedImages struct {
    images []image.Image
    signatures []ImageSignature
}

func NewImageSignature(img image.Image) ImageSignature {
    var sumr, sumg, sumb uint64
    var c color.Color

    bounds := img.Bounds()

    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            c = img.At(x, y)
            r,g,b,_ := c.RGBA()
            sumr += uint64(r)
            sumg += uint64(g)
            sumb += uint64(b)
        }
    }

    s := sumr+sumg+sumb+1

    return ImageSignature{float64(sumr)/float64(s), float64(sumg)/float64(s), float64(sumb)/float64(s)}
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

    collector := make(chan struct{img image.Image; sgn ImageSignature; err error})
    for i := range(files) {
        go indexImage(dirpath+"/"+files[i].Name(), collector)
    }
    for i := 0; i < count; i++ {
        tuple := <-collector
        indexImages.images[i] = tuple.img
        indexImages.signatures[i] = tuple.sgn
        errs[i] = tuple.err
    }

    for i := range(errs) {
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

func indexImage(filename string, out chan<- struct{img image.Image; sgn ImageSignature; err error} ) {

    file, err := os.Open(filename)
    if err != nil {
        out <- struct{img image.Image; sgn ImageSignature; err error}{nil, ImageSignature{0.0, 0.0, 0.0}, err}
        return
    }
    defer file.Close()

    img, _, err := image.Decode(file)
    if err != nil {
        out <- struct{img image.Image; sgn ImageSignature; err error}{nil, ImageSignature{0.0, 0.0, 0.0}, err}
        return
    }

    signature := NewImageSignature(img)

    out <- struct{img image.Image; sgn ImageSignature; err error}{img, signature, nil}
}

func (p *IndexedImages) FindClosest(img image.Image) image.Image {
    var mini int
    var d float64

    minIndex := math.MaxFloat64
    imageSignature := NewImageSignature(img)

    for i := range(p.images) {
        d = p.signatures[i].distance(&imageSignature)
        if minIndex > d {
            minIndex = d
            mini = i
        }
    }

    return p.images[mini]
}
