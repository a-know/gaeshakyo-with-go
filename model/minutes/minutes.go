package minutes

import (
	"appengine"
	"appengine/datastore"
	"time"

	"code.google.com/p/go-uuid/uuid"
)

type Minutes struct {
	Key *datastore.Key
	Title   string
	CreatedAt time.Time
}

func SaveAs(c appengine.Context, title string) (*datastore.Key, error) {
	key := datastore.NewKey(c, "minutes", uuid.New(), 0, nil)

	m1 := Minutes{
		Key:       key,
		Title:     title,
		CreatedAt: time.Now(),
	}

	// put
	_, err := datastore.Put(c, key, &m1)
	return key, err
}

func DescList(c appengine.Context) (minutes []Minutes, err error) {
	q := datastore.NewQuery("minutes").Order("-CreatedAt")

	_, err = q.GetAll(c, &minutes)

	return minutes, err
}

