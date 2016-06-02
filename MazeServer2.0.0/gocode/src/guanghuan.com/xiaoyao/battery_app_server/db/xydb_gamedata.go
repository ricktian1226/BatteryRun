package batterydb

import (
	//xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	//battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	//xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	"guanghuan.com/xiaoyao/superbman_server/server"
)

//func (db *BatteryDB) AddNewGame(game battery_run_net.Game) (err error) {
//	uid := game.GetUid()
//	gameid := game.GetId()

//	if uid == "" || gameid == "" || gameid == "0" {
//		return xyerror.ErrBadInputData
//	}

//	err = db.AddData(DB_TABLE_GAMEDATA, game)
//	if err != nil {
//		xylog.Error("[DB] add new game failed: %v", err)
//	}

//	err = xyerror.DBError(err)
//	return
//}

//// 查询指定用户(可为空)、游戏id的游戏记录
//func (db *BatteryDB) GetGame(game_id string, uid string) (game battery_run_net.Game, err error) {
//	if !db.IsValidGameId(game_id) {
//		err = xyerror.ErrBadInputData
//		return
//	}

//	var qstr interface{}
//	qstr = nil
//	if uid == "" {
//		qstr = bson.M{"id": game_id}
//	} else {
//		qstr = bson.M{"id": game_id, "uid": uid}
//	}
//	xylog.Debug("query:%v", qstr)

//	err = db.GetOneData(DB_TABLE_GAMEDATA, qstr, &game, mgo.Strong)
//	err = xyerror.DBError(err)
//	return
//}
//func (db *BatteryDB) UpdateGame(game battery_run_net.Game) (err error) {
//	game_id := game.GetId()
//	if !db.IsValidGameId(game_id) {
//		return xyerror.ErrBadInputData
//	}

//	err = db.UpdateData(DB_TABLE_GAMEDATA, bson.M{"id": game_id}, game)
//	if err != nil {
//		xylog.Debug("[DB] update game failed: %v", err)
//	}

//	err = xyerror.DBError(err)
//	return err
//}

//func (db *BatteryDB) IsValidGameId(game_id string) bool {
//	return game_id != "" && game_id != "0"
//}

// 游戏日志
func (db *BatteryDB) AddGameLog(gl battery.GameLog) (err error) {
	return db.AddData(xybusiness.DB_TABLE_GAME_LOG, gl)
}
