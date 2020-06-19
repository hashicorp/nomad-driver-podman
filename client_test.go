/*
Copyright 2019 Thomas Weber

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"os/user"
	"testing"
	"github.com/stretchr/testify/require"
)

func TestClient_GuessSocketPathForRootCgroupV1(t *testing.T) {
	u := user.User { Uid: "0" }
	fs := []string{"nodev	cgroup"}
	p := guessSocketPath(&u, fs)

	require.Equal(t, p, "unix://run/podman/io.podman")
}

func TestClient_GuessSocketPathForRootCgroupV2(t *testing.T) {
	u := user.User { Uid: "0" }
	fs := []string{"nodev	cgroup2"}
	p := guessSocketPath(&u, fs)

	require.Equal(t, p, "unix://run/podman/io.podman")
}

func TestClient_GuessSocketPathForUserCgroupV1(t *testing.T) {
	u := user.User { Uid: "1000" }
	fs := []string{"nodev	cgroup"}
	p := guessSocketPath(&u, fs)

	require.Equal(t, p, "unix://run/podman/io.podman")
}

func TestClient_GuessSocketPathForUserCgroupV2_1(t *testing.T) {
	u := user.User { Uid: "1000" }
	fs := []string{"nodev	cgroup2"}
	p := guessSocketPath(&u, fs)

	require.Equal(t, p, "unix://run/user/1000/podman/io.podman")
}

func TestClient_GuessSocketPathForUserCgroupV2_2(t *testing.T) {
	u := user.User { Uid: "1337" }
	fs := []string{"nodev	cgroup2"}
	p := guessSocketPath(&u, fs)

	require.Equal(t, p, "unix://run/user/1337/podman/io.podman")
}

func TestClient_GetProcFilesystems(t *testing.T) {
	procFilesystems, err := getProcFilesystems()

	require.NoError(t, err)
	require.Greater(t, len(procFilesystems), 0)
}
