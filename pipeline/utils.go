package pipeline

import (
	"encoding/binary"
	"go-meter/randnum"
)

// 使用 random package 进行随机数的生成
func Random(mask uint64) uint64 {
	randomState := randnum.RandomInit(mask)
	num := randnum.LCGRandom(randomState)
	return num
}


// 对两个数据进行XOR操作
func XOR(A uint64, B uint64) uint64{
	return A ^ B
}


// 将数据 uint64 转化为 byte 类型的数组
func UTB(dataUint uint64) []byte{
	sliceByte := make([]byte,8)
	binary.BigEndian.PutUint64(sliceByte,dataUint)
	return sliceByte
}



func ByteBlock(masterBlock *[]uint64, masterMask,fileMask,blockMask uint64) *[]byte{
	masterSeed := Random(masterMask)
	fileSeed := Random(fileMask)
	block := make([]byte, 0, 8*1024)
	for _,i := range(*masterBlock) {
		section := UTB(XOR(XOR(XOR(i,masterSeed),fileSeed),blockMask))
		block = append(block, section...)
	}
	return &block
}