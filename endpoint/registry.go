package endpoint

import "context"

type Registrar interface {
	Register(ctx context.Context, service *Instance) error
	Deregister(ctx context.Context, service *Instance) error
}
