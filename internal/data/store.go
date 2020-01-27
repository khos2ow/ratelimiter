package data

// Store respresents an interface of type of data store
// which gets implemented in different types, e.g. Redis,
// InMemory, etc.
type Store interface {
	Options() *Options
	Connect() error
	Get(key string) string
	Put(key string, value string) error
}
