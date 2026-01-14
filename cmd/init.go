package cmd

import (
	"fmt"

	"github.com/Mahmoud-Khaled-FS/zyra/internal/scaffold"
	"github.com/spf13/cobra"
)

var (
	initForce bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a Zyra project",
	Long: `Initialize a new Zyra project in the current directory.

This command creates:
  - zyra.config
  - requests/ directory
  - example request file

Safe to run multiple times.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := scaffold.BuildScaffoldOptions{
			Dir:   ".",
			Force: initForce,
		}

		if err := scaffold.BuildScaffold(opts); err != nil {
			return fmt.Errorf("init failed: %w", err)
		}

		fmt.Println("Zyra project initialized 🚀")
		return nil
	},
}

func init() {
	initCmd.Flags().BoolVarP(
		&initForce,
		"force",
		"f",
		false,
		"Overwrite existing files",
	)

	rootCmd.AddCommand(initCmd)
}
