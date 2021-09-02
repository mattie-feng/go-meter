package pipeline

import (
	"os"
	"sync"
	_ "time"
	"math"
)


//type blockData struct {
//	header map[string]interface{}
//	block []byte
//}

const MasterBlockSize = 1024

type DataFile struct {
	mutex *sync.Mutex
	wg *sync.WaitGroup
	file *os.File
}

// 单次写入
func (df *DataFile)write(blockChan chan []byte) {
	df.mutex.Lock()
	defer df.mutex.Unlock()
	content := <- blockChan
	df.file.Write(content)
	df.wg.Done()
}

func (df *DataFile)read() *[]byte{
	bs := make([]byte,4096*16)
	df.file.Read(bs)
	return &bs
}

func (df *DataFile)close(){
	df.file.Close()
}


// 处理成 block size 大小的数据块，用于单次写入，参数 rawDataSlice 存放了所有原始数据块的指针
func getBlock(rawDataSlice *[]*[]byte,chanData chan []byte,mtx *sync.Mutex,blockSize int) {
	mtx.Lock()
	defer mtx.Unlock()

	num := len(*rawDataSlice)
	dataCombined := []byte{}
	for i:=0; i<num; i++ {
		oweData := blockSize - len(dataCombined)
		if oweData > 0 {
			rawData := (*rawDataSlice)[i]
			if oweData < len(*rawData) {
				dataCombined = append(dataCombined, (*rawData)[:oweData]...)
				*rawData = (*rawData)[oweData:]
				*rawDataSlice = (*rawDataSlice)[i:]
				break
			}else {
				dataCombined = append(dataCombined, (*rawData)...)
			}
		} else {
			*rawDataSlice = (*rawDataSlice)[i:]
			break
		}
	}
	chanData <- dataCombined
}


func WriteToFile(data []*[]byte, filePath string, fileSize int, blockSize int){
	mutexWrite := &sync.Mutex{}
	mutexGetBlock := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	chanData := make(chan []byte)
	timesWrite := int(math.Ceil(float64(fileSize) / float64(blockSize))) // 向上取整

	// 生成一个要写入的文件
	fileWrite, _ := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0766)
	dataFileForWrite := DataFile{mutexWrite, wg, fileWrite}

	// 数据块处理成 block size 大小
	for i:=0; i<timesWrite ; i++{
		go getBlock(&data,chanData,mutexGetBlock,blockSize)
	}

	// 写入到文件
	for i:=0; i<timesWrite; i++{
		wg.Add(1)
		//go dataFileForWrite.write(chanData,timeChan)
		go dataFileForWrite.write(chanData)
	}
	wg.Wait()
	dataFileForWrite.close()
}