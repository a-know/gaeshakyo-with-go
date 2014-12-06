package memo

import (
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
	"appengine/taskqueue"
	"appengine/user"
	"time"

	"code.google.com/p/go-uuid/uuid"
)

type Memo struct {
	Key       *datastore.Key
	Memo      string
	Minutes   *datastore.Key
	Author    user.User
	CreatedAt time.Time
}

const ascListMemkey = "LIST_OF_MEMO"

func SaveAs(c appengine.Context, minutesKey *datastore.Key, memoString string, u *user.User) (*datastore.Key, error) {
	key := datastore.NewKey(c, "memo", uuid.New(), 0, nil)

	m1 := Memo{
		Key:       key,
		Memo:      memoString,
		Minutes:   minutesKey,
		Author:    *u,
		CreatedAt: time.Now(),
	}

	// put
	_, err := datastore.Put(c, key, &m1)
	memcache.Delete(c, ascListMemkey)

	// post taskqueue
	task := taskqueue.NewPOSTTask("/tq/IncrementMemoCount", map[string][]string{"minutesKey": {minutesKey.Encode()}})
	_, err = taskqueue.Add(c, task, "")

	return key, err
}

func AscList(c appengine.Context, minutesKey *datastore.Key) (memo []Memo, err error) {
	memcache.Gob.Get(c, ascListMemkey, &memo)

	// item not found in memcache
	if memo == nil {
		q := datastore.NewQuery("memo").Filter("Minutes =", minutesKey).Order("CreatedAt")
		_, err = q.GetAll(c, &memo)

		// put item to memcache
		mem_item := &memcache.Item{
			Key:    ascListMemkey,
			Object: memo,
		}
		memcache.Gob.Add(c, mem_item)
	}

	return memo, err
}
