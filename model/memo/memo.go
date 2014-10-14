package memo

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"time"

	"code.google.com/p/go-uuid/uuid"
)

type Memo struct {
	Key       *datastore.Key
	Memo      string
	Minutes   *datastore.Key
	Author    *user.User
	CreatedAt time.Time
}

func SaveAs(c appengine.Context, minutesKey *datastore.Key, memoString string, u *user.User) (*datastore.Key, error) {
	key := datastore.NewKey(c, "memo", uuid.New(), 0, nil)

	m1 := Memo{
		Key:       key,
		Memo:      memoString,
		Minutes:   minutesKey,
		Author:    u,
		CreatedAt: time.Now(),
	}

	// put
	_, err := datastore.Put(c, key, &m1)
	return key, err
}

func AscList(c appengine.Context, minutesKey *datastore.Key) (memo []Memo, err error) {
	q := datastore.NewQuery("memo").Filter("Minutes =", minutesKey).Order("CreatedAt")

	_, err = q.GetAll(c, &memo)

	return memo, err
}
