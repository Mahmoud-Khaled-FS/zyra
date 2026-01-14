package cmd

import (
	"fmt"

	"github.com/Mahmoud-Khaled-FS/zyra/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print Zyra version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(
			"Zyra %s\nCommit: %s\nBuilt:  %s\n",
			version.Version,
			version.Commit,
			version.Date,
		)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
