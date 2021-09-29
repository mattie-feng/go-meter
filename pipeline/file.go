package pipeline

import (
	"fmt"
	"go-meter/randnum"
	"log"
	"math"
	"os"
	"sync"
	"time"
)

type File struct {
	//wg *sync.WaitGroup
	file       *os.File
	fileSize   int
	blockNum   int
	basicBlock *basicBlock
}

func NewFile(filePath string, fileSize int, masterMask, fileMask uint64) *File {
	//file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0766)
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	randumState := randnum.RandomInit(fileMask)
	basicBlock := BasicBlockInit(masterMask, fileMask, randumState)
	blockNum := getBlockNum(fileSize)
	return &File{
		file,
		fileSize,
		blockNum,
		basicBlock,
	}
}

// 获取一个文件有多少个64K的block，向上取整
func getBlockNum(fileSize int) int {
	return int(math.Ceil(float64(fileSize) / float64(MasterBlockSize)))
}

func (f *File) WriteFile(masterBlock *[]uint64, blockSize int) {
	start := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(1)
	times := int(math.Ceil(float64(f.fileSize) / float64(blockSize))) // 向上取整
	var bufCap int
	if blockSize > 65536{
		bufCap = blockSize * 2
	} else {
		bufCap = 65536 * 2
	}
	buf := NewBuf(bufCap)
	ch := make(chan *[]byte, 2)

	// 生成数据
	go func(){
		for i:=0; i<times; i++{
			f.basicBlock.generateBlock(ch, buf, masterBlock, blockSize, f.blockNum)
		}
	}()
	// 写入数据
	go func(){
		for i:=0; i<times; i++{
			f.basicBlock.writeBlock(ch, f.file, blockSize)
		}
		wg.Done()
	}()

	wg.Wait()
	err := f.file.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(time.Since(start))
}
