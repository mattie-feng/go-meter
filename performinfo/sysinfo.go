package performinfo
//package main
import (
    //"reflect"
    "sync"
    "time"
    "errors"
    "github.com/shirou/gopsutil/cpu"
    "github.com/shirou/gopsutil/mem"
    //"github.com/shirou/gopsutil/disk"
)

var(
    ErrGetTimeFailed = errors.New("Get time failed!")
)

type IOInfo struct{
    sumIO int64
    sumMB int64
    IOMutex sync.Mutex
    iops int64
    mbps int64
    nowTime int64
    nowIOPS float64
    nowMBPS float64
}

var ioinfo = IOInfo{
    sumIO: 0,
    nowTime: 0,
}

func GetState()([]float64){

    sysInfo := make([]float64,2)

    cpuPer, _ := cpu.Percent(time.Second, false)
    sysInfo[0] = cpuPer[0]
    //fmt.Println(c)
    memInfo,_ := mem.VirtualMemory()
    //fmt.Println(memInfo.UsedPercent)
    sysInfo[1] = memInfo.UsedPercent
   // parts,_ := disk.Partitions(true)
    //diskInfo,_ := disk.Usage(parts[0].Mountpoint)
    //fmt.Println(parts)
    //fmt.Println(diskInfo.UsedPercent)
    //ioc,_ := disk.IOCounters()
    //fmt.Println(ioc)

    return sysInfo

}

func IOStart(bs int64)error{

    return nil
}

func IOEnd(bs int64)error{

    nowtimeS := time.Now().Unix()
    ioinfo.IOMutex.Lock()
    ioinfo.sumIO = ioinfo.sumIO + 1
    ioinfo.sumMB = ioinfo.sumMB + bs

    if nowtimeS != ioinfo.nowTime{
        ioinfo.nowTime = nowtimeS
        ioinfo.iops = ioinfo.sumIO
        ioinfo.mbps = ioinfo.sumMB
        ioinfo.sumIO = 0
        ioinfo.sumMB = 0
    }
    ioinfo.IOMutex.Unlock()
    return nil
}
/*Get latest second IOPS. */
func GetIOps()(int64){
    nowtimeS := time.Now().Unix()
    if nowtimeS - ioinfo.nowTime > 1{
        return 0
    }
    return ioinfo.iops
}

/*Get latest second MBPS. */
func GetMBps()(int64){

    nowtimeS := time.Now().Unix()
    if nowtimeS - ioinfo.nowTime > 1{
        return 0
    }
    mbps := ioinfo.mbps/(1024*1024)
    return mbps
}
/*func main() {

    
    IOStart(64)
    IOEnd(64)
    fmt.Println(GetIOps())

}*/

/*
 //For write function:提供给写 goroutine 的接口
 func GetIOstartTime(p *iostate, bz uint32)(err)
 func GetIOendTime(p *iostate)(err)


 //for user:提供给用户端功能的接口
 func GetIOps()(int)
 func Getmbps()(int)
 func GetState()([]float)//cpu,disk,mem等状态
 */