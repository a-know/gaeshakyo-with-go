package memo

import (
	"appengine"
	"appengine/datastore"
	"time"
)

type Memo struct {
	Key *datastore.Key
	Memo   string
	Minutes *datastore.Key
	CreatedAt time.Time
}

func SaveAs(c appengine.Context, minutesKey *datastore.Key, memoString string) (*datastore.Key, error) {
	key := datastore.NewIncompleteKey(c, "memo", nil)

	m1 := Memo{
		Key:       key,
		Memo:      memoString,
		Minutes:   minutesKey,
		CreatedAt: time.Now(),
	}

	// put
	_, err := datastore.Put(c, key, &m1)
	return key, err
}

func AscList(c appengine.Context, minutesKey *datastore.Key) (memo []Memo, err error) {
	q := datastore.NewQuery("memo").Filter("Minutes", minutesKey).Order("+CreatedAt")

	_, err = q.GetAll(c, &memo)

	return memo, err
}
