// SPDX-License-Identifier: Apache-2.0

package client

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"
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
