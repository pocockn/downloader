package store

// Store handles fetching and storing data.
type Store interface {
	Set(key string, value []byte) error
	Get(key string) ([]byte, error)
	GetAll() ([][]byte, error)
	Disconnect() error
}
