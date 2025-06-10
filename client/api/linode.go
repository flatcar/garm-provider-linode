// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/flatcar/garm-provider-linode/config"
	"github.com/linode/linodego"
	"golang.org/x/oauth2"
)

type LinodeAPI interface {
	CreateInstance(context.Context, linodego.InstanceCreateOptions) (*linodego.Instance, error)
	DeleteInstance(context.Context, int) error
	GetInstance(context.Context, int) (*linodego.Instance, error)
	ListInstances(context.Context, *linodego.ListOptions) ([]linodego.Instance, error)
}

func New(cfg *config.Config) (LinodeAPI, error) {
	if cfg == nil {
		return nil, fmt.Errorf("configuration is nil")
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validating configuration: %w", err)
	}

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cfg.Token})
	oauth2Linode := &http.Client{
		Transport: &oauth2.Transport{
			Source: tokenSource,
		},
	}

	client := linodego.NewClient(oauth2Linode)
	return &client, nil
}
