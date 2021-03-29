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
	"go.etcd.io/bbolt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var dbName = os.Getenv("DATABASE_PATH")
var Version = "dev"
var db *bbolt.DB

type ErrHandler func(http.ResponseWriter, *http.Request) error

func init() {
	var err error

	if !checkEnv() {
		log.Println("Must set database host or path")
		os.Exit(1)
	}

	db, err = bbolt.Open(dbName, 0644, nil)
	if err != nil {
		log.Printf("error opening database: %s", err)
		return
	}

}

func checkEnv() bool {
	if os.Getenv("DATABASE_PATH") == "" {
		return false
	}

	return true
}

func Serve() {
	defer db.Close()

	router := mux.NewRouter().StrictSlash(true)

	router.PathPrefix(Version)
	router.Handle("/buckets/{bucket}", ErrHandler(GetBucketByID)).Methods("GET")
	router.Handle("/buckets/{bucket}", ErrHandler(CreateBucket)).Methods("POST")
	router.Handle("/buckets/{bucket}", ErrHandler(DeleteBucketByID)).Methods("DELETE")
	router.Handle("/buckets/{bucket}/keys", ErrHandler(GetBucketKeys)).Methods("GET")
	router.Handle("/buckets/{bucket}/keys/{key}", ErrHandler(GetKVByID)).Methods("GET")
	router.Handle("/buckets/{bucket}/keys/{key}", ErrHandler(CreateKV)).Methods("POST")
	router.Handle("/buckets/{bucket}/keys/{key}", ErrHandler(DeleteKVByID)).Methods("DELETE")

	router.Use(logger)

	log.Fatal(http.ListenAndServe(":8080", router))

}

