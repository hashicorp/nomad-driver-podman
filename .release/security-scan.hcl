# Copyright IBM Corp. 2019, 2026
# SPDX-License-Identifier: MPL-2.0

binary {
  secrets    = true
  go_modules = true
  # TODO: enable osv check once dependencies are updated.
  osv       = false
  oss_index = false
  nvd       = false
}
