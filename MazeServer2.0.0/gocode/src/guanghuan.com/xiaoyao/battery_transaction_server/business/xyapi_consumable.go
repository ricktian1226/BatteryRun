// xyapi_consumable
package batteryapi

import (
	//battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

//判断玩家是否拥有消耗品道具
// uid string 玩家id
// id uint64 道具id
// amount uint32 道具数目
//return:
// isExisting bool true 拥有,false 不拥有
func (api *XYAPI) IsConsumableExisting(uid string, id uint64, amount uint32) (isExisting bool) {
	isExisting = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_CONSUMABLE).IsConsumableExisting(uid, id, amount)
	return
}

//更新玩家消耗品道具
// uid string 玩家id
// id uint64 道具id
// amount uint32 道具数目
func (api *XYAPI) updateUserDataWithConsumable(uid string, id uint64, amount uint32) (err error) {
	if api.isBeforeGameRandomGoods(id) { //如果是随机道具，需要做替换操作
		err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_CONSUMABLE).UpsertRandomConsumable(uid, id, amount)
	} else if id != BEFOREGAME_RANDOM_GOODSID { //如果是普通消耗品道具，则直接增加
		err = api.GetDB(xybusiness.BUSINESS_COLLECTION_INDEX_CONSUMABLE).AddConsumable(uid, id, amount)
	}
	return
}
