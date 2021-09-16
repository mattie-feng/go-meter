package pipeline

import (
	"bytes"
	"log"
	"os"
	"sync"
)


type basicBlock struct {
	mutexWriteBlock *sync.Mutex
	mutexGenerageBlock *sync.Mutex
	wg *sync.WaitGroup
	index int
	masterMask uint64
	fileMask uint64
	blockMask uint64
}

func BasicBlockInit(masterMask, fileMask uint64) *basicBlock{
	mutexWriteBlock := &sync.Mutex{}
	mutexGenerateBlock := &sync.Mutex{}
	wgWriteBlock := &sync.WaitGroup{}
	blockMask := Random(fileMask)

	basicblock := &basicBlock{
		mutexWriteBlock,
		mutexGenerateBlock,
		wgWriteBlock,
		0,
		masterMask,
		fileMask,
		blockMask,
	}

	return basicblock
}


// 每次生成 64KB block（根据 master block 大小）
func (b *basicBlock)generageBlock(ch chan *[]byte, buffer *bytes.Buffer ,masterBlock *[]uint64, blockSize, blockNum int) {
	b.mutexGenerageBlock.Lock()
	defer b.mutexGenerageBlock.Unlock()

	for buffer.Len() < blockSize {
		if b.index < blockNum {
			dataMaster := ByteBlock(masterBlock, b.masterMask, b.fileMask, b.blockMask)
			buffer.Write(*dataMaster)
			newBlockMask := Random(b.blockMask)
			b.blockMask = newBlockMask // 生成之后更改BlockMask
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
func (b *basicBlock)writeBlock(ch chan *[]byte, file *os.File) {
	b.mutexWriteBlock.Lock()
	defer b.mutexWriteBlock.Unlock()

	block := <- ch
	_, err := file.Write(*block)
	if err != nil{
		log.Fatal(err)
	}
	b.wg.Done()
}