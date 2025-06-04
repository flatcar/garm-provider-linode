// SPDX-License-Identifier: Apache-2.0

package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	// Token used to authenticate the Linode HTTP client.
	Token string `toml:"token"`
}

// New returns a new config
func New(cfgFile string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(cfgFile, &config); err != nil {
		return nil, fmt.Errorf("decoding config: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}

	return &config, nil
}

func (c *Config) Validate() error {
	if c.Token == "" {
		return fmt.Errorf("token needs to be set")
	}

	return nil
}
