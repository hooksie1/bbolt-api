package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go.etcd.io/bbolt"
	"net/http"
)

//TODO: Implement GetBucketByID handler
func GetBucketByID(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	data := make(map[string]string)
	db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(vars["id"]))

		cursor := bucket.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			data[string(k)] = string(v)
		}

		return nil
	})

	resp, err := json.Marshal(data)
	if err != nil {
		return NewHTTPError(err, 500, "error marshaling response")
	}

	w.Write(resp)

	return nil
}

func CreateBucket(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucket([]byte(vars["id"]))
		if err != nil {
			return NewHTTPError(err, 500, "error creating bucket")
		}

		return nil

	})
	return nil
}

//TODO: Implement DeleteBucketByID handler
func DeleteBucketByID(w http.ResponseWriter, r *http.Request) error {

	return nil
}