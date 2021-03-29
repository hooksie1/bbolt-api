/*
Copyright © 2021 John Hooks john@hooks.technology

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go.etcd.io/bbolt"
	"net/http"
)

type Record struct {
	Data string `json:"data"`
}

func getBucketKeys(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	data := make(map[string]string)
	if err := db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(vars["bucket"]))
		cursor := bucket.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			data[string(k)] = string(v)
		}

		return nil
	}); err != nil {
		return err
	}

	resp, err := json.Marshal(data)
	if err != nil {
		return NewHTTPError(err, 500, "error marshaling response")
	}

	w.Write(resp)

	return nil
}

func getKVByID(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	var record Record
	if err := db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(vars["bucket"]))
		value := bucket.Get([]byte(vars["key"]))
		if value == nil {
			return NewHTTPError(nil, 404, "kv not found")
		}

		record.Data = string(value)

		return nil
	}); err != nil {
		return err
	}

	err := json.NewEncoder(w).Encode(record)
	if err != nil {
		return NewHTTPError(err, 500, "error encoding json data")
	}

	return nil
}

func createKV(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	record, err := keyData(r.Body)
	if err != nil {
		return NewHTTPError(err, 500, "error decoding JSON data")
	}
	db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(vars["bucket"]))
		if err := bucket.Put([]byte(vars["key"]), []byte(record.Data)); err != nil {
			return NewHTTPError(err, 500, "error writing key/value pair")
		}

		return nil
	})

	return nil
}

func deleteKVByID(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(vars["bucket"]))
		if err := bucket.Delete([]byte(vars["key"])); err != nil {
			return NewHTTPError(err, 500, "error deleting key")
		}

		return nil
	})

	return nil
}