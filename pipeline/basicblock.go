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
	mu   sync.Mutex
}

func NewBuf(cap int) (buf *Buf) {
	return &Buf{data: make([]byte, 0, cap)}
}

func (b *Buf) EnBuf(v []byte) {
	b.mu.Lock()
	b.data = append(b.data, v...)
	b.mu.Unlock()
}

func (b *Buf) DeBuf(size int) []byte {
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
	index       int
	blockSeed   uint64
	randomStata *randnum.RandomState
}

func BasicBlockInit(rs *randnum.RandomState) *basicBlock {
	blockSeed := randnum.LCGRandom(rs)
	basicblock := &basicBlock{
		0,
		blockSeed,
		rs,
	}
	return basicblock
}

func (b *basicBlock) generateBlock(ch chan *[]uint64, buf *Buf, block *[]uint64, blockSize, blockNum int) {
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

func (b *basicBlock) writeBlock(ch chan *[]uint64, file *os.File, blockSize int) {
	block := <-ch
	data := UTB(*block)
	_, err := file.Write(*block)
	if err != nil {
		log.Fatal(err)
	}
	performinfo.IOEnd(int64(blockSize))
}
