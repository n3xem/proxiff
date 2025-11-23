package main

import (
	"os"

	"github.com/n3xem/proxiff/cmd/proxiff/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
