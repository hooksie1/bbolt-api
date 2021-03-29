package server

import (
	"encoding/json"
	"io"
)

func jsonData(r io.Reader) (map[string]string, error){
	kv := make(map[string]string)
	err := json.NewDecoder(r).Decode(&kv)
	if err != nil {
		return nil, err
	}

	return kv, nil

}
