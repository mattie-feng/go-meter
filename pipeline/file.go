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
	masterSeed uint64
	fileSeed   uint64
	basicBlock *basicBlock
}

func NewFile(filePath string, fileSize int, masterSeed, fileSeed uint64) *File {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	randumState := randnum.RandomInit(fileSeed)
	basicBlock := BasicBlockInit(randumState)
	return &File{
		file,
		fileSize,
		masterSeed,
		fileSeed,
		basicBlock,
	}
}

// 获取一个文件有多少个64K的block，向上取整
func getBlockNum(fileSize int) int {
	return int(math.Ceil(float64(fileSize) / float64(MasterBlockSize)))
}

func getXORBlock(masterBlock *[]uint64, masterSeed, fileSeed uint64) *[]uint64 {
	block := make([]uint64, 8*1024)
	for i, data := range *masterBlock {
		v := data ^ masterSeed ^ fileSeed
		block[i] = v
	}
	return &block
}

func (f *File) WriteFile(masterBlock *[]uint64, blockSize int) {
	start := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(1)
	blockNum := getBlockNum(f.fileSize)
	block := getXORBlock(masterBlock, f.masterSeed, f.fileSeed)
	times := int(math.Ceil(float64(f.fileSize) / float64(blockSize))) // 向上取整

	var bufCap int
	if blockSize > 65536 {
		bufCap = blockSize * 2
	} else {
		bufCap = 65536 * 2
	}
	buf := NewBuf(bufCap)
	ch := make(chan *[]uint64, 2)

	// 生成数据
	go func() {
		for i := 0; i < times; i++ {
			f.basicBlock.generateBlock(ch, buf, block, blockSize, blockNum)
		}
	}()
	// 写入数据
	go func() {
		for i := 0; i < times; i++ {
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
