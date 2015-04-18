package minutes

import (
	"net/http"

	"appengine"
	"appengine/datastore"
	"appengine/memcache"
	"appengine/user"
	"google.golang.org/appengine/urlfetch"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	newappengine "google.golang.org/appengine"
	"google.golang.org/cloud"
	"google.golang.org/cloud/storage"

	"model/memo"
)

type Minutes struct {
	Key       *datastore.Key
	Title     string
	Author    user.User
	MemoCount int
	CreatedAt time.Time
	UpdatedAt time.Time
}

const descListMemkey = "LIST_OF_MINUTES"
const bucket = "gaeshakyo-with-go-bucket"

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

func QueryForUpdateMemoCount(c appengine.Context) (minutesKeyList []*datastore.Key, err error) {
	now := time.Now()
	before24hours := now.AddDate(0, 0, -1)
	before48hours := now.AddDate(0, 0, -2)
	q := datastore.NewQuery("minutes").Filter("UpdatedAt <", before24hours).Filter("UpdatedAt >", before48hours).KeysOnly()
	minutesKeyList, err = q.GetAll(c, nil)
	return
}

func UpdateMemoCount(c appengine.Context, minutesKey *datastore.Key) (err error) {
	var count int
	q := datastore.NewQuery("memo").Filter("Minutes =", minutesKey)
	count, err = q.Count(c)
	if err != nil {
		return err
	}

	var m Minutes
	err = datastore.Get(c, minutesKey, &m)
	if err != nil {
		return err
	}

	m.MemoCount = count
	m.UpdatedAt = time.Now()

	_, err = datastore.Put(c, minutesKey, &m)
	if err != nil {
		return err
	}

	memcache.Delete(c, descListMemkey)
	return
}

func Delete(c appengine.Context, minutesKey *datastore.Key) (err error) {
	memcache.Delete(c, descListMemkey)
	// 議事録に紐付くメモの取得
	var memoKeyList []*datastore.Key
	q := datastore.NewQuery("memo").Filter("Minutes =", minutesKey).KeysOnly()
	memoKeyList, err = q.GetAll(c, nil)

	err = datastore.DeleteMulti(c, memoKeyList)
	if err != nil {
		return err
	}
	err = datastore.Delete(c, minutesKey)
	if err != nil {
		return err
	}
	return
}

func ExportAsTsv(r *http.Request, c appengine.Context, minutes Minutes) (fileName string, err error) {
	// Minutes タイトルの取得
	fileName = minutes.Title
	// Minutes に紐付く Memo の取得
	memo_list, memolist_err := memo.AscList(c, minutes.Key)
	if memolist_err != nil {
		return fileName, memolist_err
	}

	// see https://cloud.google.com/appengine/docs/go/googlecloudstorageclient/getstarted
	hc := &http.Client{}
	ctx := newappengine.NewContext(r)
	hc.Transport = &oauth2.Transport{
		Source: google.AppEngineTokenSource(ctx, storage.ScopeFullControl),
		Base:   &urlfetch.Transport{Context: ctx},
	}

	cloud_context := cloud.WithContext(ctx, newappengine.AppID(ctx), hc)

	wc := storage.NewWriter(cloud_context, bucket, fileName)
	wc.ContentType = "text/tab-separated-values"

	// 全 Memo の内容を書き出し
	var content string
	for _, memo := range memo_list {
		content = "\"" + memo.CreatedAt.String() + "\"\t\"" + memo.Author.Email + "\"\t\"" + memo.Memo + "\"\n"
		if _, write_err := wc.Write([]byte(content)); write_err != nil {
			return fileName, write_err
		}
	}

	if close_err := wc.Close(); close_err != nil {
		return fileName, close_err
	}

	return fileName, nil
}

func GetTsvUrl(r *http.Request, c appengine.Context, fileName string) (url string, err error) {
	// see http://godoc.org/google.golang.org/cloud/storage#SignedURL
	hc := &http.Client{}
	// ctx := cloud.NewContext(appengine.AppID(c), hc)
	ctx := newappengine.NewContext(r)
	hc.Transport = &oauth2.Transport{
		Source: google.AppEngineTokenSource(ctx, storage.ScopeFullControl),
		Base:   &urlfetch.Transport{Context: ctx},
	}

	url, err = storage.SignedURL(bucket, fileName, &storage.SignedURLOptions{Expires: time.Now().AddDate(1, 0, 0)}) // 1年後に失効
	return
}
