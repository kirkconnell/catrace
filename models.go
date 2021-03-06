package catrace

import (
	"time"

	"appengine"
	"appengine/datastore"
)

type Image struct {
	OriginalUrl string
	CloudUrl    string
	Category    string
	Views       int
	Votes       int
	Rank        int
	CreatedAt   time.Time
	VotedAt     time.Time
}

func (i *Image) Save(c appengine.Context) bool {
	// TODO: Figure out what todo with the Key
	_, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Image", nil), i)
	if err != nil {
		c.Errorf("Saving Image to Datastore failed with error: %v", err.Error())
		return false
	}
	return true
}

func ImageByOriginalUrl(c appengine.Context, url string) (img Image, err error) {
	q := datastore.NewQuery("Image").Filter("OriginalUrl =", url).KeysOnly()
	t := q.Limit(1).Run(c)
	_, err = t.Next(&img)
	return img, err
}
