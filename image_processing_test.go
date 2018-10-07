package photomosaic

import (
    "image"
    "image/color"
    "image/draw"
    "testing"
)

func TestColorFractionsReturnsCorrectNumbers(t *testing.T) {
    input := image.NewNRGBA(image.Rect(0,0, 600, 150))

    blue := color.NRGBA{0, 0, 255, 255}
    green := color.NRGBA{0, 255, 0, 255}
    red := color.NRGBA{255, 0, 0, 255}

    draw.Draw(input, image.Rect(0, 0, 200, 150), &image.Uniform{blue}, image.ZP, draw.Src)
    draw.Draw(input, image.Rect(200, 0, 400, 150), &image.Uniform{green}, image.ZP, draw.Src)
    draw.Draw(input, image.Rect(400, 0, 600, 150), &image.Uniform{red}, image.ZP, draw.Src)

    r, g, b := ColorFractions(input)

    if r != 1.0/3 || g != 1.0/3 || b != 1.0/3 {
        t.Errorf("The concentration of red, blue or green color is not correct r %f, g %f, b %f", r,g,b)
    }
}


func TestColorFractionsReturnsCorrectNumbersOnlyBlueAndRed(t *testing.T) {
    input := image.NewNRGBA(image.Rect(0,0, 600, 150))

    blue := color.NRGBA{0, 0, 255, 255}
    red := color.NRGBA{255, 0, 0, 255}
    black := color.NRGBA{0, 0, 0, 255}

    draw.Draw(input, image.Rect(0, 0, 200, 150), &image.Uniform{blue}, image.ZP, draw.Src)
    draw.Draw(input, image.Rect(200, 0, 400, 150), &image.Uniform{black}, image.ZP, draw.Src)
    draw.Draw(input, image.Rect(400, 0, 600, 150), &image.Uniform{red}, image.ZP, draw.Src)

    r, g, b := ColorFractions(input)

    if r != 1.0/2 || g != 0.0 || b != 1.0/2 {
        t.Errorf("The concentration of red, blue or green color is not correct r %f, g %f, b %f", r,g,b)
    }
}

