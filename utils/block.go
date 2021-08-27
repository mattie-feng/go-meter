package utils

import (
	"fmt"
	"os"
	"strconv"
)

const (
	_ = iota
	KB int = 1 << (10*iota)
	MB
	GB
	TB
)
var UNIT = []string{"KB","kb","MB","mb","GB","gb","TB","tb"}

func isSizeUnit(unitString string) bool{
	for _,unit := range UNIT {
		if unitString == unit {
			return true
		}
	}
	return false
}


func ParseSize(size string) int{
	if len(size) < 2{
		fmt.Println("输入大小有误")
		os.Exit(1)
	}
	sizeNum, _ := strconv.Atoi(size[:len(size)-2]) // 暂时先不做err判断
	sizeUnit := size[len(size)-2:]
	if !isSizeUnit(sizeUnit){
		fmt.Println("输入单位有误，请重新输入(\"KB\",\"kb\",\"MB\",\"mb\",\"GB\",\"gb\",\"TB\",\"tb\")")
	}

	switch sizeUnit {
	case "KB","kb":
		sizeNum = sizeNum * KB
	case "MB","mb":
		sizeNum = sizeNum * MB
	case "GB","gb":
		sizeNum = sizeNum * GB
	case "TB","tb":
		sizeNum = sizeNum * TB
	default:
		sizeNum = sizeNum * KB
	}


	return sizeNum

}
