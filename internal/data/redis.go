package data

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
)

// Redis holds the client and options for connecting
// and working against Redis
type Redis struct {
	options *Options
	client  *redis.Client
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

// Options returns the current options passed to Redis
func (r *Redis) Options() *Options {
	return r.options
}

// Connect connects to Redis database
func (r *Redis) Connect() error {
	logrus.Info("Connecting to Redis at ", r.options.RedisURL, ":", r.options.RedisPort)
	pong, err := r.client.Ping().Result()
	if err != nil {
		return err
	}
	logrus.Info("Connected to Redis, ping: ", pong)
	return nil
}

// Has checks if the time slot bucket 'key' exists
func (r *Redis) Has(key string) bool {
	return true
}

// Get hits count from Redis based on bucket key and resource name
func (r *Redis) Get(key string, resource string) (int, error) {
	return 0, nil
}

// Create bucket in Redis for specitic time slot and resource name
func (r *Redis) Create(key string, resource string) error {
	return nil
}

// Add new hit for particular 'resource' in particular time bucket 'key'
func (r *Redis) Add(key string, resource string, now time.Time) (int, error) {
	return 0, nil
}

// Delete a specific 'key' time bucket
func (r Redis) Delete(key string) error {
	return nil
}
