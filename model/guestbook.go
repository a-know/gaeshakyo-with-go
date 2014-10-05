package model

import (
	"appengine"
	"appengine/datastore"
	"time"
)

type Guestbook struct {
	Message   string
	CreatedAt time.Time
}

func Save(c appengine.Context, message string) (err error) {
	g1 := Guestbook{
		Message:   message,
		CreatedAt: time.Now(),
	}

	// put
	key := datastore.NewIncompleteKey(c, "guestbook", nil)
	key, err = datastore.Put(c, key, &g1)
	return err
}

func DescList(c appengine.Context) (guestbooks []Guestbook) {
	q := datastore.NewQuery("guestbook")

	_, _ = q.GetAll(c, &guestbooks)

	return guestbooks
}
