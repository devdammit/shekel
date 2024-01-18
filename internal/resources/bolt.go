package resources

import (
	"encoding/binary"
	"time"

	"go.etcd.io/bbolt"
)

var BoltRootBucket = []byte("root")

type Bolt struct {
	*bbolt.DB

	path string
}

func NewBolt(path string) *Bolt {
	return &Bolt{
		path: path,
	}
}

func (s *Bolt) GetName() string {
	return "bolt"
}

func (s *Bolt) Start() error {
	db, err := bbolt.Open(s.path, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}

	s.DB = db

	return db.Update(func(tx *bbolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists(BoltRootBucket)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *Bolt) Stop() error {
	return s.DB.Close()
}

func Itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
