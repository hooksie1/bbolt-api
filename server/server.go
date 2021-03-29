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
	"log"
	"net/http"
	"time"
)



type ClientError interface {
	Error() string
	Body() ([]byte, error)
	Headers() (int, map[string]string)
}

type HTTPError struct {
	Cause   error  `json:"-"`
	Details string `json:"details"`
	Status  int    `json:"-"`
}

func (e *HTTPError) Error() string {
	if e.Cause == nil {
		return e.Details
	}
	return e.Details + " : " + e.Cause.Error()
}

func (e *HTTPError) Body() ([]byte, error) {
	body, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling response: %v", err)
	}
	return body, nil
}

func (e *HTTPError) Headers() (int, map[string]string) {
	return e.Status, map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}
}

func NewHTTPError(err error, status int, detail string) error {
	return &HTTPError{
		Cause:   err,
		Details: detail,
		Status:  status,
	}
}

// logger logs the endpoint requested and times how long the request takes.
func logger(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s %s %s",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}

func (fn ErrHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	err := fn(w, r)
	if err == nil {
		return
	}

	clientError, ok := err.(ClientError)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := clientError.Body()
	if err != nil {
		log.Printf("An error ocurred: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	status, headers := clientError.Headers()
	for k, v := range headers {
		w.Header().Set(k, v)
	}

	w.WriteHeader(status)
	w.Write(body)

}

