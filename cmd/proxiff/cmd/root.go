package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "proxiff",
	Short: "A Diffy-like HTTP proxy tool for comparing responses from two servers",
	Long: `Proxiff is a HTTP proxy tool that forwards requests to two different servers
(newer and current), compares their responses, and logs any differences.

It uses a gRPC-based plugin system for customizable comparison logic.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}
