// Copyright IBM Corp. 2019, 2025
// SPDX-License-Identifier: MPL-2.0

package ci

import (
	"os"
)

func TestSELinux() bool {
	return os.Getenv("CI_SELINUX") == "1"
}
