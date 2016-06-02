package bussiness

import (
	batterydb "guanghuan.com/xiaoyao/battery_maintenance_server/db"
	xyconf "guanghuan.com/xiaoyao/common/conf"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

// // 解密反序列化请求消息
// func Before(req *http.Request, pbmsg proto.Message) (err error) {
// 	reqData, err := ioutil.ReadAll(req.Body)
// 	if err != nil {
// 		xylog.ErrorNoId("read req err:%v", err.Error())
// 		return
// 	}
// 	defer req.Body.Close()

// 	reqData, err = crypto.Decrypt(reqData)
// 	if err != nil {
// 		xylog.Error("Error decrypt :%s", err.Error())
// 	} else {
// 		err = xyencoder.PbDecode(reqData, pbmsg)
// 		if err != nil {
// 			xylog.Error("Error Unmarshal :%s", err.Error())
// 		}
// 	}
// 	xylog.DebugNoId("reqdata :%v", pbmsg)
// 	return
// }

// 加密序列化pb消息
// func After(pbmsg proto.Message) (resp string) {
// 	data, err := xyencoder.PbEncode(pbmsg)
// 	if err != xyerror.ErrOK {
// 		xylog.ErrorNoId("xyencoder.PbEncode failed : %v", err)
// 		return
// 	}
// 	//进行加密
// 	data, err = crypto.Encrypt(data)
// 	if err != xyerror.ErrOK {
// 		xylog.ErrorNoId("crypto.Encrypt failed : %v", err)
// 	}
// 	resp = string(data)
// 	return
// }

type XYAPI struct {
	platform battery.PLATFORM_TYPE
}

var apiConfigUtil xyconf.ApiConfigUtil

//创建业务实例对象
func NewXYAPI() *XYAPI {
	return &XYAPI{
		platform: battery.PLATFORM_TYPE_PLATFORM_TYPE_UNKOWN,
	}
}

func (api *XYAPI) GetCommonDB(index int) *batterydb.BatteryDB {
	dbInterface := xybusiness.DefBusinessDBSessionManager.Get(index, battery.PLATFORM_TYPE_PLATFORM_TYPE_UNKOWN)
	//xybusiness.DefBusinessDBSessionManager.Print()
	//if dbInterface == nil {
	//	xylog.ErrorNoId("DefBusinessDBSessionManager.Get index(%d) platform(%v) is nil", index, battery.PLATFORM_TYPE_PLATFORM_TYPE_UNKOWN)
	//	return nil
	//} else {
	return dbInterface.(*batterydb.BatteryDB)
	//}
}
