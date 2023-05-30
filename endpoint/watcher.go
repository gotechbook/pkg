package endpoint

type Watcher interface {
	Next() ([]*Instance, error)
	Stop() error
}
