package cmd

import (
	"fmt"
	"go-meter/pipeline"
	"strconv"
	"sync"

	"github.com/robfig/cron"
	"github.com/spf13/cobra"
)

// copyCmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "copy",
	Long:  `Copy`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Start to write files...")
		checkInputArgs()
		// number := InputArgs.Lineage[1] - InputArgs.Lineage[0] + 1
		wg := &sync.WaitGroup{}
		// wg.Add(int(number))
		wg.Add(1)

		c := cron.New()
		c.AddFunc("@every 1s", func() {
			printPerfor()
		})
		c.Start()

		// for i := InputArgs.Lineage[0]; i <= InputArgs.Lineage[1]; i++ {
		// 	go WF(i, wg)
		// }
		for i := uint(0); i <= 0; i++ {
			go WF(i, wg)
		}

		wg.Wait()
		fmt.Println("Finish to write files...")
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)
}
func WF(i uint, wg *sync.WaitGroup) {
	masterBlock := pipeline.MasterBlockInit()
	fileMask := uint64(i)
	fileSize, _ := strconv.Atoi(InputArgs.TotalSize)
	blockSize, _ := strconv.Atoi(InputArgs.BlockSize)
	filename := InputArgs.FilePath + "/" + strconv.FormatUint(fileMask, 10)
	file := pipeline.NewFile(filename, fileSize, InputArgs.MasterMask, fileMask)
	file.WriteFile(masterBlock, blockSize)
	wg.Done()
}
