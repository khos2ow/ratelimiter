package data

// Options respresents the available options to be passed
// to data store implementations. These can be set with
// flags through the CLI or environment variable to Docker
// container.
type Options struct {
	RedisURL      string
	RedisPort     int
	RedisPassword string
}
