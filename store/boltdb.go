package store

import (
	"fmt"
	"log"

	bolt "go.etcd.io/bbolt"
)

// Bolt implements the store interface and holds a connection to the BoltDB.
type Bolt struct {
	Client *bolt.DB
	bucket string
}

// ConnectBolt opens a connection to Bolt and creates the bucket if it doesn't exist.
func ConnectBolt(bucket string) (Store, error) {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create new Bolt instance: %w", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &Bolt{Client: db, bucket: bucket}, nil
}

// Set sets key value.
func (r *Bolt) Set(key string, value []byte) error {
	return r.Client.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(r.bucket))
		return b.Put([]byte(key), value)
	})
}

// Get gets a value from Bolt.
func (r *Bolt) Get(key string) ([]byte, error) {
	log.Printf("fetching key %s from bolt", key)
	var result []byte
	if err := r.Client.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(r.bucket))
		result = b.Get([]byte(key))
		return nil
	}); err != nil {
		return nil, fmt.Errorf("unable to fetch key %s : %w", key, err)
	}

	return result, nil
}

// GetAll will fetch all records from Bolt.
func (r *Bolt) GetAll() ([][]byte, error) {
	var results [][]byte

	if err := r.Client.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(r.bucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", r.bucket)
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			results = append(results, v)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return results, nil
}

// Disconnect will disconnect the Bolt connection.
func (r *Bolt) Disconnect() error {
	return r.Client.Close()
}
