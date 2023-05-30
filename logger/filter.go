package logger

// FilterOption is filter option.
type FilterOption func(*Filter)

const fuzzyStr = "***"

type Filter struct {
	logger Logger
	level  Level
	key    map[interface{}]struct{}
	value  map[interface{}]struct{}
	filter func(level Level, keyValue ...interface{}) bool
}

func NewFilter(logger Logger, opts ...FilterOption) *Filter {
	options := Filter{
		logger: logger,
		key:    make(map[interface{}]struct{}),
		value:  make(map[interface{}]struct{}),
	}
	for _, o := range opts {
		o(&options)
	}
	return &options
}

func (f *Filter) Log(level Level, keyValue ...interface{}) error {
	if level < f.level {
		return nil
	}
	var prefixKv []interface{} // contains the slice of arguments defined as prefixes during the log initialization
	l, ok := f.logger.(*logger)
	if ok && len(l.prefix) > 0 {
		prefixKv = make([]interface{}, 0, len(l.prefix))
		prefixKv = append(prefixKv, l.prefix...)
	}

	if f.filter != nil && (f.filter(level, prefixKv...) || f.filter(level, keyValue...)) {
		return nil
	}

	if len(f.key) > 0 || len(f.value) > 0 {
		for i := 0; i < len(keyValue); i += 2 {
			v := i + 1
			if v >= len(keyValue) {
				continue
			}
			if _, ok := f.key[keyValue[i]]; ok {
				keyValue[v] = fuzzyStr
			}
			if _, ok := f.value[keyValue[v]]; ok {
				keyValue[v] = fuzzyStr
			}
		}
	}
	return f.logger.Log(level, keyValue...)
}

// FilterLevel with filter level.
func FilterLevel(level Level) FilterOption {
	return func(opts *Filter) {
		opts.level = level
	}
}

// FilterKey with filter key.
func FilterKey(key ...string) FilterOption {
	return func(o *Filter) {
		for _, v := range key {
			o.key[v] = struct{}{}
		}
	}
}

// FilterValue with filter value.
func FilterValue(value ...string) FilterOption {
	return func(o *Filter) {
		for _, v := range value {
			o.value[v] = struct{}{}
		}
	}
}

// FilterFunc with filter func.
func FilterFunc(f func(level Level, keyValue ...interface{}) bool) FilterOption {
	return func(o *Filter) {
		o.filter = f
	}
}
