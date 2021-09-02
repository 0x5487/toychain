package main

import (
	"os"
	"toychain/cmd/node"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "root",
	Short: "choose instance to run: node",
	Long:  ``,
}

func main() {
	rootCmd.AddCommand(node.RunCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
