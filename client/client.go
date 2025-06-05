// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"golang.org/x/oauth2"

	"github.com/cloudbase/garm-provider-common/cloudconfig"
	"github.com/cloudbase/garm-provider-common/params"
	"github.com/cloudbase/garm-provider-common/util"
	"github.com/linode/linodego"

	"github.com/flatcar/garm-provider-linode/config"
)

type LinodeClient struct {
	Client *linodego.Client
	config *config.Config
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

	return &LinodeClient{
		Client: &client,
		config: cfg,
	}, nil
}

func (c *LinodeClient) CreateInstance(ctx context.Context, bootstrapParams params.BootstrapInstance) (*linodego.Instance, error) {
	tools, err := util.GetTools(bootstrapParams.OSType, bootstrapParams.OSArch, bootstrapParams.Tools)
	if err != nil {
		return nil, fmt.Errorf("getting tools: %w", err)
	}

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
			fmt.Sprintf("pool=%s", bootstrapParams.PoolID),
		},
		Type: bootstrapParams.Flavor,
	}

	instance, err := c.Client.CreateInstance(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("creating instance: %w", err)
	}

	return instance, nil
}

func (c *LinodeClient) DeleteInstance(ctx context.Context, ID string) error {
	// TODO: Consider case where ID is the label (i.e name) of the instance.
	i, err := strconv.Atoi(ID)
	if err != nil {
		return fmt.Errorf("converting ID string to ID int: %w", err)
	}

	if err := c.Client.DeleteInstance(ctx, i); err != nil {
		return fmt.Errorf("deleting instance from Linode API: %w", err)
	}

	return nil
}

func (c *LinodeClient) GetInstance(ctx context.Context, ID string) (*linodego.Instance, error) {
	// TODO: Consider case where ID is the label (i.e name) of the instance.
	i, err := strconv.Atoi(ID)
	if err != nil {
		return nil, fmt.Errorf("converting ID string to ID int: %w", err)
	}

	instance, err := c.Client.GetInstance(ctx, i)
	if err != nil {
		return nil, fmt.Errorf("getting instance from Linode API: %w", err)
	}

	return instance, nil
}

func (c *LinodeClient) ListInstances(ctx context.Context, poolID string) ([]linodego.Instance, error) {
	f := map[string]string{
		"tags": fmt.Sprintf("pool=%s", poolID),
	}
	filter, err := json.Marshal(f)
	if err != nil {
		return nil, fmt.Errorf("marshalling filter: %w", err)
	}

	instances, err := c.Client.ListInstances(ctx, &linodego.ListOptions{
		Filter: string(filter),
	})
	if err != nil {
		return nil, fmt.Errorf("getting instances list from Linode API: %w", err)
	}

	return instances, nil
}
