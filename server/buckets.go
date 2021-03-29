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
	"github.com/gorilla/mux"
	"go.etcd.io/bbolt"
	"net/http"
)

func GetBucketByID(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	if err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(vars["bucket"]))
		if b == nil {
			return NewHTTPError(nil, 404, "bucket not found")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func CreateBucket(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	if err := db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(vars["bucket"]))
		if err != nil {
			return NewHTTPError(err, 500, "error creating bucket")
		}

		return nil

	}); err != nil {
		return err
	}
	return nil
}

func DeleteBucketByID(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	if err := db.Update(func(tx *bbolt.Tx) error {
		if err := tx.DeleteBucket([]byte(vars["bucket"])); err != nil {
			return NewHTTPError(err, 500, "error deleting bucket")
		}

		return nil
	}); err != nil {
		return err
	}
	return nil
}