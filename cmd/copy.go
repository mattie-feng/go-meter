package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// copyCmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "copy",
	Long:  `Copy`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("copyCmd")
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)
}
