package plugin

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"github.com/n3xem/proxiff/comparator"
	"github.com/n3xem/proxiff/plugin/proto"
	"google.golang.org/grpc"
)

// Handshake is used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "PROXIFF_PLUGIN",
	MagicCookieValue: "comparator",
}

// PluginMap is the map of plugins we can dispense.
var PluginMap = map[string]plugin.Plugin{
	"comparator": &ComparatorPlugin{},
}

// ComparatorPlugin is the implementation of plugin.Plugin
type ComparatorPlugin struct {
	plugin.Plugin
	// Impl is the actual implementation of the comparator
	Impl comparator.Comparator
}

func (p *ComparatorPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterComparatorPluginServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *ComparatorPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: proto.NewComparatorPluginClient(c)}, nil
}
