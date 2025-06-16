// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/cloudbase/garm-provider-common/cloudconfig"
	"github.com/cloudbase/garm-provider-common/params"
	"github.com/cloudbase/garm-provider-common/util"
	"github.com/linode/linodego"

	"github.com/flatcar/garm-provider-linode/client/api"
	"github.com/flatcar/garm-provider-linode/config"
)

const (
	TagPool       = "garm-pool-id"
	TagController = "garm-controller-id"
)

type Linode struct {
	api    api.LinodeAPI
	config *config.Config
	id     string
}

// New returns a new Linode client.
func New(cfg *config.Config, a api.LinodeAPI, controllerID string) (*Linode, error) {
	return &Linode{
		api:    a,
		config: cfg,
		id:     controllerID,
	}, nil
}

func (c *Linode) CreateInstance(ctx context.Context, bootstrapParams params.BootstrapInstance) (*linodego.Instance, error) {
	tools, err := util.GetTools(bootstrapParams.OSType, bootstrapParams.OSArch, bootstrapParams.Tools)
	if err != nil {
		return nil, fmt.Errorf("getting tools: %w", err)
	}

	extraSpecs, err := extraSpecsFromBootstrapData(bootstrapParams)
	if err != nil {
		return nil, fmt.Errorf("getting extra specs: %w", err)
	}

	bootstrapParams.UserDataOptions.ExtraPackages = extraSpecs.ExtraPackages

	userData, err := cloudconfig.GetCloudConfig(bootstrapParams, tools, bootstrapParams.Name)
	if err != nil {
		return nil, fmt.Errorf("generating userdata: %w", err)
	}

	password, err := createRandomRootPassword()
	if err != nil {
		return nil, fmt.Errorf("generating root password: %w", err)
	}

	booted := true

	opts := linodego.InstanceCreateOptions{
		Booted: &booted,
		Image:  bootstrapParams.Image,
		Label:  bootstrapParams.Name,
		Metadata: &linodego.InstanceMetadataOptions{
			UserData: base64.StdEncoding.EncodeToString([]byte(userData)),
		},
		Region:   c.config.Region,
		RootPass: password,
		Tags: []string{
			fmt.Sprintf("%s=%s", TagPool, bootstrapParams.PoolID),
			fmt.Sprintf("%s=%s", TagController, c.id),
		},
		Type: bootstrapParams.Flavor,
	}

	instance, err := c.api.CreateInstance(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("creating instance: %w", err)
	}

	// We wait for the instance to be provisioned, booted and running.
	if err := waitUntilReady(5*time.Minute, 5*time.Second, func() (bool, error) {
		instance, err = c.api.GetInstance(ctx, instance.ID)
		if err != nil {
			return false, fmt.Errorf("getting instance: %w", err)
		}

		return instance.Status == linodego.InstanceRunning, nil
	}); err != nil {
		return nil, fmt.Errorf("getting instance running: %w", err)
	}

	return instance, nil
}

func (c *Linode) DeleteInstance(ctx context.Context, ID string) error {
	var id int
	id, err := strconv.Atoi(ID)
	if err != nil {
		i, err := c.GetInstanceID(ctx, ID)
		if err != nil {
			return fmt.Errorf("getting instance ID by its name: %w", err)
		}

		id = i
	}

	if err := c.api.DeleteInstance(ctx, id); err != nil {
		return fmt.Errorf("deleting instance from Linode API: %w", err)
	}

	return nil
}

func (c *Linode) GetInstance(ctx context.Context, ID string) (*linodego.Instance, error) {
	var id int
	id, err := strconv.Atoi(ID)
	if err != nil {
		i, err := c.GetInstanceID(ctx, ID)
		if err != nil {
			return nil, fmt.Errorf("getting instance ID by its name: %w", err)
		}

		id = i
	}

	instance, err := c.api.GetInstance(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting instance from Linode API: %w", err)
	}

	return instance, nil
}

func (c *Linode) GetInstanceID(ctx context.Context, name string) (int, error) {
	f := map[string]string{
		"label": name,
	}
	filter, err := json.Marshal(f)
	if err != nil {
		return -1, fmt.Errorf("marshalling filter: %w", err)
	}

	instances, err := c.api.ListInstances(ctx, &linodego.ListOptions{
		Filter: string(filter),
	})
	if err != nil {
		return -1, fmt.Errorf("listing instances from the API: %w", err)
	}

	if len(instances) == 0 {
		return -1, fmt.Errorf("no instances matching this name: %s", name)
	}

	return instances[0].ID, nil
}

func (c *Linode) ListInstances(ctx context.Context, poolID string) ([]linodego.Instance, error) {
	f := map[string]string{
		"tags": fmt.Sprintf("%s=%s", TagPool, poolID),
	}
	filter, err := json.Marshal(f)
	if err != nil {
		return nil, fmt.Errorf("marshalling filter: %w", err)
	}

	instances, err := c.api.ListInstances(ctx, &linodego.ListOptions{
		Filter: string(filter),
	})
	if err != nil {
		return nil, fmt.Errorf("getting instances list from Linode API: %w", err)
	}

	return instances, nil
}

func (c *Linode) RemoveAllInstances(ctx context.Context) error {
	f := map[string]string{
		"tags": fmt.Sprintf("%s=%s", TagController, c.id),
	}
	filter, err := json.Marshal(f)
	if err != nil {
		return fmt.Errorf("marshalling filter: %w", err)
	}

	instances, err := c.api.ListInstances(ctx, &linodego.ListOptions{
		Filter: string(filter),
	})
	if err != nil {
		return fmt.Errorf("getting instances list from Linode API: %w", err)
	}

	for _, instance := range instances {
		id := strconv.Itoa(instance.ID)
		if err := c.DeleteInstance(ctx, id); err != nil {
			return fmt.Errorf("deleting instance %s: %w", id, err)
		}
	}

	return nil
}
