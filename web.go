package catrace

import (
	"fmt"
	"math/rand"
	"net/http"

	"appengine"
	"appengine/datastore"

	"text/template"
)

var templates = make(map[string]*template.Template)

func init() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/catmin", adminHandler)

	for _, tmpl := range []string{"catmin"} {
		templates[tmpl] = template.Must(template.ParseFiles("templates/" + tmpl + ".html"))
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, cat race!")
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	offset1 := getRandomOffset(c)
	offset2 := 0

	for {
		offset2 = getRandomOffset(c)
		if offset1 != offset2 {
			break
		}
	}

	img1 := getImage(c, offset1)
	img2 := getImage(c, offset2)

	obs := make(map[string]interface{})
	obs["Image1"] = img1
	obs["Image2"] = img2

	renderTemplate(w, obs, "catmin")
}

// func voteHandler(w http.ResponseWriter, r *http.Request) {
//  c := appengine.NewContext(r)
//  url := r.FormValue("image_url")
//
//  img := getImage(c, url)
//  img.Votes += 1
//  img.Save(c)
//}

func renderTemplate(w http.ResponseWriter, obs map[string]interface{}, tmpl string) {
	if err := templates[tmpl].Execute(w, obs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getImage(c appengine.Context, offset int) (img Image) {
	q := datastore.NewQuery("Image").Order("CreatedAt").Offset(offset).Limit(1)
	t := q.Run(c)

	_, err := t.Next(&img)

	if err != nil {
		c.Errorf("Error retrieving cat from datastore: %v", err.Error())
	}

	return
}

func getRandomOffset(c appengine.Context) int {
	count, err := datastore.NewQuery("Image").Count(c)
	if err != nil {
		c.Errorf("Counting rows in datastore failed: %v", err.Error())
	}

	if count == 0 {
		c.Errorf("No images have been downloaded")
		return 0
	}
	return rand.Int() % count
}
