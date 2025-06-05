// SPDX-License-Identifier: Apache-2.0

package client

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

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
