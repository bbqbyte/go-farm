package pbconfig

import (
	"testing"
)

func TestDir(t *testing.T) {
	t.Log(HomeDir())
}

func TestExpand(t *testing.T) {
	t.Log(Expand("~/.m2"))
}

func TestExpandValueEnv(t *testing.T) {
	t.Log(ExpandValueEnv("$say good ${GOPATHZ||/usr/local/go} haha"))
}