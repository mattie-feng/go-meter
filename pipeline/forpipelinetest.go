package pipeline

import (
	"sync"
	"os"
	"fmt"
	"bytes"
)

// For test
func GetData() []*[]byte{
	mutexW := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	readPath := "/Users/vince/go/src/go-meter/64k.txt"
	fileR, _ := os.OpenFile(readPath, os.O_RDWR|os.O_CREATE, 0766)
	dataFileForRead := DataFile{mutexW, wg, fileR}
	masterBlock := dataFileForRead.read()

	dataAll := [](*[]byte){}
	for i:=0; i<64; i++ {
		dataOne := (*masterBlock)[:1024]
		dataAll = append(dataAll, &dataOne)
		*masterBlock = (*masterBlock)[1024:]
	}
	return dataAll
}


// For test
func Compare(fileAPath string,fileBPath string){
	fileA, _ := os.OpenFile(fileAPath, os.O_RDWR|os.O_CREATE, 0766)
	fileB, _ := os.OpenFile(fileBPath, os.O_RDWR|os.O_CREATE, 0766)

	bsA := make([]byte,4096*16)
	bsB := make([]byte,4096*16)

	fileA.Read(bsA)
	fileB.Read(bsB)

	fmt.Println(len(bsA))
	fmt.Println(bytes.Equal(bsA,bsB))

	//err:=os.Remove(fileAPath)
	//if err!=nil{
	//	fmt.Println(err)
	//}
}
