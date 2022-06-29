package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ExecConfig contains the configuration of an exec session
type ExecConfig struct {
	// Command the the command that will be invoked in the exec session.
	// Must not be empty.
	Command []string `json:"Cmd"`
	// DetachKeys are keys that will be used to detach from the exec
	// session.
	DetachKeys string `json:"DetachKeys,omitempty"`
	// Environment is a set of environment variables that will be set for
	// the first process started by the exec session.
	Environment map[string]string `json:"Env,omitempty"`
	// The user, and optionally, group to run the exec process inside the container.
	// Format is one of: user, user:group, uid, or uid:gid."
	User string `json:"User,omitempty"`
	// WorkDir is the working directory for the first process that will be
	// launched by the exec session.
	// If set to "" the exec session will be started in / within the
	// container.
	WorkDir string `json:"WorkingDir,omitempty"`
	// Tty is whether the exec session will allocate a pseudoterminal.
	Tty bool `json:"Tty,omitempty"`
	// AttachStdin is whether the STDIN stream will be forwarded to the exec
	// session's first process when attaching. Only available if Terminal is
	// false.
	AttachStdin bool `json:"AttachStdin,omitempty"`
	// AttachStdout is whether the STDOUT stream will be forwarded to the
	// exec session's first process when attaching. Only available if
	// Terminal is false.
	AttachStdout bool `json:"AttachStdout,omitempty"`
	// AttachStderr is whether the STDERR stream will be forwarded to the
	// exec session's first process when attaching. Only available if
	// Terminal is false.
	AttachStderr bool `json:"AttachStderr,omitempty"`
	// Privileged is whether the exec session will be privileged - that is,
	// will be granted additional capabilities.
	Privileged bool `json:"Privileged,omitempty"`
}

// ExecSessionResponse contains the ID of a newly created exec session
type ExecSessionResponse struct {
	ID string
}

// ExecCreate creates an exec session to run a command inside a running container
func (c *API) ExecCreate(ctx context.Context, name string, config ExecConfig) (string, error) {

	jsonString, err := json.Marshal(config)
	if err != nil {
		return "", err
	}

	res, err := c.Post(ctx, fmt.Sprintf("/v1.0.0/libpod/containers/%s/exec", name), bytes.NewBuffer(jsonString))
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(res.Body)
		return "", fmt.Errorf("cannot create exec session, status code: %d: %s", res.StatusCode, body)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	execResponse := &ExecSessionResponse{}
	err = json.Unmarshal(body, execResponse)
	if err != nil {
		return "", err
	}
	return execResponse.ID, err
}
