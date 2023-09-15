// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf6muxserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var _ tfprotov6.ProviderServer = &muxServer{}

// muxServer is a gRPC server implementation that stands in front of other
// gRPC servers, routing requests to them as if they were a single server. It
// should always be instantiated by calling NewMuxServer().
type muxServer struct {
	// Routing for data source types
	dataSources map[string]tfprotov6.ProviderServer

	// Provider schema is cached during GetProviderSchema for
	// ValidateProviderConfig equality checking.
	providerSchema *tfprotov6.Schema

	// Routing for resource types
	resources map[string]tfprotov6.ProviderServer

	// Resource capabilities are cached during GetProviderSchema
	resourceCapabilities map[string]*tfprotov6.ServerCapabilities

	// Underlying servers for requests that should be handled by all servers
	servers []tfprotov6.ProviderServer
}

// ProviderServer is a function compatible with tf6server.Serve.
func (s muxServer) ProviderServer() tfprotov6.ProviderServer {
	return &s
}

// NewMuxServer returns a muxed server that will route gRPC requests between
// tfprotov6.ProviderServers specified. When the GetProviderSchema RPC of each
// is called, there is verification that the overall muxed server is compatible
// by ensuring:
//
//   - All provider schemas exactly match
//   - All provider meta schemas exactly match
//   - Only one provider implements each managed resource
//   - Only one provider implements each data source
func NewMuxServer(_ context.Context, servers ...func() tfprotov6.ProviderServer) (*muxServer, error) {
	result := muxServer{
		dataSources:          make(map[string]tfprotov6.ProviderServer),
		resources:            make(map[string]tfprotov6.ProviderServer),
		resourceCapabilities: make(map[string]*tfprotov6.ServerCapabilities),
		servers:              make([]tfprotov6.ProviderServer, 0, len(servers)),
	}

	for _, server := range servers {
		result.servers = append(result.servers, server())
	}

	return &result, nil
}
