// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/flatcar/garm-provider-linode/client"
	"github.com/flatcar/garm-provider-linode/config"
	"github.com/linode/linodego"

	execution "github.com/cloudbase/garm-provider-common/execution/v0.1.0"
	"github.com/cloudbase/garm-provider-common/params"
)

var (
	_       execution.ExternalProvider = &linodeProvider{}
	Version                            = "v0.0.0-unknown"
	status                             = map[linodego.InstanceStatus]params.InstanceStatus{
		linodego.InstanceRunning:      params.InstanceRunning,
		linodego.InstanceOffline:      params.InstanceStopped,
		linodego.InstanceDeleting:     params.InstanceDeleting,
		linodego.InstanceProvisioning: params.InstancePendingCreate,
		linodego.InstanceBooting:      params.InstanceCreating,
		linodego.InstanceShuttingDown: params.InstanceStopped,
		linodego.InstanceRebooting:    params.InstanceStatusUnknown,
		linodego.InstanceMigrating:    params.InstanceStatusUnknown,
		linodego.InstanceRebuilding:   params.InstanceStatusUnknown,
		linodego.InstanceCloning:      params.InstanceStatusUnknown,
		linodego.InstanceRestoring:    params.InstanceStatusUnknown,
		linodego.InstanceResizing:     params.InstanceStatusUnknown,
	}
)

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
	instances, err := p.cli.ListInstances(poolID)
	if err != nil {
		return nil, fmt.Errorf("listing instances: %w", err)
	}

	res := make([]params.ProviderInstance, len(instances))
	for i, instance := range instances {
		// Extra safety net if the instance status does not exist.
		instanceStatus := params.InstanceStatusUnknown
		if s, ok := status[instance.Status]; ok {
			instanceStatus = s
		}

		inst := params.ProviderInstance{
			ProviderID: strconv.Itoa(instance.ID),
			Name:       instance.Label,
			Status:     instanceStatus,
			// TODO: Add OSType, OSName, OSVersion and OSArch.
		}

		// Best effort to get the public IP.
		ipv4s := instance.IPv4
		if len(ipv4s) > 0 {
			inst.Addresses = []params.Address{
				{
					Type:    params.PublicAddress,
					Address: ipv4s[0].String(),
				},
			}
		}

		res[i] = inst
	}

	return res, nil
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
