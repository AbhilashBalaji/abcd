package db_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/AbhilashBalaji/abcd/db"
)

func setKey(t *testing.T, d *db.Database, key, value string) {
	t.Helper()
	if err := d.SetKey(key, []byte(value)); err != nil {
		t.Fatalf("SetKey(%q,%q)failed: %v", key, value, err)
	}
}

func getKey(t *testing.T, d *db.Database, key string) string {
	t.Helper()
	value, err := d.GetKey(key)

	if err != nil {

		t.Fatalf("GetKey(%q)failed: %v", key, err)

	}
	return string(value)
}

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

	setKey(t, db, "name", "Jeff")

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

func TestDeleteExtraKeys(t *testing.T) {
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

	setKey(t, db, "name", "Jeff")
	setKey(t, db, "foo", "bar")

	value := getKey(t, db, "name")
	if value != "Jeff" {
		t.Errorf(`Unexpected value for key "bruh": got %q , want %q`, value, "moment")
	}

	if err := db.DeleteExtraKeys(func(name string) bool { return name == "name" }); err != nil {
		t.Fatalf("Could not delete extra keys : %v", err)
	}
	if value := getKey(t, db, "foo"); value != "bar" {
		t.Errorf(`unexpected value for key "name": got %q , want %q`, value, "")
	}

}
