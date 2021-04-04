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
	"fmt"
	"github.com/gorilla/mux"
	"go.etcd.io/bbolt"
	"net/http"
)

type BucketList struct {
	Buckets []string `json:"buckets"`
}

func getBuckets(w http.ResponseWriter, r *http.Request) error {
	var buckets BucketList

	if err := db.View(func(tx *bbolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bbolt.Bucket) error {
			b := []string{string(name)}
			buckets.Buckets = append(buckets.Buckets, b...)
			return nil
		})
	}); err != nil {
		return fmt.Errorf("error getting buckets: %s", err)
	}

	if err := json.NewEncoder(w).Encode(&buckets); err != nil {
		return fmt.Errorf("error encoding json bucket data: %s", err)
	}


	return nil
}
func getBucketByID(w http.ResponseWriter, r *http.Request) error {
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

func createBucket(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	if err := db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(vars["bucket"]))
		if err != nil {
			return fmt.Errorf("error creating bucket: %s", err)
		}

		return nil

	}); err != nil {
		return err
	}
	return nil
}

func deleteBucketByID(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	if err := db.Update(func(tx *bbolt.Tx) error {
		if err := tx.DeleteBucket([]byte(vars["bucket"])); err != nil {
			return fmt.Errorf("error deleting bucket: %s", err)
		}

		return nil
	}); err != nil {
		return err
	}
	return nil
}