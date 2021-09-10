package main

import "go-meter/pipeline"

func main() {
	data := pipeline.GetData()
	pipeline.WriteEntireFile(data,"/Users/vince/go/src/go-meter/wfile.txt",64*1024,333)
	pipeline.Compare("/Users/vince/go/src/go-meter/wfile.txt","/Users/vince/go/src/go-meter/64k.txt")
}