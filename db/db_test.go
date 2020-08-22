package db_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/AbhilashBalaji/abcd/db"
)

func createTempDb(t *testing.T, readOnly bool) *db.Database {
	t.Helper()
	f, err := ioutil.TempFile(os.TempDir(), "abcdDB")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}
	name := f.Name()
	f.Close()
	t.Cleanup(func() { os.Remove(name) })

	db, closeFunc, err := db.NewDatabase(name, readOnly)
	if err != nil {
		t.Fatalf("could not create a new database :%v", err)
	}
	t.Cleanup(func() { closeFunc() })
	return db
}

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
	db := createTempDb(t, false)

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
	k, v, err := db.GetNextKeyForReplication()
	if err != nil {
		t.Fatalf(`Unexpected error for GetNextKeyForReplication(): %v`, err)
	}
	if !bytes.Equal(k, []byte("bruh")) || !bytes.Equal(v, []byte("moment")) {
		t.Errorf(`GetNextKeyForReplication(): got %q, %q; want %q, %q`, k, v, "bruh", "moment")
	}

}
func TestDeleteReplicationKey(t *testing.T) {
	db := createTempDb(t, false)

	setKey(t, db, "party", "Great")

	k, v, err := db.GetNextKeyForReplication()
	if err != nil {
		t.Fatalf(`Unexpected error for GetNextKeyForReplication(): %v`, err)
	}

	if !bytes.Equal(k, []byte("party")) || !bytes.Equal(v, []byte("Great")) {
		t.Errorf(`GetNextKeyForReplication(): got %q, %q; want %q, %q`, k, v, "party", "Great")
	}

	if err := db.DeleteReplicationKey([]byte("party"), []byte("Bad")); err == nil {
		t.Fatalf(`DeleteReplicationKey("party", "Bad"): got nil error, want non-nil error`)
	}

	if err := db.DeleteReplicationKey([]byte("party"), []byte("Great")); err != nil {
		t.Fatalf(`DeleteReplicationKey("party", "Great"): got %q, want nil error`, err)
	}

	k, v, err = db.GetNextKeyForReplication()
	if err != nil {
		t.Fatalf(`Unexpected error for GetNextKeyForReplication(): %v`, err)
	}

	if k != nil || v != nil {
		t.Errorf(`GetNextKeyForReplication(): got %v, %v; want nil, nil`, k, v)
	}
}
func TestDeleteExtraKeys(t *testing.T) {
	db := createTempDb(t, false)

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

func TestSetReadOnly(t *testing.T) {
	db := createTempDb(t, true)

	if err := db.SetKey("bruh", []byte("moment")); err == nil {
		t.Fatalf("SetKey(%q, %q): got nil error, want non-nil error", []byte("moment"), err)
	}
}
