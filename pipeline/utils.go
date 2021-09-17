package pipeline

import (
	"encoding/binary"
)

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


func XORBlock(masterBlock *[]uint64, masterSeed,fileSeed,blockSeed uint64) *[]byte{
	block := make([]byte, 0, 8*1024)
	for _,i := range(*masterBlock) {
		section := UTB(XOR(XOR(XOR(i,masterSeed),fileSeed),blockSeed))
		block = append(block, section...)
	}
	return &block
}