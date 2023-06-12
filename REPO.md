# REPO

Some notes on repository maintenence specific to `nomad-driver-podman`.

### Update Go version

Simply update `.go-version` - all other scripts should reference this file to
determine the Go build version. Update go.mod if it's a new major version.

### Update CHANGELOG.md

When adding a new feature or fixing a bug, update the `CHANGELOG.md` file to
make a note of the change. The description should provide a link to the primary
PR that includes the modification.
