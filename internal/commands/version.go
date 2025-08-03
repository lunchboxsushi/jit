package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of jit",
	Long:  `Display the current version of the jit CLI tool`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("jit version %s\n", cmd.Root().Version)
	},
}

// GetVersionCmd returns the version command
func GetVersionCmd() *cobra.Command {
	return versionCmd
}
