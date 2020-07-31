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
	"os/exec"
	"os/user"

	"path/filepath"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad-driver-podman/iopodman"
	"github.com/hashicorp/nomad/client/taskenv"
	"github.com/hashicorp/nomad/helper/freeport"
	"github.com/hashicorp/nomad/helper/testlog"
	"github.com/hashicorp/nomad/helper/uuid"
	"github.com/hashicorp/nomad/nomad/structs"
	"github.com/hashicorp/nomad/plugins/drivers"
	dtestutil "github.com/hashicorp/nomad/plugins/drivers/testutils"
	tu "github.com/hashicorp/nomad/testutil"
	"github.com/stretchr/testify/require"
	"github.com/varlink/go/varlink"
)

var (
	// busyboxLongRunningCmd is a busybox command that runs indefinitely, and
	// ideally responds to SIGINT/SIGTERM.  Sadly, busybox:1.29.3 /bin/sleep doesn't.
	busyboxLongRunningCmd = []string{"nc", "-l", "-p", "3000", "127.0.0.1"}
	varlinkSocketPath     = ""
)

func init() {
	user, _ := user.Current()
	procFilesystems, err := getProcFilesystems()

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	socketPath := guessSocketPath(user, procFilesystems)

	varlinkSocketPath = socketPath
}

func createBasicResources() *drivers.Resources {
	res := drivers.Resources{
		NomadResources: &structs.AllocatedTaskResources{
			Memory: structs.AllocatedMemoryResources{
				MemoryMB: 100,
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

	d := NewPodmanDriver(testlog.HCLogger(t)).(*Driver)
	d.podmanClient = newPodmanClient()
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
	case <-time.After(time.Duration(tu.TestMultiplier()*2) * time.Second):
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

	// Create a container of the same name but don't start it. This mimics
	// the case of dockerd getting restarted and stopping containers while
	// Nomad is watching them.
	containerName := BuildContainerName(task)
	createOpts := iopodman.Create{
		Args: []string{
			taskCfg.Image,
			"sleep",
			"5",
		},
		Name: &containerName,
	}

	podman := newPodmanClient()
	_, err := podman.CreateContainer(createOpts)
	require.NoError(t, err)

	_, _, err = d.StartTask(task)
	defer d.DestroyTask(task.ID, true)
	require.NoError(t, err)

	require.NoError(t, d.WaitUntilStarted(task.ID, 5*time.Second))
	require.NoError(t, d.DestroyTask(task.ID, true))
}

func TestPodmanDriver_Start_Wait_AllocDir(t *testing.T) {
	if !tu.IsCI() {
		t.Parallel()
	}

	exp := []byte{'w', 'i', 'n'}
	file := "output.txt"
	allocDir := "/mnt/alloc"

	taskCfg := newTaskConfig("", []string{
		"sh",
		"-c",
		fmt.Sprintf(`echo -n %s > $%s/%s; sleep 1`,
			string(exp), taskenv.AllocDir, file),
	})
	taskCfg.Volumes = []string{fmt.Sprintf("alloc/:%s", allocDir)}
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
	case <-time.After(time.Duration(tu.TestMultiplier()*2) * time.Second):
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
	case <-time.After(time.Duration(tu.TestMultiplier()*2) * time.Second):
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

	// Attempt to wait
	waitCh, err := d.WaitTask(context.Background(), task.ID)
	require.NoError(t, err)

	select {
	case <-waitCh:
	case <-time.After(time.Duration(tu.TestMultiplier()*2) * time.Second):
		t.Fatalf("Container did not exit in time")
	}

	tasklog := readLogfile(t, task)
	require.Contains(t, tasklog, check)

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

	// Attempt to wait
	waitCh, err := d.WaitTask(context.Background(), task.ID)
	require.NoError(t, err)

	select {
	case <-waitCh:
	case <-time.After(time.Duration(tu.TestMultiplier()*2) * time.Second):
		t.Fatalf("Container did not exit in time")
	}

	// check if the hostname was visible in the container
	tasklog := readLogfile(t, task)
	require.Contains(t, tasklog, shouldHaveHostname)

}

// check port_map feature
func TestPodmanDriver_PortMap(t *testing.T) {
	if !tu.IsCI() {
		t.Parallel()
	}
	ports := freeport.MustTake(2)
	defer freeport.Return(ports)

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

	inspectData := inspectContainer(t, containerName)

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
	initPath := "/usr/libexec/podman/catatonit"
	_, err := os.Stat(initPath)
	if os.IsNotExist(err) {
		path, err := exec.LookPath("catatonit")
		if err != nil {
			t.Skip("Skipping --init test because catatonit is not installed")
			return
		}
		initPath = path
	}

	taskCfg := newTaskConfig("", []string{
		// print pid 1 filename to stdout
		"realpath",
		"/proc/1/exe",
	})
	// enable --init
	taskCfg.Init = true
	taskCfg.InitPath = initPath

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

	// Attempt to wait
	waitCh, err := d.WaitTask(context.Background(), task.ID)
	require.NoError(t, err)

	select {
	case <-waitCh:
	case <-time.After(time.Duration(tu.TestMultiplier()*2) * time.Second):
		t.Fatalf("Container did not exit in time")
	}

	// podman maps init process to /dev/init, so we should see this
	tasklog := readLogfile(t, task)
	require.Contains(t, tasklog, "/dev/init")

}

// test oom flag propagation
func TestPodmanDriver_OOM(t *testing.T) {
	if !tu.IsCI() {
		t.Parallel()
	}

	taskCfg := newTaskConfig("", []string{
		// Incrementally creates a bigger and bigger variable.
		"sh",
		"-c",
		"tail /dev/zero",
	})

	// only enable init if catatonit is installed
	initPath := "/usr/libexec/podman/catatonit"
	_, err := os.Stat(initPath)
	if os.IsNotExist(err) {
		path, err := exec.LookPath("catatonit")
		if err != nil {
			t.Skip("Skipping oom test because catatonit is not installed")
			return
		}
		initPath = path
	}
	taskCfg.Init = true
	taskCfg.InitPath = initPath

	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "oom",
		AllocID:   uuid.Generate(),
		Resources: createBasicResources(),
	}
	// limit memory to 10MB to trigger oom soon enough
	task.Resources.NomadResources.Memory.MemoryMB = 10
	require.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	_, _, err = d.StartTask(task)
	require.NoError(t, err)

	defer d.DestroyTask(task.ID, true)

	// Attempt to wait
	waitCh, err := d.WaitTask(context.Background(), task.ID)
	require.NoError(t, err)

	select {
	case res := <-waitCh:
		require.False(t, res.Successful(), "Should have failed because of oom but was successful")
		require.True(t, res.OOMKilled, "OOM Flag not set")
		require.Contains(t, res.Err.Error(), "OOM killer")
	case <-time.After(time.Duration(tu.TestMultiplier()*2) * time.Second):
		t.Fatalf("Container did not exit in time")
	}
}

// check setting a user for the task
func TestPodmanDriver_User(t *testing.T) {
	if !tu.IsCI() {
		t.Parallel()
	}

	taskCfg := newTaskConfig("", []string{
		// print our username to stdout
		"whoami",
	})

	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "user",
		AllocID:   uuid.Generate(),
		Resources: createBasicResources(),
	}
	// use "www-data" as a user for our test, it's part of the busybox image
	task.User = "www-data"
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
		// should have a exitcode=0 result
		require.True(t, res.Successful())
	case <-time.After(time.Duration(tu.TestMultiplier()*2) * time.Second):
		t.Fatalf("Container did not exit in time")
	}

	// see if stdout was populated with the "whoami" output
	tasklog := readLogfile(t, task)
	require.Contains(t, tasklog, "www-data")

}

// test memory/swap options
func TestPodmanDriver_Swap(t *testing.T) {
	if !tu.IsCI() {
		t.Parallel()
	}

	taskCfg := newTaskConfig("", busyboxLongRunningCmd)
	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "swap",
		AllocID:   uuid.Generate(),
		Resources: createBasicResources(),
	}
	// limit memory to 50MB
	task.Resources.NomadResources.Memory.MemoryMB = 50
	// but reserve 40MB
	taskCfg.MemoryReservation = "40m"
	// and allow mem+swap of 100MB (= 50 MB Swap)
	taskCfg.MemorySwap = "100m"
	// set a swappiness of 60
	taskCfg.MemorySwappiness = 60
	require.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	containerName := BuildContainerName(task)

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
	case <-time.After(time.Duration(tu.TestMultiplier()*2) * time.Second):
	}
	// inspect container to learn about the actual podman limits
	inspectData := inspectContainer(t, containerName)

	// see if the configured values are set correctly
	require.Equal(t, int64(52428800), inspectData.HostConfig.Memory)
	require.Equal(t, int64(41943040), inspectData.HostConfig.MemoryReservation)
	require.Equal(t, int64(104857600), inspectData.HostConfig.MemorySwap)

	procFilesystems, err := getProcFilesystems()
	if err == nil {
		cgroupv2 := isCGroupV2(procFilesystems)
		if cgroupv2 == false {
			require.Equal(t, int64(60), inspectData.HostConfig.MemorySwappiness)
		}
	}
}

// check tmpfs mounts
func TestPodmanDriver_Tmpfs(t *testing.T) {
	if !tu.IsCI() {
		t.Parallel()
	}

	taskCfg := newTaskConfig("", []string{
		// print our username to stdout
		"sh",
		"-c",
		"mount|grep tmpfs",
	})
	taskCfg.Tmpfs = []string{
		"/tmpdata1",
		"/tmpdata2",
	}

	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "tmpfs",
		AllocID:   uuid.Generate(),
		Resources: createBasicResources(),
	}
	// use "www-data" as a user for our test, it's part of the busybox image
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
	case <-time.After(time.Duration(tu.TestMultiplier()*2) * time.Second):
		t.Fatalf("Container did not exit in time")
	}

	// see if tmpfs was propagated to podman
	inspectData := inspectContainer(t, containerName)
	expectedFilesystem := map[string]string{
		"/tmpdata1": "rw,rprivate,nosuid,nodev,tmpcopyup",
		"/tmpdata2": "rw,rprivate,nosuid,nodev,tmpcopyup",
	}
	require.Exactly(t, expectedFilesystem, inspectData.HostConfig.Tmpfs)

	// see if stdout was populated with expected "mount" output
	tasklog := readLogfile(t, task)
	require.Contains(t, tasklog, " tmpfs on /tmpdata1 type tmpfs ")
	require.Contains(t, tasklog, " tmpfs on /tmpdata2 type tmpfs ")
}

// check default capabilities
func TestPodmanDriver_DefaultCaps(t *testing.T) {
	if !tu.IsCI() {
		t.Parallel()
	}

	taskCfg := newTaskConfig("", busyboxLongRunningCmd)
	inspectData := startDestroyInspect(t, taskCfg, nil, "defaultcaps")

	// a default container should not have SYS_TIME
	require.NotContains(t, inspectData.EffectiveCaps, "CAP_SYS_TIME")
	// a default container gets MKNOD cap
	require.Contains(t, inspectData.EffectiveCaps, "CAP_MKNOD")
}

// check modified capabilities (CapAdd/CapDrop)
func TestPodmanDriver_Caps(t *testing.T) {
	if !tu.IsCI() {
		t.Parallel()
	}

	taskCfg := newTaskConfig("", busyboxLongRunningCmd)
	// 	cap_add = [
	//     "SYS_TIME",
	//   ]
	taskCfg.CapAdd = []string{"SYS_TIME"}
	// 	cap_drop = [
	//     "MKNOD",
	//   ]
	taskCfg.CapDrop = []string{"MKNOD"}

	inspectData := startDestroyInspect(t, taskCfg, nil, "caps")

	// we added SYS_TIME, so we should see it in inspect
	require.Contains(t, inspectData.EffectiveCaps, "CAP_SYS_TIME")
	// we dropped CAP_MKNOD, so we should NOT see it in inspect
	require.NotContains(t, inspectData.EffectiveCaps, "CAP_MKNOD")
}

// check dns server configuration
func TestPodmanDriver_Dns(t *testing.T) {
	if !tu.IsCI() {
		t.Parallel()
	}

	taskCfg := newTaskConfig("", []string{
		"cat",
		"/etc/resolv.conf",
	})
	// config {
	//   dns = [
	//     "1.1.1.1"
	//   ]
	// }
	taskCfg.Dns = []string{"1.1.1.1"}

	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "dns",
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
		// should have a exitcode=0 result
		require.True(t, res.Successful())
	case <-time.After(time.Duration(tu.TestMultiplier()*2) * time.Second):
		t.Fatalf("Container did not exit in time")
	}

	// see if stdout was populated with the correct output
	tasklog := readLogfile(t, task)
	require.Contains(t, tasklog, "nameserver 1.1.1.1")

}

// TestPodmanDriver_NetworkMode asserts we can specify different network modes
// Default podman cni subnet 10.88.0.0/16
func TestPodmanDriver_NetworkMode(t *testing.T) {
	if !tu.IsCI() {
		t.Parallel()
	}

	testCases := []struct {
		mode     string
		expected string
		gateway  string
	}{
		{
			mode:    "host",
			gateway: "",
		},
		{
			// https://github.com/containers/libpod/issues/6618
			// bridge mode info is not fully populated for inspect command
			// so we can only make certain assertions
			mode:    "bridge",
			gateway: "10.88.0.1",
		},
		{
			// slirp4netns information not populated by podman and
			// is not supported for root containers
			// https://github.com/containers/libpod/issues/6097
			mode: "slirp4netns",
		},
		{
			mode:    "none",
			gateway: "",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s_mode_%s", t.Name(), tc.mode), func(t *testing.T) {

			taskCfg := newTaskConfig("", busyboxLongRunningCmd)
			taskCfg.NetworkMode = tc.mode

			task := &drivers.TaskConfig{
				ID:        uuid.Generate(),
				Name:      fmt.Sprintf("network_mode_%s", tc.mode),
				AllocID:   uuid.Generate(),
				Resources: createBasicResources(),
			}

			require.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

			d := podmanDriverHarness(t, nil)
			defer d.MkAllocDir(task, true)()

			containerName := BuildContainerName(task)
			_, _, err := d.StartTask(task)
			require.NoError(t, err)

			defer d.DestroyTask(task.ID, true)

			require.NoError(t, d.WaitUntilStarted(task.ID, time.Duration(tu.TestMultiplier()*3)*time.Second))

			inspectData := inspectContainer(t, containerName)
			if tc.mode == "host" {
				require.Equal(t, "host", inspectData.HostConfig.NetworkMode)
			}
			require.Equal(t, tc.gateway, inspectData.NetworkSettings.Gateway)
		})
	}
}

// TestTestPodmanDriver_IsolatedEnv ensures that env vars do not leak
// into other tasks
func TestPodmanDriver_IsolatedEnv(t *testing.T) {
	if !tu.IsCI() {
		t.Parallel()
	}

	taskCfg1 := newTaskConfig("", busyboxLongRunningCmd)
	cb := func(d *drivers.TaskConfig) {
		if d.Env == nil {
			d.Env = make(map[string]string)
		}
		d.Env["foo"] = "bar"
	}
	inspectData := startDestroyInspect(t, taskCfg1, cb, "env1")
	require.NotNil(t, inspectData)
	require.Contains(t, inspectData.Config.Env, "foo=bar")

	taskCfg2 := newTaskConfig("", busyboxLongRunningCmd)
	cb2 := func(d *drivers.TaskConfig) {
		if d.Env == nil {
			d.Env = make(map[string]string)
		}
		d.Env["env2"] = "bar2"
	}
	inspectData2 := startDestroyInspect(t, taskCfg2, cb2, "env2")

	require.NotContains(t, inspectData2.Config.Env, "foo=bar")
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
	busyboxImageID := "docker://docker.io/library/busybox:latest"

	image := busyboxImageID

	return TaskConfig{
		Image: image,
		// LoadImage: loadImage,
		Command: command[0],
		Args:    command[1:],
	}
}

func inspectContainer(t *testing.T, containerName string) iopodman.InspectContainerData {
	ctx := context.Background()
	varlinkConnection, err := getPodmanConnection(ctx)
	require.NoError(t, err)

	defer varlinkConnection.Close()

	inspectJSON, err := iopodman.InspectContainer().Call(ctx, varlinkConnection, containerName)
	require.NoError(t, err)

	var inspectData iopodman.InspectContainerData
	require.NoError(t, json.Unmarshal([]byte(inspectJSON), &inspectData))

	return inspectData
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
	varlinkConnection, err := varlink.NewConnection(ctx, varlinkSocketPath)
	return varlinkConnection, err
}

func newPodmanClient() *PodmanClient {
	testLogger := hclog.New(&hclog.LoggerOptions{
		Name:  "testClient",
		Level: hclog.LevelFromString("DEBUG"),
	})
	client := &PodmanClient{
		ctx:               context.Background(),
		logger:            testLogger,
		varlinkSocketPath: varlinkSocketPath,
	}
	return client
}

func getPodmanDriver(t *testing.T, harness *dtestutil.DriverHarness) *Driver {
	driver, ok := harness.Impl().(*Driver)
	require.True(t, ok)
	return driver
}

// helper to start, destroy and inspect a long running container
func startDestroyInspect(t *testing.T, taskCfg TaskConfig, driverCB func(d *drivers.TaskConfig), taskName string) iopodman.InspectContainerData {
	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      taskName,
		AllocID:   uuid.Generate(),
		Resources: createBasicResources(),
	}

	if driverCB != nil {
		driverCB(task)
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
	case <-time.After(time.Duration(tu.TestMultiplier()*2) * time.Second):
	}
	// inspect container to learn about the actual podman limits
	inspectData := inspectContainer(t, containerName)
	return inspectData
}
