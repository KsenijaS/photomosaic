package main

import (
	"fmt"
	"github.com/KsenijaS/photomosaic"
	"github.com/disintegration/imaging"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type user struct {
	UserName string
	Password []byte
	First    string
	Last     string
}

var tpl *template.Template
var dbUsers = map[string]user{}      // user ID, user
var dbSessions = map[string]string{} // session ID, user ID

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

	http.HandleFunc("/", index)
	http.HandleFunc("/mosaic", mosaic(indexedImages))
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/guest", guest)
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("./public"))))
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, req *http.Request) {
	u := getUser(w, req)
	tpl.ExecuteTemplate(w, "index.gohtml", u)
}

func mosaic(i photomosaic.IndexedImages) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var filename string
		var newimg string

		if !alreadyLoggedIn(r) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

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
		err := tpl.ExecuteTemplate(w, "mosaic.gohtml", newimg)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func signup(w http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	// process form submission
	if req.Method == http.MethodPost {

		// get form values
		un := req.FormValue("username")
		p := req.FormValue("password")
		f := req.FormValue("firstname")
		l := req.FormValue("lastname")

		// username taken?
		if _, ok := dbUsers[un]; ok {
			http.Error(w, "Username already taken", http.StatusForbidden)
			return
		}

		bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// store user in dbUsers
		u := user{un, bs, f, l}
		dbUsers[un] = u

		// redirect
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(w, "signup.gohtml", nil)
}

func login(w http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {
		un := req.FormValue("username")
		p := req.FormValue("password")

		u, ok := dbUsers[un]
		if !ok {
			http.Error(w, "Username and/or password is incorrect", http.StatusForbidden)
		}

		err := bcrypt.CompareHashAndPassword(u.Password, []byte(p))
		if err != nil {
			http.Error(w, "Username and/or password is incorrect", http.StatusForbidden)
		}

		sID, _ := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		http.SetCookie(w, c)
		dbSessions[c.Value] = un

		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "login.gohtml", nil)
}

func logout(w http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	c, _ := req.Cookie("session")
	delete(dbSessions, c.Value)

	c.Value = ""
	c.MaxAge = -1
	http.SetCookie(w, c)

	http.Redirect(w, req, "/login", http.StatusSeeOther)
}

func guest(w http.ResponseWriter, req *http.Request) {
	g := "guest"
	if alreadyLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	dbUsers[g] = user{g, []byte(g), g, g}
	sID, _ := uuid.NewV4()
	c := &http.Cookie{
		Name:  "session",
		Value: sID.String(),
	}
	http.SetCookie(w, c)
	dbSessions[c.Value] = g

	http.Redirect(w, req, "/", http.StatusSeeOther)
}
