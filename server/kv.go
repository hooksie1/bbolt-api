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