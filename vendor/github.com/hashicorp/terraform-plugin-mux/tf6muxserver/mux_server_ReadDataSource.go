// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf6muxserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
)

// ReadDataSource calls the ReadDataSource method, passing `req`, on the
// provider that returned the data source specified by req.TypeName in its
// schema.
func (s muxServer) ReadDataSource(ctx context.Context, req *tfprotov6.ReadDataSourceRequest) (*tfprotov6.ReadDataSourceResponse, error) {
	rpc := "ReadDataSource"
	ctx = logging.InitContext(ctx)
	ctx = logging.RpcContext(ctx, rpc)
	server, ok := s.dataSources[req.TypeName]

	if !ok {
		return nil, fmt.Errorf("%q isn't supported by any servers", req.TypeName)
	}

	ctx = logging.Tfprotov6ProviderServerContext(ctx, server)
	logging.MuxTrace(ctx, "calling downstream server")

	return server.ReadDataSource(ctx, req)
}
