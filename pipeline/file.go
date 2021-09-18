package pipeline

import (
	"go-meter/randnum"
	"log"
	"math"
	"os"
)

type File struct {
	//wg *sync.WaitGroup
	file       *os.File
	fileSize   int
	blockNum   int
	basicBlock *basicBlock
}

func NewFile(filePath string, fileSize int, masterMask, fileMask uint64) *File {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0766)
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
	times := int(math.Ceil(float64(f.fileSize) / float64(blockSize))) // 向上取整
	buf := make([]byte,blockSize*2)
	ch := make(chan *[]byte, 2)

	//生成数据
	for i := 0; i < times; i++ {
		go f.basicBlock.generateBlock(ch, &buf, masterBlock, blockSize, f.blockNum)
	}

	// 写入到文件
	for i := 0; i < times; i++ {
		f.basicBlock.wg.Add(1)
		go f.basicBlock.writeBlock(ch, f.file, blockSize)

	}
	f.basicBlock.wg.Wait()
	err := f.file.Close()
	if err != nil {
		log.Fatal(err)
	}
}
