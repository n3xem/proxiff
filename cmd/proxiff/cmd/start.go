package cmd

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/hashicorp/go-plugin"
	"github.com/n3xem/proxiff/comparator"
	pluginpkg "github.com/n3xem/proxiff/plugin"
	"github.com/n3xem/proxiff/proxy"
	"github.com/spf13/cobra"
)

var (
	newerURL   string
	currentURL string
	port       string
	pluginPath string
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the proxiff proxy server",
	Long: `Start the proxiff proxy server that forwards requests to both newer and current servers,
compares their responses, and logs any differences.`,
	RunE: runStart,
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().StringVar(&newerURL, "newer", "", "URL of the newer server (required)")
	startCmd.Flags().StringVar(&currentURL, "current", "", "URL of the current server (required)")
	startCmd.Flags().StringVar(&port, "port", "8080", "Port to listen on")
	startCmd.Flags().StringVar(&pluginPath, "plugin", "", "Path to comparator plugin binary (optional)")

	startCmd.MarkFlagRequired("newer")
	startCmd.MarkFlagRequired("current")
}

func runStart(cmd *cobra.Command, args []string) error {
	// Setup logger with JSON format for Kubernetes environments
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	var comp comparator.Comparator
	var client *plugin.Client

	// Load plugin if specified, otherwise use builtin plugin
	if pluginPath != "" {
		logger.Info("loading comparator plugin",
			slog.String("plugin_path", pluginPath),
		)
		pluginComp, pluginClient, err := pluginpkg.LoadComparatorPlugin(pluginPath)
		if err != nil {
			return fmt.Errorf("failed to load plugin: %w", err)
		}
		client = pluginClient
		comp = pluginComp
		logger.Info("plugin loaded successfully",
			slog.String("plugin_path", pluginPath),
		)
	} else {
		// Use builtin plugin by default
		logger.Info("using builtin SimpleComparator plugin")
		pluginComp, pluginClient, err := pluginpkg.StartBuiltinPlugin()
		if err != nil {
			return fmt.Errorf("failed to start builtin plugin: %w", err)
		}
		client = pluginClient
		comp = pluginComp
		logger.Info("builtin plugin started successfully")
	}

	if client != nil {
		defer client.Kill()
	}

	// Create proxy
	p := proxy.NewProxy(newerURL, currentURL, comp, logger)

	// Start HTTP server
	addr := ":" + port
	logger.Info("starting proxiff server",
		slog.String("addr", addr),
		slog.String("newer_server", newerURL),
		slog.String("current_server", currentURL),
	)

	if err := http.ListenAndServe(addr, p); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
