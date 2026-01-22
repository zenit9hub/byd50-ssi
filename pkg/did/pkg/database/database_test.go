package database

import (
	"os"
	"testing"
)

func TestInitialize(t *testing.T) {
	t.Setenv("LEVELDB_PATH", "/tmp/foo.db")
	_ = os.RemoveAll("/tmp/foo.db")

	db, err := Initialize()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = db.Close()
		_ = os.RemoveAll("/tmp/foo.db")
	}()

	if err := db.Put([]byte("k"), []byte("v"), nil); err != nil {
		t.Fatal(err)
	}
	got, err := db.Get([]byte("k"), nil)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "v" {
		t.Fatalf("unexpected value: %s", string(got))
	}
}

func TestInitializeError(t *testing.T) {
	t.Setenv("LEVELDB_PATH", "/dev/null/leveldb")

	if _, err := Initialize(); err == nil {
		t.Fatal("expected error opening invalid db path")
	}
}
