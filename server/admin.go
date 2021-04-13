package server

import (
	"encoding/json"
	"go.etcd.io/bbolt"
	"net/http"
	"strconv"
)

type Statistics struct {
	Reads int `json:"total_read_transactions"`
	Writes int `json:"total_writes"`
	Time float64 `json:"total_disk_write_duration"`
}

func backupDB(w http.ResponseWriter, r *http.Request) error {
	err := db.View(func(tx *bbolt.Tx) error {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", `attachment; filename="db-backup.db"`)
		w.Header().Set("Content-Length", strconv.Itoa(int(tx.Size())))
		_, err := tx.WriteTo(w)
		return err
	})
	if err != nil {
		return NewHTTPError(err, 500, "error creating backup")
	}

	return nil
}

func getDBStats(w http.ResponseWriter, r *http.Request) error {
	stats := db.Stats()

	statistics := Statistics{
		Reads: stats.TxN,
		Writes: stats.TxStats.Write,
		Time: stats.TxStats.WriteTime.Seconds(),
	}

	json.NewEncoder(w).Encode(&statistics)

	return nil
}