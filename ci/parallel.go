package ci

import (
	"os"
	"testing"
)

func Parallel(t *testing.T) {
	if os.Getenv("CI") == "" {
		t.Parallel()
	}
}
