// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/nomad-driver-podman/api"
	"github.com/hashicorp/nomad-driver-podman/ci"
	"github.com/hashicorp/nomad/client/taskenv"
	"github.com/hashicorp/nomad/helper/testlog"
	"github.com/hashicorp/nomad/helper/uuid"
	"github.com/hashicorp/nomad/nomad/structs"
	"github.com/hashicorp/nomad/plugins/base"
	"github.com/hashicorp/nomad/plugins/drivers"
	dtestutil "github.com/hashicorp/nomad/plugins/drivers/testutils"
	spec "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/shoenig/test/must"
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
			CPUPeriod:        100000,
			CPUQuota:         100000,
			CPUShares:        500,
			MemoryLimitBytes: 256 * 1024 * 1024,
			PercentTicks:     float64(500) / float64(2000),
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

	// Set ClientHttpTimeout before calling SetConfig() because it affects the
	// HTTP client that is created.
	if v, ok := cfg["ClientHttpTimeout"]; ok {
		if sv, ok := v.(string); ok {
			pluginConfig.ClientHttpTimeout = sv
		}
	}

	// Configure the Plugin with named socket
	if v, ok := cfg["Socket"]; ok {
		if sv, ok := v.([]PluginSocketConfig); ok {
			pluginConfig.Socket = sv
		}
	}

	if err := base.MsgPackEncode(&baseConfig.PluginConfig, &pluginConfig); err != nil {
		t.Error("Unable to encode plugin config", err)
	}

	d := NewPodmanDriver(logger).(*Driver)
	must.NoError(t, d.SetConfig(&baseConfig))
	d.buildFingerprint()
	d.config.Volumes.Enabled = true
	if enforce, err := os.ReadFile("/sys/fs/selinux/enforce"); err == nil {
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
	if v, ok := cfg["extra_labels"]; ok {
		if bv, ok := v.([]string); ok {
			d.config.ExtraLabels = bv
		}
	}

	harness := dtestutil.NewDriverHarness(t, d)

	return harness
}

func TestPodmanDriver_PingPodman(t *testing.T) {
	ci.Parallel(t)

	d := podmanDriverHarness(t, nil)
	version, err := getPodmanDriver(t, d).defaultPodman.Ping(context.Background())
	must.NoError(t, err)
	must.NotEq(t, "", version)
}

func TestPodmanDriver_Start_NoImage(t *testing.T) {
	ci.Parallel(t)

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
	must.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, false)
	defer cleanup()

	_, _, err := d.StartTask(task)
	must.Error(t, err)
	must.ErrorContains(t, err, "image name required")

	_ = d.DestroyTask(task.ID, true)
}

// start a long running container
func TestPodmanDriver_Start_Wait(t *testing.T) {
	ci.Parallel(t)

	taskCfg := newTaskConfig("", busyboxLongRunningCmd)
	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "start_wait",
		AllocID:   uuid.Generate(),
		Resources: createBasicResources(),
	}
	must.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	_, _, err := d.StartTask(task)
	must.NoError(t, err)

	defer func() {
		_ = d.DestroyTask(task.ID, true)
	}()

	// Attempt to wait
	waitCh, err := d.WaitTask(context.Background(), task.ID)
	must.NoError(t, err)

	select {
	case <-waitCh:
		t.Fatalf("wait channel should not have received an exit result")
	case <-time.After(10 * time.Second):
	}
}

// test a short-living container
func TestPodmanDriver_Start_WaitFinish(t *testing.T) {
	ci.Parallel(t)

	taskCfg := newTaskConfig("", []string{"echo", "hello"})
	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "start_waitfinish",
		AllocID:   uuid.Generate(),
		Resources: createBasicResources(),
	}
	must.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	_, _, err := d.StartTask(task)
	must.NoError(t, err)

	defer func() {
		_ = d.DestroyTask(task.ID, true)
	}()

	// Attempt to wait
	waitCh, err := d.WaitTask(context.Background(), task.ID)
	must.NoError(t, err)

	select {
	case res := <-waitCh:
		must.True(t, res.Successful())
	case <-time.After(20 * time.Second):
		must.Unreachable(t, must.Sprint("timeout"))
	}
}

// TestPodmanDriver_Start_StoppedContainer asserts that Nomad will detect a
// stopped task container, remove it, and start a new container.
//
// See https://github.com/hashicorp/nomad/issues/3419
func TestPodmanDriver_Start_StoppedContainer(t *testing.T) {
	ci.Parallel(t)

	taskCfg := newTaskConfig("", []string{"sleep", "5"})
	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "start_stoppedContainer",
		AllocID:   uuid.Generate(),
		Resources: createBasicResources(),
	}
	must.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

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

	_, err := getPodmanDriver(t, d).defaultPodman.ContainerCreate(context.Background(), createOpts)
	must.NoError(t, err)

	_, _, err = d.StartTask(task)

	defer func() {
		_ = d.DestroyTask(task.ID, true)
	}()

	must.NoError(t, err)

	must.NoError(t, d.WaitUntilStarted(task.ID, 5*time.Second))
	must.NoError(t, d.DestroyTask(task.ID, true))
}

func TestPodmanDriver_Start_Wait_AllocDir(t *testing.T) {
	ci.Parallel(t)

	exp := []byte{'w', 'i', 'n'}
	file := "output.txt"

	taskCfg := newTaskConfig("", []string{
		"sh",
		"-c",
		fmt.Sprintf(`echo -n %s > $%s/%s`,
			string(exp), taskenv.AllocDir, file),
	})
	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "start_wait_allocDir",
		AllocID:   uuid.Generate(),
		Resources: createBasicResources(),
	}
	must.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	_, _, err := d.StartTask(task)
	must.NoError(t, err)

	defer func() {
		_ = d.DestroyTask(task.ID, true)
	}()

	// Attempt to wait
	waitCh, err := d.WaitTask(context.Background(), task.ID)
	must.NoError(t, err)

	select {
	case res := <-waitCh:
		must.True(t, res.Successful())
	case <-time.After(20 * time.Second):
		must.Unreachable(t, must.Sprint("timeout"))
	}

	// Check that data was written to the shared alloc directory.
	outputFile := filepath.Join(task.TaskDir().SharedAllocDir, file)
	act, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Couldn't read expected output: %v", err)
	}

	if !reflect.DeepEqual(act, exp) {
		t.Fatalf("Command outputted %v; want %v", act, exp)
	}
}

// check if container is destroyed if gc.container=true
func TestPodmanDriver_GC_Container_on(t *testing.T) {
	ci.Parallel(t)

	taskCfg := newTaskConfig("", busyboxLongRunningCmd)
	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "gc_container_on",
		AllocID:   uuid.Generate(),
		Resources: createBasicResources(),
	}
	must.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	containerName := BuildContainerName(task)

	_, _, err := d.StartTask(task)
	must.NoError(t, err)

	defer func() {
		_ = d.DestroyTask(task.ID, true)
	}()

	// Attempt to wait
	waitCh, err := d.WaitTask(context.Background(), task.ID)
	must.NoError(t, err)

	select {
	case <-waitCh:
		t.Fatalf("wait channel should not have received an exit result")
	case <-time.After(10 * time.Second):
	}

	_ = d.DestroyTask(task.ID, true)

	// see if the container does not exist (404)
	_, err = getPodmanDriver(t, d).defaultPodman.ContainerStats(context.Background(), containerName)
	must.ErrorIs(t, err, api.ContainerNotFound)
}

// check if container is destroyed if gc.container=false
func TestPodmanDriver_GC_Container_off(t *testing.T) {
	ci.Parallel(t)

	taskCfg := newTaskConfig("", busyboxLongRunningCmd)
	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "gc_container_off",
		AllocID:   uuid.Generate(),
		Resources: createBasicResources(),
	}
	must.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	opts := make(map[string]interface{})
	opts["GC.Container"] = false

	d := podmanDriverHarness(t, opts)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	containerName := BuildContainerName(task)

	_, _, err := d.StartTask(task)
	must.NoError(t, err)

	defer func() {
		_ = d.DestroyTask(task.ID, true)
	}()

	// Attempt to wait
	waitCh, err := d.WaitTask(context.Background(), task.ID)
	must.NoError(t, err)

	select {
	case <-waitCh:
		t.Fatalf("wait channel should not have received an exit result")
	case <-time.After(10 * time.Second):
	}

	_ = d.DestroyTask(task.ID, true)

	// see if the stopped container can be inspected
	_, err = getPodmanDriver(t, d).defaultPodman.ContainerInspect(context.Background(), containerName)
	must.NoError(t, err)

	// and cleanup after ourself
	err = getPodmanDriver(t, d).defaultPodman.ContainerDelete(context.Background(), containerName, true, true)
	must.NoError(t, err)
}

// Check log_opt=journald logger
func TestPodmanDriver_logJournald(t *testing.T) {
	ci.Parallel(t)

	stdoutMagic := uuid.Generate()
	stderrMagic := uuid.Generate()

	// Added a "sleep 1" because there is a race condition between WaitTask and the ending of the container
	// that makes the test flaky. If the container stops running before the WaitTask, the container logs are
	// gone when we query podman and this test fails.
	// Adding a sleep "reliably" gives an edge to the WaitTask to start running before the container exits.
	taskCfg := newTaskConfig("", []string{
		"sh",
		"-c",
		fmt.Sprintf("echo %s; 1>&2 echo %s; sleep 1", stdoutMagic, stderrMagic),
	})
	taskCfg.Logging.Driver = "journald"
	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "logJournald",
		AllocID:   uuid.Generate(),
		Resources: createBasicResources(),
	}
	must.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	_, _, err := d.StartTask(task)
	must.NoError(t, err)

	defer func() {
		_ = d.DestroyTask(task.ID, true)
	}()

	// Attempt to wait
	waitCh, err := d.WaitTask(context.Background(), task.ID)
	must.NoError(t, err)

	select {
	case <-waitCh:
	case <-time.After(10 * time.Second):
		t.Fatalf("Container did not exit in time")
	}

	stdoutLog := readStdoutLog(t, task)
	must.StrContains(t, stdoutLog, stdoutMagic)
	must.StrNotContains(t, stdoutLog, stderrMagic)

	stderrLog := readStderrLog(t, task)
	must.StrContains(t, stderrLog, stderrMagic)
	must.StrNotContains(t, stderrLog, stdoutMagic)
}

// Check log_opt=nomad logger
func TestPodmanDriver_logNomad(t *testing.T) {
	ci.Parallel(t)

	stdoutMagic := uuid.Generate()
	stderrMagic := uuid.Generate()

	taskCfg := newTaskConfig("", []string{
		"sh",
		"-c",
		fmt.Sprintf("echo %s; 1>&2 echo %s", stdoutMagic, stderrMagic),
	})
	taskCfg.Logging.Driver = "nomad"
	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "logNomad",
		AllocID:   uuid.Generate(),
		Resources: createBasicResources(),
	}
	must.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	_, _, err := d.StartTask(task)
	must.NoError(t, err)

	defer func() {
		_ = d.DestroyTask(task.ID, true)
	}()

	// Attempt to wait
	waitCh, err := d.WaitTask(context.Background(), task.ID)
	must.NoError(t, err)

	select {
	case <-waitCh:
	case <-time.After(10 * time.Second):
		t.Fatalf("Container did not exit in time")
	}

	// log_driver=nomad combines both streams into stdout, so we will find both
	// magic values in the same stream
	stdoutLog := readStdoutLog(t, task)
	must.StrContains(t, stdoutLog, stdoutMagic)
	must.StrContains(t, stdoutLog, stderrMagic)
}

// Check default plugin logger
// TODO: Can we make this nicer?
func TestPodmanDriver_logNomadDefault(t *testing.T) {
	ci.Parallel(t)

	stdoutMagic := uuid.Generate()
	stderrMagic := uuid.Generate()

	taskCfg := newTaskConfig("", []string{
		"sh",
		"-c",
		fmt.Sprintf("echo %s; 1>&2 echo %s", stdoutMagic, stderrMagic),
	})
	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "logNomad",
		AllocID:   uuid.Generate(),
		Resources: createBasicResources(),
	}
	must.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	_, _, err := d.StartTask(task)
	must.NoError(t, err)

	defer func() {
		_ = d.DestroyTask(task.ID, true)
	}()

	// Attempt to wait
	waitCh, err := d.WaitTask(context.Background(), task.ID)
	must.NoError(t, err)

	select {
	case <-waitCh:
	case <-time.After(10 * time.Second):
		t.Fatalf("Container did not exit in time")
	}

	// log_driver=nomad combines both streams into stdout, so we will find both
	// magic values in the same stream
	stdoutLog := readStdoutLog(t, task)
	must.StrContains(t, stdoutLog, stdoutMagic)
	must.StrContains(t, stdoutLog, stderrMagic)
}

// check if extra labels make it into logs
func TestPodmanDriver_ExtraLabels(t *testing.T) {
	ci.Parallel(t)

	taskCfg := newTaskConfig("", []string{
		"sh",
		"-c",
		"sleep 1",
	})

	task := &drivers.TaskConfig{
		ID:            uuid.Generate(),
		Name:          "my-task",
		TaskGroupName: "my-group",
		AllocID:       uuid.Generate(),
		Resources:     createBasicResources(),
	}
	must.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	opts := map[string]interface{}{
		"extra_labels": []string{
			"task*",
		},
	}
	d := podmanDriverHarness(t, opts)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	_, _, err := d.StartTask(task)
	must.NoError(t, err)

	defer func() {
		_ = d.DestroyTask(task.ID, true)
	}()

	// Attempt to wait
	waitCh, err := d.WaitTask(context.Background(), task.ID)
	must.NoError(t, err)

	containerName := BuildContainerName(task)
	inspectData, err := getPodmanDriver(t, d).defaultPodman.ContainerInspect(context.Background(), containerName)
	must.NoError(t, err)

	select {
	case <-waitCh:
	case <-time.After(20 * time.Second):
		t.Fatalf("Container did not exit in time")
	}

	expectedLabels := map[string]string{
		"com.hashicorp.nomad.alloc_id":        task.AllocID,
		"com.hashicorp.nomad.task_group_name": "my-group",
		"com.hashicorp.nomad.task_name":       "my-task",
	}
	must.MapEq(t, expectedLabels, inspectData.Config.Labels)
}

// check hostname task config options
func TestPodmanDriver_Hostname(t *testing.T) {
	ci.Parallel(t)

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
	must.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	_, _, err := d.StartTask(task)
	must.NoError(t, err)

	defer func() {
		_ = d.DestroyTask(task.ID, true)
	}()

	// Attempt to wait
	waitCh, err := d.WaitTask(context.Background(), task.ID)
	must.NoError(t, err)

	select {
	case <-waitCh:
	case <-time.After(10 * time.Second):
		t.Fatalf("Container did not exit in time")
	}

	// check if the hostname was visible in the container
	tasklog := readStdoutLog(t, task)
	must.StrContains(t, tasklog, shouldHaveHostname)
}

func TestPodmanDriver_setExtraHosts(t *testing.T) {
	ci.Parallel(t)

	cases := []struct {
		hosts []string
		exp   error
	}{{
		hosts: nil,
		exp:   nil,
	}, {
		hosts: []string{"example.com"},
		exp:   errors.New("cannot use \"example.com\" as extra_hosts (must be host:ip)"),
	}, {
		hosts: []string{"example.com:10.0.0.1"},
		exp:   nil,
	}, {
		hosts: []string{"ip6.example.com:[4321::1234]"},
		exp:   nil,
	}}

	for _, tc := range cases {
		opts := new(api.SpecGenerator)
		err := setExtraHosts(tc.hosts, opts)
		must.Eq(t, tc.exp, err)
	}
}

// check port_map feature
func TestPodmanDriver_PortMap(t *testing.T) {
	ci.Parallel(t)

	ports := ci.PortAllocator.Grab(2)

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
	must.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	containerName := BuildContainerName(task)

	_, _, err := d.StartTask(task)
	must.NoError(t, err)

	defer func() {
		_ = d.DestroyTask(task.ID, true)
	}()

	inspectData, err := getPodmanDriver(t, d).defaultPodman.ContainerInspect(context.Background(), containerName)
	must.NoError(t, err)

	// Verify that the port environment variables are set
	must.SliceContains(t, inspectData.Config.Env, "NOMAD_PORT_main=8888")
	must.SliceContains(t, inspectData.Config.Env, "NOMAD_PORT_REDIS=6379")

	// Verify that the correct ports are bound
	expectedPortBindings := map[string][]api.InspectHostPort{
		"8888/tcp": {{
			HostIP:   "127.0.0.1",
			HostPort: strconv.Itoa(ports[0]),
		}},
		"8888/udp": {{
			HostIP:   "127.0.0.1",
			HostPort: strconv.Itoa(ports[0]),
		}},
		"6379/tcp": {{
			HostIP:   "127.0.0.1",
			HostPort: strconv.Itoa(ports[1]),
		}},
		"6379/udp": {{
			HostIP:   "127.0.0.1",
			HostPort: strconv.Itoa(ports[1]),
		}},
		// FIXME: REDIS UDP
	}

	must.MapEq(t, expectedPortBindings, inspectData.HostConfig.PortBindings)
}

func TestPodmanDriver_Ports(t *testing.T) {
	ci.Parallel(t)

	ports := ci.PortAllocator.Grab(2)

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
	must.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	containerName := BuildContainerName(task)

	_, _, err := d.StartTask(task)
	must.NoError(t, err)

	defer func() {
		_ = d.DestroyTask(task.ID, true)
	}()

	inspectData, err := getPodmanDriver(t, d).defaultPodman.ContainerInspect(context.Background(), containerName)
	must.NoError(t, err)
	must.SliceContains(t, inspectData.Config.Env, fmt.Sprintf("NOMAD_PORT_redis=%d", ports[0]))

	must.MapLen(t, 4, inspectData.HostConfig.PortBindings)
	expectedPortBindings := map[string][]api.InspectHostPort{
		fmt.Sprintf("%d/tcp", 8888): {{
			HostIP:   "127.0.0.1",
			HostPort: strconv.Itoa(ports[0]),
		}},
		fmt.Sprintf("%d/udp", 8888): {{
			HostIP:   "127.0.0.1",
			HostPort: strconv.Itoa(ports[0]),
		}},
		fmt.Sprintf("%d/tcp", ports[1]): {{
			HostIP:   "127.0.0.1",
			HostPort: strconv.Itoa(ports[1]),
		}},
		fmt.Sprintf("%d/udp", ports[1]): {{
			HostIP:   "127.0.0.1",
			HostPort: strconv.Itoa(ports[1]),
		}},
	}
	must.MapEq(t, expectedPortBindings, inspectData.HostConfig.PortBindings)
}

func TestPodmanDriver_Ports_MissingFromGroup(t *testing.T) {
	ci.Parallel(t)

	ports := ci.PortAllocator.Grab(1)

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
	must.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	_, _, err := d.StartTask(task)
	must.ErrorContains(t, err, `Port "missing" not found, check network stanza`)
}

func TestPodmanDriver_Ports_MissingDriverConfig(t *testing.T) {
	ci.Parallel(t)

	ports := ci.PortAllocator.Grab(1)

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
	must.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	_, _, err := d.StartTask(task)
	must.Error(t, err)
	must.ErrorContains(t, err, "No ports defined in network stanza")
}

func TestPodmanDriver_Ports_WithPortMap(t *testing.T) {
	ci.Parallel(t)

	ports := ci.PortAllocator.Grab(1)

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
	must.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	_, _, err := d.StartTask(task)
	must.Error(t, err)
	must.ErrorContains(t, err, "Invalid port declaration; use of port_map and ports")
}

// check --init with default path
func TestPodmanDriver_Init(t *testing.T) {
	ci.Parallel(t)

	// only test --init if catatonit is installed
	initPath, err := exec.LookPath("catatonit")
	if os.IsNotExist(err) {
		t.Skip("Skipping --init test because catatonit is not installed")
	}
	must.NoError(t, err)

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
	must.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	_, _, err = d.StartTask(task)
	must.NoError(t, err)

	defer func() {
		_ = d.DestroyTask(task.ID, true)
	}()

	// Attempt to wait
	waitCh, err := d.WaitTask(context.Background(), task.ID)
	must.NoError(t, err)

	select {
	case <-waitCh:
	case <-time.After(10 * time.Second):
		t.Fatalf("Container did not exit in time")
	}

	apiVersion := getPodmanDriver(t, d).defaultPodman.GetAPIVersion()

	tasklog := readStdoutLog(t, task)
	// podman maps init process to /run/podman-init >= 4.2 and /dev/init < 4.2
	versionChange, _ := version.NewVersion("4.2")
	parsedApiVersion, err := version.NewVersion(apiVersion)
	if err != nil {
		t.Fatalf("Not a valid version: %s", apiVersion)
	}

	if parsedApiVersion.LessThan(versionChange) {
		must.StrContains(t, tasklog, "/dev/init")
	} else {
		must.StrContains(t, tasklog, "/run/podman-init")
	}

}

// test oom flag propagation
func TestPodmanDriver_OOM(t *testing.T) {
	ci.Parallel(t)

	t.Skip("Skipping oom test because of podman cgroup v2 bugs")

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
		path, pathErr := exec.LookPath("catatonit")
		if pathErr != nil {
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
	must.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	_, _, err = d.StartTask(task)
	must.NoError(t, err)

	defer func() {
		_ = d.DestroyTask(task.ID, true)
	}()

	// Attempt to wait
	waitCh, err := d.WaitTask(context.Background(), task.ID)
	must.NoError(t, err)

	select {
	case res := <-waitCh:
		must.False(t, res.Successful(), must.Sprint("Should have failed because of oom but was successful"))
		must.True(t, res.OOMKilled, must.Sprint("OOM Flag not set"))
		must.ErrorContains(t, res.Err, "OOM killer")
	case <-time.After(20 * time.Second):
		must.Unreachable(t, must.Sprint("Container did not exit in time"))
	}
}

// check setting a user for the task
func TestPodmanDriver_User(t *testing.T) {
	// if os.Getuid() != 0 {
	// 	t.Skip("Skipping User test ")
	// }
	ci.Parallel(t)

	taskCfg := newTaskConfig("", []string{
		// print our username to stdout
		"sh",
		"-c",
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
	must.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	_, _, err := d.StartTask(task)
	must.NoError(t, err)

	defer func() {
		_ = d.DestroyTask(task.ID, true)
	}()

	// Attempt to wait
	waitCh, err := d.WaitTask(context.Background(), task.ID)
	must.NoError(t, err)

	select {
	case res := <-waitCh:
		// should have a exitcode=0 result
		must.True(t, res.Successful())
	case <-time.After(10 * time.Second):
		must.Unreachable(t, must.Sprint("Container did not exit in time"))
	}

	// see if stdout was populated with the "whoami" output
	tasklog := readStdoutLog(t, task)
	must.StrContains(t, tasklog, "www-data")

}

func TestPodmanDriver_Device(t *testing.T) {
	ci.Parallel(t)

	taskCfg := newTaskConfig("", []string{
		// print our username to stdout
		"sh",
		"-c",
		"sleep 1; ls -l /dev/net/tun",
	})

	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "device",
		AllocID:   uuid.Generate(),
		Resources: createBasicResources(),
	}
	taskCfg.Devices = []string{"/dev/net/tun"}
	must.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	_, _, err := d.StartTask(task)
	must.NoError(t, err)

	defer func() {
		_ = d.DestroyTask(task.ID, true)
	}()

	// Attempt to wait
	waitCh, err := d.WaitTask(context.Background(), task.ID)
	must.NoError(t, err)

	select {
	case res := <-waitCh:
		// should have a exitcode=0 result
		must.True(t, res.Successful())
	case <-time.After(10 * time.Second):
		t.Fatalf("Container did not exit in time")
	}

	// see if stdout was populated with the "whoami" output
	tasklog := readStdoutLog(t, task)
	must.StrContains(t, tasklog, "dev/net/tun")

}

// test memory/swap options
func TestPodmanDriver_Swap(t *testing.T) {
	ci.Parallel(t)

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
	// set shm_size of 100m
	taskCfg.ShmSize = "100m"
	must.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	containerName := BuildContainerName(task)

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	_, _, err := d.StartTask(task)
	must.NoError(t, err)

	defer func() {
		_ = d.DestroyTask(task.ID, true)
	}()

	// Attempt to wait
	waitCh, err := d.WaitTask(context.Background(), task.ID)
	must.NoError(t, err)

	select {
	case <-waitCh:
		t.Fatalf("wait channel should not have received an exit result")
	case <-time.After(10 * time.Second):
	}
	// inspect container to learn about the actual podman limits
	inspectData, err := getPodmanDriver(t, d).defaultPodman.ContainerInspect(context.Background(), containerName)
	must.NoError(t, err)

	// see if the configured values are set correctly
	must.Eq(t, 52428800, inspectData.HostConfig.Memory)
	must.Eq(t, 41943040, inspectData.HostConfig.MemoryReservation)
	must.Eq(t, 104857600, inspectData.HostConfig.MemorySwap)
	must.Eq(t, 104857600, inspectData.HostConfig.ShmSize)

	if !getPodmanDriver(t, d).defaultPodman.IsCgroupV2() {
		must.Eq(t, 60, inspectData.HostConfig.MemorySwappiness)
	}
}

// check tmpfs mounts
func TestPodmanDriver_Tmpfs(t *testing.T) {
	ci.Parallel(t)

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
	must.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	containerName := BuildContainerName(task)
	_, _, err := d.StartTask(task)
	must.NoError(t, err)

	defer func() {
		_ = d.DestroyTask(task.ID, true)
	}()

	// Attempt to wait
	waitCh, err := d.WaitTask(context.Background(), task.ID)
	must.NoError(t, err)

	select {
	case <-waitCh:
	case <-time.After(10 * time.Second):
		t.Fatalf("Container did not exit in time")
	}

	// see if tmpfs was propagated to podman
	inspectData, err := getPodmanDriver(t, d).defaultPodman.ContainerInspect(context.Background(), containerName)
	must.NoError(t, err)

	expectedFilesystem := map[string]string{
		"/tmpdata1": "rw,rprivate,nosuid,nodev,tmpcopyup",
		"/tmpdata2": "rw,rprivate,nosuid,nodev,tmpcopyup",
	}
	must.MapEq(t, expectedFilesystem, inspectData.HostConfig.Tmpfs)

	// see if stdout was populated with expected "mount" output
	tasklog := readStdoutLog(t, task)
	must.StrContains(t, tasklog, " on /tmpdata1 type tmpfs ")
	must.StrContains(t, tasklog, " on /tmpdata2 type tmpfs ")
}

// check mount options
func TestPodmanDriver_Mount(t *testing.T) {
	ci.Parallel(t)

	taskCfg := newTaskConfig("", []string{
		// print our username to stdout
		"sh",
		"-c",
		"mount|grep check",
	})
	taskCfg.Volumes = []string{
		// explicitly check that we can have more then one option
		"/tmp:/checka:ro,shared",
		"/tmp:/checkb:private",
		"/tmp:/checkc",
	}

	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "Mount",
		AllocID:   uuid.Generate(),
		Resources: createBasicResources(),
	}
	// use "www-data" as a user for our test, it's part of the busybox image
	must.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	containerName := BuildContainerName(task)
	_, _, err := d.StartTask(task)
	must.NoError(t, err)

	defer func() {
		_ = d.DestroyTask(task.ID, true)
	}()

	// Attempt to wait
	waitCh, err := d.WaitTask(context.Background(), task.ID)
	must.NoError(t, err)

	select {
	case <-waitCh:
	case <-time.After(10 * time.Second):
		t.Fatalf("Container did not exit in time")
	}

	// see if options where correctly sent to podman
	inspectData, err := getPodmanDriver(t, d).defaultPodman.ContainerInspect(context.Background(), containerName)
	must.NoError(t, err)

	aok := false
	bok := false
	cok := false

	// this part is a bit verbose but the exact mount options and
	// their order depend on the target os
	// so we need to dissect each result line
	for _, bind := range inspectData.HostConfig.Binds {
		if strings.HasPrefix(bind, "/tmp:/check") {
			prefix := bind[0:13]
			opts := strings.Split(bind[13:], ",")
			if prefix == "/tmp:/checka:" {
				must.SliceContains(t, opts, "ro")
				must.SliceContains(t, opts, "shared")
				aok = true
			}
			if prefix == "/tmp:/checkb:" {
				must.SliceContains(t, opts, "rw")
				must.SliceContains(t, opts, "private")
				bok = true
			}
			if prefix == "/tmp:/checkc:" {
				must.SliceContains(t, opts, "rw")
				must.SliceContains(t, opts, "rprivate")
				cok = true
			}
		}
	}
	must.True(t, aok, must.Sprint("checka not ok"))
	must.True(t, bok, must.Sprint("checkb not ok"))
	must.True(t, cok, must.Sprint("checkc not ok"))

	// see if stdout was populated with expected "mount" output
	tasklog := readStdoutLog(t, task)
	must.StrContains(t, tasklog, " on /checka type ")
	must.StrContains(t, tasklog, " on /checkb type ")
	must.StrContains(t, tasklog, " on /checkc type ")
}

// check default capabilities
func TestPodmanDriver_DefaultCaps(t *testing.T) {
	ci.Parallel(t)

	taskCfg := newTaskConfig("", busyboxLongRunningCmd)
	inspectData := startDestroyInspect(t, taskCfg, "defaultcaps")

	// a default container should not have SYS_TIME
	must.SliceNotContains(t, inspectData.EffectiveCaps, "CAP_SYS_TIME")
	// a default container gets CHOWN cap
	must.SliceContains(t, inspectData.EffectiveCaps, "CAP_CHOWN")
}

// check default process label
func TestPodmanDriver_DefaultProcessLabel(t *testing.T) {
	ci.Parallel(t)

	taskCfg := newTaskConfig("", busyboxLongRunningCmd)
	inspectData := startDestroyInspect(t, taskCfg, "defaultprocesslabel")

	// TODO: skip SELinux check in CI since it's not supported yet.
	// https://github.com/hashicorp/nomad-driver-podman/pull/139#issuecomment-1192929834
	if ci.TestSELinux() {
		// a default container gets "disable" process label
		must.StrContains(t, inspectData.ProcessLabel, "label=disable")
	}
}

// check modified capabilities (CapAdd/CapDrop/SelinuxOpts)
func TestPodmanDriver_Caps(t *testing.T) {
	ci.Parallel(t)

	taskCfg := newTaskConfig("", busyboxLongRunningCmd)
	// 	cap_add = [
	//     "SYS_TIME",
	//   ]
	taskCfg.CapAdd = []string{"SYS_TIME"}
	// 	cap_drop = [
	//     "MKNOD",
	//   ]
	taskCfg.CapDrop = []string{"CHOWN"}
	// 	selinux_opts = [
	//     "disable",
	//   ]
	taskCfg.SelinuxOpts = []string{"disable"}

	inspectData := startDestroyInspect(t, taskCfg, "caps")

	// we added SYS_TIME, so we should see it in inspect
	must.SliceContains(t, inspectData.EffectiveCaps, "CAP_SYS_TIME")
	// we dropped CAP_CHOWN, so we should NOT see it in inspect
	must.SliceNotContains(t, inspectData.EffectiveCaps, "CAP_CHOWN")

	// TODO: skip SELinux check in CI since it's not supported yet.
	// https://github.com/hashicorp/nomad-driver-podman/pull/139#issuecomment-1192929834
	if ci.TestSELinux() {
		// we added "disable" process label, so we should see it in inspect
		must.StrContains(t, inspectData.ProcessLabel, "label=disable")
	}
}

// check security_opt option
func TestPodmanDriver_SecurityOpt(t *testing.T) {
	taskCfg := newTaskConfig("", busyboxLongRunningCmd)
	// add a security_opt
	taskCfg.SecurityOpt = []string{"no-new-privileges"}
	inspectData := startDestroyInspect(t, taskCfg, "securityopt")
	// and compare it
	must.SliceContains(t, inspectData.HostConfig.SecurityOpt, "no-new-privileges")
}

// check enabled tty option
func TestPodmanDriver_Tty(t *testing.T) {
	ci.Parallel(t)

	taskCfg := newTaskConfig("", busyboxLongRunningCmd)
	taskCfg.Tty = true
	inspectData := startDestroyInspect(t, taskCfg, "tty")
	must.True(t, inspectData.Config.Tty)
}

// check labels option
func TestPodmanDriver_Labels(t *testing.T) {
	ci.Parallel(t)

	taskCfg := newTaskConfig("", busyboxLongRunningCmd)
	taskCfg.Labels = map[string]string{"nomad": "job"}
	inspectData := startDestroyInspect(t, taskCfg, "labels")
	expectedLabels := map[string]string{"nomad": "job"}
	must.MapEq(t, expectedLabels, inspectData.Config.Labels)
}

// check enabled privileged option
func TestPodmanDriver_Privileged(t *testing.T) {
	ci.Parallel(t)

	taskCfg := newTaskConfig("", busyboxLongRunningCmd)
	taskCfg.Privileged = true
	inspectData := startDestroyInspect(t, taskCfg, "privileged")
	must.True(t, inspectData.HostConfig.Privileged)
}

// check apparmor default value
func TestPodmanDriver_AppArmorDefault(t *testing.T) {
	ci.Parallel(t)

	d := podmanDriverHarness(t, nil)

	// Skip test if apparmor is not available
	if !getPodmanDriver(t, d).defaultPodman.IsAppArmorEnabled() {
		t.Skip("Skipping AppArmor test ")
	}

	defaultCfg := newTaskConfig("", busyboxLongRunningCmd)
	defaultInspectResult := startDestroyInspect(t, defaultCfg, "aa-default")
	must.StrContains(t, defaultInspectResult.AppArmorProfile, "containers-default")
}

// check apparmor option
func TestPodmanDriver_AppArmorUnconfined(t *testing.T) {
	ci.Parallel(t)

	d := podmanDriverHarness(t, nil)

	// Skip test if apparmor is not available
	if !getPodmanDriver(t, d).defaultPodman.IsAppArmorEnabled() {
		t.Skip("Skipping AppArmor test ")
	}

	unconfinedCfg := newTaskConfig("", busyboxLongRunningCmd)
	unconfinedCfg.ApparmorProfile = "unconfined"
	unconfinedInspectResult := startDestroyInspect(t, unconfinedCfg, "aa-unconfined")
	must.Eq(t, "unconfined", unconfinedInspectResult.AppArmorProfile)
}

// check ulimit option
func TestPodmanDriver_Ulimit(t *testing.T) {
	ci.Parallel(t)

	taskCfg := newTaskConfig("", busyboxLongRunningCmd)
	taskCfg.Ulimit = map[string]string{"nproc": "4096", "nofile": "2048:4096"}
	inspectData := startDestroyInspect(t, taskCfg, "ulimits")

	nofileLimit := api.InspectUlimit{
		Name: "RLIMIT_NOFILE",
		Soft: uint64(2048),
		Hard: uint64(4096),
	}
	nprocLimit := api.InspectUlimit{
		Name: "RLIMIT_NPROC",
		Soft: uint64(4096),
		Hard: uint64(4096),
	}

	must.Len(t, 2, inspectData.HostConfig.Ulimits)
	if inspectData.HostConfig.Ulimits[0].Name == "RLIMIT_NOFILE" {
		must.Eq(t, nofileLimit, inspectData.HostConfig.Ulimits[0])
		must.Eq(t, nprocLimit, inspectData.HostConfig.Ulimits[1])
	} else {
		must.Eq(t, nofileLimit, inspectData.HostConfig.Ulimits[1])
		must.Eq(t, nprocLimit, inspectData.HostConfig.Ulimits[0])
	}
}

// check pids_limit option
func TestPodmanDriver_PidsLimit(t *testing.T) {
	ci.Parallel(t)

	taskCfg := newTaskConfig("", busyboxLongRunningCmd)
	// set a random pids limit
	taskCfg.PidsLimit = int64(100 + rand.Intn(100))
	inspectData := startDestroyInspect(t, taskCfg, "pidsLimit")
	// and compare it
	must.Eq(t, taskCfg.PidsLimit, inspectData.HostConfig.PidsLimit)
}

// check enabled readonly_rootfs option
func TestPodmanDriver_ReadOnlyFilesystem(t *testing.T) {
	ci.Parallel(t)

	taskCfg := newTaskConfig("", busyboxLongRunningCmd)
	taskCfg.ReadOnlyRootfs = true
	inspectData := startDestroyInspect(t, taskCfg, "readonly_rootfs")
	must.True(t, inspectData.HostConfig.ReadonlyRootfs)
}

// check userns mode configuration (single value)
func TestPodmanDriver_UsernsMode(t *testing.T) {
	ci.Parallel(t)

	// TODO: run test once CI can use rootless Podman.
	t.Skip("Test suite requires rootful Podman")

	taskCfg := newTaskConfig("", busyboxLongRunningCmd)
	taskCfg.UserNS = "host"
	inspectData := startDestroyInspect(t, taskCfg, "userns")

	must.Eq(t, "host", inspectData.HostConfig.UsernsMode)
}

// check userns mode configuration (parsed value)
func TestPodmanDriver_UsernsModeParsed(t *testing.T) {
	ci.Parallel(t)

	// TODO: run test once CI can use rootless Podman.
	t.Skip("Test suite requires rootful Podman")

	taskCfg := newTaskConfig("", busyboxLongRunningCmd)
	taskCfg.UserNS = "keep-id:uid=200,gid=210"
	inspectData := startDestroyInspect(t, taskCfg, "userns")
	must.Eq(t, "keep-id:uid=200,gid=210", inspectData.HostConfig.UsernsMode)
}

// check dns server configuration
func TestPodmanDriver_DNS(t *testing.T) {
	ci.Parallel(t)

	taskCfg := newTaskConfig("", []string{
		"sh",
		"-c",
		"cat /etc/resolv.conf",
	})
	// network {
	//   dns {
	//     servers = ["1.1.1.1"]
	// 	   searches = ["internal.corp"]
	//     options = ["ndots:2"]
	//   }
	// }
	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "dns",
		AllocID:   uuid.Generate(),
		Resources: createBasicResources(),
		DNS: &drivers.DNSConfig{
			Servers:  []string{"1.1.1.1"},
			Searches: []string{"internal.corp"},
			Options:  []string{"ndots:2"},
		},
	}
	must.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	_, _, err := d.StartTask(task)
	must.NoError(t, err)

	defer func() {
		_ = d.DestroyTask(task.ID, true)
	}()

	// Attempt to wait
	waitCh, err := d.WaitTask(context.Background(), task.ID)
	must.NoError(t, err)

	select {
	case res := <-waitCh:
		// should have a exitcode=0 result
		must.True(t, res.Successful())
	case <-time.After(10 * time.Second):
		t.Fatalf("Container did not exit in time")
	}

	// see if stdout was populated with the correct output
	tasklog := readStdoutLog(t, task)
	must.StrContains(t, tasklog, "nameserver 1.1.1.1")
	must.StrContains(t, tasklog, "search internal.corp")
	must.StrContains(t, tasklog, "options ndots:2")

}

// TestPodmanDriver_NetworkMode asserts we can specify different network modes
// Default podman cni subnet 10.88.0.0/16
func TestPodmanDriver_NetworkModes(t *testing.T) {
	ci.Parallel(t)

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

			must.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

			d := podmanDriverHarness(t, nil)
			defer d.MkAllocDir(task, true)()

			containerName := BuildContainerName(task)
			_, _, err := d.StartTask(task)
			must.NoError(t, err)

			defer func() {
				_ = d.DestroyTask(task.ID, true)
			}()

			must.NoError(t, d.WaitUntilStarted(task.ID, 20*time.Second))

			inspectData, err := getPodmanDriver(t, d).defaultPodman.ContainerInspect(context.Background(), containerName)
			must.NoError(t, err)
			if tc.mode == "host" {
				must.Eq(t, "host", inspectData.HostConfig.NetworkMode)
			}
			must.Eq(t, tc.gateway, inspectData.NetworkSettings.Gateway)
		})
	}
}

// let a task join NetworkNS of another container via network_mode=container:
func TestPodmanDriver_NetworkMode_Container(t *testing.T) {
	ci.Parallel(t)

	allocId := uuid.Generate()

	// we're running "nc" on localhost here
	mainTaskCfg := newTaskConfig("", []string{
		"nc",
		"-l",
		"-p",
		"6748",
		"-s",
		"localhost",
	})
	mainTask := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "maintask",
		AllocID:   allocId,
		Resources: createBasicResources(),
	}
	must.NoError(t, mainTask.EncodeConcreteDriverConfig(&mainTaskCfg))

	// we're running a second task in same networkNS and invoke netstat in it
	sidecarTaskCfg := newTaskConfig("", []string{
		"sh",
		"-c",
		"netstat -tulpen",
	})
	// join maintask network
	sidecarTaskCfg.NetworkMode = "container:maintask-" + allocId
	sidecarTask := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "sidecar",
		AllocID:   allocId,
		Resources: createBasicResources(),
	}
	must.NoError(t, sidecarTask.EncodeConcreteDriverConfig(&sidecarTaskCfg))

	mainHarness := podmanDriverHarness(t, nil)
	mainCleanup := mainHarness.MkAllocDir(mainTask, true)
	defer mainCleanup()

	_, _, err := mainHarness.StartTask(mainTask)
	must.NoError(t, err)
	defer func() {
		_ = mainHarness.DestroyTask(mainTask.ID, true)
	}()

	sidecarHarness := podmanDriverHarness(t, nil)
	sidecarCleanup := sidecarHarness.MkAllocDir(sidecarTask, true)
	defer sidecarCleanup()

	_, _, err = sidecarHarness.StartTask(sidecarTask)
	must.NoError(t, err)
	defer func() {
		_ = sidecarHarness.DestroyTask(sidecarTask.ID, true)
	}()

	// Attempt to wait
	waitCh, err := sidecarHarness.WaitTask(context.Background(), sidecarTask.ID)
	must.NoError(t, err)

	select {
	case res := <-waitCh:
		// should have a exitcode=0 result
		must.True(t, res.Successful())
	case <-time.After(10 * time.Second):
		t.Fatalf("Sidecar did not exit in time")
	}

	// see if stdout was populated with the correct output
	tasklog := readStdoutLog(t, sidecarTask)
	must.StrContains(t, tasklog, "127.0.0.1:6748")
}

// let a task joint NetorkNS of another container via network_mode=task:
func TestPodmanDriver_NetworkMode_Task(t *testing.T) {
	ci.Parallel(t)

	allocId := uuid.Generate()

	// we're running "nc" on localhost here
	mainTaskCfg := newTaskConfig("", []string{
		"nc",
		"-l",
		"-p",
		"6748",
		"-s",
		"localhost",
	})
	mainTask := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "maintask",
		AllocID:   allocId,
		Resources: createBasicResources(),
	}
	must.NoError(t, mainTask.EncodeConcreteDriverConfig(&mainTaskCfg))

	// we're running a second task in same networkNS and invoke netstat in it
	sidecarTaskCfg := newTaskConfig("", []string{
		"sh",
		"-c",
		"netstat -tulpen",
	})
	// join maintask network
	sidecarTaskCfg.NetworkMode = "task:maintask"
	sidecarTask := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "sidecar",
		AllocID:   allocId,
		Resources: createBasicResources(),
	}
	must.NoError(t, sidecarTask.EncodeConcreteDriverConfig(&sidecarTaskCfg))

	mainHarness := podmanDriverHarness(t, nil)
	mainCleanup := mainHarness.MkAllocDir(mainTask, true)
	defer mainCleanup()

	_, _, err := mainHarness.StartTask(mainTask)
	must.NoError(t, err)
	defer func() {
		_ = mainHarness.DestroyTask(mainTask.ID, true)
	}()

	sidecarHarness := podmanDriverHarness(t, nil)
	sidecarCleanup := sidecarHarness.MkAllocDir(sidecarTask, true)
	defer sidecarCleanup()

	_, _, err = sidecarHarness.StartTask(sidecarTask)
	must.NoError(t, err)
	defer func() {
		_ = sidecarHarness.DestroyTask(sidecarTask.ID, true)
	}()

	// Attempt to wait
	waitCh, err := sidecarHarness.WaitTask(context.Background(), sidecarTask.ID)
	must.NoError(t, err)

	select {
	case res := <-waitCh:
		// should have a exitcode=0 result
		must.True(t, res.Successful())
	case <-time.After(10 * time.Second):
		t.Fatalf("Sidecar did not exit in time")
	}

	// see if stdout was populated with the correct output
	tasklog := readStdoutLog(t, sidecarTask)
	must.StrContains(t, tasklog, "127.0.0.1:6748")
}

// test kill / signal support
func TestPodmanDriver_SignalTask(t *testing.T) {
	ci.Parallel(t)

	taskCfg := newTaskConfig("", busyboxLongRunningCmd)
	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "signal_task",
		AllocID:   uuid.Generate(),
		Resources: createBasicResources(),
	}
	must.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	_, _, err := d.StartTask(task)
	must.NoError(t, err)

	defer func() {
		_ = d.DestroyTask(task.ID, true)
	}()

	time.Sleep(300 * time.Millisecond)
	// try to send non-existing singal, should yield an error
	must.Error(t, d.SignalTask(task.ID, "FOO"))
	time.Sleep(300 * time.Millisecond)
	// send a friendly CTRL+C to busybox to stop the container
	must.NoError(t, d.SignalTask(task.ID, "SIGINT"))

	// Attempt to wait
	waitCh, err := d.WaitTask(context.Background(), task.ID)
	must.NoError(t, err)

	select {
	case res := <-waitCh:
		must.Eq(t, 130, res.ExitCode)
	case <-time.After(60 * time.Second):
		t.Fatalf("Container did not exit in time")
	}
}

func TestPodmanDriver_Sysctl(t *testing.T) {
	ci.Parallel(t)

	// Test fails if not root
	if os.Geteuid() != 0 {
		t.Skip("Skipping Sysctl test")
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
	must.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	_, _, err := d.StartTask(task)
	must.NoError(t, err)

	defer func() {
		_ = d.DestroyTask(task.ID, true)
	}()

	// Attempt to wait
	waitCh, err := d.WaitTask(context.Background(), task.ID)
	must.NoError(t, err)

	select {
	case <-waitCh:
	case <-time.After(20 * time.Second):
		t.Fatalf("Container did not exit in time")
	}

	tasklog := readStdoutLog(t, task)
	must.StrContains(t, tasklog, "net.core.somaxconn = 12321")

}

// Make sure we can pull and start "non-latest" containers
func TestPodmanDriver_Pull(t *testing.T) {
	ci.Parallel(t)

	type imageTestCase struct {
		Image    string
		TaskName string
	}

	testCases := []imageTestCase{
		{Image: "busybox:unstable", TaskName: "pull_tag"},
		{Image: "busybox", TaskName: "pull_non_tag"},
	}
	switch runtime.GOARCH {
	case "arm64":
		testCases = append(testCases, imageTestCase{Image: "busybox@sha256:5acba83a746c7608ed544dc1533b87c737a0b0fb730301639a0179f9344b1678", TaskName: "pull_digest"})
	case "amd64":
		testCases = append(testCases, imageTestCase{Image: "busybox@sha256:ce98b632acbcbdf8d6fdc50d5f91fea39c770cd5b3a2724f52551dde4d088e96", TaskName: "pull_digest"})
	default:
		t.Fatalf("Unexpected architecture %s", runtime.GOARCH)
	}

	for _, testCase := range testCases {
		startDestroyInspectImage(t, testCase.Image, testCase.TaskName)
	}
}

func startDestroyInspectImage(t *testing.T, image string, taskName string) {
	taskCfg := newTaskConfig(image, busyboxLongRunningCmd)
	inspectData := startDestroyInspect(t, taskCfg, taskName)

	d := podmanDriverHarness(t, nil)

	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      taskName,
		AllocID:   uuid.Generate(),
		Resources: createBasicResources(),
	}
	driver := getPodmanDriver(t, d)
	imageID, err := driver.createImage(image, &TaskAuthConfig{}, false, false, driver.defaultPodman, 5*time.Minute, task)
	must.NoError(t, err)
	must.Eq(t, imageID, inspectData.Image)
}

// TestPodmanDriver_Pull_Timeout verifies that the task image_pull_timeout
// configuration is respected and that it can be set to higher value than the
// driver's client_http_timeout.
//
// This test starts a proxy on port 5000 and requires the machine running the
// test to allow an insecure registry at localhost:5000. To run this test add
// the following to /etc/containers/registries.conf:
//
// [[registry]]
// location = "localhost:5000"
// insecure = true
func TestPodmanDriver_Pull_Timeout(t *testing.T) {
	ci.Parallel(t)

	// Check if the machine running the test is properly configured to run this
	// test.
	expectedRegistry := `[[registry]]
location = "localhost:5000"
insecure = true`

	content, err := os.ReadFile("/etc/containers/registries.conf")
	must.NoError(t, err)
	if !strings.Contains(string(content), expectedRegistry) {
		t.Skip("Skipping test because /etc/containers/registries.conf doesn't have insecure registry config for localhost:5000")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a listener on port 5000.
	l, err := net.Listen("tcp", "127.0.0.1:5000")
	must.NoError(t, err)

	// Proxy requests to quay.io except when trying to pull a busybox image.
	parsedURL, err := url.Parse("https://quay.io")
	must.NoError(t, err)

	proxy := httputil.NewSingleHostReverseProxy(parsedURL)
	proxy.ModifyResponse = func(resp *http.Response) error {
		if strings.Contains(resp.Request.URL.Path, "busybox") {
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(5 * time.Second):
				return fmt.Errorf("expected image pull to timeout")
			}
		}
		return nil
	}
	// Create a new HTTP test server but don't start it so we can use our
	// own listener on port 5000.
	ts := httptest.NewUnstartedServer(proxy)
	ts.Listener.Close()
	ts.Listener = l
	ts.Start()
	defer ts.Close()

	// Create the test harness and a sample task.
	d := podmanDriverHarness(t, map[string]interface{}{
		"ClientHttpTimeout": "1s",
	})
	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "test",
		AllocID:   uuid.Generate(),
		Resources: createBasicResources(),
	}

	// Time how long it takes to pull the image.
	now := time.Now()

	resultCh := make(chan error)
	go func() {
		// Pull image using our proxy.
		image := "localhost:5000/quay/busybox:latest"
		driver := getPodmanDriver(t, d)
		_, err = driver.createImage(image, &TaskAuthConfig{}, false, true, driver.defaultPodman, 3*time.Second, task)
		resultCh <- err
	}()

	err = <-resultCh
	cancel()

	// Verify that the timeout happened after client_http_timeout.
	must.Greater(t, 2*time.Second, time.Since(now))
	must.ErrorContains(t, err, "context deadline exceeded")
}

func Test_createImage(t *testing.T) {
	ci.Parallel(t)

	testCases := []struct {
		Image     string
		Reference string
	}{
		{Image: "busybox:musl", Reference: "docker.io/library/busybox:musl"},
		{Image: "docker://busybox:latest", Reference: "docker.io/library/busybox:latest"},
		{Image: "docker.io/library/busybox", Reference: "docker.io/library/busybox:latest"},
	}

	for _, testCase := range testCases {
		createInspectImage(t, testCase.Image, testCase.Reference)
	}
}

func Test_createImageArchives(t *testing.T) {
	ci.Parallel(t)

	doesNotExist := func(filepath string) bool {
		_, err := os.Stat(filepath)
		if errors.Is(err, os.ErrNotExist) {
			return true
		}
		must.NoError(t, err)
		return false
	}

	if doesNotExist("/tmp/oci-archive") || doesNotExist("/tmp/docker-archive") {
		t.Skip("Skipping image archive test. Missing prepared archive file(s).")
	}

	archiveDir := os.Getenv("ARCHIVE_DIR")
	if archiveDir == "" {
		t.Skip("Skipping image archive test. Missing \"ARCHIVE_DIR\" environment variable")
	}

	testCases := []struct {
		Image     string
		Reference string
	}{
		{
			Image:     fmt.Sprintf("oci-archive:%s/oci-archive", archiveDir),
			Reference: "docker.io/library/alpine:latest",
		},
		{
			Image:     fmt.Sprintf("docker-archive:%s/docker-archive", archiveDir),
			Reference: "docker.io/library/alpine:latest",
		},
	}

	for _, testCase := range testCases {
		createInspectImage(t, testCase.Image, testCase.Reference)
	}
}

func createInspectImage(t *testing.T, image, reference string) {
	d := podmanDriverHarness(t, nil)

	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      "inspectImage",
		AllocID:   uuid.Generate(),
		Resources: createBasicResources(),
	}
	driver := getPodmanDriver(t, d)
	idTest, err := driver.createImage(image, &TaskAuthConfig{}, false, false, driver.defaultPodman, 5*time.Minute, task)
	must.NoError(t, err)

	idRef, err := driver.defaultPodman.ImageInspectID(context.Background(), reference)
	must.NoError(t, err)
	must.Eq(t, idRef, idTest)
}

func Test_setTaskCpuset(t *testing.T) {
	ci.Parallel(t)

	t.Run("empty", func(t *testing.T) {
		sysResources := &drivers.LinuxResources{CpusetCpus: ""}
		taskCPU := new(spec.LinuxCPU)
		cfg := TaskConfig{}
		err := setCPUResources(cfg, sysResources, taskCPU)
		must.NoError(t, err)
		must.Eq(t, "", taskCPU.Cpus)
	})

	t.Run("non-empty", func(t *testing.T) {
		sysResources := &drivers.LinuxResources{CpusetCpus: "2,5-8"}
		taskCPU := new(spec.LinuxCPU)
		cfg := TaskConfig{}
		err := setCPUResources(cfg, sysResources, taskCPU)
		must.NoError(t, err)
		must.Eq(t, "2,5-8", taskCPU.Cpus)
	})
}

func Test_cpuLimits(t *testing.T) {
	ci.Parallel(t)

	numCores := runtime.NumCPU()
	cases := []struct {
		name            string
		systemResources drivers.LinuxResources
		hardLimit       bool
		period          uint64
		expectedQuota   int64
		expectedPeriod  uint64
	}{
		{
			name: "no hard limit",
			systemResources: drivers.LinuxResources{
				PercentTicks: 1.0,
				CPUPeriod:    100000,
			},
			hardLimit:      false,
			expectedQuota:  100000,
			expectedPeriod: 100000,
		},
		{
			name: "hard limit max quota",
			systemResources: drivers.LinuxResources{
				PercentTicks: 1.0,
				CPUPeriod:    100000,
			},
			hardLimit:      true,
			expectedQuota:  100000,
			expectedPeriod: 100000,
		},
		{
			name: "hard limit, half quota",
			systemResources: drivers.LinuxResources{
				PercentTicks: 0.5,
				CPUPeriod:    100000,
			},
			hardLimit:      true,
			expectedQuota:  50000,
			expectedPeriod: 100000,
		},
		{
			name: "hard limit, custom period",
			systemResources: drivers.LinuxResources{
				PercentTicks: 1.0,
				CPUPeriod:    100000,
			},
			hardLimit:      true,
			period:         20000,
			expectedQuota:  20000,
			expectedPeriod: 20000,
		},
		{
			name: "hard limit, half quota, custom period",
			systemResources: drivers.LinuxResources{
				PercentTicks: 0.5,
				CPUPeriod:    100000,
			},
			hardLimit:      true,
			period:         20000,
			expectedQuota:  10000,
			expectedPeriod: 20000,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			taskCPU := &spec.LinuxCPU{}
			cfg := TaskConfig{
				CPUHardLimit: c.hardLimit,
				CPUCFSPeriod: c.period,
			}
			err := setCPUResources(cfg, &c.systemResources, taskCPU)
			must.Nil(t, err)

			if c.hardLimit {
				must.NotNil(t, taskCPU.Quota)
				must.Eq(t, c.expectedQuota*int64(numCores), *taskCPU.Quota)

				must.NotNil(t, taskCPU.Period)
				must.Eq(t, c.expectedPeriod, *taskCPU.Period)
			} else {
				must.Nil(t, taskCPU.Quota)
				must.Nil(t, taskCPU.Period)
			}
		})
	}
}

func Test_memoryLimits(t *testing.T) {
	ci.Parallel(t)

	cases := []struct {
		name         string
		memResources drivers.MemoryResources
		reservation  string
		expectedHard int64
		expectedSoft int64
	}{
		{
			name: "plain",
			memResources: drivers.MemoryResources{
				MemoryMB: 20,
			},
			expectedHard: 20 * 1024 * 1024,
			expectedSoft: 0,
		},
		{
			name: "memory oversubscription",
			memResources: drivers.MemoryResources{
				MemoryMB:    20,
				MemoryMaxMB: 30,
			},
			expectedHard: 30 * 1024 * 1024,
			expectedSoft: 20 * 1024 * 1024,
		},
		{
			name: "plain but using memory reservations",
			memResources: drivers.MemoryResources{
				MemoryMB: 20,
			},
			reservation:  "10m",
			expectedHard: 20 * 1024 * 1024,
			expectedSoft: 10 * 1024 * 1024,
		},
		{
			name: "oversubscription but with specifying memory reservation",
			memResources: drivers.MemoryResources{
				MemoryMB:    20,
				MemoryMaxMB: 30,
			},
			reservation:  "10m",
			expectedHard: 30 * 1024 * 1024,
			expectedSoft: 10 * 1024 * 1024,
		},
		{
			name: "oversubscription but with specifying high memory reservation",
			memResources: drivers.MemoryResources{
				MemoryMB:    20,
				MemoryMaxMB: 30,
			},
			reservation:  "25m",
			expectedHard: 30 * 1024 * 1024,
			expectedSoft: 20 * 1024 * 1024,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			hard, soft, err := memoryLimits(c.memResources, c.reservation)
			must.NoError(t, err)

			if c.expectedHard > 0 {
				must.NotNil(t, hard)
				must.Eq(t, c.expectedHard, *hard)
			} else {
				must.Nil(t, hard)
			}

			if c.expectedSoft > 0 {
				must.NotNil(t, soft)
				must.Eq(t, c.expectedSoft, *soft)
			} else {
				must.Nil(t, soft)
			}
		})
	}
}

func Test_parseImage(t *testing.T) {
	ci.Parallel(t)

	digest := "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	testCases := []struct {
		Input     string
		Name      string
		Transport string
	}{
		{Input: "quay.io/repo/busybox:glibc", Name: "quay.io/repo/busybox:glibc", Transport: "docker"},
		{Input: "docker.io/library/root", Name: "docker.io/library/root:latest", Transport: "docker"},
		{Input: "docker://root", Name: "docker.io/library/root:latest", Transport: "docker"},
		{Input: "user/repo@" + digest, Name: "docker.io/user/repo@" + digest, Transport: "docker"},
		{Input: "user/repo:tag", Name: "docker.io/user/repo:tag", Transport: "docker"},
		{Input: "url:5000/repo", Name: "url:5000/repo:latest", Transport: "docker"},
		{Input: "url:5000/repo:tag", Name: "url:5000/repo:tag", Transport: "docker"},
		{Input: "oci-archive:path:tag", Name: "path:tag", Transport: "oci-archive"},
		{Input: "docker-archive:path:image:tag", Name: "path:docker.io/library/image:tag", Transport: "docker-archive"},
	}
	for _, testCase := range testCases {
		ref, err := parseImage(testCase.Input)
		must.NoError(t, err)
		must.Eq(t, testCase.Transport, ref.Transport().Name())
		if ref.Transport().Name() == "docker" {
			must.Eq(t, testCase.Name, ref.DockerReference().String())
		} else {
			must.Eq(t, testCase.Name, ref.StringWithinTransport())
		}

	}
}

func Test_namedSocketIsDefault(t *testing.T) {
	ci.Parallel(t)

	defaultSocket := PluginSocketConfig{Name: "default", SocketPath: "unix://run/podman/podman.sock"}
	sockets := make([]PluginSocketConfig, 1)
	sockets[0] = defaultSocket

	d := podmanDriverHarness(t, map[string]interface{}{
		"Socket": sockets,
	})

	// Get back our podman socket via its name
	api, err := getPodmanDriver(t, d).getPodmanClient("default")
	must.NoError(t, err)
	must.True(t, api.IsDefaultClient())
}

// Test that omitting a name in the Socket definition gives it the name "default"
func Test_unnamedSocketIsRenamedToDefault(t *testing.T) {
	ci.Parallel(t)

	defaultSocket := PluginSocketConfig{Name: "", SocketPath: "unix://run/podman/podman.sock"}
	sockets := make([]PluginSocketConfig, 1)
	sockets[0] = defaultSocket

	d := podmanDriverHarness(t, map[string]interface{}{
		"Socket": sockets,
	})

	_, err := getPodmanDriver(t, d).getPodmanClient("default")
	must.NoError(t, err)
}

// Test that if there's one socket, it becomes the default socket (regardless of name)
func Test_namedSocketBecomesDefaultSocket(t *testing.T) {
	ci.Parallel(t)

	socketPath := "unix://run/podman/podman.sock"
	if os.Geteuid() != 0 {
		socketPath = fmt.Sprintf("unix://run/user/%d/podman/podman.sock", os.Geteuid())
	}

	defaultSocket := PluginSocketConfig{Name: "podmanSock", SocketPath: socketPath}
	sockets := make([]PluginSocketConfig, 1)
	sockets[0] = defaultSocket

	d := podmanDriverHarness(t, map[string]interface{}{
		"Socket": sockets,
	})

	api, err := getPodmanDriver(t, d).getPodmanClient("podmanSock")
	must.NoError(t, err)
	must.True(t, api.IsDefaultClient())

	fingerprint := getPodmanDriver(t, d).buildFingerprint()
	must.MapContainsKey(t, fingerprint.Attributes, "driver.podman.socketName")

	originalName, ok := fingerprint.Attributes["driver.podman.socketName"].GetString()
	must.True(t, ok)
	must.Eq(t, "podmanSock", originalName)
	isDefaultAttr, ok := fingerprint.Attributes["driver.podman.defaultPodman"].GetBool()
	must.True(t, ok)
	must.True(t, isDefaultAttr)
}

// A socket with non alpha numeric ([a-zA-Z0-9_]) chars gets cleaned before being used
// in attributes.
func Test_socketNameIsCleanedInAttributes(t *testing.T) {
	ci.Parallel(t)

	socketPath := "unix://run/podman/podman.sock"
	if os.Geteuid() != 0 {
		socketPath = fmt.Sprintf("unix://run/user/%d/podman/podman.sock", os.Geteuid())
	}

	// We need two sockets because the default socket will get prefix driver.podman. without the socket name.
	defaultSocket := PluginSocketConfig{Name: "default", SocketPath: socketPath}
	socketWithSpaces := PluginSocketConfig{Name: "a socket with spaces", SocketPath: socketPath}
	sockets := make([]PluginSocketConfig, 2)
	sockets[0] = defaultSocket
	sockets[1] = socketWithSpaces

	d := podmanDriverHarness(t, map[string]interface{}{
		"Socket": sockets,
	})

	fingerprint := getPodmanDriver(t, d).buildFingerprint()
	must.MapContainsKey(t, fingerprint.Attributes, "driver.podman.a_socket_with_spaces.socketName")

	originalName, ok := fingerprint.Attributes["driver.podman.a_socket_with_spaces.socketName"].GetString()
	must.True(t, ok)
	must.Eq(t, "a socket with spaces", originalName)
}

// read a tasks stdout logfile into a string, fail on error
func readStdoutLog(t *testing.T, task *drivers.TaskConfig) string {
	logfile := filepath.Join(filepath.Dir(task.StdoutPath), fmt.Sprintf("%s.stdout.0", task.Name))
	stdout, err := os.ReadFile(logfile)
	must.NoError(t, err)
	return string(stdout)
}

// read a tasks stderr logfile into a string, fail on error
func readStderrLog(t *testing.T, task *drivers.TaskConfig) string {
	logfile := filepath.Join(filepath.Dir(task.StderrPath), fmt.Sprintf("%s.stderr.0", task.Name))
	stderr, err := os.ReadFile(logfile)
	must.NoError(t, err)
	return string(stderr)
}

func newTaskConfig(image string, command []string) TaskConfig {
	if len(image) == 0 {
		image = "docker.io/library/busybox:latest"
	}

	return TaskConfig{
		Image: image,
		// LoadImage: loadImage,
		Command:          command[0],
		Args:             command[1:],
		ImagePullTimeout: "5m",
	}
}

func getPodmanDriver(t *testing.T, harness *dtestutil.DriverHarness) *Driver {
	driver, ok := harness.Impl().(*Driver)
	must.True(t, ok)
	return driver
}

// helper to start, destroy and inspect a long running container
func startDestroyInspect(t *testing.T, taskCfg TaskConfig, taskName string) api.InspectContainerData {
	task := &drivers.TaskConfig{
		ID:        uuid.Generate(),
		Name:      taskName,
		AllocID:   uuid.Generate(),
		Resources: createBasicResources(),
	}
	must.NoError(t, task.EncodeConcreteDriverConfig(&taskCfg))

	d := podmanDriverHarness(t, nil)
	cleanup := d.MkAllocDir(task, true)
	defer cleanup()

	containerName := BuildContainerName(task)

	_, _, err := d.StartTask(task)
	must.NoError(t, err)

	defer func() {
		_ = d.DestroyTask(task.ID, true)
	}()

	// Attempt to wait
	waitCh, err := d.WaitTask(context.Background(), task.ID)
	must.NoError(t, err)

	select {
	case <-waitCh:
		t.Fatalf("wait channel should not have received an exit result")
	case <-time.After(20 * time.Second):
	}
	driver := getPodmanDriver(t, d)
	inspectData, err := driver.defaultPodman.ContainerInspect(context.Background(), containerName)
	must.NoError(t, err)

	return inspectData
}
