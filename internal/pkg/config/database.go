package config

import bolt "go.etcd.io/bbolt"

var database *bolt.DB

func (c *Config) GetBoltDB() *bolt.DB {
	if database == nil {
		db, err := bolt.Open(c.Database.Filepath, 0600, nil)

		if err != nil {
			panic(err)
		}

		database = db
	}

	return database
}