/*
command for inatall cobra
export GO111MODULE=on & export GOPROXY=https://goproxy.cn
go mod init go-meter & go run go.mod
go get github.com/spf13/cobra
go get github.com/spf13/viper
*/
package main

import (
	"fmt"
	"go-meter/cmd"
	"regexp"
)

// Check the format of size
func checkSize(size string) bool {
	// fmt.Println(size)
	str := `^([0-9.]+)(K|M|G|T)(i?B)?$`
	r := regexp.MustCompile(str)
	matchsBool := r.MatchString(size)
	return matchsBool
}

func main() {
	cmd.Execute()
	fmt.Println("main:", cmd.InputArgs)
	if !checkSize(cmd.InputArgs.BlockSize) {
		fmt.Println("Please input correct block size.")
	}
	if !checkSize(cmd.InputArgs.TotalSize) {
		fmt.Println("Please input correct total size.")
	}
	if len(cmd.InputArgs.Lineage) > 2 {
		fmt.Println("Please input correct Lineage.")
	}
	if cmd.InputArgs.MasterMask < 0 {
		fmt.Println("Mastermask is not negative.")
	}
}
