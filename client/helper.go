// SPDX-License-Identifier: Apache-2.0

package client

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cloudbase/garm-provider-common/cloudconfig"
	"github.com/cloudbase/garm-provider-common/params"
)

type extraSpecs struct {
	//ExtraPackages to install on the VM.
	ExtraPackages []string `json:"extra_packages,omitempty" jsonschema:"description=Extra packages to install on the VM."`
	// The Cloudconfig struct from common package
	cloudconfig.CloudConfigSpec
}

// createRandomRootPassword for instances.
// (imported from github.com/linode/terraform-provider/linode)
func createRandomRootPassword() (string, error) {
	rawRootPass := make([]byte, 50)
	if _, err := rand.Read(rawRootPass); err != nil {
		return "", fmt.Errorf("generating random password: %w", err)
	}

	rootPass := base64.StdEncoding.EncodeToString(rawRootPass)

	return rootPass, nil
}

// waitUntilReady tests the checkFunction until it succeeds.
func waitUntilReady(timeout, delay time.Duration, checkFunction func() (bool, error)) error {
	after := time.After(timeout)
	for {
		select {
		case <-after:
			return fmt.Errorf("time limit exceeded")
		default:
		}

		done, err := checkFunction()
		if err != nil {
			return err
		}

		if done {
			break
		}

		time.Sleep(delay)
	}
	return nil
}

func extraSpecsFromBootstrapData(data params.BootstrapInstance) (extraSpecs, error) {
	if len(data.ExtraSpecs) == 0 {
		return extraSpecs{}, nil
	}

	var spec extraSpecs
	if err := json.Unmarshal(data.ExtraSpecs, &spec); err != nil {
		return extraSpecs{}, fmt.Errorf("unmarshalling extra_specs: %w", err)
	}

	return spec, nil
}
