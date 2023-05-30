package logger

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"sync"
)

var _ Logger = (*stdLogger)(nil)

type stdLogger struct {
	log  *log.Logger
	pool *sync.Pool
}

func NewStdLogger(w io.Writer) Logger {
	return &stdLogger{log: log.New(w, "", 0), pool: &sync.Pool{New: func() interface{} { return new(bytes.Buffer) }}}
}

func (s *stdLogger) Log(level Level, keyValue ...interface{}) error {
	if len(keyValue) == 0 {
		return nil
	}
	if (len(keyValue) & 1) == 1 {
		keyValue = append(keyValue, "KEY VALUE UNPAIRED")
	}
	buf := s.pool.Get().(*bytes.Buffer)

	buf.WriteString(fmt.Sprintf("[%-5s]", level.String()))

	for i := 0; i < len(keyValue); i += 2 {
		_, _ = fmt.Fprintf(buf, " [%s=%v]", keyValue[i], keyValue[i+1])
	}
	_ = s.log.Output(4, buf.String())
	buf.Reset()
	s.pool.Put(buf)
	return nil
}

func (s *stdLogger) Close() error {
	return nil
}
