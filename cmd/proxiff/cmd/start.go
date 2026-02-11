package cmd

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/n3xem/proxiff/comparator"
	"github.com/n3xem/proxiff/proxy"
	"github.com/spf13/cobra"
)

var (
	newerURL   string
	currentURL string
	port       string
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

	startCmd.MarkFlagRequired("newer")
	startCmd.MarkFlagRequired("current")
}

func runStart(cmd *cobra.Command, args []string) error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	comp := comparator.NewSimpleComparator()

	p := proxy.NewProxy(newerURL, currentURL, comp, logger)

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
