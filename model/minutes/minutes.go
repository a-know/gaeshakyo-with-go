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

func Saveas(c appengine.Context, title string) (*datastore.Key, error) {
	g1 := Minutes{
		Title:   title,
		CreatedAt: time.Now(),
	}

	// put
	key := datastore.NewIncompleteKey(c, "guestbook", nil)
	_, err := datastore.Put(c, key, &g1)
	return key, err
}

func AscList(c appengine.Context) (minutes []Minutes, err error) {
	q := datastore.NewQuery("minutes").Order("+CreatedAt")

	_, err = q.GetAll(c, &minutes)

	return minutes, err
}

