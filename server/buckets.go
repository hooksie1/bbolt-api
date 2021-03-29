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