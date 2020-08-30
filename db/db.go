package db

import (
	"bytes"
	"errors"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

var defaultBucket = []byte("Default")
var replicaBucket = []byte("Default-replica")

// Database type using bolt
type Database struct {
	db       *bolt.DB
	readOnly bool
}

// NewDatabase returns a DB instance
func NewDatabase(dbPath string, readOnly bool) (db *Database, closeFunc func() error, err error) {
	boltDb, err := bolt.Open(dbPath, 0666, nil)

	if err != nil {
		return nil, nil, err
	}
	// Disk sync flushes data from disc to DB , good idea to keep it
	// if you're not cool with loosing data ; it does speed it up tho if you dont
	// boltDb.NoSync = true

	db = &Database{db: boltDb, readOnly: readOnly}
	closeFunc = boltDb.Close

	if err := db.createDefaultBuckets(); err != nil {
		closeFunc()
		return nil, nil, fmt.Errorf("creating deafult bucker : %w", err)
	}

	return db, closeFunc, nil
}

// SetKey sets the key to requested value or returns error
func (d *Database) SetKey(key string, value []byte) error {
	if d.readOnly {
		return errors.New("Read Only mode")
	}
	return d.db.Update(func(tx *bolt.Tx) error {
		if err := tx.Bucket(defaultBucket).Put([]byte(key), value); err != nil {
			return err
		}
		return tx.Bucket(replicaBucket).Put([]byte(key), value)
	})

}

// SetKeyOnReplica sets the key to the requested value into the default database and does not write
// to the replication queue.
// This method is intended to be used only on replicas.
func (d *Database) SetKeyOnReplica(key string, value []byte) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(defaultBucket).Put([]byte(key), value)
	})
}

func copyByteSlice(b []byte) []byte {
	if b == nil {
		return nil
	}
	res := make([]byte, len(b))
	copy(res, b)
	return res
}

// GetNextKeyForReplication returns k,v for keys that
// have changed but not applied to Replica.
// If there  are no keys,
func (d *Database) GetNextKeyForReplication() (key, value []byte, err error) {
	err = d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(replicaBucket)
		k, v := b.Cursor().First()
		key = copyByteSlice(k)
		value = copyByteSlice(v)
		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return key, value, nil
}

// GetKey gets the key to requested value or returns error
func (d *Database) GetKey(key string) ([]byte, error) {
	var result []byte
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucket)
		result = copyByteSlice(b.Get([]byte(key)))
		return nil

	})

	if err == nil {
		return result, nil
	}
	return nil, err

}

func (d *Database) createDefaultBuckets() error {
	return d.db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(defaultBucket)); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists([]byte(replicaBucket)); err != nil {
			return err
		}
		return nil
	})
}

// DeleteReplicationKey deletes the key from the replication queue
// if the value matches the contents or if the key is already absent.
func (d *Database) DeleteReplicationKey(key, value []byte) (err error) {
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(replicaBucket)

		v := b.Get(key)
		if v == nil {
			return errors.New("key does not exist")
		}

		if !bytes.Equal(v, value) {
			return errors.New("value does not match")
		}

		return b.Delete(key)
	})
}

//DeleteExtraKeys deletes keys that do not belong in current shard (on new shard)
func (d *Database) DeleteExtraKeys(isExtra func(string) bool) error {
	var keys []string
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucket)
		return b.ForEach(func(k, v []byte) error {
			ks := string(k)
			if isExtra(ks) {
				keys = append(keys, string(k))
			}
			return nil
		})
	})
	if err != nil {
		return err
	}
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucket)

		for _, k := range keys {
			if err := b.Delete([]byte(k)); err != nil {
				return err
			}
		}
		return nil
	})
}
