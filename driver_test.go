package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"path/filepath"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad-driver-podman/api"
	"github.com/hashicorp/nomad/client/taskenv"
	"github.com/hashicorp/nomad/helper/freeport"
	"github.com/hashicorp/nomad/helper/testlog"
	"github.com/hashicorp/nomad/helper/uuid"
	"github.com/hashicorp/nomad/nomad/structs"
	"github.com/hashicorp/nomad/plugins/base"
	"github.com/hashicorp/nomad/plugins/drivers"
	dtestutil "github.com/hashicorp/nomad/plugins/drivers/testutils"
	tu "github.com/hashicorp/nomad/testutil"
	"github.com/stretchr/testify/require"
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

	logger := testlog.HCLogger(t)
	if testing.Verbose() {
		logger.SetLevel(hclog.Trace)
	} else {
		logger.SetLevel(hclog.Info)
	}

	baseConfig := base.Config{}
	pluginConfig := PluginConfig{}
	if err := base.MsgPackEncode(&baseConfig.PluginConfig, &pluginConfig); err != nil {
		t.Error("Unable to encode plugin config", err)
	}

	d := NewPodmanDriver(logger).(*Driver)
	d.SetConfig(&baseConfig)
	d.buildFingerprint()
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
	createOpts := api.SpecGenerator{}
	createOpts.ContainerBasicConfig.Name = containerName
	createOpts.ContainerStorageConfig.Image = taskCfg.Image
	createOpts.ContainerBasicConfig.Command = []string{
		"sleep",
		"5",
	}

	_, err := getPodmanDriver(t, d).podman.ContainerCreate(context.Background(), createOpts)
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
	case <-time.After(time.Duration(tu.TestMultiplier()*2) * time.Second):
	}

	d.DestroyTask(task.ID, true)

	// see if the container does not exist (404)
	_, err = getPodmanDriver(t, d).podman.ContainerStats(context.Background(), containerName)
	require.Error(t, err, api.ContainerNotFound)
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

	// see if the stopped container can be inspected
	_, err = getPodmanDriver(t, d).podman.ContainerInspect(context.Background(), containerName)
	require.NoError(t, err)

	// and cleanup after ourself
	err = getPodmanDriver(t, d).podman.ContainerDelete(context.Background(), containerName, true, true)
	require.NoError(t, err)
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

	inspectData, err := getPodmanDriver(t, d).podman.ContainerInspect(context.Background(), containerName)
	require.NoError(t, err)

	// Verify that the port environment variables are set
	require.Contains(t, inspectData.Config.Env, "NOMAD_PORT_main=8888")
	require.Contains(t, inspectData.Config.Env, "NOMAD_PORT_REDIS=6379")

	// Verify that the correct ports are bound
	expectedPortBindings := map[string][]api.InspectHostPort{
		"8888/tcp": []api.InspectHostPort{
			api.InspectHostPort{
				HostIP:   "127.0.0.1",
				HostPort: strconv.Itoa(ports[0]),
			},
		},
		"8888/udp": []api.InspectHostPort{
			api.InspectHostPort{
				HostIP:   "127.0.0.1",
				HostPort: strconv.Itoa(ports[0]),
			},
		},
		"6379/tcp": []api.InspectHostPort{
			api.InspectHostPort{
				HostIP:   "127.0.0.1",
				HostPort: strconv.Itoa(ports[1]),
			},
		},
		"6379/udp": []api.InspectHostPort{
			api.InspectHostPort{
				HostIP:   "127.0.0.1",
				HostPort: strconv.Itoa(ports[1]),
			},
		},
		// FIXME: REDIS UDP
	}

	require.Exactly(t, expectedPortBindings, inspectData.HostConfig.PortBindings)

}

func TestPodmanDriver_Ports(t *testing.T) {
	if !tu.IsCI() {
		t.Parallel()
	}

	ports := freeport.MustTake(2)
	defer freeport.Return(ports)

	taskCfg := newTaskConfig("", busyboxLongRunningCmd)

	taskCfg.Ports = []string{
		"redis",
		"other",
	}

	hostIP := "127.0.0.1"

	resources := createBasicResources()
	resources.Ports = &structs.AllocatedPorts{
		{
			Label:  "redis",
			HostIP: hostIP,
			To:     8888,
			Value:  ports[0],
		},
		{
			Label:  "other",
			HostIP: hostIP,
			Value:  ports[1],
		},
	}

	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "ports",
		AllocID:   uuid.Generate(),
		Resources: resources,
	}

	task.Resources.NomadResources.Networks = []*structs.NetworkResource{
		{
			IP: "127.0.0.1",
			DynamicPorts: []structs.Port{
				{Label: "redis", Value: ports[0]},
				{Label: "other", Value: ports[1]},
			},
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

	inspectData, err := getPodmanDriver(t, d).podman.ContainerInspect(context.Background(), containerName)
	require.NoError(t, err)

	require.Contains(t, inspectData.Config.Env, fmt.Sprintf("NOMAD_PORT_redis=%d", ports[0]))

	require.Len(t, inspectData.HostConfig.PortBindings, 4)
	expectedPortBindings := map[string][]api.InspectHostPort{
		fmt.Sprintf("%d/tcp", 8888): []api.InspectHostPort{
			api.InspectHostPort{
				HostIP:   "127.0.0.1",
				HostPort: strconv.Itoa(ports[0]),
			},
		},
		fmt.Sprintf("%d/udp", 8888): []api.InspectHostPort{
			api.InspectHostPort{
				HostIP:   "127.0.0.1",
				HostPort: strconv.Itoa(ports[0]),
			},
		},
		fmt.Sprintf("%d/tcp", ports[1]): []api.InspectHostPort{
			api.InspectHostPort{
				HostIP:   "127.0.0.1",
				HostPort: strconv.Itoa(ports[1]),
			},
		},
		fmt.Sprintf("%d/udp", ports[1]): []api.InspectHostPort{
			api.InspectHostPort{
				HostIP:   "127.0.0.1",
				HostPort: strconv.Itoa(ports[1]),
			},
		},
	}
	require.Exactly(t, expectedPortBindings, inspectData.HostConfig.PortBindings)
}

func TestPodmanDriver_Ports_MissingFromGroup(t *testing.T) {
	if !tu.IsCI() {
		t.Parallel()
	}

	ports := freeport.MustTake(1)
	defer freeport.Return(ports)

	taskCfg := newTaskConfig("", busyboxLongRunningCmd)

	taskCfg.Ports = []string{
		"redis",
		"missing",
	}

	hostIP := "127.0.0.1"

	resources := createBasicResources()
	resources.Ports = &structs.AllocatedPorts{
		{
			Label:  "redis",
			HostIP: hostIP,
			Value:  ports[0],
		},
	}

	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "ports",
		AllocID:   uuid.Generate(),
		Resources: resources,
	}

	task.Resources.NomadResources.Networks = []*structs.NetworkResource{
		{
			IP: "127.0.0.1",
			DynamicPorts: []structs.Port{
				{Label: "redis", Value: ports[0]},
			},
		},
	}
	require.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	_, _, err := d.StartTask(task)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Port \"missing\" not found, check network stanza")
}

func TestPodmanDriver_Ports_MissingDriverConfig(t *testing.T) {
	if !tu.IsCI() {
		t.Parallel()
	}

	ports := freeport.MustTake(1)
	defer freeport.Return(ports)

	taskCfg := newTaskConfig("", busyboxLongRunningCmd)
	taskCfg.Ports = []string{
		"redis",
	}

	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "ports",
		AllocID:   uuid.Generate(),
		Resources: createBasicResources(),
	}

	task.Resources.NomadResources.Networks = []*structs.NetworkResource{
		{
			IP: "127.0.0.1",
			DynamicPorts: []structs.Port{
				{Label: "redis", Value: ports[0]},
			},
		},
	}
	require.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	_, _, err := d.StartTask(task)
	require.Error(t, err)
	require.Contains(t, err.Error(), "No ports defined in network stanza")
}

func TestPodmanDriver_Ports_WithPortMap(t *testing.T) {
	if !tu.IsCI() {
		t.Parallel()
	}

	ports := freeport.MustTake(1)
	defer freeport.Return(ports)

	taskCfg := newTaskConfig("", busyboxLongRunningCmd)

	taskCfg.Ports = []string{
		"redis",
	}

	taskCfg.PortMap = map[string]int{
		"main":  8888,
		"REDIS": 6379,
	}

	hostIP := "127.0.0.1"

	resources := createBasicResources()
	resources.Ports = &structs.AllocatedPorts{
		{
			Label:  "redis",
			HostIP: hostIP,
			To:     8888,
			Value:  ports[0],
		},
	}

	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "ports",
		AllocID:   uuid.Generate(),
		Resources: resources,
	}

	task.Resources.NomadResources.Networks = []*structs.NetworkResource{
		{
			IP: "127.0.0.1",
			DynamicPorts: []structs.Port{
				{Label: "redis", Value: ports[0]},
			},
		},
	}
	require.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	_, _, err := d.StartTask(task)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Invalid port declaration; use of port_map and ports")
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

	t.Skip("Skipping oom test because of podman cgroup v2 bugs")

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
	// limit memory to 5MB to trigger oom soon enough
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
	case <-time.After(time.Duration(tu.TestMultiplier()*3) * time.Second):
		t.Fatalf("Container did not exit in time")
	}
}

// check setting a user for the task
func TestPodmanDriver_User(t *testing.T) {
	// if os.Getuid() != 0 {
	// 	t.Skip("Skipping User test ")
	// }
	if !tu.IsCI() {
		t.Parallel()
	}

	taskCfg := newTaskConfig("", []string{
		// print our username to stdout
		"sh",
		"-c",
		"sleep 1; whoami",
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
	inspectData, err := getPodmanDriver(t, d).podman.ContainerInspect(context.Background(), containerName)
	require.NoError(t, err)

	// see if the configured values are set correctly
	require.Equal(t, int64(52428800), inspectData.HostConfig.Memory)
	require.Equal(t, int64(41943040), inspectData.HostConfig.MemoryReservation)
	require.Equal(t, int64(104857600), inspectData.HostConfig.MemorySwap)

	if !getPodmanDriver(t, d).cgroupV2 {
		require.Equal(t, int64(60), inspectData.HostConfig.MemorySwappiness)
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
	inspectData, err := getPodmanDriver(t, d).podman.ContainerInspect(context.Background(), containerName)
	require.NoError(t, err)

	expectedFilesystem := map[string]string{
		"/tmpdata1": "rw,rprivate,nosuid,nodev,tmpcopyup",
		"/tmpdata2": "rw,rprivate,nosuid,nodev,tmpcopyup",
	}
	require.Exactly(t, expectedFilesystem, inspectData.HostConfig.Tmpfs)

	// see if stdout was populated with expected "mount" output
	tasklog := readLogfile(t, task)
	require.Contains(t, tasklog, " on /tmpdata1 type tmpfs ")
	require.Contains(t, tasklog, " on /tmpdata2 type tmpfs ")
}

// check default capabilities
func TestPodmanDriver_DefaultCaps(t *testing.T) {
	taskCfg := newTaskConfig("", busyboxLongRunningCmd)
	inspectData := startDestroyInspect(t, taskCfg, "defaultcaps")

	// a default container should not have SYS_TIME
	require.NotContains(t, inspectData.EffectiveCaps, "CAP_SYS_TIME")
	// a default container gets CHOWN cap
	require.Contains(t, inspectData.EffectiveCaps, "CAP_CHOWN")
}

// check modified capabilities (CapAdd/CapDrop)
func TestPodmanDriver_Caps(t *testing.T) {
	taskCfg := newTaskConfig("", busyboxLongRunningCmd)
	// 	cap_add = [
	//     "SYS_TIME",
	//   ]
	taskCfg.CapAdd = []string{"SYS_TIME"}
	// 	cap_drop = [
	//     "MKNOD",
	//   ]
	taskCfg.CapDrop = []string{"CHOWN"}

	inspectData := startDestroyInspect(t, taskCfg, "caps")

	// we added SYS_TIME, so we should see it in inspect
	require.Contains(t, inspectData.EffectiveCaps, "CAP_SYS_TIME")
	// we dropped CAP_CHOWN, so we should NOT see it in inspect
	require.NotContains(t, inspectData.EffectiveCaps, "CAP_CHOWN")
}

// check enabled tty option
func TestPodmanDriver_Tty(t *testing.T) {
	taskCfg := newTaskConfig("", busyboxLongRunningCmd)
	taskCfg.Tty = true
	inspectData := startDestroyInspect(t, taskCfg, "tty")

	require.True(t, inspectData.Config.Tty)
}

// check dns server configuration
func TestPodmanDriver_Dns(t *testing.T) {
	if !tu.IsCI() {
		t.Parallel()
	}

	taskCfg := newTaskConfig("", []string{
		"sh",
		"-c",
		"sleep 1; cat /etc/resolv.conf",
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

			inspectData, err := getPodmanDriver(t, d).podman.ContainerInspect(context.Background(), containerName)
			require.NoError(t, err)
			if tc.mode == "host" {
				require.Equal(t, "host", inspectData.HostConfig.NetworkMode)
			}
			require.Equal(t, tc.gateway, inspectData.NetworkSettings.Gateway)
		})
	}
}

// test kill / signal support
func TestPodmanDriver_SignalTask(t *testing.T) {
	if !tu.IsCI() {
		t.Parallel()
	}

	taskCfg := newTaskConfig("", busyboxLongRunningCmd)
	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "signal_task",
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

	time.Sleep(300 * time.Millisecond)
	// try to send non-existing singal, should yield an error
	require.Error(t, d.SignalTask(task.ID, "FOO"))
	time.Sleep(300 * time.Millisecond)
	// send a friendly CTRL+C to busybox to stop the container
	require.NoError(t, d.SignalTask(task.ID, "SIGINT"))

	// Attempt to wait
	waitCh, err := d.WaitTask(context.Background(), task.ID)
	require.NoError(t, err)

	select {
	case res := <-waitCh:
		require.Equal(t, 130, res.ExitCode, "Should have got exit code 130")
	case <-time.After(time.Duration(tu.TestMultiplier()*5) * time.Second):
		t.Fatalf("Container did not exit in time")
	}
}

func TestPodmanDriver_Sysctl(t *testing.T) {
	if !tu.IsCI() {
		t.Parallel()
	}

	// set a uncommon somaxconn value and echo the effective
	// in-container value
	taskCfg := newTaskConfig("", []string{
		"sysctl",
		"net.core.somaxconn",
	})
	taskCfg.Sysctl = map[string]string{"net.core.somaxconn": "12321"}
	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "sysctl",
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
	require.Contains(t, tasklog, "net.core.somaxconn = 12321")

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

func getPodmanDriver(t *testing.T, harness *dtestutil.DriverHarness) *Driver {
	driver, ok := harness.Impl().(*Driver)
	require.True(t, ok)
	return driver
}

// helper to start, destroy and inspect a long running container
func startDestroyInspect(t *testing.T, taskCfg TaskConfig, taskName string) api.InspectContainerData {
	if !tu.IsCI() {
		t.Parallel()
	}

	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      taskName,
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
	inspectData, err := getPodmanDriver(t, d).podman.ContainerInspect(context.Background(), containerName)
	require.NoError(t, err)

	return inspectData
}
