package minutes

import (
	"appengine"
	"appengine/datastore"
	"time"
)

type Minutes struct {
	Title   string
	CreatedAt time.Time
}

func SaveAs(c appengine.Context, title string) (*datastore.Key, error) {
	m1 := Minutes{
		Title:   title,
		CreatedAt: time.Now(),
	}

	// put
	key := datastore.NewIncompleteKey(c, "minutes", nil)
	_, err := datastore.Put(c, key, &m1)
	return key, err
}

func DescList(c appengine.Context) (minutes []Minutes, err error) {
	q := datastore.NewQuery("minutes").Order("-CreatedAt")

	_, err = q.GetAll(c, &minutes)

	return minutes, err
}

