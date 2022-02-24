package store

import bolt "go.etcd.io/bbolt"

type FileDB struct {
	DB *bolt.DB
}

func NewStore(db *bolt.DB) (*FileDB, error) {
	fileDB := &FileDB{DB: db}

	err := fileDB.init()

	if err != nil {
		return nil, err
	}

	return fileDB, nil
}

func (f *FileDB) init() error {
	return f.DB.Update(func (tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(ChannelBucketKey))

		if err != nil {
			return err
		}

		return nil
	})
}