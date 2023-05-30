package middleware

import (
	"sort"
	"strings"
)

type Matcher interface {
	Use(ms ...Middleware)
	Add(selector string, ms ...Middleware)
	Matcher(operation string) []Middleware
}

type matcher struct {
	prefix   []string
	defaults []Middleware
	match    map[string][]Middleware
}

func (m *matcher) Use(ms ...Middleware) {
	m.defaults = ms
}

func (m *matcher) Add(selector string, ms ...Middleware) {
	if strings.HasSuffix(selector, "*") {
		selector = strings.TrimSuffix(selector, "*")
		m.prefix = append(m.prefix, selector)
		// sort the prefix:
		//  - /foo/bar
		//  - /foo
		sort.Slice(m.prefix, func(i, j int) bool {
			return m.prefix[i] > m.prefix[j]
		})
	}
	m.match[selector] = ms
}

func (m *matcher) Matcher(operation string) []Middleware {
	ms := make([]Middleware, 0, len(m.defaults))
	if len(m.defaults) > 0 {
		ms = append(ms, m.defaults...)
	}
	if next, ok := m.match[operation]; ok {
		return append(ms, next...)
	}
	for _, prefix := range m.prefix {
		if strings.HasPrefix(operation, prefix) {
			return append(ms, m.match[prefix]...)
		}
	}
	return ms
}

func NewMatcher() Matcher {
	return &matcher{
		match: make(map[string][]Middleware),
	}
}
