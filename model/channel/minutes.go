package channel

import (
	"appengine/datastore"
	"time"
)

type MinutesChannel struct {
	Key        *datastore.Key
	Token      string
	MinutesKey *datastore.Key
	CreatedAt  time.Time
}
