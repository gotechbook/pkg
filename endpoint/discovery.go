package endpoint

import "context"

type Discovery interface {
	GetService(ctx context.Context, name string) ([]*Instance, error)
	Watch(ctx context.Context, name string) (Watcher, error)
}
