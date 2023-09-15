// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf6muxserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
	"github.com/hashicorp/terraform-plugin-mux/internal/tf6dynamicvalue"
)

// ValidateProviderConfig calls the ValidateProviderConfig method on each server
// in order, passing `req`. Response diagnostics are appended from all servers.
// Response PreparedConfig must be equal across all servers with nil values
// skipped.
func (s muxServer) ValidateProviderConfig(ctx context.Context, req *tfprotov6.ValidateProviderConfigRequest) (*tfprotov6.ValidateProviderConfigResponse, error) {
	rpc := "ValidateProviderConfig"
	ctx = logging.InitContext(ctx)
	ctx = logging.RpcContext(ctx, rpc)
	var resp *tfprotov6.ValidateProviderConfigResponse

	for _, server := range s.servers {
		ctx = logging.Tfprotov6ProviderServerContext(ctx, server)
		logging.MuxTrace(ctx, "calling downstream server")

		res, err := server.ValidateProviderConfig(ctx, req)

		if err != nil {
			return resp, fmt.Errorf("error from %T validating provider config: %w", server, err)
		}

		if res == nil {
			continue
		}

		if resp == nil {
			resp = res
			continue
		}

		if len(res.Diagnostics) > 0 {
			// This could implement Diagnostic deduplication if/when
			// implemented upstream.
			resp.Diagnostics = append(resp.Diagnostics, res.Diagnostics...)
		}

		// Do not check equality on missing PreparedConfig or unset PreparedConfig
		if res.PreparedConfig == nil {
			continue
		}

		equal, err := tf6dynamicvalue.Equals(s.providerSchema.ValueType(), res.PreparedConfig, resp.PreparedConfig)

		if err != nil {
			return nil, fmt.Errorf("unable to compare ValidateProviderConfig PreparedConfig responses: %w", err)
		}

		if !equal {
			return nil, fmt.Errorf("got different ValidateProviderConfig PreparedConfig response from multiple servers, not sure which to use")
		}

		resp.PreparedConfig = res.PreparedConfig
	}

	return resp, nil
}
