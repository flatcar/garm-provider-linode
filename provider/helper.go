// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"strconv"

	"github.com/cloudbase/garm-provider-common/params"
	"github.com/linode/linodego"
)

var (
	status = map[linodego.InstanceStatus]params.InstanceStatus{
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

// instanceLinodeToGarm takes care of converting a Linode instance
// to a Garm instance.
func instanceLinodeToGarm(in *linodego.Instance) params.ProviderInstance {
	if in == nil {
		return params.ProviderInstance{}
	}

	// Extra safety net if the instance status does not exist.
	instanceStatus := params.InstanceStatusUnknown
	if s, ok := status[in.Status]; ok {
		instanceStatus = s
	}

	out := params.ProviderInstance{
		ProviderID: strconv.Itoa(in.ID),
		Name:       in.Label,
		Status:     instanceStatus,
		// TODO: Add OSType, OSName, OSVersion and OSArch.
	}

	// Best effort to get the public IP.
	ipv4s := in.IPv4
	if len(ipv4s) > 0 {
		out.Addresses = []params.Address{
			{
				Type:    params.PublicAddress,
				Address: ipv4s[0].String(),
			},
		}
	}

	return out
}
