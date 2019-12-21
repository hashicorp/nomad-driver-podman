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
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/consul/lib/freeport"
	"github.com/hashicorp/nomad/client/taskenv"
	ctestutil "github.com/hashicorp/nomad/client/testutil"
	"github.com/hashicorp/nomad/helper/testlog"
	"github.com/hashicorp/nomad/helper/uuid"
	"github.com/hashicorp/nomad/nomad/structs"
	"github.com/hashicorp/nomad/plugins/drivers"
	dtestutil "github.com/hashicorp/nomad/plugins/drivers/testutils"
	tu "github.com/hashicorp/nomad/testutil"
	"github.com/pascomnet/nomad-driver-podman/iopodman"
	"github.com/stretchr/testify/require"
	"github.com/varlink/go/varlink"
)

var (
	// busyboxLongRunningCmd is a busybox command that runs indefinitely, and
	// ideally responds to SIGINT/SIGTERM.  Sadly, busybox:1.29.3 /bin/sleep doesn't.
	busyboxLongRunningCmd = []string{"nc", "-l", "-p", "3000", "127.0.0.1"}
)

func createBasicResources() *drivers.Resources {
	res := drivers.Resources{
		NomadResources: &structs.AllocatedTaskResources{
			Memory: structs.AllocatedMemoryResources{
				// MemoryMB: 256,
			},
			Cpu: structs.AllocatedCpuResources{
				CpuShares: 250,
			},
		},
		LinuxResources: &drivers.LinuxResources{
			CPUShares:        512,
			MemoryLimitBytes: 256 * 1024 * 1024,
		},
	}
	return &res
}

// podmanDriverHarness wires up everything needed to launch a task with a podman driver.
// A driver plugin interface and cleanup function is returned
func podmanDriverHarness(t *testing.T, cfg map[string]interface{}) *dtestutil.DriverHarness {

	ctestutil.RequireRoot(t)

	d := NewPodmanDriver(testlog.HCLogger(t)).(*Driver)
	d.config.Volumes.Enabled = true
	if enforce, err := ioutil.ReadFile("/sys/fs/selinux/enforce"); err == nil {
		if string(enforce) == "1" {
			d.logger.Info("Enabling SelinuxLabel")
			d.config.Volumes.SelinuxLabel = "z"
		}
	}
	d.config.GC.Container = true
	if v, ok := cfg["GC.Container"]; ok {
		if bv, ok := v.(bool); ok {
			d.config.GC.Container = bv
		}
	}

	harness := dtestutil.NewDriverHarness(t, d)

	return harness
}

func TestPodmanDriver_Start_NoImage(t *testing.T) {
	if !tu.IsCI() {
		t.Parallel()
	}

	taskCfg := TaskConfig{
		Command: "echo",
		Args:    []string{"foo"},
	}
	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "echo",
		AllocID:   uuid.Generate(),
		Resources: createBasicResources(),
	}
	require.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, false)
	defer cleanup()

	_, _, err := d.StartTask(task)
	require.Error(t, err)
	require.Contains(t, err.Error(), "image name required")

	d.DestroyTask(task.ID, true)

}

// start a long running container
func TestPodmanDriver_Start_Wait(t *testing.T) {
	if !tu.IsCI() {
		t.Parallel()
	}

	taskCfg := newTaskConfig("", busyboxLongRunningCmd)
	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "start_wait",
		AllocID:   uuid.Generate(),
		Resources: createBasicResources(),
	}
	require.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	_, _, err := d.StartTask(task)
	require.NoError(t, err)

	defer d.DestroyTask(task.ID, true)

	// Attempt to wait
	waitCh, err := d.WaitTask(context.Background(), task.ID)
	require.NoError(t, err)

	select {
	case <-waitCh:
		t.Fatalf("wait channel should not have received an exit result")
	case <-time.After(time.Duration(tu.TestMultiplier()*1) * time.Second):
	}
}

// test a short-living container
func TestPodmanDriver_Start_WaitFinish(t *testing.T) {
	if !tu.IsCI() {
		t.Parallel()
	}

	taskCfg := newTaskConfig("", []string{"echo", "hello"})
	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "start_waitfinish",
		AllocID:   uuid.Generate(),
		Resources: createBasicResources(),
	}
	require.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	_, _, err := d.StartTask(task)
	require.NoError(t, err)

	defer d.DestroyTask(task.ID, true)

	// Attempt to wait
	waitCh, err := d.WaitTask(context.Background(), task.ID)
	require.NoError(t, err)

	select {
	case res := <-waitCh:
		if !res.Successful() {
			require.Fail(t, "ExitResult should be successful: %v", res)
		}
	case <-time.After(time.Duration(tu.TestMultiplier()*5) * time.Second):
		require.Fail(t, "timeout")
	}
}

/*
// TestPodmanDriver_Start_StoppedContainer asserts that Nomad will detect a
// stopped task container, remove it, and start a new container.
//
// See https://github.com/hashicorp/nomad/issues/3419
func TestPodmanDriver_Start_StoppedContainer(t *testing.T) {
	if !tu.IsCI() {
		t.Parallel()
	}

	taskCfg := newTaskConfig("", []string{"sleep", "5"})
	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "start_stoppedContainer",
		AllocID:   uuid.Generate(),
		Resources: createBasicResources(),
	}
	require.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	ctx := context.Background()
	varlinkConnection, err := getPodmanConnection(ctx)
	require.NoError(t, err)

	defer varlinkConnection.Close()

	// Create a container of the same name but don't start it. This mimics
	// the case of dockerd getting restarted and stopping containers while
	// Nomad is watching them.
	containerName := task.Name+"-"+task.ID //strings.Replace(task.ID, "/", "_", -1)
	createOpts := iopodman.Create{
		Args:       []string{
			taskCfg.Image,
			"sleep",
			"5",
		},
		Name:       &containerName,
	}

	_, err = iopodman.CreateContainer().Call(ctx, varlinkConnection, createOpts)
	require.NoError(t, err)

	_, _, err = d.StartTask(task)
	defer d.DestroyTask(task.ID, true)
	require.NoError(t, err)

	require.NoError(t, d.WaitUntilStarted(task.ID, 5*time.Second))
	require.NoError(t, d.DestroyTask(task.ID, true))
}
*/

func TestPodmanDriver_Start_Wait_AllocDir(t *testing.T) {
	if !tu.IsCI() {
		t.Parallel()
	}

	exp := []byte{'w', 'i', 'n'}
	file := "output.txt"

	taskCfg := newTaskConfig("", []string{
		"sh",
		"-c",
		fmt.Sprintf(`echo -n %s > $%s/%s; sleep 1`,
			string(exp), taskenv.AllocDir, file),
	})
	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "start_wait_allocDir",
		AllocID:   uuid.Generate(),
		Resources: createBasicResources(),
	}
	require.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	_, _, err := d.StartTask(task)
	require.NoError(t, err)

	defer d.DestroyTask(task.ID, true)

	// Attempt to wait
	waitCh, err := d.WaitTask(context.Background(), task.ID)
	require.NoError(t, err)

	select {
	case res := <-waitCh:
		if !res.Successful() {
			require.Fail(t, fmt.Sprintf("ExitResult should be successful: %v", res))
		}
	case <-time.After(time.Duration(tu.TestMultiplier()*5) * time.Second):
		require.Fail(t, "timeout")
	}

	// Check that data was written to the shared alloc directory.
	outputFile := filepath.Join(task.TaskDir().SharedAllocDir, file)
	act, err := ioutil.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Couldn't read expected output: %v", err)
	}

	if !reflect.DeepEqual(act, exp) {
		t.Fatalf("Command outputted %v; want %v", act, exp)
	}
}

// check if container is destroyed if gc.container=true
func TestPodmanDriver_GC_Container_on(t *testing.T) {
	if !tu.IsCI() {
		t.Parallel()
	}

	taskCfg := newTaskConfig("", busyboxLongRunningCmd)
	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "gc_container_on",
		AllocID:   uuid.Generate(),
		Resources: createBasicResources(),
	}
	require.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	containerName := BuildContainerName(task)

	_, _, err := d.StartTask(task)
	require.NoError(t, err)

	defer d.DestroyTask(task.ID, true)

	// Attempt to wait
	waitCh, err := d.WaitTask(context.Background(), task.ID)
	require.NoError(t, err)

	select {
	case <-waitCh:
		t.Fatalf("wait channel should not have received an exit result")
	case <-time.After(time.Duration(tu.TestMultiplier()*1) * time.Second):
	}

	d.DestroyTask(task.ID, true)

	ctx := context.Background()
	varlinkConnection, err := getPodmanConnection(ctx)
	require.NoError(t, err)

	defer varlinkConnection.Close()

	// ... 1 means container could not be found (is deleted)
	exists, err := iopodman.ContainerExists().Call(ctx, varlinkConnection, containerName)
	require.NoError(t, err)
	require.Equal(t, 1, int(exists))
}

// check if container is destroyed if gc.container=false
func TestPodmanDriver_GC_Container_off(t *testing.T) {
	if !tu.IsCI() {
		t.Parallel()
	}

	taskCfg := newTaskConfig("", busyboxLongRunningCmd)
	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "gc_container_off",
		AllocID:   uuid.Generate(),
		Resources: createBasicResources(),
	}
	require.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	opts := make(map[string]interface{})
	opts["GC.Container"] = false

	d := podmanDriverHarness(t, opts)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	containerName := BuildContainerName(task)

	_, _, err := d.StartTask(task)
	require.NoError(t, err)

	defer d.DestroyTask(task.ID, true)

	// Attempt to wait
	waitCh, err := d.WaitTask(context.Background(), task.ID)
	require.NoError(t, err)

	select {
	case <-waitCh:
		t.Fatalf("wait channel should not have received an exit result")
	case <-time.After(time.Duration(tu.TestMultiplier()*1) * time.Second):
	}

	d.DestroyTask(task.ID, true)

	ctx := context.Background()
	varlinkConnection, err := getPodmanConnection(ctx)
	require.NoError(t, err)

	defer varlinkConnection.Close()

	// ... 0 means container could be found (still exists)
	exists, err := iopodman.ContainerExists().Call(ctx, varlinkConnection, containerName)
	require.NoError(t, err)
	require.Equal(t, 0, int(exists))

	// and cleanup after ourself
	iopodman.RemoveContainer().Call(ctx, varlinkConnection, containerName, true, true)
}

// Check stdout/stderr logging
func TestPodmanDriver_Stdout(t *testing.T) {
	if !tu.IsCI() {
		t.Parallel()
	}

	check := uuid.Generate()

	taskCfg := newTaskConfig("", []string{
		"sh",
		"-c",
		"echo " + check,
	})
	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "stdout",
		AllocID:   uuid.Generate(),
		Resources: createBasicResources(),
	}
	require.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	_, _, err := d.StartTask(task)
	require.NoError(t, err)

	defer d.DestroyTask(task.ID, true)

	tasklog := readLogfile(t, task)
	require.Contains(t, tasklog, check)

	// Attempt to wait
	waitCh, err := d.WaitTask(context.Background(), task.ID)
	require.NoError(t, err)

	select {
	case <-waitCh:
		t.Fatalf("wait channel should not have received an exit result")
	case <-time.After(time.Duration(tu.TestMultiplier()*1) * time.Second):
	}

}

// check hostname task config options
func TestPodmanDriver_Hostname(t *testing.T) {
	if !tu.IsCI() {
		t.Parallel()
	}

	taskCfg := newTaskConfig("", []string{
		// print hostname to stdout
		"hostname",
	})
	shouldHaveHostname := "host_" + uuid.Generate()

	// populate Hostname option in task configuration
	taskCfg.Hostname = shouldHaveHostname

	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "stdout",
		AllocID:   uuid.Generate(),
		Resources: createBasicResources(),
	}
	require.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	_, _, err := d.StartTask(task)
	require.NoError(t, err)

	defer d.DestroyTask(task.ID, true)

	// check if the hostname was visible in the container
	tasklog := readLogfile(t, task)
	require.Contains(t, tasklog, shouldHaveHostname)

	// Attempt to wait
	waitCh, err := d.WaitTask(context.Background(), task.ID)
	require.NoError(t, err)

	select {
	case <-waitCh:
		t.Fatalf("wait channel should not have received an exit result")
	case <-time.After(time.Duration(tu.TestMultiplier()*1) * time.Second):
	}

}

// check port_map feature
func TestPodmanDriver_PortMap(t *testing.T) {
	if !tu.IsCI() {
		t.Parallel()
	}
	ports := freeport.GetT(t, 2)

	taskCfg := newTaskConfig("", busyboxLongRunningCmd)

	taskCfg.PortMap = map[string]int{
		"main":  8888,
		"REDIS": 6379,
	}

	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "portmap",
		AllocID:   uuid.Generate(),
		Resources: createBasicResources(),
	}
	task.Resources.NomadResources.Networks = []*structs.NetworkResource{
		{
			IP:            "127.0.0.1",
			ReservedPorts: []structs.Port{{Label: "main", Value: ports[0]}},
			DynamicPorts:  []structs.Port{{Label: "REDIS", Value: ports[1]}},
		},
	}
	require.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	containerName := BuildContainerName(task)

	_, _, err := d.StartTask(task)
	require.NoError(t, err)

	defer d.DestroyTask(task.ID, true)

	inspectJSON := inspectContainer(t, containerName)
	var inspectData iopodman.InspectContainerData
	require.NoError(t, json.Unmarshal([]byte(inspectJSON), &inspectData))

	// Verify that the port environment variables are set
	require.Contains(t, inspectData.Config.Env, "NOMAD_PORT_main=8888")
	require.Contains(t, inspectData.Config.Env, "NOMAD_PORT_REDIS=6379")

	// Verify that the correct ports are bound
	expectedPortBindings := map[string][]iopodman.InspectHostPort{
		"8888/tcp": []iopodman.InspectHostPort{
			iopodman.InspectHostPort{
				HostIP:   "127.0.0.1",
				HostPort: strconv.Itoa(ports[0]),
			},
		},
		"8888/udp": []iopodman.InspectHostPort{
			iopodman.InspectHostPort{
				HostIP:   "127.0.0.1",
				HostPort: strconv.Itoa(ports[0]),
			},
		},
		"6379/tcp": []iopodman.InspectHostPort{
			iopodman.InspectHostPort{
				HostIP:   "127.0.0.1",
				HostPort: strconv.Itoa(ports[1]),
			},
		},
		"6379/udp": []iopodman.InspectHostPort{
			iopodman.InspectHostPort{
				HostIP:   "127.0.0.1",
				HostPort: strconv.Itoa(ports[1]),
			},
		},
		// FIXME: REDIS UDP
	}

	require.Exactly(t, expectedPortBindings, inspectData.HostConfig.PortBindings)

	// fmt.Printf("Inspect %v",inspectData.HostConfig.PortBindings)

}

// check --init with default path
func TestPodmanDriver_Init(t *testing.T) {
	if !tu.IsCI() {
		t.Parallel()
	}

	// only test --init if catatonit is installed
	_, err := os.Stat("/usr/libexec/podman/catatonit")
	if os.IsNotExist(err) {
		t.Skip("Skipping --init test because catatonit is not installed")
		return
	}

	taskCfg := newTaskConfig("", []string{
		// print pid 1 filename to stdout
		"realpath",
		"/proc/1/exe",
	})
	// enable --init
	taskCfg.Init = true

	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "init",
		AllocID:   uuid.Generate(),
		Resources: createBasicResources(),
	}
	require.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	_, _, err = d.StartTask(task)
	require.NoError(t, err)

	defer d.DestroyTask(task.ID, true)

	// podman maps init process to /dev/init, so we should see this
	tasklog := readLogfile(t, task)
	require.Contains(t, tasklog, "/dev/init")

}

// read a tasks logfile into a string, fail on error
func readLogfile(t *testing.T, task *drivers.TaskConfig) string {
	logfile := filepath.Join(filepath.Dir(task.StdoutPath), fmt.Sprintf("%s.stdout.0", task.Name))
	// Get the stdout of the process and assert that it's not empty
	stdout, err := ioutil.ReadFile(logfile)
	require.NoError(t, err)
	return string(stdout)
}

func newTaskConfig(variant string, command []string) TaskConfig {
	// busyboxImageID is the ID stored in busybox.tar
	busyboxImageID := "docker://busybox"
	// busyboxImageID := "busybox:1.29.3"

	image := busyboxImageID

	return TaskConfig{
		Image: image,
		// LoadImage: loadImage,
		Command: command[0],
		Args:    command[1:],
	}
}

func inspectContainer(t *testing.T, containerName string) string {
	ctx := context.Background()
	varlinkConnection, err := getPodmanConnection(ctx)
	require.NoError(t, err)

	defer varlinkConnection.Close()

	inspectJSON, err := iopodman.InspectContainer().Call(ctx, varlinkConnection, containerName)
	require.NoError(t, err)

	return inspectJSON
}

func getContainer(t *testing.T, containerName string) iopodman.Container {
	ctx := context.Background()
	varlinkConnection, err := getPodmanConnection(ctx)
	require.NoError(t, err)

	defer varlinkConnection.Close()

	container, err := iopodman.GetContainer().Call(ctx, varlinkConnection, containerName)
	require.NoError(t, err)

	return container
}

func getPodmanConnection(ctx context.Context) (*varlink.Connection, error) {
	// FIXME: a parameter for the socket would be nice
	varlinkConnection, err := varlink.NewConnection(ctx, "unix://run/podman/io.podman")
	return varlinkConnection, err
}
