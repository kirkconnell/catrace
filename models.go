package models

import (
    "time"

    "appengine"
    "appengine/datastore"
)

type Image struct {
    OriginalUrl  string
    CloudUrl     string
    Category     string
    Views        int
    Votes        int
    Rank         int
    CreatedAt    time.Time
    VotedAt      time.Time
}

func (i *Image) Save(c appengine.Context) bool {
    // TODO: Figure out what todo with the Key
    _, err := datastore.Put(c, datastore.NewIncompleteKey(c, "image", nil), &i)
    if err != nil {
        c.Errorf("Saving Image to Datastore failed with error: %v", err.Error())
        return false
    }
    return true
} 

