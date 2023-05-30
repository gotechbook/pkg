package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gotechbook/pkg/endpoint"
	"math/rand"
	"time"
)

var _ endpoint.Registrar = (*Registry)(nil)

type Registry struct {
	opts   *options
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

func New(c *clientv3.Client, opts ...Option) *Registry {
	op := &options{
		ctx:       context.Background(),
		namespace: "/microservices",
		ttl:       15 * time.Second,
		maxRetry:  5,
	}

	for _, o := range opts {
		o(op)
	}

	return &Registry{
		opts:   op,
		client: c,
		kv:     clientv3.NewKV(c),
	}
}

func (r *Registry) Register(ctx context.Context, service *endpoint.Instance) error {
	key := fmt.Sprintf("%s/%s/%s", r.opts.namespace, service.Name, service.ID)
	value, err := json.Marshal(service)
	if err != nil {
		return err
	}
	if r.lease != nil {
		r.lease.Close()
	}
	r.lease = clientv3.NewLease(r.client)
	leaseID, err := r.registerWithKV(ctx, key, string(value))
	if err != nil {
		return err
	}
	go r.heartBeat(r.opts.ctx, leaseID, key, string(value))
	return nil
}

func (r *Registry) Deregister(ctx context.Context, service *endpoint.Instance) error {
	defer func() {
		if r.lease != nil {
			r.lease.Close()
		}
	}()
	key := fmt.Sprintf("%s/%s/%s", r.opts.namespace, service.Name, service.ID)
	_, err := r.client.Delete(ctx, key)
	return err
}

func (r *Registry) GetService(ctx context.Context, name string) ([]*endpoint.Instance, error) {
	key := fmt.Sprintf("%s/%s", r.opts.namespace, name)
	resp, err := r.kv.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	items := make([]*endpoint.Instance, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		si, err := unmarshal(kv.Value)
		if err != nil {
			return nil, err
		}
		if si.Name != name {
			continue
		}
		items = append(items, si)
	}
	return items, nil
}

func (r *Registry) Watch(ctx context.Context, name string) (endpoint.Watcher, error) {
	key := fmt.Sprintf("%s/%s", r.opts.namespace, name)
	return newWatcher(ctx, key, name, r.client)
}

func (r *Registry) registerWithKV(ctx context.Context, k, v string) (clientv3.LeaseID, error) {
	grant, err := r.lease.Grant(ctx, int64(r.opts.ttl.Seconds()))
	if err != nil {
		return 0, err
	}
	_, err = r.client.Put(ctx, k, v, clientv3.WithLease(grant.ID))
	if err != nil {
		return 0, err
	}
	return grant.ID, nil
}

func (r *Registry) heartBeat(ctx context.Context, leaseID clientv3.LeaseID, key string, value string) {
	curLeaseID := leaseID
	kac, err := r.client.KeepAlive(ctx, leaseID)
	if err != nil {
		curLeaseID = 0
	}
	rand.Seed(time.Now().Unix())

	for {
		if curLeaseID == 0 {
			// try to registerWithKV
			var retreat []int
			for retryCnt := 0; retryCnt < r.opts.maxRetry; retryCnt++ {
				if ctx.Err() != nil {
					return
				}
				// prevent infinite blocking
				idChan := make(chan clientv3.LeaseID, 1)
				errChan := make(chan error, 1)
				cancelCtx, cancel := context.WithCancel(ctx)
				go func() {
					defer cancel()
					id, registerErr := r.registerWithKV(cancelCtx, key, value)
					if registerErr != nil {
						errChan <- registerErr
					} else {
						idChan <- id
					}
				}()

				select {
				case <-time.After(3 * time.Second):
					cancel()
					continue
				case <-errChan:
					continue
				case curLeaseID = <-idChan:
				}

				kac, err = r.client.KeepAlive(ctx, curLeaseID)
				if err == nil {
					break
				}
				retreat = append(retreat, 1<<retryCnt)
				time.Sleep(time.Duration(retreat[rand.Intn(len(retreat))]) * time.Second)
			}
			if _, ok := <-kac; !ok {
				// retry failed
				return
			}
		}

		select {
		case _, ok := <-kac:
			if !ok {
				if ctx.Err() != nil {
					// channel closed due to context cancel
					return
				}
				// need to retry registration
				curLeaseID = 0
				continue
			}
		case <-r.opts.ctx.Done():
			return
		}
	}
}

func marshal(si *endpoint.Instance) (string, error) {
	data, err := json.Marshal(si)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func unmarshal(data []byte) (si *endpoint.Instance, err error) {
	err = json.Unmarshal(data, &si)
	return
}
