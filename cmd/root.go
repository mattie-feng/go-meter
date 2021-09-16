package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var InputArgs struct {
	Lineage    []uint
	BlockSize  string
	TotalSize  string
	MasterMask uint64
	FilePath   string
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
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().UintSliceVarP(&InputArgs.Lineage, "lineage", "l", nil, "Start Lineage,End Lineage")
	rootCmd.PersistentFlags().StringVarP(&InputArgs.TotalSize, "tsize", "t", "", "Total Size")
	rootCmd.PersistentFlags().StringVarP(&InputArgs.BlockSize, "bsize", "b", "", "Block Size")
	rootCmd.PersistentFlags().StringVarP(&InputArgs.FilePath, "path", "p", "", "FilePath")
	rootCmd.PersistentFlags().Uint64VarP(&InputArgs.MasterMask, "mask", "m", 0, "Master Mask")

	viper.BindPFlag("TotalSize", rootCmd.PersistentFlags().Lookup("tsize"))
	viper.BindPFlag("BlockSize", rootCmd.PersistentFlags().Lookup("bsize"))
	viper.BindPFlag("MasterMask", rootCmd.PersistentFlags().Lookup("mask"))
	viper.BindPFlag("FilePath", rootCmd.PersistentFlags().Lookup("path"))
	viper.BindPFlag("Lineage", rootCmd.PersistentFlags().Lookup("lineage"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigName("go-meter")
	viper.SetConfigType("yaml")
	// viper.SetConfigFile("go-meter.yaml")

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	var err error
	if err = viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	// if _, ok := err.(viper.ConfigFileNotFoundError); ok {
	// 	// Config file not found; ignore error if desired
	// 	fmt.Println("no such config file")
	// } else {
	// 	// Config file was found but another error was produced
	// 	fmt.Println("read config error")
	// }

	// read config from yaml file
	viper.UnmarshalKey("TotalSize", &InputArgs.TotalSize)
	viper.UnmarshalKey("BlockSize", &InputArgs.BlockSize)
	viper.UnmarshalKey("MasterMask", &InputArgs.MasterMask)
	viper.UnmarshalKey("FilePath", &InputArgs.FilePath)
	viper.UnmarshalKey("Lineage", &InputArgs.Lineage)

	// fmt.Println("Total Size:", viper.GetString("TotalSize"))

}
