package cmd

import (
	"github.com/spf13/cobra"
)

var InputArgs struct {
	LineAge    []int
	BlockSize  string
	TotalSize  string
	MasterMask int
	Path       string
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-meter",
	Short: "Tool of copy and compare",
	Long:  `go-meter is a tool of copy and compare.`,
	Args:  cobra.ExactArgs(1),
	// Run: func(cmd *cobra.Command, args []string) {

	// },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().IntSliceVarP(&InputArgs.LineAge, "lineage", "l", nil, "Start Lineage,End Lineage")
	rootCmd.PersistentFlags().StringVarP(&InputArgs.TotalSize, "tsize", "t", "", "Total Size")
	rootCmd.PersistentFlags().StringVarP(&InputArgs.BlockSize, "bsize", "b", "", "Block Size")
	rootCmd.PersistentFlags().StringVarP(&InputArgs.Path, "path", "p", "", "Path")
	rootCmd.PersistentFlags().IntVarP(&InputArgs.MasterMask, "mask", "m", 0, "Master Mask")
	// rootCmd.MarkFlagRequired("lineage")

}
