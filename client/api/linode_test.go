// SPDX-License-Identifier: Apache-2.0

package api_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/flatcar/garm-provider-linode/client/api"
	"github.com/flatcar/garm-provider-linode/config"
)

func TestCreateAPI(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		_, err := api.New(
			&config.Config{
				Token: "foo",
			},
		)
		require.NoError(t, err)
	})

	t.Run("Failure without token", func(t *testing.T) {
		_, err := api.New(
			&config.Config{},
		)
		require.ErrorContains(t, err, "validating configuration: token needs to be set")
	})

	t.Run("Failure without configuration", func(t *testing.T) {
		_, err := api.New(nil)
		require.ErrorContains(t, err, "configuration is nil")
	})
}
