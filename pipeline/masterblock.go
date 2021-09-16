package pipeline

import (
	"go-meter/randnum"
)

const MasterBlockSize = 64 * 1024
const RandomTimes = 8 * 1024

func MasterBlockInit() *[]uint64 {
	MasterBlock := make([]uint64,0)
	for i := 0; i<RandomTimes; i++ {
		randomState := randnum.RandomInit(uint64(i))
		section := randnum.LCGRandom(randomState)
		MasterBlock = append(MasterBlock, section)
	}
	return &MasterBlock
}

