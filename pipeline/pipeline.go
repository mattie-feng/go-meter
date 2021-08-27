package pipeline

import (
	"bytes"
	"fmt"
	"go-meter/utils"
	"os"
	"sync"
	"time"
)


type blockData struct {
	header map[string]interface{}
	block []byte
}



type DataFile struct {
	mutex *sync.Mutex
	wg *sync.WaitGroup
	f *os.File
}

func (df *DataFile)write(blockChan chan []byte) {
	df.mutex.Lock()
	defer df.mutex.Unlock()
	content := <- blockChan
	start_1 := time.Now()
	df.f.Write(content)
	df.wg.Done()

	start_2 := time.Since(start_1)
	fmt.Println("singlg cost=[%s]",start_2)

}


func (df *DataFile)close(){
	df.f.Close()
}





func (df *DataFile)read() *[]byte{
	bs := make([]byte,4096*16)
	//for {
	//	count, err := df.f.Read(bs)
	//	if err == io.EOF{
	//		fmt.Println()
	//		fmt.Println("数据读取完毕。。")
	//		break
	//	}
	//	fmt.Println()
	//	fmt.Println(string(bs[:count]))
	//	fmt.Println(len(bs[:count]))
	//}
	df.f.Read(bs)
	return &bs

}




func bsSwitch(blockSize string, fileSize string, masterBlock *[]byte,){
	bytesBS := utils.ParseSize(blockSize)
	bytesFS := utils.ParseSize(fileSize)

	lenBlock := bytesFS/bytesBS
	BsChan := make(chan []byte,10)

	fmt.Println(lenBlock,BsChan)

	//for i := 0; i < lenBlock; i++{
	//	if if bytesBS > 65536 {
	//		go func(c chan []byte){
	//			BsChana <- getBlockS()
	//		}(BsChan)
	//	} else {
	//		go func(c chan []byte){
	//			BsChana <- getBlockS()
	//		}(BsChan)
	//	}
	//}

}

// 给一个有缓存空间的chan，准备好一个个blocksize的数据块，写入这个chan，全部写完则关闭这个chanel
// 写的goroutine 从这个chanel读数据并写入（可以开启一定数量(全部blocksize数量)的goroutine）
// 一个 channel 是用来存储bs 大于 64K的，缓存的大小


func getBlockS(bsChan chan []byte, masterBlock *[]byte, blockSize int){
	// 获取小于64K的block



	//for {
	//
	//	count, err := df.f.Read(bs)
	//	if err == io.EOF{
	//		fmt.Println()
	//		fmt.Println("数据读取完毕。。")
	//		break
	//	}
	//	fmt.Println()
	//	fmt.Println(string(bs[:count]))
	//	fmt.Println(len(bs[:count]))
	//}

	//*masterBlock[:]


}



func getBlocktest(dataSource *[]byte,chanData chan []byte,mtx *sync.Mutex,blockSize int) {
	mtx.Lock()
	defer mtx.Unlock()
	chanData <- (*dataSource)[:blockSize]
	*dataSource = (*dataSource)[blockSize:]
}


func writeToSlice(bytesA *[]byte,c1 chan []byte,mtx *sync.Mutex,wg *sync.WaitGroup){
	mtx.Lock()
	defer mtx.Unlock()
	part := <- c1
	*bytesA = append(*bytesA, part...)
	wg.Done()
}


func WriteToSingleFile(filePath string, blockSize int) {
	mutexW := &sync.Mutex{}
	mutexR := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	chanData := make(chan []byte)

	// 展示取代数据的获取 64 master block
	readPath := "/Users/vince/go/src/go-meter/64k.txt"
	fileR, _ := os.OpenFile(readPath, os.O_RDWR|os.O_CREATE, 0766)
	dataFileForRead := DataFile{mutexW, wg, fileR}
	masterBlock := dataFileForRead.read()
	timesWrite := len(*masterBlock) / blockSize

	// 生成一个要写入的文件
	fileW, _ := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0766)
	dataFileForWrite := DataFile{mutexW, wg, fileW}


	start := time.Now()
	// 数据块处理成 block size 大小
	for i:=0; i<timesWrite; i++{
		go getBlocktest(masterBlock,chanData,mutexR,blockSize)
	}

	// 写入到文件
	for i:=0; i<timesWrite; i++{
		wg.Add(1)
		go dataFileForWrite.write(chanData)
	}
	wg.Wait()
	//dataFileForWrite.close()

	cost := time.Since(start)
	fmt.Printf("cost=[%s]",cost)
}

func Compare(fileAPath string,fileBPath string){
	fileA, _ := os.OpenFile(fileAPath, os.O_RDWR|os.O_CREATE, 0766)
	fileB, _ := os.OpenFile(fileBPath, os.O_RDWR|os.O_CREATE, 0766)

	bsA := make([]byte,4096*16)
	bsB := make([]byte,4096*16)

	fileA.Read(bsA)
	fileB.Read(bsB)

	fmt.Println(bytes.Equal(bsA,bsB))

}

