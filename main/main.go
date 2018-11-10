package main

import (
	"fmt"
	"github.com/KsenijaS/photomosaic"
	"github.com/disintegration/imaging"
	"html/template"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	log.Print("Generating index...")
	indexedImages, err := photomosaic.NewIndexedImages("sample_images")
	if err != nil {
		log.Fatalf("Error while indexing images: %v", err)
	}
	log.Print("Images indexed successfully.")

	http.HandleFunc("/", index(indexedImages))
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("./public"))))
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}

func index(i photomosaic.IndexedImages) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var filename string
		var newimg string

		if r.Method == http.MethodPost {
			f, h, err := r.FormFile("nf")
			if err != nil {
				fmt.Println(err)
			}
			defer f.Close()

			filename = h.Filename
			newimg = "new" + filename
			dst, err := os.Create(filepath.Join("./public/pics/", newimg))
			if err != nil {
				fmt.Println(err)
			}
			defer dst.Close()

			f.Seek(0, 0)
			img, err := jpeg.Decode(f)
			if err != nil {
				fmt.Println(err)
			}

			img = imaging.Resize(img, 1200, 1200, imaging.Lanczos)
			myimg := img.(photomosaic.SubbableImage)
			imgmtx, err := photomosaic.GenerateImageMatrixFromImages(myimg, i)
			if err != nil {
				log.Fatalf("Cannot create Image Matrix %s", err)
			}

			jpeg.Encode(dst, imgmtx.Img, &jpeg.Options{jpeg.DefaultQuality})

		}
		err := tpl.ExecuteTemplate(w, "index.gohtml", newimg)
		if err != nil {
			fmt.Println(err)
		}
	}
}
