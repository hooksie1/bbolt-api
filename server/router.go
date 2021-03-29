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
	"fmt"
	"go.etcd.io/bbolt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var dbName = os.Getenv("DATABASE_PATH")
var db *bbolt.DB

type errHandler func(http.ResponseWriter, *http.Request) error


func checkEnv() bool {
	if os.Getenv("DATABASE_PATH") == "" || os.Getenv("SERVER_PORT") == "" {
		return false
	}

	return true
}

func Serve() {
	var err error

	if !checkEnv() {
		log.Println("Must set database path and server port")
		os.Exit(1)
	}

	db, err = bbolt.Open(dbName, 0644, nil)
	if err != nil {
		log.Printf("error opening database: %s", err)
		return
	}

	defer db.Close()


	router := mux.NewRouter().StrictSlash(true)
	apiRouter := router.PathPrefix("/v1").Subrouter().StrictSlash(true)

	apiRouter.Handle("/buckets/{bucket}", errHandler(getBucketByID)).Methods("GET")
	apiRouter.Handle("/buckets/{bucket}", errHandler(createBucket)).Methods("POST")
	apiRouter.Handle("/buckets/{bucket}", errHandler(deleteBucketByID)).Methods("DELETE")
	apiRouter.Handle("/buckets/{bucket}/keys", errHandler(getBucketKeys)).Methods("GET")
	apiRouter.Handle("/buckets/{bucket}/keys/{key}", errHandler(getKVByID)).Methods("GET")
	apiRouter.Handle("/buckets/{bucket}/keys/{key}", errHandler(createKV)).Methods("POST")
	apiRouter.Handle("/buckets/{bucket}/keys/{key}", errHandler(deleteKVByID)).Methods("DELETE")

	apiRouter.Use(logger)


	port := fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))
	log.Fatal(http.ListenAndServe(port, router))

}

