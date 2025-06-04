// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/linode/linodego"

	"github.com/flatcar/garm-provider-linode/config"
)

type LinodeClient struct {
	Client *linodego.Client
}

// New returns a new Linode client.
func New(cfg *config.Config, controllerID string) (*LinodeClient, error) {
	if cfg == nil {
		return nil, fmt.Errorf("configuration is nil")
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validating configuration: %w", err)
	}

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cfg.Token})
	oauth2Client := &http.Client{
		Transport: &oauth2.Transport{
			Source: tokenSource,
		},
	}

	client := linodego.NewClient(oauth2Client)

	return &LinodeClient{Client: &client}, nil
}

}
