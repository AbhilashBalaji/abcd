package db_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/AbhilashBalaji/abcd/db"
)

func TestGetSet(t *testing.T) {
	f, err := ioutil.TempFile(os.TempDir(), "abcdDB")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}
	name := f.Name()
	defer f.Close()
	defer os.Remove(name)
	db, closeFunc, err := db.NewDatabase(name)
	if err != nil {
		t.Fatalf("could not create a new database :%v", err)
	}
	defer closeFunc()

	if err := db.SetKey("bruh", []byte("moment")); err != nil {
		t.Fatalf("could not write key : %v", err)
	}

	value, err := db.GetKey("bruh")
	if err != nil {
		t.Fatalf(`could not get the key "bruh" : %v`, err)
	}
	if !bytes.Equal(value, []byte("moment")) {
		t.Errorf(`Unexpected value for key "bruh": got %q , want %q`, value, "moment")
	}

}
