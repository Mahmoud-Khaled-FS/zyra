package cmd

import (
	"fmt"
	"os"

	"github.com/Mahmoud-Khaled-FS/zyra/internal/zyra"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [path]",
	Short: "Run zyra file",
	Args:  cobra.ExactArgs(1), // require exactly 1 argument
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("path is required")
		}

		if _, err := os.Stat(args[0]); err != nil {
			return fmt.Errorf("invalid path: %s", args[0])
		}

		cfg, err := cmd.Flags().GetString("config")
		if err != nil {
			return err
		}

		noTest, err := cmd.Flags().GetBool("no-test")
		if err != nil {
			return err
		}

		if cfg != "" {
			if _, err := os.Stat(cfg); err != nil {
				return fmt.Errorf("invalid path: %s", cfg)
			}
		}

		path := args[0]
		err = zyra.Run(zyra.RunOption{
			Path:       path,
			ConfigPath: cfg,
			NoTest:     noTest,
		})

		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	runCmd.Flags().StringP("config", "c", "", "config file path")
	runCmd.Flags().Bool("no-test", false, "skip test execution")
	rootCmd.AddCommand(runCmd)
}
