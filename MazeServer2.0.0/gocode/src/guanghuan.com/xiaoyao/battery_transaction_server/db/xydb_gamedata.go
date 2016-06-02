package batterydb

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	xylog "guanghuan.com/xiaoyao/common/log"
	"guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

//增加游戏记录
func (db *BatteryDB) AddNewGame(game battery_run_net.Game) (err error) {
	uid := game.GetUid()
	gameid := game.GetId()

	if uid == "" || gameid == "" || gameid == "0" {
		return xyerror.ErrBadInputData
	}

	err = db.AddData(xybusiness.DB_TABLE_GAMEDATA, game)
	if err != xyerror.ErrOK {
		xylog.Error(uid, "[DB] add new game failed: %v", err)
	}

	err = xyerror.DBError(err)
	return
}

// 查询指定用户(可为空)、游戏id的游戏记录
func (db *BatteryDB) GetGame(game_id string) (game battery_run_net.Game, err error) {
	if !db.IsValidGameId(game_id) {
		err = xyerror.ErrBadInputData
		return
	}

	condition := bson.M{"id": game_id}
	selector := bson.M{"_id": 0, "xxx_unrecognized": 0}

	err = db.GetOneData(xybusiness.DB_TABLE_GAMEDATA, condition, selector, &game, mgo.Strong)
	return
}

//刷新游戏状态
func (db *BatteryDB) UpdateGame(game battery_run_net.Game) (err error) {
	game_id := game.GetId()
	if !db.IsValidGameId(game_id) {
		return xyerror.ErrBadInputData
	}

	err = db.UpdateData(xybusiness.DB_TABLE_GAMEDATA, bson.M{"id": game_id}, game)
	if err != xyerror.ErrOK {
		xylog.ErrorNoId("[DB] update game failed: %v", err)
	}

	err = xyerror.DBError(err)
	return err
}

//游戏id是否可用
// game_id string
//return:
// true 有效, false 无效
func (db *BatteryDB) IsValidGameId(game_id string) bool {
	return game_id != "" && game_id != "0"
}

//增加游戏日志
func (db *BatteryDB) AddGameLog(gl battery.GameLog) (err error) {
	return db.AddData(xybusiness.DB_TABLE_GAME_LOG, gl)
}
