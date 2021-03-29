package server

import (
	"go.etcd.io/bbolt"
	"net/http"
	"strconv"
)

func backupDB(w http.ResponseWriter, r *http.Request) error {
	err := db.View(func(tx *bbolt.Tx) error {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", `attachment; filename="piggy-backup.db"`)
		w.Header().Set("Content-Length", strconv.Itoa(int(tx.Size())))
		_, err := tx.WriteTo(w)
		return err
	})
	if err != nil {
		return NewHTTPError(err, 500, "error creating backup")
	}

	return nil
}
