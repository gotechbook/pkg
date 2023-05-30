package utils

import (
	"github.com/pkg/errors"
	"io"
	"os"
)

func ReadFile(f string) ([]byte, error) {
	kf, err := os.Open(f)
	if err != nil {
		return nil, errors.Wrapf(err, "read %s error", f)
	}

	return io.ReadAll(kf)
}
