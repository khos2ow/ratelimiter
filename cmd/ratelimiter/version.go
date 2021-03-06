package ratelimiter

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/khos2ow/ratelimiter/internal/version"
)

var versionCmd = &cobra.Command{
	Args:  cobra.NoArgs,
	Use:   "version",
	Short: "Print the version of ratelimiter",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(fmt.Sprintf("ratelimiter version %s\n", version.String()))
	},
}
