// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf6muxserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
)

// UpgradeResourceState calls the UpgradeResourceState method, passing `req`,
// on the provider that returned the resource specified by req.TypeName in its
// schema.
func (s muxServer) UpgradeResourceState(ctx context.Context, req *tfprotov6.UpgradeResourceStateRequest) (*tfprotov6.UpgradeResourceStateResponse, error) {
	rpc := "UpgradeResourceState"
	ctx = logging.InitContext(ctx)
	ctx = logging.RpcContext(ctx, rpc)
	server, ok := s.resources[req.TypeName]

	if !ok {
		return nil, fmt.Errorf("%q isn't supported by any servers", req.TypeName)
	}

	ctx = logging.Tfprotov6ProviderServerContext(ctx, server)
	logging.MuxTrace(ctx, "calling downstream server")

	return server.UpgradeResourceState(ctx, req)
}
