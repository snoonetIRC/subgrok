package store

import bolt "go.etcd.io/bbolt"

const ChannelBucketKey = "channel"

func (f *FileDB) ToggleSubscription(channel string, subreddit string, subscribed bool) error {
	return f.DB.Update(func(tx *bolt.Tx) error {
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

			return channelBucket.ForEach(func(k, v []byte) error {
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

	return subscriptions, err
}

func (f *FileDB) GetChannels() ([]string, error) {
	subscriptions, err := f.GetSubscriptions()

	if err != nil {
		return nil, err
	}

	var channels []string

	for channel, subreddits := range subscriptions {
		for _, subscribed := range subreddits {
			if subscribed {
				channels = append(channels, channel)
				break
			}
		}
	}

	return channels, nil
}

func (f *FileDB) GetChannelSubscriptions(name string) ([]string, error) {
	var subscriptions []string

	err := f.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(ChannelBucketKey))

		if bucket == nil {
			panic("channel bucket doesn't exist")
		}

		channelBucket := bucket.Bucket([]byte(name))

		if channelBucket == nil {
			return nil
		}

		return channelBucket.ForEach(func(k, v []byte) error {
			subscriptionString := string(v)

			if subscriptionString == "true" {
				subscriptions = append(subscriptions, string(k))
			}

			return nil
		})
	})

	return subscriptions, err
}

func (f *FileDB) GetChannelSubscriptionCount(name string) (int, error) {
	subscriptions, err := f.GetChannelSubscriptions(name)

	if err != nil {
		return 0, err
	}

	return len(subscriptions), nil
}
