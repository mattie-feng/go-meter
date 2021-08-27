/*
command for inatall cobra
export GO111MODULE=on & export GOPROXY=https://goproxy.cn
go get github.com/spf13/cobra
// go get github.com/spf13/viper
*/
package main

import (
	"fmt"
	"go-meter/cmd"
)

type inputArgs struct {
	lineAge    [2]int
	blockSize  string
	totalSize  string
	masterMask int
	path       string
}

func main() {
	cmd.Execute()
	fmt.Println(cmd.Lineage)
}
