package batteryapi

//选择角色(返回选择确认信息（直接返回）
//购买角色(返回购买确认信息（成功则扣除相关消耗：碎片、金币、宝石），失败则提示错误信息)
//获取角色列表（发送n个角色信息列表（是否选择，等级，是否已购买））
//返回游戏结算结果(在游戏结算中应该有了加进去就行)

import (
//xylog "guanghuan.com/xiaoyao/common/log"
//battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
)

//func (api *XYAPI) OperationRoleInfoConfigListRequest(req *battery.RoleInfoConfigListRequest, resp *battery.RoleInfoConfigListResponse) (err error) {

//	xylog.Debug("[RoleInfoConfigListRequest] cmd start")
//	defer xylog.Debug("[RoleInfoConfigListRequest] cmd done")
//	//xylog.Debug("req : %v", req)

//	configList := req.GetRoleInfoConfigList()
//	xylog.Debug("configList length : %d", len(configList))
//	xylog.Debug("configList : %v", configList)
//	if nil != configList {
//		//先重置所有数据为无效
//		err = api.GetDB().ResetAllRoleInfoConfigDataInvalid()
//		//添加数据
//		for _, config := range configList {
//			xylog.Debug("config : %v", config)
//			err = api.GetDB().UpsertRoleConfig(*config)
//			if err != nil {
//				return err
//			}
//		}
//	}

//	mallList := req.GetRoleLevelMallItemList()
//	if mallList != nil && len(mallList) > 0 {
//		xylog.Debug("mallList length : %d", len(mallList))
//		xylog.Debug("mallList : %v", mallList)
//		//先重置所有数据为无效
//		//添加数据
//		for _, mallItem := range mallList {
//			err = api.GetDB().AddRoleLevelInfoToGoodsTB(*mallItem)
//			xylog.Debug("err : %v", err)
//			if err != nil {
//				return err
//			}
//		}
//	}
//	return err
//}
