// randnum包的主要功能是：
// （1）产生master mask、file mask和block mask
// （2）产生master block

package randnum

import (
	"errors"
	"crypto/sha512"
	"strconv"
	"strings"
)

// 测试数据
// var InputArgs struct {
// 	LineAge []int
// 	BlockSize string
// 	TotalSize string
// 	MasterMask int
// 	Path string
// }

// BufSize:maste block size暂定64k
// HashSize:产生MasterBlock需要调用sha512的次数
// 常量
const(
	BufSize = 65536
	HashSize = 1024
)

// xxState:0表示未完成,1表示已完成
// Fmstate:filemask state
// BmState:blockmast state
// MmState:mastermask state
// MbState:masterblock state
// EncrymMask:用于保存Encrypted MasterMask
// FileMask:用于保存FileMask
// BlockMask:用于保存BlockMask
// MasterBlock:用于保存MasterMask
type RandomState struct {
	FmState int
	BmState int
	MmState int
	MbState int
	EncrymMask int
	FileMask map[int]int
	BlockMask map[int]int
	MasterBlock []byte
}

// totalSize:由用户输入的TotalSize
// lineAge:由用户输入的lineage
// 调用random函数，根据不同的seed产生master mask
// 产生filemask和blockmask
// filemask = seed(Starting lineage ~ Ending lineage)
// blcokmask = seed(……seed(filemask))
// 返回一个randstate
func RandomInit(mseed int, totalSize string, lineAge []int) (*RandomState, error) {
	var maskErr,fileErr,blockErr error
	var randstate RandomState
	var totSize int

	randstate.FileMask = make(map[int]int)
	randstate.BlockMask = make(map[int]int)

	randstate.MmState = 0
	randstate.FmState = 0
	randstate.BmState = 0

	randstate.EncrymMask, maskErr = gmRandom(mseed)

	if maskErr != nil {
		return &randstate, errors.New("encrymMask create error")
	}else{
		randstate.MmState = 1
	}

	totSize, err := StringToint(totalSize)
	if err != nil {
		return &randstate, err
	} 

	for i := lineAge[0]; i <= lineAge[1]; i++ {
		randstate.FileMask[i], fileErr = gmRandom(i)
		// fmt.Println("i", randstate.fileMask[i], fileErr)
		blMask := randstate.FileMask[i]
		n := totSize / BufSize
		for j := 0; j < n; j++ {
			blMask, blockErr = gmRandom(blMask)
			randstate.BlockMask[i * n + j] = blMask
			// fmt.Println(i * n + j, randstate.blockMask[i * n + j], berr)
		}
	}

	if fileErr != nil && blockErr != nil {
		return &randstate, errors.New("filemask or blockmask create error")
	}else {
		randstate.FmState = 1
		randstate.BmState = 1
		return &randstate, nil
	}

}

// random算法：根据mSeed产生对应的随机数
// 随机数 = (((mSeed + 1) << 3) ^ (mSeed + 1)) % 127
// 产生随机数goMask后返回
func gmRandom(mSeed int) (int, error) {
	mSeed1 := mSeed + 1
	mSeed2 := mSeed1 << 3
	goMask := (mSeed1 ^ mSeed2) % 127

	if goMask < 0 {
		return 0, errors.New("filemask or blockmask create error")
	}else {
		return goMask, nil
	}
}

// SHA512
// 计算哈希值，返回一个长度为64的数组
func getHashcode(message string) [64]byte {
	hashBytes := sha512.Sum512([]byte(message))
	return hashBytes
}

// 产生master block:
// writeSize用于计算master block size
// random number是由sha512(0~1023)的结果拼接而成
func RandomBlock(randState *RandomState) (int, error) {
	var writeSize int
	// var hashCode bytes.Buffer
	var hashNumber [64]byte
	var masterNumber []byte
	randState.MasterBlock = make([]byte, 64)

	masterNumber = make([]byte, 64)

	randState.MbState = 0

	hashNumber = getHashcode(strconv.Itoa(0))
	copy(masterNumber,hashNumber[:])
	copy(randState.MasterBlock,hashNumber[:])

	for i := 1; i < HashSize; i++ {
		hashNumber = getHashcode(strconv.Itoa(i))
		copy(masterNumber,hashNumber[:])
		randState.MasterBlock = append(randState.MasterBlock, masterNumber...)
	}
	writeSize = len(randState.MasterBlock)

	if writeSize != BufSize {
		return writeSize, errors.New("write error")
	}else{
		randState.MbState = 1
		return writeSize, nil
	}

}

// 将String转换成int类型
// 主要是将用户输入的TotalSize转换成byte
func StringToint(strSize string)(int, error) {
	strSize2 := strings.ToUpper(strSize)
	intSize, isOk := strconv.Atoi(strSize[0:len(strSize) - 1])

	switch strSize2[len(strSize2) - 1] {
	case 'K':
		intSize = intSize * 1024
	case 'M':
		intSize = intSize * 1024 * 1024
	case 'G':
		intSize = intSize * 1024 * 1024 * 1024
	case 'T':
		intSize = intSize * 1024 * 1024 * 1024 * 1024
	default:
		return 0, errors.New("string to int error")
	}

	return intSize, isOk
}