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
	"log"
	"net/http"
	"os"

	"github.com/spf13/viper"

	"github.com/gorilla/mux"
)

var dbName = os.Getenv("DATABASE_PATH")
var Version = "dev"

type ErrHandler func(http.ResponseWriter, *http.Request) error

func checkEnv() bool {
	if os.Getenv("DATABASE_PATH") == "" && viper.GetBool("devMode") {
		return false
	}

	return true
}

func Serve() {

	if !checkEnv() {
		log.Println("Must set database host or path")
		os.Exit(1)
	}

	router := mux.NewRouter().StrictSlash(true)

	apiRouter := router.PathPrefix(Version).Subrouter().StrictSlash(true)
	apiRouter.Handle("/buckets", ErrHandler(GetBuckets)).Methods("GET")
	apiRouter.Handle("/buckets/{id}", ErrHandler(GetBucketByID)).Methods("GET")
	apiRouter.Handle("/buckets/{id}", ErrHandler(CreateBucket)).Methods("POST")
	apiRouter.Handle("/buckets/{id}", ErrHandler(DeleteBucketByID)).Methods("DELETE")
	apiRouter.Handle("/buckets/{id}/{key}", ErrHandler(GetKVByID)).Methods("GET")
	apiRouter.Handle("/buckets/{id}/{key}", ErrHandler(CreateKV)).Methods("POST")
	apiRouter.Handle("/buckets/{id}/{key}", ErrHandler(DeleteKVByID)).Methods("DELETE")

	apiRouter.Use(logger)

	log.Fatal(http.ListenAndServe(":8080", router))

}

