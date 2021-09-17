package pipeline

import (
	"bytes"
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

// 每次生成 64KB block（根据 master block 大小）
func (b *basicBlock)generateBlock(ch chan *[]byte, buffer *bytes.Buffer ,masterBlock *[]uint64, blockSize, blockNum int) {
	b.mutexGenerageBlock.Lock()
	defer b.mutexGenerageBlock.Unlock()

	for buffer.Len() < blockSize {
		if b.index < blockNum {
			data64K := XORBlock(masterBlock, b.masterSeed, b.fileSeed, b.blockSeed)
			buffer.Write(*data64K)
			b.blockSeed = randnum.LCGRandom(b.randomStata) // 更新为下一个BlockSeed
			b.index++
		} else {
			data := buffer.Next(buffer.Len())
			ch <- &data
			return
		}
	}
	data := buffer.Next(blockSize)
	ch <- &data
	return
}

// 单次写入(blocksize)
func (b *basicBlock) writeBlock(ch chan *[]byte, file *os.File, blockSize int) {
	b.mutexWriteBlock.Lock()
	defer b.mutexWriteBlock.Unlock()

	block := <-ch
	performinfo.IOStart(int64(blockSize))
	_, err := file.Write(*block)
	if err != nil {
		log.Fatal(err)
	}
	b.wg.Done()
}
