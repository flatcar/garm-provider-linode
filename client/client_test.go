// SPDX-License-Identifier: Apache-2.0

package client_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/cloudbase/garm-provider-common/params"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/flatcar/garm-provider-linode/client"
	"github.com/flatcar/garm-provider-linode/config"
)

const (
	MockCreateInstance = "create_instance"
	MockDeleteInstance = "delete_instance"
	MockGetInstance    = "get_instance"
	MockListInstances  = "list_instances"
)

type call struct {
	name string
	args any
}

func ptr[T any](v T) *T {
	return &v
}

type mockLinode struct {
	calls          []call
	createInstance func(context.Context, linodego.InstanceCreateOptions) (*linodego.Instance, error)
	deleteInstance func(context.Context, int) error
	getInstance    func(context.Context, int) (*linodego.Instance, error)
	listInstances  func(context.Context, *linodego.ListOptions) ([]linodego.Instance, error)
}

func (m *mockLinode) CreateInstance(ctx context.Context, opts linodego.InstanceCreateOptions) (*linodego.Instance, error) {
	m.calls = append(m.calls, call{name: MockCreateInstance, args: opts})
	if m.createInstance != nil {
		return m.createInstance(ctx, opts)
	}

	return nil, nil
}

func (m *mockLinode) DeleteInstance(ctx context.Context, ID int) error {
	m.calls = append(m.calls, call{name: MockDeleteInstance, args: ID})
	if m.deleteInstance != nil {
		return m.deleteInstance(ctx, ID)
	}

	return nil
}

func (m *mockLinode) GetInstance(ctx context.Context, ID int) (*linodego.Instance, error) {
	m.calls = append(m.calls, call{name: MockGetInstance, args: ID})
	if m.getInstance != nil {
		return m.getInstance(ctx, ID)
	}

	return nil, nil
}

func (m *mockLinode) ListInstances(ctx context.Context, opts *linodego.ListOptions) ([]linodego.Instance, error) {
	m.calls = append(m.calls, call{name: MockListInstances, args: opts})
	if m.listInstances != nil {
		return m.listInstances(ctx, opts)
	}

	return nil, nil
}

func TestCreateInstance(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		m := &mockLinode{
			calls: []call{},
			createInstance: func(ctx context.Context, opts linodego.InstanceCreateOptions) (*linodego.Instance, error) {
				return &linodego.Instance{
					ID:     9876,
					Label:  "test-instance",
					Status: linodego.InstanceBooting,
				}, nil
			},
			getInstance: func(ctx context.Context, ID int) (*linodego.Instance, error) {
				return &linodego.Instance{
					ID:     9876,
					Label:  "test-instance",
					Status: linodego.InstanceRunning,
				}, nil
			},
		}

		client, err := client.New(
			&config.Config{
				Token: "foo",
			},
			m,
			"1234",
		)
		require.NoError(t, err)

		i, err := client.CreateInstance(t.Context(), params.BootstrapInstance{
			Name:          "test-instance",
			InstanceToken: "test-token",
			OSArch:        params.Amd64,
			OSType:        params.Linux,
			Flavor:        "m1.micro",
			Image:         "ubuntu-20.04",
			Tools: []params.RunnerApplicationDownload{
				{
					OS:                ptr("linux"),
					Architecture:      ptr("x64"),
					DownloadURL:       ptr("http://test.com"),
					Filename:          ptr("runner.tar.gz"),
					SHA256Checksum:    ptr("sha256:1123"),
					TempDownloadToken: ptr("test-token"),
				},
			},
			ExtraSpecs: json.RawMessage(`{
				"extra_packages": ["curl"]
			}`),
			PoolID: "test-pool",
		})
		require.NoError(t, err)

		assert.NotNil(t, i)

		require.Len(t, m.calls, 2)

		c := m.calls[0]
		assert.Equal(t, c.name, MockCreateInstance)
		opts, ok := c.args.(linodego.InstanceCreateOptions)
		require.True(t, ok)
		assert.Equal(t, opts.Tags, []string{"pool=test-pool"})
		require.NotNil(t, opts.Metadata)
		assert.NotEmpty(t, opts.Metadata.UserData)

		c = m.calls[1]
		assert.Equal(t, c.name, MockGetInstance)

		ID, ok := c.args.(int)
		require.True(t, ok)
		assert.Equal(t, ID, 9876)
	})
}

func TestDeleteInstance(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		m := &mockLinode{
			calls: []call{},
		}

		client, err := client.New(
			&config.Config{
				Token: "foo",
			},
			m,
			"1234",
		)
		require.NoError(t, err)

		err = client.DeleteInstance(t.Context(), "9876")
		require.NoError(t, err)

		require.Len(t, m.calls, 1)
		c := m.calls[0]
		assert.Equal(t, c.name, MockDeleteInstance)

		opts, ok := c.args.(int)
		require.True(t, ok)
		assert.Equal(t, opts, 9876)
	})

	t.Run("Fail from API", func(t *testing.T) {
		m := &mockLinode{
			calls: []call{},
			deleteInstance: func(ctx context.Context, ID int) error {
				return fmt.Errorf("random error from the API")
			},
		}

		client, err := client.New(
			&config.Config{
				Token: "foo",
			},
			m,
			"1234",
		)
		require.NoError(t, err)

		err = client.DeleteInstance(t.Context(), "9876")
		assert.ErrorContains(t, err, "deleting instance from Linode API: random error from the API")

		require.Len(t, m.calls, 1)
		c := m.calls[0]
		assert.Equal(t, c.name, MockDeleteInstance)

		opts, ok := c.args.(int)
		require.True(t, ok)
		assert.Equal(t, opts, 9876)
	})

	t.Run("Fail from ID not being an ID", func(t *testing.T) {
		m := &mockLinode{
			calls: []call{},
		}

		client, err := client.New(
			&config.Config{
				Token: "foo",
			},
			m,
			"1234",
		)
		require.NoError(t, err)

		err = client.DeleteInstance(t.Context(), "foo")
		assert.ErrorContains(t, err, "converting ID string to ID int: strconv.Atoi: parsing \"foo\": invalid syntax")

		require.Len(t, m.calls, 0)
	})
}

func TestGetInstance(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		m := &mockLinode{
			calls: []call{},
			getInstance: func(ctx context.Context, ID int) (*linodego.Instance, error) {
				return &linodego.Instance{
					ID: 1234,
				}, nil
			},
		}

		client, err := client.New(
			&config.Config{
				Token: "foo",
			},
			m,
			"1234",
		)
		require.NoError(t, err)

		i, err := client.GetInstance(t.Context(), "9876")
		require.NoError(t, err)
		assert.Equal(t, i.ID, 1234)

		require.Len(t, m.calls, 1)
		c := m.calls[0]
		assert.Equal(t, c.name, MockGetInstance)

		opts, ok := c.args.(int)
		require.True(t, ok)
		assert.Equal(t, opts, 9876)
	})

	t.Run("Fail from API", func(t *testing.T) {
		m := &mockLinode{
			calls: []call{},
			getInstance: func(ctx context.Context, ID int) (*linodego.Instance, error) {
				return nil, fmt.Errorf("random error from the API")
			},
		}

		client, err := client.New(
			&config.Config{
				Token: "foo",
			},
			m,
			"1234",
		)
		require.NoError(t, err)

		i, err := client.GetInstance(t.Context(), "9876")
		require.Nil(t, i)
		assert.ErrorContains(t, err, "getting instance from Linode API: random error from the API")

		require.Len(t, m.calls, 1)
		c := m.calls[0]
		assert.Equal(t, c.name, MockGetInstance)

		opts, ok := c.args.(int)
		require.True(t, ok)
		assert.Equal(t, opts, 9876)
	})

	t.Run("Fail from ID not being an ID", func(t *testing.T) {
		m := &mockLinode{
			calls: []call{},
			getInstance: func(ctx context.Context, ID int) (*linodego.Instance, error) {
				return nil, fmt.Errorf("random error from the API")
			},
		}

		client, err := client.New(
			&config.Config{
				Token: "foo",
			},
			m,
			"1234",
		)
		require.NoError(t, err)

		i, err := client.GetInstance(t.Context(), "foo")
		require.Nil(t, i)
		assert.ErrorContains(t, err, "converting ID string to ID int: strconv.Atoi: parsing \"foo\": invalid syntax")

		require.Len(t, m.calls, 0)
	})
}

func TestListInstances(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		m := &mockLinode{
			calls: []call{},
			listInstances: func(ctx context.Context, opts *linodego.ListOptions) ([]linodego.Instance, error) {
				return []linodego.Instance{
					{
						ID: 1234,
					},
					{
						ID: 5678,
					},
				}, nil
			},
		}

		client, err := client.New(
			&config.Config{
				Token: "foo",
			},
			m,
			"1234",
		)
		require.NoError(t, err)

		i, err := client.ListInstances(t.Context(), "1234")
		require.NoError(t, err)
		assert.Equal(t, len(i), 2)

		require.Len(t, m.calls, 1)
		c := m.calls[0]
		assert.Equal(t, c.name, MockListInstances)

		opts, ok := c.args.(*linodego.ListOptions)
		require.True(t, ok)
		assert.Equal(t, opts.Filter, `{"tags":"pool=1234"}`)
	})

	t.Run("Fail from API", func(t *testing.T) {
		m := &mockLinode{
			calls: []call{},
			listInstances: func(ctx context.Context, opts *linodego.ListOptions) ([]linodego.Instance, error) {
				return nil, fmt.Errorf("random error from the API")
			},
		}

		client, err := client.New(
			&config.Config{
				Token: "foo",
			},
			m,
			"1234",
		)
		require.NoError(t, err)

		i, err := client.ListInstances(t.Context(), "1234")
		require.Nil(t, i)
		assert.ErrorContains(t, err, "getting instances list from Linode API: random error from the API")

		require.Len(t, m.calls, 1)
		c := m.calls[0]
		assert.Equal(t, c.name, MockListInstances)

		opts, ok := c.args.(*linodego.ListOptions)
		require.True(t, ok)
		assert.Equal(t, opts.Filter, `{"tags":"pool=1234"}`)
	})
}
