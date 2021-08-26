/*
command for inatall cobra
export GO111MODULE=on & export GOPROXY=https://goproxy.cn
go get github.com/spf13/cobra
// go get github.com/spf13/viper
*/
package main

import "go-meter/cmd"

func main() {
	cmd.Execute()
}
