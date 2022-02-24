package store

import bolt "go.etcd.io/bbolt"
import "github.com/davecgh/go-spew/spew"

const ChannelBucketKey = "channel"

func (f *FileDB) ToggleSubscription(channel string, subreddit string, subscribed bool) error {
	return f.DB.Update(func (tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(ChannelBucketKey))

		if bucket == nil {
			panic("channel bucket doesn't exist")
		}

		channelSubscriptions, err := bucket.CreateBucketIfNotExists([]byte(channel))

		if err != nil {
			return err
		}

		// TODO: store something sensible instead of a string representation of
		// a boolean
		subscribedString := "true"

		if !subscribed {
			subscribedString = "false"
		}

		return channelSubscriptions.Put([]byte(subreddit), []byte(subscribedString))
	})
}

func (f *FileDB) GetSubscriptions() (map[string]map[string]bool, error) {
	subscriptions := make(map[string]map[string]bool)

	err := f.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(ChannelBucketKey))

		if bucket == nil {
			panic("channel bucket doesn't exist")
		}

		return bucket.ForEach(func(k, v []byte) error {
			channel := string(k)
			channelBucket := bucket.Bucket(k)

			if channelBucket == nil {
				return nil
			}

			return channelBucket.ForEach(func (k, v []byte) error {
				subreddit := string(k)
				subscriptionString := string(v)

				if subscriptions[channel] == nil {
					subscriptions[channel] = make(map[string]bool)
				}

				if subscriptionString == "true" {
					subscriptions[channel][subreddit] = true
				}

				return nil
			})
		})
	})

	spew.Dump(subscriptions)

	return subscriptions, err
}