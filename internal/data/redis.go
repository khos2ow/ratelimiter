package data

import (
	"fmt"

	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
)

// Redis holds the client and options for connecting
// and working against Redis
type Redis struct {
	options *Options
	client  *redis.Client
}

// Options returns the current options passed to Redis
func (r Redis) Options() *Options {
	return r.options
}

// Connect connects to Redis database
func (r Redis) Connect() error {
	logrus.Info("Connecting to Redis at ", r.options.RedisURL, ":", r.options.RedisPort)
	pong, err := r.client.Ping().Result()
	if err != nil {
		return err
	}
	logrus.Info("Connected to Redis, ping: ", pong)
	return nil
}

// Get gets a value from Redis based on provided key
func (r Redis) Get(key string) string {
	return ""
}

// Put saves a new key/value to Redis
func (r Redis) Put(key string, value string) error {
	return nil
}

// NewRedis returns new Redis client with provided
// URL and port and password through CLI options
func NewRedis(options *Options) *Redis {
	return &Redis{
		options: options,
		client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", options.RedisURL, options.RedisPort),
			Password: options.RedisPassword,
			DB:       0, // use default DB
		}),
	}
}
