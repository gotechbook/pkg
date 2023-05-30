package eth

import (
	"fmt"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	s := GenerateKey()
	fmt.Println(s)
}
