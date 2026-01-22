package database

import (
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"os"
)

// Initialize the database
func Initialize() (*leveldb.DB, error) {
	path := os.Getenv("LEVELDB_PATH")
	if path == "" {
		path = "/tmp/foo.db"
	}
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}
	return db, nil
}
