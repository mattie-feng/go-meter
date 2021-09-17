package pipeline

import (
	"go-meter/randnum"
)

const MasterBlockSize = 64 * 1024
const RandomTimes = 8 * 1024

// 显示初始化 MasterBlock 使用的Seed是1
func MasterBlockInit() *[]uint64 {
	rs := randnum.RandomInit(1)
	MasterBlock := make([]uint64,0)
	for i := 0; i<RandomTimes; i++ {
		section := randnum.LCGRandom(rs)
		MasterBlock = append(MasterBlock, section)
	}
	return &MasterBlock
}

