package main

import (
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/n3xem/proxiff/comparator"
	pluginpkg "github.com/n3xem/proxiff/plugin"
)

// StatusOnlyComparator only compares HTTP status codes
type StatusOnlyComparator struct{}

func (s *StatusOnlyComparator) Compare(newer, current *comparator.Response) *comparator.Result {
	result := &comparator.Result{
		Match:   newer.StatusCode == current.StatusCode,
		Newer:   newer,
		Current: current,
	}

	if !result.Match {
		result.Difference = fmt.Sprintf("Status code differs: newer=%d, current=%d",
			newer.StatusCode, current.StatusCode)
	}

	return result
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     nil,
		JSONFormat: true,
	})

	// Create our custom comparator
	comp := &StatusOnlyComparator{}

	var pluginMap = map[string]plugin.Plugin{
		"comparator": &pluginpkg.ComparatorPlugin{Impl: comp},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: pluginpkg.Handshake,
		Plugins:         pluginMap,
		GRPCServer:      plugin.DefaultGRPCServer,
		Logger:          logger,
	})
}
