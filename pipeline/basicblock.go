package pipeline

import (
	"go-meter/performinfo"
	"go-meter/randnum"
	"log"
	"os"
	"sync"
)

type Buf struct {
	data []byte
	mu sync.Mutex
}

func NewBuf(cap int)(buf *Buf){
	return &Buf{data: make([]byte,0,cap)}
}

func (b *Buf) EnBuf(v []byte){
	b.mu.Lock()
	b.data = append(b.data, v...)
	b.mu.Unlock()
}

func (b *Buf) DeBuf(size int) []byte{
	b.mu.Lock()
	if len(b.data) == 0 {
		b.mu.Unlock()
		return nil
	}
	v := b.data[:size]
	b.data = b.data[size:]
	b.mu.Unlock()
	return v
}


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


func (b *basicBlock)generateBlock(ch chan *[]byte, buf *Buf ,masterBlock *[]uint64, blockSize, blockNum,times int) {
	for i := 0; i < times; i++ {
		for len(buf.data) < blockSize {
			if b.index < blockNum {
				dataMaster := XORBlock(masterBlock, b.masterSeed, b.fileSeed, b.blockSeed)
				buf.EnBuf(*dataMaster)
				b.blockSeed = randnum.LCGRandom(b.randomStata) // 更新为下一个BlockSeed
				b.index++
			} else {
				ch <- &buf.data
				return
			}
		}
		block := buf.DeBuf(blockSize)
		ch <- &block
	}
}


func (b *basicBlock) writeBlock(ch chan *[]byte, file *os.File, blockSize,times int) {
	for i :=0; i < times; i++ {
		block := <-ch
		_, err := file.Write(*block)
		if err != nil {
			log.Fatal(err)
		}
		performinfo.IOEnd(int64(blockSize))
	}
	b.wg.Done()
}

