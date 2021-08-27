package main

import "go-meter/pipeline"

func main() {

	pipeline.WriteToSingleFile("/Users/vince/go/src/go-meter/wfile.txt",4096*4)
	pipeline.Compare("/Users/vince/go/src/go-meter/wfile.txt","/Users/vince/go/src/go-meter/64k.txt")
}