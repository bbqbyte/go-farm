package iputil

import (
	"testing"
)

func TestSingleIP(t *testing.T) {
	f := New(Options{
		AllowedIPs:     []string{"222.25.118.1"},
		DefaultBlocked: true,
	})

	assertEqual(t, f.Allowed("222.25.118.1"), true)
	assertEqual(t, f.Blocked("222.25.118.2"), true)
}
