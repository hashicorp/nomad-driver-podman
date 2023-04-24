package ci

import (
	"os"
)

func TestSELinux() bool {
	return os.Getenv("CI_SELINUX") == "1"
}
