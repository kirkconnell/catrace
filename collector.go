package catrace

import (
	"fmt"
	"strings"
    "net/http"
    "io/ioutil"
    "time"

    "appengine"
    "appengine/datastore"
    "appengine/taskqueue"
    "appengine/urlfetch"
)

func init() {
	http.HandleFunc("/datacollector", collectUrls)
    http.HandleFunc("/worker", worker)
}

func collectUrls(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    url := "http://catoverflow.com/api/query?offset=0&limit=1000"
    
    client := urlfetch.Client(c)    
    resp, err := client.Get(url)
	if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
    	return
    }
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
    	return
    }
    
    //Split the body by new lines to get the url for each image.
    s := string(body)
    
    urls := strings.Fields(s)
    for _, u := range urls {
    	t := taskqueue.NewPOSTTask("/worker", map[string][]string{"url": {u}})
        if _, err := taskqueue.Add(c, t, ""); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    }
}

func worker(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    url := r.FormValue("url")
    
    q := datastore.NewQuery("Image").Filter("OriginalUrl =", url).KeysOnly()
    t := q.Limit(1).Run(c)
    var img Image
    _, err := t.Next(&img)
    fmt.Printf("Image is: %v", img)
    
	if err == datastore.Done {
		img := new(Image)
		img.OriginalUrl = url
		img.Category = "cats"
		img.CreatedAt = time.Now()
		img.Views = 0
		img.Votes = 0
		img.Rank = 0
		img.Save(c)
	} else if err != nil {
        c.Errorf("Fetching Images from data store failed: %v", err.Error())
        return
    }
}

