package photomosaic

import (
	"github.com/disintegration/imaging"
	"log"
)

func GenerateImageMatrixFromImages(img SubbableImage, i IndexedImages) (*ImageMatrix, error) {
	imgmtx, err := NewImageMatrix(img)
	if err != nil {
		log.Fatalf("Cannot create Image Matrix %s", err)
		return nil, err
	}

	outimgs := imgmtx.SubImages()
	for k := range outimgs {
		outimgs[k] = imaging.Resize(i.FindClosest(outimgs[k]), 30, 30, imaging.Lanczos)
	}

	imgmtx.generateImageMatrix(outimgs)

	return imgmtx, nil
}
