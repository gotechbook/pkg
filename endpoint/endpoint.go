package endpoint

import (
	"fmt"
	"net/url"
)

func NewEndpoint(scheme, host string) *url.URL {
	return &url.URL{Scheme: scheme, Host: host}
}

func ParseEndpoint(endpoints []string, scheme string) (string, error) {
	for _, e := range endpoints {
		u, err := url.Parse(e)
		if err != nil {
			return "", err
		}
		if u.Scheme == scheme {
			return u.Host, nil
		}
	}
	return "", nil
}

func Scheme(scheme string, isSecure bool) string {
	if isSecure {
		return fmt.Sprintf("%ss", scheme)
	}
	return scheme
}
