// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf6muxserver

import "github.com/hashicorp/terraform-plugin-go/tfprotov6"

// serverSupportsPlanDestroy returns true if the given ServerCapabilities is not
// nil and enables the PlanDestroy capability.
func serverSupportsPlanDestroy(capabilities *tfprotov6.ServerCapabilities) bool {
	if capabilities == nil {
		return false
	}

	return capabilities.PlanDestroy
}
