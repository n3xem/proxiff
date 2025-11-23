package plugin

import (
	"fmt"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/n3xem/proxiff/comparator"
)

// LoadComparatorPlugin loads a comparator plugin from the given path
func LoadComparatorPlugin(pluginPath string) (comparator.Comparator, *plugin.Client, error) {
	// Create an hclog.Logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: nil, // Suppress output
		Level:  hclog.Error,
	})

	// We're a host. Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: Handshake,
		Plugins:         PluginMap,
		Cmd:             exec.Command(pluginPath),
		Logger:          logger,
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolGRPC,
		},
	})

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		client.Kill()
		return nil, nil, fmt.Errorf("failed to connect to plugin: %w", err)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("comparator")
	if err != nil {
		client.Kill()
		return nil, nil, fmt.Errorf("failed to dispense plugin: %w", err)
	}

	// We should have a Comparator now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	comp := raw.(comparator.Comparator)

	return comp, client, nil
}
