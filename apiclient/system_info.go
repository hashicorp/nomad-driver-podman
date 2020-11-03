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

package apiclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// SystemInfo returns information on the system and libpod configuration
func (c *APIClient) SystemInfo(ctx context.Context) (Info, error) {

	var infoData Info

	// the libpod/info endpoint seems to have some trouble
	// using "compat" endpoint and minimal struct
	// until podman returns proper data.
	res, err := c.Get(ctx, fmt.Sprintf("/%s/libpod/info", PODMAN_API_VERSION))
	if err != nil {
		return infoData, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return infoData, fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return infoData, err
	}
	err = json.Unmarshal(body, &infoData)
	if err != nil {
		return infoData, err
	}

	return infoData, nil
}

// -------------------------------------------------------------------------------------------------------
// structs copied from https://github.com/containers/podman/blob/master/libpod/define/info.go
//
// some unused parts are modified/commented out to not pull more dependencies
//
// some fields are reordert to make the linter happy (bytes maligned complains)
// -------------------------------------------------------------------------------------------------------

// Info is the overall struct that describes the host system
// running libpod/podman
type Info struct {
	Host       *HostInfo              `json:"host"`
	Store      *StoreInfo             `json:"store"`
	Registries map[string]interface{} `json:"registries"`
	Version    Version                `json:"version"`
}

//HostInfo describes the libpod host
type HostInfo struct {
	Arch           string           `json:"arch"`
	BuildahVersion string           `json:"buildahVersion"`
	CgroupManager  string           `json:"cgroupManager"`
	CGroupsVersion string           `json:"cgroupVersion"`
	Conmon         *ConmonInfo      `json:"conmon"`
	CPUs           int              `json:"cpus"`
	Distribution   DistributionInfo `json:"distribution"`
	EventLogger    string           `json:"eventLogger"`
	Hostname       string           `json:"hostname"`
	// IDMappings     IDMappings             `json:"idMappings,omitempty"`
	Kernel       string                 `json:"kernel"`
	MemFree      int64                  `json:"memFree"`
	MemTotal     int64                  `json:"memTotal"`
	OCIRuntime   *OCIRuntimeInfo        `json:"ociRuntime"`
	OS           string                 `json:"os"`
	RemoteSocket *RemoteSocket          `json:"remoteSocket,omitempty"`
	Rootless     bool                   `json:"rootless"`
	RuntimeInfo  map[string]interface{} `json:"runtimeInfo,omitempty"`
	Slirp4NetNS  SlirpInfo              `json:"slirp4netns,omitempty"`
	SwapFree     int64                  `json:"swapFree"`
	SwapTotal    int64                  `json:"swapTotal"`
	Uptime       string                 `json:"uptime"`
	Linkmode     string                 `json:"linkmode"`
}

// RemoteSocket describes information about the API socket
type RemoteSocket struct {
	Path   string `json:"path,omitempty"`
	Exists bool   `json:"exists,omitempty"`
}

// SlirpInfo describes the slirp executable that
// is being being used.
type SlirpInfo struct {
	Executable string `json:"executable"`
	Package    string `json:"package"`
	Version    string `json:"version"`
}

// IDMappings describe the GID and UID mappings
// type IDMappings struct {
// 	GIDMap []idtools.IDMap `json:"gidmap"`
// 	UIDMap []idtools.IDMap `json:"uidmap"`
// }

// DistributionInfo describes the host distribution
// for libpod
type DistributionInfo struct {
	Distribution string `json:"distribution"`
	Version      string `json:"version"`
}

// ConmonInfo describes the conmon executable being used
type ConmonInfo struct {
	Package string `json:"package"`
	Path    string `json:"path"`
	Version string `json:"version"`
}

// OCIRuntimeInfo describes the runtime (crun or runc) being
// used with podman
type OCIRuntimeInfo struct {
	Name    string `json:"name"`
	Package string `json:"package"`
	Path    string `json:"path"`
	Version string `json:"version"`
}

// StoreInfo describes the container storage and its
// attributes
type StoreInfo struct {
	ConfigFile      string                 `json:"configFile"`
	ContainerStore  ContainerStore         `json:"containerStore"`
	GraphDriverName string                 `json:"graphDriverName"`
	GraphOptions    map[string]interface{} `json:"graphOptions"`
	GraphRoot       string                 `json:"graphRoot"`
	GraphStatus     map[string]string      `json:"graphStatus"`
	ImageStore      ImageStore             `json:"imageStore"`
	RunRoot         string                 `json:"runRoot"`
	VolumePath      string                 `json:"volumePath"`
}

// ImageStore describes the image store.  Right now only the number
// of images present
type ImageStore struct {
	Number int `json:"number"`
}

// ContainerStore describes the quantity of containers in the
// store by status
type ContainerStore struct {
	Number  int `json:"number"`
	Paused  int `json:"paused"`
	Running int `json:"running"`
	Stopped int `json:"stopped"`
}

// Version is an output struct for API
type Version struct {
	APIVersion string
	Version    string
	GoVersion  string
	GitCommit  string
	BuiltTime  string
	Built      int64
	OsArch     string
}
