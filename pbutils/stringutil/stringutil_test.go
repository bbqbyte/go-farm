package stringutil

import (
	"fmt"
	"testing"
)

func TestRandString(t *testing.T) {
	fmt.Printf(RandomString(16))
}