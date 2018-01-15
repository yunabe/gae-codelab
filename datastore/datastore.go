package datastore

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"time"

	cloudds "cloud.google.com/go/datastore"
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

func handleConflictNew(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	rn, err := rand.Int(rand.Reader, big.NewInt(1000*1000))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	keyStr := fmt.Sprintf("randkey%d", rn.Int64())
	key := datastore.NewKey(ctx, "Entity", keyStr, 0, nil)
	done := make(chan error)
	for i := 0; i < 2; i++ {
		go func(i int) {
			done <- datastore.RunInTransaction(ctx, func(ctx context.Context) error {
				var entity Entity
				err := datastore.Get(ctx, key, &entity)
				if err != nil && err != datastore.ErrNoSuchEntity {
					return err
				}
				if err == nil {
					return fmt.Errorf("%s already exists", keyStr)
				}
				time.Sleep(time.Second)
				datastore.Put(ctx, key, &Entity{Value: int64(i)})
				return nil
			}, &datastore.TransactionOptions{
				Attempts: 1,
			})
		}(i)
	}
	w.Header().Add("Content-Type", "text/plain")
	for i := 0; i < 2; i++ {
		if err := <-done; err != nil {
			log.Errorf(ctx, "%v", err)
			fmt.Fprintln(w, err)
		} else {
			fmt.Fprintln(w, "OK!")
		}
	}
}

// It seems cloud datastore API also works in GAE SE.
func handleCloudDatastoreAPI(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	client, err := cloudds.NewClient(ctx, "yunabe-codelab")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = client.Put(ctx, cloudds.NameKey("Entity", "cloudapi", nil), &Entity{
		Name:  "cloud-api-entity",
		Value: 1234,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func RegisterHandlers() {
	http.HandleFunc("/datastore/store_entity", handleStoreEntity)
	http.HandleFunc("/datastore/conflict_new", handleConflictNew)
	http.HandleFunc("/datastore/cloud_api", handleCloudDatastoreAPI)
}
