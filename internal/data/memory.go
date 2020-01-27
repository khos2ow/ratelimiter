package data

// InMemory holds the options and for with in memory
// key-value data store
type InMemory struct {
	options *Options
}

// Options returns the current options passed to InMemory
func (i InMemory) Options() *Options {
	return i.options
}

// Connect connects to InMemory database
func (i InMemory) Connect() error {
	return nil
}

// Get gets a value from InMemory based on provided key
func (i InMemory) Get(key string) string {
	return ""
}

// Put saves a new key/value to InMemory
func (i InMemory) Put(key string, value string) error {
	return nil
}

// NewInMemory returns new InMemory client
func NewInMemory(options *Options) *InMemory {
	return &InMemory{
		options: options,
	}
}
