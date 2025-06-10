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

var (
	_       execution.ExternalProvider = &linodeProvider{}
	Version                            = "v0.0.0-unknown"
)

type linodeProvider struct {
	cfg          *config.Config
	cli          *client.Linode
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
	instance, err := p.cli.CreateInstance(ctx, bootstrapParams)
	if err != nil {
		return params.ProviderInstance{}, fmt.Errorf("creating the instance: %w", err)
	}

	inst := instanceLinodeToGarm(instance)

	return inst, nil
}

// Delete instance will delete the instance in a provider.
func (p *linodeProvider) DeleteInstance(ctx context.Context, ID string) error {
	if err := p.cli.DeleteInstance(ctx, ID); err != nil {
		return fmt.Errorf("deleting instance: %w", err)
	}

	return nil
}

// GetInstance will return details about one instance.
func (p *linodeProvider) GetInstance(ctx context.Context, ID string) (params.ProviderInstance, error) {
	instance, err := p.cli.GetInstance(ctx, ID)
	if err != nil {
		return params.ProviderInstance{}, fmt.Errorf("getting instance: %w", err)
	}

	inst := instanceLinodeToGarm(instance)

	return inst, nil
}

// GetVersion returns the version of the provider.
func (p *linodeProvider) GetVersion(ctx context.Context) string {
	return Version
}

// ListInstances will list all instances for a provider.
func (p *linodeProvider) ListInstances(ctx context.Context, poolID string) ([]params.ProviderInstance, error) {
	instances, err := p.cli.ListInstances(ctx, poolID)
	if err != nil {
		return nil, fmt.Errorf("listing instances: %w", err)
	}

	res := make([]params.ProviderInstance, len(instances))
	for i, instance := range instances {
		res[i] = instanceLinodeToGarm(&instance)
	}

	return res, nil
}

// RemoveAllInstances will remove all instances created by this provider.
func (p *linodeProvider) RemoveAllInstances(ctx context.Context) error {
	// TODO: Implements p.cli.RemoveAllInstances(ctx).
	return nil
}

// Stop shuts down the instance.
func (p *linodeProvider) Stop(ctx context.Context, instance string, force bool) error {
	// TODO: Implements p.cli.StopInstance(ctx, instance, force).
	return nil
}

// Start boots up an instance.
func (p *linodeProvider) Start(ctx context.Context, instance string) error {
	return nil
}
