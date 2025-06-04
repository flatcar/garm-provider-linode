// SPDX-License-Identifier: Apache-2.0

package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/flatcar/garm-provider-linode/config"
)

func TestCredentials_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *config.Config
		wantErr bool
	}{
		{
			name: "valid",
			config: &config.Config{
				Token: "foo",
			},
			wantErr: false,
		},
		{
			name:    "invalid (missing token)",
			config:  &config.Config{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				assert.Error(t, tt.config.Validate())
			} else {
				assert.NoError(t, tt.config.Validate())
			}
		})
	}
}
