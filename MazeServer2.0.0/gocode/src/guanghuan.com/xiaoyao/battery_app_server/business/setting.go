package batteryapi

import (
	xyversion "guanghuan.com/xiaoyao/common/version"
)

var (
	DefVersion          = *xyversion.New(1, 0, 0) // 1.00.00
	DefMinClientVersion = DefVersion
	DefMaxClientVersion = DefVersion
)

const (
	GAME_OP_START  int32 = 0
	GAME_OP_END    int32 = 1
	GAME_OP_UPLOAD int32 = 2
)

const (
	ACCOUNT_OP_LOGIN            = 0
	ACCOUNT_OP_NEW              = 1
	ACCOUNT_OP_UPDATE           = 2
	ACCOUNT_OP_UPDATE_DEVICE_ID = 3
)

const (
	DIAMOND_OP_OUT  = 0  // 消耗
	DIAMOND_OP_GAIN = 1  // 非购买获得
	DIAMOND_OP_IAP  = 10 // 游戏内购买
	DIAMOND_OP_OAP  = 20 // 游戏外购买
)

//怪物id对应的下标
var DefMonsterId2IndexMap = NewMonsterId2IndexMap()
var DefMonsterIndex2IdMap = NewMonsterIndex2IdMap()

func NewMonsterId2IndexMap() (m map[string]int) {
	m = make(map[string]int, 3)
	m["20001"] = 0
	m["20002"] = 1
	m["20003"] = 2
	return
}

func NewMonsterIndex2IdMap() (m map[int]string) {
	m = make(map[int]string, 3)
	m[0] = "20001"
	m[1] = "20002"
	m[2] = "20003"
	return
}

var DefGainItemId2IndexMap = NewGainItemId2IndexMap()
var DefGainItemIndex2IdMap = NewGainItemIndex2IdMap()

func NewGainItemId2IndexMap() (m map[string]int) {
	m = make(map[string]int, 15)
	m["30000"] = 0
	m["40000"] = 1
	m["50000"] = 2
	m["50001"] = 3
	m["50002"] = 4
	m["50003"] = 5
	m["50004"] = 6
	m["50005"] = 7
	m["50006"] = 8
	m["50007"] = 9
	m["50008"] = 10
	m["60000"] = 11
	m["70000"] = 12
	m["80000"] = 13
	m["90000"] = 14
	return
}

func NewGainItemIndex2IdMap() (m map[int]string) {
	m = make(map[int]string, 15)
	m[0] = "30000"
	m[1] = "40000"
	m[2] = "50000"
	m[3] = "50001"
	m[4] = "50002"
	m[5] = "50003"
	m[6] = "50004"
	m[7] = "50005"
	m[8] = "50006"
	m[9] = "50007"
	m[10] = "50008"
	m[11] = "60000"
	m[12] = "70000"
	m[13] = "80000"
	m[14] = "90000"
	return
}
