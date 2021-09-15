// randnum包的主要功能是：产生random number

package randnum

// 保存seed
type randomState struct {
	randomSeed interface{}
}

// 初始化Random state
func RandomInit(randSeed uint64)(*randomState) {
	randState := new(randomState)
	randState.randomSeed = randSeed
	return randState
}

// X(n+1) = (aXn + b) mod c
// lcgresult =  (randState.randomSeed.(uint64) * 25214903917 + 11) % (1 << 48)
// lcgresult的结果位数不定，所以对它取余
// 最后返回的随机数是8位的，例如：61806731
func LCGRandom(randState *randomState)(uint64) {

	// 随机算法
	lcgRand := (randState.randomSeed.(uint64) * 25214903917 + 11) % (1 << 48)
	lcgResult := lcgRand % 100000000
	if lcgResult == 0 {
		lcgResult = lcgRand / 1000000
	}

	// 对于位数不够的，右边补0
	for lcgResult < 10000000 {
		lcgResult = lcgResult * 10
	}

	return lcgResult
}