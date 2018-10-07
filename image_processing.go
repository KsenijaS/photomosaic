package photomosaic

import (
    "image"
    "image/color"
)

func ColorFractions(img image.Image) (rPix float32, gPix float32, bPix float32)  {
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

    s := sumr+sumg+sumb

    return float32(sumr)/float32(s), float32(sumg)/float32(s), float32(sumb)/float32(s)
}
