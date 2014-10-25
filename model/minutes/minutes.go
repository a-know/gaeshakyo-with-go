package minutes

import (
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
	"appengine/user"
	"time"

	"code.google.com/p/go-uuid/uuid"
)

type Minutes struct {
	Key       *datastore.Key
	Title     string
	Author    user.User
	CreatedAt time.Time
}

func SaveAs(c appengine.Context, title string, u *user.User) (*datastore.Key, error) {
	key := datastore.NewKey(c, "minutes", uuid.New(), 0, nil)

	m1 := Minutes{
		Key:       key,
		Title:     title,
		Author:    *u,
		CreatedAt: time.Now(),
	}

	// put
	_, err := datastore.Put(c, key, &m1)
	memcache.Delete(c, descListMemkey)
	return key, err
}

const descListMemkey = "LIST_OF_MINUTES"

func DescList(c appengine.Context) (minutes []Minutes, err error) {
	memcache.Gob.Get(c, descListMemkey, &minutes)

	// item not found in memcache
	if minutes == nil {
		q := datastore.NewQuery("minutes").Order("-CreatedAt")
		_, err = q.GetAll(c, &minutes)

		// put item to memcache
		mem_item := &memcache.Item{
			Key:    descListMemkey,
			Object: minutes,
		}
		memcache.Gob.Add(c, mem_item)
	}

	return minutes, err
}
