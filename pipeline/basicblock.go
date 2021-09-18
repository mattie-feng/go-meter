package pipeline

import (
	"go-meter/performinfo"
	"go-meter/randnum"
	"log"
	"os"
	"sync"
)

type basicBlock struct {
	mutexWriteBlock *sync.Mutex
	mutexGenerageBlock *sync.Mutex
	wg *sync.WaitGroup
	index int
	masterSeed uint64
	fileSeed uint64
	blockSeed uint64
	randomStata *randnum.RandomState
}

func BasicBlockInit(masterSeed, fileSeed uint64, rs *randnum.RandomState) *basicBlock{
	mutexWriteBlock := &sync.Mutex{}
	mutexGenerateBlock := &sync.Mutex{}
	wgWriteBlock := &sync.WaitGroup{}
	blockSeed := randnum.LCGRandom(rs)
	basicblock := &basicBlock{
		mutexWriteBlock,
		mutexGenerateBlock,
		wgWriteBlock,
		0,
		masterSeed,
		fileSeed,
		blockSeed,
		rs,
	}
	return basicblock
}

// 每次生成64K数据，处理成block size大小
func (b *basicBlock)generateBlock(ch chan *[]byte, buffer *[]byte ,masterBlock *[]uint64, blockSize, blockNum int) {
	b.mutexGenerageBlock.Lock()
	defer b.mutexGenerageBlock.Unlock()
	for len(*buffer) < blockSize {
		if b.index < blockNum {
			dataMaster := XORBlock(masterBlock, b.masterSeed, b.fileSeed, b.blockSeed)
			*buffer = append(*buffer, *dataMaster...)
			b.blockSeed = randnum.LCGRandom(b.randomStata) // 更新为下一个BlockSeed
			b.index++
		} else {
			ch <- buffer
			return
		}
	}
	data := (*buffer)[:blockSize]
	*buffer = (*buffer)[blockSize:]
	ch <- &data
	return
}

// 单次写入(blocksize)
func (b *basicBlock) writeBlock(ch chan *[]byte, file *os.File, blockSize int) {
	b.mutexWriteBlock.Lock()
	defer b.mutexWriteBlock.Unlock()

	block := <-ch
	_, err := file.Write(*block)
	if err != nil {
		log.Fatal(err)
	}
	performinfo.IOEnd(int64(blockSize))
	b.wg.Done()
}
