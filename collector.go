package catrace

import (
	"strings"
    "net/http"
    "bytes"
    "io/ioutil"
    "time"

    "appengine"
    "appengine/datastore"
    "appengine/taskqueue"
)

func init() {
	http.HandleFunc("/datacollector", collectUrls)
    http.HandleFunc("/worker", worker)
}

func collectUrls(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    url := "http://catoverflow.com/api/query?offset=0&limit=1000"
    
    resp, err := http.Get(url)
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
    n := bytes.Index(body, []byte{0})
    s := string(body[:n])
    
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
    if t == nil {
    	img := new(Image)
    	img.OriginalUrl = url
    	img.Category = "cats"
    	img.CreatedAt = time.Now()
    	img.Save(c)
    }
}

