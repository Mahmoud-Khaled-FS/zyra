package cmd

import (
	"github.com/spf13/cobra"

	"github.com/Mahmoud-Khaled-FS/zyra/internal/zyra"
)

var (
	listCount   bool
	listJSON    bool
	listAbs     bool
	listPattern string
)

var listCmd = &cobra.Command{
	Use:   "list [path]",
	Short: "List discovered Zyra test files",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := "."
		if len(args) == 1 {
			path = args[0]
		}

		return zyra.ListZyraFiles(zyra.ListZyraFilesOptions{
			Path:        path,
			ListCount:   listCount,
			ListJSON:    listJSON,
			ListPattern: listPattern,
		})
	},
}

func init() {
	listCmd.Flags().BoolVar(&listCount, "count", false, "Print only total number of test files")
	listCmd.Flags().BoolVar(&listJSON, "json", false, "Output as JSON")
	listCmd.Flags().BoolVar(&listAbs, "abs", false, "Show absolute paths")
	listCmd.Flags().StringVar(&listPattern, "pattern", "", "Filter files by name pattern")
	rootCmd.AddCommand(listCmd)
}
