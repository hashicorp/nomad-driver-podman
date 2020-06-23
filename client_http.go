package main

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/go-hclog"
	"github.com/pascomnet/nomad-driver-podman/client"
	"github.com/pascomnet/nomad-driver-podman/client/pods"
	"github.com/pascomnet/nomad-driver-podman/models"
)

type PodmanClientV2 struct {
	client client.ProvidesaContainerCompatibleInterfaceExperimental
	log    hclog.Logger
}

func NewPodmanClientV2(ctx context.Context, logger hclog.Logger) (*PodmanClientV2, error) {
	c := client.NewHTTPClient(strfmt.Default)

	return &PodmanClientV2{
		client: c,
		log:    logger,
	}, nil
}

func (c *PodmanClientV2) GetContainerStats(containerID string) (models.PodStatsReport, error) {

	params := pods.NewStatsPodParams()
	params.NamesOrIDs = []string{containerID}
	r, err := c.client.Pods.StatsPod(params)
	if err != nil {
		return nil, err
	}
	return r.Payload.
}
