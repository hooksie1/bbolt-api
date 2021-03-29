package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go.etcd.io/bbolt"
	"net/http"
)

//TODO: Implement GetKVByID handler
func GetKVByID(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	data := make(map[string]string)
	db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(vars["id"]))
		value := bucket.Get([]byte(vars["key"]))
		if value == nil {
			return NewHTTPError(nil, 404, "kv not found")
		}

		data[vars["key"]] = string(value)

		return nil
	})

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return NewHTTPError(err, 500, "error encoding json data")
	}

	return nil
}

//TODO: Implement CreateKV handler
func CreateKV(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	kv, err := jsonData(r.Body)
	if err != nil {
		return NewHTTPError(err, 500, "error decoding JSON data")
	}
	db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(vars["id"]))
		for k, v := range kv {
			err := bucket.Put([]byte(k), []byte(v))
			if err != nil {
				return NewHTTPError(err, 500, "error writing to database")
			}
		}

		return nil
	})

	return nil
}

//TODO: Implement DeleteKVByID handler
func DeleteKVByID(w http.ResponseWriter, r *http.Request) error {

	return nil
}