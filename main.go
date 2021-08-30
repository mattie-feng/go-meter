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
	"regexp"
)

func checkSize(size string) bool {
	// fmt.Println(size)
	str := `^([0-9.]+)(K|M|G|T)(i?B)?$`
	r := regexp.MustCompile(str)
	matchsBool := r.MatchString(size)
	return matchsBool
}

func main() {
	cmd.Execute()
	fmt.Println(cmd.InputArgs)
	if !checkSize(cmd.InputArgs.BlockSize) {
		fmt.Println("Please input correct block size.")
	}
	if !checkSize(cmd.InputArgs.TotalSize) {
		fmt.Println("Please input correct total size.")
	}
	if len(cmd.InputArgs.LineAge) > 2 {
		fmt.Println("Please input correct lineAge.")
	}
	if cmd.InputArgs.MasterMask < 0 {
		fmt.Println("Mastermask is not negative.")
	}
}
