package memo

import (
	"appengine/datastore"
	"time"
)

type Memo struct {
	Memo   string
	Minutes *datastore.Key
	CreatedAt time.Time
}
