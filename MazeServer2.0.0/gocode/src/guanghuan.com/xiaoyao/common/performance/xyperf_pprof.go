// xyperf_pprof
// 性能采集操作接口
// 无需重启服务，直接通过操作接口动态收集性能相关信息
package xyperf

import (
	"fmt"
	"guanghuan.com/xiaoyao/common/log"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"time"
)

const PPROF_SUBJECT = "pprof" //pprof nats消息subject

const (
	StrPProfOperation_Mem                = "mem"          //pprof.WriteHeapProfile
	StrPProfOperation_LookUpHeap         = "heap"         //pprof.Lookup("heap")
	StrPProfOperation_LookUpGoruntine    = "goroutine"    //pprof.Lookup("goroutine")
	StrPProfOperation_LookUpThreadCreate = "threadcreate" //pprof.Lookup("threadcreate")
	StrPProfOperation_LookUpBlock        = "block"        //pprof.Lookup("block")
	StrPProfOperation_CPUStart           = "cpustart"     //pprof.StartCpuProfile()
	StrPProfOperation_CPUStop            = "cpustop"      //pprof.StopCpuProfile()
	StrPProfOperation_FreeMem            = "freemem"      //pprof.StopCpuProfile()
)

var (
	cpuFd, memFd, heapFd, goroutineFd, otherFd *os.File //cpu/内存/堆内存/其他 的文件句柄
	cpuFile, memFile                           string   //cpu/内存采集信息的文件名
)

//校验下文件句柄
func checkFd(fd **os.File) {
	if (*fd) != nil {
		(*fd).Close()
		(*fd) = nil
	}
}

//初始化pprof配置信息文件句柄
func InitPProf(path, appName string, dcId, nodeId int) (err error) {

	//cpu 统计信息
	cpuFile = fmt.Sprintf("%s/pprof_cpu_%s_%d_%d.log", path, appName, dcId, nodeId)

	//内存统计信息
	memFile = fmt.Sprintf("%s/pprof_mem_%s_%d_%d.log", path, appName, dcId, nodeId)

	//对内存统计信息
	checkFd(&heapFd)
	file := fmt.Sprintf("%s/pprof_heap_%s_%d_%d.log", path, appName, dcId, nodeId)
	if heapFd, err = os.Create(file); err != nil {
		xylog.ErrorNoId("create %s failed : %v", file, err)
		return
	}

	//对goroutine统计信息
	checkFd(&goroutineFd)
	file = fmt.Sprintf("%s/pprof_goruntine_%s_%d_%d.log", path, appName, dcId, nodeId)
	if goroutineFd, err = os.Create(file); err != nil {
		xylog.ErrorNoId("create %s failed : %v", file, err)
		return
	}

	//其他统计信息：heap/threadcreate/block
	checkFd(&otherFd)
	file = fmt.Sprintf("%s/pprof_other_%s_%d_%d.log", path, appName, dcId, nodeId)
	if otherFd, err = os.Create(file); err != nil {
		xylog.ErrorNoId("create %s failed : %v", file, err)
		return
	}

	return
}

//性能统计操作函数
//func OperationPProf(operation int) {
func OperationPProf(op string) {

	switch op {

	case StrPProfOperation_Mem:
		xylog.InfoNoId("pprof mem")
		//runtime.GC() //先强制做下垃圾回收
		//n, err := memFd.WriteString(fmt.Sprintf("%s WriteHeapProfile\n", time.Now().String())) //写入时间戳信息
		//if n <= 0 || err != nil {
		//	xylog.Error("write to pprof mem file failed : n(%d), err(%v)", n, err)
		//	return
		//}
		checkFd(&memFd)
		_, err := os.Stat(memFile)
		if err == nil || os.IsExist(err) { //文件存在，则改名
			fileNew := fmt.Sprintf("%s.%s", memFile, time.Now().Format("2006-01-06_15_04_05"))
			err = os.Rename(memFile, fileNew)
			if err != nil {
				xylog.ErrorNoId("rename %s to %s failed : %v", memFile, fileNew, err)
				return
			}
		}

		if memFd, err = os.Create(memFile); err != nil {
			xylog.ErrorNoId("create %s failed : %v", memFile, err)
			return
		}
		defer checkFd(&memFd)

		err = pprof.WriteHeapProfile(memFd) //写入heap信息
		if err != nil {
			xylog.ErrorNoId("pprof.WriteHeapProfile(memFd) failed : err(%v)", err)
			return
		}

	case StrPProfOperation_LookUpHeap:
		xylog.InfoNoId("pprof lookup heap")
		runtime.GC() //先强制做下垃圾回收
		p := pprof.Lookup("heap")

		n, err := heapFd.WriteString(fmt.Sprintf("%s lookup heap\n", time.Now().String())) //写入时间戳信息
		if n <= 0 || err != nil {
			xylog.ErrorNoId("write heap to pprof mem file failed : n(%d), err(%v)", n, err)
			return
		}

		err = p.WriteTo(heapFd, 2)
		if err != nil {
			xylog.ErrorNoId("write heap to pprof mem file failed : err(%v)", err)
			return
		}

	case StrPProfOperation_LookUpGoruntine:
		xylog.InfoNoId("pprof lookup goroutine")

		p := pprof.Lookup("goroutine")

		n, err := goroutineFd.WriteString(fmt.Sprintf("%s lookup goroutine\n", time.Now().String())) //写入时间戳信息
		if n <= 0 || err != nil {
			xylog.ErrorNoId("write goroutine to pprof goroutine file failed : n(%d), err(%v)", n, err)
			return
		}

		err = p.WriteTo(goroutineFd, 2)
		if err != nil {
			xylog.ErrorNoId("write goroutines to pprof goroutine file failed : err(%v)", err)
			return
		}

	case StrPProfOperation_LookUpThreadCreate:
		xylog.InfoNoId("pprof lookup threadcreate")

		p := pprof.Lookup("threadcreate")

		n, err := otherFd.WriteString(fmt.Sprintf("%s lookup threadcreate\n", time.Now().String())) //写入时间戳信息
		if n <= 0 || err != nil {
			xylog.ErrorNoId("write threadcreate to pprof other file failed : n(%d), err(%v)", n, err)
			return
		}

		err = p.WriteTo(otherFd, 2)
		if err != nil {
			xylog.ErrorNoId("write threadcreate to pprof other file failed : err(%v)", err)
			return
		}

	case StrPProfOperation_LookUpBlock:
		xylog.InfoNoId("pprof lookup block")
		p := pprof.Lookup("block")
		n, err := otherFd.WriteString(fmt.Sprintf("%s lookup block\n", time.Now().String())) //写入时间戳信息
		if n <= 0 || err != nil {
			xylog.ErrorNoId("write block to pprof other file failed : n(%d), err(%v)", n, err)
			return
		}

		err = p.WriteTo(otherFd, 2)
		if err != nil {
			xylog.ErrorNoId("write block to pprof other file failed : err(%v)", err)
			return
		}

	case StrPProfOperation_CPUStart: //每次采集都会生成新的文件。避免多次采集的内容混在一起
		xylog.InfoNoId("pprof cpu start")
		checkFd(&cpuFd)
		_, err := os.Stat(cpuFile)
		if err == nil || os.IsExist(err) { //文件存在，则改名
			fileNew := fmt.Sprintf("%s.%s", cpuFile, time.Now().Format("2006-01-06_15_04_05"))
			err = os.Rename(cpuFile, fileNew)
			if err != nil {
				xylog.ErrorNoId("rename %s to %s failed : %v", cpuFile, fileNew, err)
				return
			}
		}

		if cpuFd, err = os.Create(cpuFile); err != nil {
			xylog.ErrorNoId("create %s failed : %v", cpuFile, err)
			return
		}

		err = pprof.StartCPUProfile(cpuFd)
		if err != nil {
			xylog.ErrorNoId("pprof.StartCPUProfile failed : err(%v)", err)
			return
		}

	case StrPProfOperation_CPUStop:
		xylog.InfoNoId("pprof cpu stop")
		//n, err := otherFd.WriteString(fmt.Sprintf("%s pprof cpu stop\n", time.Now().String())) //写入时间戳信息
		//if n <= 0 || err != nil {
		//	xylog.Error("write pprof cpu start to pprof cpu file failed : n(%d), err(%v)", n, err)
		//	return
		//}

		pprof.StopCPUProfile()
		checkFd(&cpuFd)
	case StrPProfOperation_FreeMem:
		xylog.InfoNoId("pprof free os memory")
		debug.FreeOSMemory()

	default:
		xylog.ErrorNoId("Unkown PProfOperation : %d", op)
	}

	return
}

//关闭文件句柄
func FiniPProf() {
	checkFd(&cpuFd)
	checkFd(&memFd)
	checkFd(&otherFd)
}
