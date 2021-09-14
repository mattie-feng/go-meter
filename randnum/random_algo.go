// randnum包的主要功能是：产生random number

package randnum

import (
	"errors"
)

// 保存上一次的Random number
type RandomState struct {
	preRandom interface{}
}

// seed范围 0-281474976710655
// X(n+1) = (aXn + b) mod c
// lcgresult = (randState.preRandom * lcgseed) & ((1 << 48) - 1)
// lcgresult的结果位数不定，所以对它取余
// 最后返回的随机数是8位的，例如：61806731
func LCGRandom(lcgSeed uint64, randState *RandomState)(uint64, error) {

	lcgSeed = lcgSeed ^ ((1 << 48) - 1)
	if lcgSeed == 0 {
		return 0,errors.New("Seed Out of Range")
	}

	// 给randState.preRandom赋初值
	if randState.preRandom == nil {
		randState.preRandom = uint64(25214903917)
	}

	// 产生random number
	lcgResult := (randState.preRandom.(uint64) * lcgSeed) & ((1 << 48) - 1)
	lcgResultmod := lcgResult % 100000000
	if lcgResultmod == 0 {
		lcgResultmod = lcgResult / 1000000
	}

	// 对于位数不够的，右边补0
	for lcgResultmod < 10000000 {
		lcgResultmod = lcgResultmod * 10
	}

	// 用于保存本次生成的random_number
	randState.preRandom = lcgResultmod * 25214903917

	return lcgResultmod,nil
}