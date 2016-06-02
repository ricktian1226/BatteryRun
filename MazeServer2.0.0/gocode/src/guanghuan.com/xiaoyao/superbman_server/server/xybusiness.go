// xybusiness
package xybusiness

import (
//"io"
//"os"
//"time"

//"guanghuan.com/xiaoyao/common/log"
//"guanghuan.com/xiaoyao/common/util"
)

//uid取后9位十进制整型作为transaction node和channel node的路由依据
//高5位作为transaction node路由依据
//低4位作为channel node路由依据
const BASE_UID_BAND = 10000

////pid校验失败，进程退出时间间隔
//const PIDFAILED_EXIT_SECONDS int64 = 10

//var closer io.Closer

//// InitPid 判断进程锁文件是否已经被锁，如果已经被锁，则认定已经有相同业务进程存在，本进程直接退出
//func InitPid(file string) {

//	////判断是否已经加锁
//	//locked, err := xyutil.IsLocked(file)
//	//xylog.DebugNoId("xyutil.IsLocked(%s) result : locked(%v), err(%v)", file, locked, err)
//	//if locked {
//	//	goto ErrHandle
//	//}

//	//加锁
//	var err error
//	closer, err = xyutil.Lock(file)
//	xylog.DebugNoId("xyutil.Lock(%s) result : err(%v)", file, err)
//	if err != nil {
//		goto ErrHandle
//	}

//	return

//ErrHandle:
//	xylog.ErrorNoId("InitPid failed : %v, process will exit in %d seconds", err, PIDFAILED_EXIT_SECONDS)
//	time.Sleep(time.Second * time.Duration(PIDFAILED_EXIT_SECONDS))
//	os.Exit(0)

//}
