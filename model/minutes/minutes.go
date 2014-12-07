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
	MemoCount int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

const descListMemkey = "LIST_OF_MINUTES"

func SaveAs(c appengine.Context, title string, u *user.User) (*datastore.Key, error) {
	key := datastore.NewKey(c, "minutes", uuid.New(), 0, nil)

	m1 := Minutes{
		Key:       key,
		Title:     title,
		Author:    *u,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// put
	_, err := datastore.Put(c, key, &m1)
	if err != nil {
		return key, err
	}
	memcache.Delete(c, descListMemkey)
	return key, err
}

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

func IncrementMemoCount(c appengine.Context, minutesKey *datastore.Key) error {
	return datastore.RunInTransaction(c, func(c appengine.Context) error {
		var m Minutes
		err := datastore.Get(c, minutesKey, &m)
		if err != nil {
			return err
		}
		m.MemoCount++
		m.UpdatedAt = time.Now()
		_, err = datastore.Put(c, minutesKey, &m)
		memcache.Delete(c, descListMemkey)
		return err
	}, nil)
}
