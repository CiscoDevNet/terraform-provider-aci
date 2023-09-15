// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf6muxserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
)

// ValidateDataResourceConfig calls the ValidateDataResourceConfig method, passing
// `req`, on the provider that returned the data source specified by
// req.TypeName in its schema.
func (s muxServer) ValidateDataResourceConfig(ctx context.Context, req *tfprotov6.ValidateDataResourceConfigRequest) (*tfprotov6.ValidateDataResourceConfigResponse, error) {
	rpc := "ValidateDataResourceConfig"
	ctx = logging.InitContext(ctx)
	ctx = logging.RpcContext(ctx, rpc)
	server, ok := s.dataSources[req.TypeName]

	if !ok {
		return nil, fmt.Errorf("%q isn't supported by any servers", req.TypeName)
	}

	ctx = logging.Tfprotov6ProviderServerContext(ctx, server)
	logging.MuxTrace(ctx, "calling downstream server")

	return server.ValidateDataResourceConfig(ctx, req)
}
