// SPDX-License-Identifier: Apache-2.0

package client

import "github.com/flatcar/garm-provider-linode/config"

type LinodeClient struct{}

// New returns a new Linode client.
func New(cfg *config.Config, controllerID string) (*LinodeClient, error) {
	return &LinodeClient{}, nil

}
