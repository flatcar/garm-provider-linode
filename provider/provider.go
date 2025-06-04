// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/flatcar/garm-provider-linode/client"
	"github.com/flatcar/garm-provider-linode/config"

	execution "github.com/cloudbase/garm-provider-common/execution/v0.1.0"
	"github.com/cloudbase/garm-provider-common/params"
)

var _ execution.ExternalProvider = &linodeProvider{}

var Version = "v0.0.0-unknown"

type linodeProvider struct {
	cfg          *config.Config
	cli          *client.LinodeClient
	controllerID string
}

func New(configPath, controllerID string) (execution.ExternalProvider, error) {
	conf, err := config.New(configPath)
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}

	cli, err := client.New(conf, controllerID)
	if err != nil {
		return nil, fmt.Errorf("getting client: %w", err)
	}

	return &linodeProvider{
		cfg:          conf,
		controllerID: controllerID,
		cli:          cli,
	}, nil
}

// CreateInstance creates a new compute instance in the provider.
func (p *linodeProvider) CreateInstance(ctx context.Context, bootstrapParams params.BootstrapInstance) (params.ProviderInstance, error) {
	return params.ProviderInstance{}, nil
}

// Delete instance will delete the instance in a provider.
func (p *linodeProvider) DeleteInstance(ctx context.Context, instance string) error {
	return nil
}

// GetInstance will return details about one instance.
func (p *linodeProvider) GetInstance(ctx context.Context, instance string) (params.ProviderInstance, error) {
	return params.ProviderInstance{}, nil
}

// ListInstances will list all instances for a provider.
func (p *linodeProvider) ListInstances(ctx context.Context, poolID string) ([]params.ProviderInstance, error) {
	return nil, nil
}

// RemoveAllInstances will remove all instances created by this provider.
func (p *linodeProvider) RemoveAllInstances(ctx context.Context) error {
	return nil
}

// Stop shuts down the instance.
func (p *linodeProvider) Stop(ctx context.Context, instance string, force bool) error {
	return nil
}

// Start boots up an instance.
func (p *linodeProvider) Start(ctx context.Context, instance string) error {
	return nil
}

// GetVersion returns the version of the provider.
func (p *linodeProvider) GetVersion(ctx context.Context) string {
	return Version
}
