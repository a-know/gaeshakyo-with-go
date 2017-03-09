package memo

import (
	"appengine"
	"appengine/channel"
	"appengine/datastore"
	"appengine/memcache"
	"appengine/taskqueue"
	"appengine/user"
	"encoding/json"
	"net/url"
	"time"

	"github.com/pborman/uuid"
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
	task := taskqueue.NewPOSTTask("/tq/IncrementMemoCount", url.Values{
		"minutesKey": {minutesKey.Encode()},
	})
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

func PushNotification(c appengine.Context, memoKey *datastore.Key) error {
	var m Memo
	err := datastore.Get(c, memoKey, &m)
	if err != nil {
		return err
	}

	q := datastore.NewQuery("minutes_channel").Filter("MinutesKey =", m.Minutes).KeysOnly()
	minutesChannelKeyList, q_err := q.GetAll(c, nil)
	if q_err != nil {
		return q_err
	}

	// push 通知する内容 = 追加された memo の内容、なので、それをあらかじめ json 化しておく
	js, js_err := json.Marshal(m)
	if js_err != nil {
		return js_err
	}

	for _, minutesChannelKey := range minutesChannelKeyList {
		send_err := channel.SendJSON(c, minutesChannelKey.StringID(), string(js))
		if send_err != nil {
			c.Errorf("failed to push notification: %v", send_err)
			return send_err
		}
	}
	return nil
}
