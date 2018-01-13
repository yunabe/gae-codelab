package datastore

import (
	"net/http"
	"strconv"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type Entity struct {
	Name    string
	Value   int64
	Created time.Time
}

func StoreEntity(ctx context.Context, name string, value int64) error {
	e := &Entity{
		Name:    name,
		Value:   value,
		Created: time.Now(),
	}
	key := datastore.NewIncompleteKey(ctx, "Entity", nil)
	// key := datastore.Key{}
	key, err := datastore.Put(ctx, key, e)

	log.Infof(ctx, "new key = %#v", key)
	return err
}

func handleStoreEntity(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	name := r.FormValue("name")
	if name == "" {
		name = "World"
	}
	value, err := strconv.ParseInt(r.FormValue("value"), 0, 64)
	if err != nil {
		log.Errorf(ctx, "Failed to parse value: %v", err)
	}
	if err := StoreEntity(ctx, name, value); err != nil {
		log.Errorf(ctx, "Failed to store an entity: %v", err)
	}
}

func RegisterHandlers() {
	http.HandleFunc("/datastore/store_entity", handleStoreEntity)
}
