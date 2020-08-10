package db

import (
	"fmt"

	bolt "go.etcd.io/bbolt"
)

var defaultBucket = []byte("Default")

// Database type using bolt
type Database struct {
	db *bolt.DB
}

// NewDatabase returns a DB instance
func NewDatabase(dbPath string) (db *Database, closeFunc func() error, err error) {
	boltDb, err := bolt.Open(dbPath, 0666, nil)

	if err != nil {
		return nil, nil, err
	}
	// Disk sync flushes data from disc to DB , good idea to keep it
	// if you're not cool with loosing data ; it does speed it up tho if you dont
	// boltDb.NoSync = true

	db = &Database{db: boltDb}
	closeFunc = boltDb.Close

	if err := db.createDefaultBucket(); err != nil {
		closeFunc()
		return nil, nil, fmt.Errorf("creating deafult bucker : %w", err)
	}

	return db, closeFunc, nil
}

// SetKey sets the key to requested value or returns error
func (d *Database) SetKey(key string, value []byte) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		// tx.WriteTo()
		b := tx.Bucket(defaultBucket)
		return b.Put([]byte(key), value)
	})

}

// GetKey gets the key to requested value or returns error
func (d *Database) GetKey(key string) ([]byte, error) {
	var result []byte
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucket)
		result = b.Get([]byte(key))
		return nil

	})

	if err == nil {
		return result, nil
	}
	return nil, err

}

func (d *Database) createDefaultBucket() error {
	return d.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(defaultBucket))
		return err
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
