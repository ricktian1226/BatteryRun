// bussiness
package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	proto "code.google.com/p/goprotobuf/proto"
	//apns "github.com/anachronistic/apns"
	httplib "github.com/astaxie/beego/httplib"

	//"guanghuan.com/xiaoyao/common/apn"
	xyencoder "guanghuan.com/xiaoyao/common/encoding"
	xylog "guanghuan.com/xiaoyao/common/log"
	//xynatsservice "guanghuan.com/xiaoyao/common/service/nats"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	crypto "guanghuan.com/xiaoyao/superbman_server/crypto"
	"guanghuan.com/xiaoyao/superbman_server/error"
	//"guanghuan.com/xiaoyao/superbman_server/server"
	"guanghuan.com/xiaoyao/common/service/timer"

	//apnsserver "guanghuan.com/xiaoyao/battery_apns_server/business"
	transactionbusiness "guanghuan.com/xiaoyao/battery_transaction_server/business"
)

const (
	CLIENT_VERSION  = 20000
	REQUEST_TIMEOUT = 2 * 1000 * 1000
)

type SPECIFIC_FUNC func(string) error

//网络交互函数
func test_op(req proto.Message, resp proto.Message, s string, rwtimeout time.Duration) (err error) {
	var (
		data      []byte
		api_url   string = "http://" + DefConfig.Url + s
		timeBegin        = time.Now()
	)

	data, err = xyencoder.PbEncode(req)
	if err != nil {
		xylog.Error("Error Encoding: %s", err.Error())
		return err
	}
	//xylog.Debug("reqdata after PbEncode: [%d][%v]\n", len(data), data)

	data, err = crypto.Encrypt(data)
	if err != nil {
		xylog.Error("Error Encrypt: %s", err.Error())
		return err
	}
	//xylog.Debug("reqdata after Encrypt: [%d][%v]\n", len(data), data)

	//fmt.Println(api_url)
	request := httplib.Post(api_url)
	request.Body(data)
	request.SetTimeout(time.Duration(DefConfig.ConnTimeout)*time.Second, time.Duration(DefConfig.RwTimeout)*time.Second)

	//timeBegin1 := time.Now()
	data, err = request.Bytes()

	xylog.DebugNoId("error : %v", err)
	if err != nil {
		xylog.Error("Error Send: %s", err.Error())
		return
	}

	//time.Sleep(240 * time.Second)

	//xylog.Debug("request.Bytes cost %d ms", time.Since(timeBegin1).Nanoseconds()/int64(time.Millisecond))
	//xylog.Debug("respdata after Send: [%d][%v]\n", len(data), data)

	data, err = crypto.Decrypt(data)
	if err != nil {
		xylog.Error("Error Decrypt: %s", err.Error())
		return
	}
	//xylog.Debug("respdata after Decrypt: [%d][%v]\n", len(data), data)

	err = xyencoder.PbDecode(data, resp)
	if err != nil {
		xylog.Error("Error Decoding: %s", err.Error())
		return err
	}

	if err == nil {
		xylog.DebugNoId("response : %v", resp)
	} else {
		xylog.ErrorNoId("response : %v, error : %v", resp, err)
	}
	//xylog.Debug("test_op end %d nano seconds", time.Now().UnixNano())

	xylog.DebugNoId("test_op cost %d ms", time.Since(timeBegin).Nanoseconds()/int64(time.Millisecond))

	return err
}

//性能测试公共函数
func test_pprof_op(function SPECIFIC_FUNC, done, errSum *uint64, id, sum uint64) (err error) {
	c := make(chan bool, DefConfig.FlowControl)
	for i := uint64(0); i < sum; i++ {
		c <- true
		j := i + id*sum

		go func() {
			//fmt.Printf("#%d result : %v\n", j, test_pprof_sub_op(function, done, errSum, j))
			test_pprof_sub_op(function, done, errSum, j)
			<-c
		}()

	}

	return
}

func test_pprof_sub_op(function SPECIFIC_FUNC, done, errSum *uint64, j uint64) (err error) {
	sid := fmt.Sprintf("sina_weibo_%d", j)
	var uid string
	uid, err = test_specific_login(sid,
		fmt.Sprintf("nameaaaa_%d", j),
		"cb7a0d0b5cd15a4b3128e3c015d92c47f72661493fae12e6cd03082623992213",
		fmt.Sprintf("iconUrl_aaaa%d", j),
		int32(CLIENT_VERSION),
		battery.ID_SOURCE_SRC_SINA_WEIBO)
	if err != nil {
		xylog.ErrorNoId("%v", err)
		*errSum++
		return
	} else {
		*done++
	}

	xylog.DebugNoId("uid(%s)  %v", uid, function)

	if function != nil {
		err = function(uid)
		if err != nil {
			xylog.ErrorNoId("%v", err)
			*errSum++
			return
		} else {
			*done++
		}
	}

	return
}

//打印性能测试进度函数
func test_print_progress(done, errSum *uint64, sum uint64) {
	timer := time.NewTicker(time.Minute)
	for {
		select {
		case <-timer.C:
			fmt.Printf("%s %d/%d/%d\n", time.Now().String(), *done, *errSum, sum)
		}
	}
	return
}

func test_login() (err error) {
	// transaction 7
	//sid := "5607087217"
	//name := "地瓜999牌"
	//devId := ""
	//iconUrl := "http://tp2.sinaimg.cn/5607087217/50/5729564554/1"

	// transaction 5
	//sid := "2844846670"
	//name := "ricktian"
	//devId := "01ca006c02d7bd6cedf4f91e98796ef64eb2600f2ed34e4585f8bedcd4477beb"
	//iconUrl := "http://tp3.sinaimg.cn/2844846670/50/5635390186/1"

	sid := "guest1"
	name := "bitch"
	devId := "01ca006c02d7bd6cedf4f91e98796ef64eb2600f2ed34e4585f8bedcd4477beb"
	iconUrl := ""
	//source := battery.ID_SOURCE_SRC_SINA_WEIBO
	source := battery.ID_SOURCE_SRC_TOKEN
	version := int32(20000)
	_, err = test_specific_login(sid, name, devId, iconUrl, version, source)

	return
}

func test_specific_login(sid, name, devId, iconUrl string, version int32, source battery.ID_SOURCE) (uid string, err error) {
	time.Sleep(time.Millisecond * 2)

	req := &battery.LoginRequest{
		LoginId: &battery.TPID{
			Id:     proto.String(sid),
			Source: source.Enum(),
			Name:   proto.String(name),
		},
		DeviceId: &battery.TPID{
			Id: proto.String(devId),
		},
		IconUrl:      proto.String(iconUrl),
		Version:      proto.Int32(version),
		PlatformType: battery.PLATFORM_TYPE_PLATFORM_TYPE_IOS.Enum(),
	}

	resp := &battery.LoginResponse{}

	err = test_op(req, resp, API_URI_LOGIN, REQUEST_TIMEOUT)

	//xylog.DebugNoId("resp : %v", resp)

	return resp.GetData().GetUid(), err
}

func test_pprof_login(done, errSum *uint64, id, sum uint64) (err error) {
	err = test_pprof_op(nil, done, errSum, id, sum)
	return
}

func test_bind() error {
	sid := "2844846671"
	name := "ricktian1"
	uid := "1442828654322275370081"
	source := battery.ID_SOURCE_SRC_SINA_WEIBO
	platform := battery.PLATFORM_TYPE_PLATFORM_TYPE_IOS
	return test_specific_bind(sid, name, uid, source, platform)
}

func test_specific_bind(sid, name, uid string, source battery.ID_SOURCE, platform battery.PLATFORM_TYPE) (err error) {
	req := &battery.BindRequest{
		Uid: proto.String(uid),
		Target: &battery.TPID{
			Id:     proto.String(sid),
			Source: source.Enum(),
			Name:   proto.String(name),
		},
		PlatformType: platform.Enum(),
	}

	resp := &battery.BindResponse{}

	err = test_op(req, resp, API_URI_BIND, REQUEST_TIMEOUT)

	xylog.DebugNoId("resp : %v", resp)

	return
}

func test_announcement_querybytime() (err error) {
	err = test_specific_announcement_querybytime("1433213960479312378820")
	return
}

func test_specific_announcement_querybytime(uid string) (err error) {
	cmd := battery.OP_CMD_Query
	req := &battery.AnnouncementRequest{
		Uid: proto.String(uid),
		Cmd: &cmd,
	}

	resp := &battery.AnnouncementResponse{}

	err = test_op(req, resp, API_URI_ANNOUNCEMENT, REQUEST_TIMEOUT)

	return
}

func test_pprof_announcement_querybytime(done, errSum *uint64, id, sum uint64) (err error) {
	err = test_pprof_op(test_specific_announcement_querybytime, done, errSum, id, sum)
	return
}

//todo iap 性能测试
func test_iap() (err error) {
	uids := []string{"139752879435073500001621"}
	for i, n := 0, len(uids); i < n; i++ {
		err = test_iapSub(uids[i])
	}
	return
}

func test_iapSub(uid string) (err error) {

	req := &battery.OrderVerifyRequest{
		//OrderId:     proto.String(fmt.Sprintf("%s.%04d", uid, i)),
		OrderId:     proto.String("123456"),
		Uid:         proto.String(uid),
		ReceiptData: proto.String("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"),
	}

	xylog.Debug("orderid : %s", "123456")

	resp := &battery.OrderVerifyResponse{}
	err = test_op(req, resp, API_URI_IAP_VERIFY, REQUEST_TIMEOUT)

	return
}

func test_query_goods() (err error) {
	err = test_specific_query_goods("1444917620835502213059")
	return
}

func test_pprof_query_goods(done, errSum *uint64, id, sum uint64) (err error) {
	err = test_pprof_op(test_specific_query_goods, done, errSum, id, sum)
	return
}

func test_specific_query_goods(uid string) (err error) {
	req := &battery.QueryGoodsRequest{
		Uid:         proto.String(uid),
		MallType:    battery.MallType_Mall_Main.Enum(),
		MallSubType: battery.MallSubType_MallSubType_Recommend.Enum(),
	}

	resp := &battery.QueryGoodsResponse{}

	err = test_op(req, resp, API_URI_GOODS_QUERY, REQUEST_TIMEOUT)

	return
}

func test_buy_good() (err error) {
	err = test_specific_buy_good("1445257345918880074081")
	return
}

func test_pprof_buy_good(done, errSum *uint64, id, sum uint64) (err error) {
	err = test_pprof_op(test_specific_buy_good, done, errSum, id, sum)
	return
}

func test_specific_buy_good(uid string) (err error) {
	req := &battery.BuyGoodsRequest{
		Uid:     proto.String(uid),
		GoodsId: proto.Uint64(110020000),
	}
	resp := &battery.BuyGoodsResponse{}

	err = test_op(req, resp, API_URI_GOODS_BUY, REQUEST_TIMEOUT)

	return
}

//func test_giftquery() (err error) {
//	req := &battery.QueryStaminaGiftRequest{
//		Uid: proto.String("1400911655578417395161"),
//	}

//	xylog.Debug("gift query uid(%s)", req.GetUid())

//	resp := &battery.QueryStaminaGiftResponse{}

//	err = test_op(req, resp, API_URI_GIFT_QUERY, 1000*1000)

//	return
//}

func test_query_frienddata() (err error) {
	err = test_query_specific_frienddata("1445257345918880074081")
	return
}

func test_query_specific_frienddata(uid string) (err error) {
	src := battery.ID_SOURCE_SRC_SINA_WEIBO

	req := &battery.QueryFriendsDataRequest{
		Uid:    proto.String(uid),
		Source: &src,
	}

	req.Sids = make([]string, 0)
	req.Sids = append(req.Sids, "sina_weibo_0")
	//req.Sids = append(req.Sids, "2066228185")
	//req.Sids = append(req.Sids, "2844846670")

	resp := &battery.QueryFriendsDataResponse{}
	err = test_op(req, resp, API_URI_GET_FRIEND_GAMEDATA, 1000*1000)

	return
}

func test_pprof_query_frienddata(done, errSum *uint64, id, sum uint64) (err error) {
	err = test_pprof_op(test_query_specific_frienddata, done, errSum, id, sum)
	return
}

func test_prop_res_request() (err error) {
	req := &battery.QueryPropResRequest{
		Uid: proto.String("1411695993946224940940"),
	}
	resp := &battery.QueryPropResResponse{}
	err = test_op(req, resp, API_URI_PROP_RES_QUERY, 1000*1000)
	return
}

func test_sys_lotto() (err error) {
	var (
		uid     = "1445257345918880074081"
		lottoid uint64
	)

	err = test_specific_sys_lotto_request_initial(uid, &lottoid)
	if err != xyerror.ErrOK {
		return
	}

	err = test_specific_sys_lotto_request_commit(uid, &lottoid)
	return
}

func test_sys_lotto_request_initial(uid string, lottoid *uint64) (err error) {
	req := &battery.LottoRequest{
		Uid:              proto.String(uid),
		Cmd:              battery.LottoCmd_Lotto_Initial.Enum(),
		ForceRefreshSlot: proto.Bool(true),
	}
	resp := &battery.LottoResponse{}
	err = test_op(req, resp, API_URI_LOTTO_OP, 1000*1000)
	*lottoid = resp.Stuff.GetLottoid()
	return
}

func test_specific_sys_lotto_request_initial(uid string, lottoid *uint64) (err error) {
	cmd := battery.LottoCmd_Lotto_Initial
	req := &battery.LottoRequest{
		Uid: proto.String(uid),
		Cmd: &cmd,
		//SerialNum: proto.Int32(1),
		ForceRefreshSlot: proto.Bool(true),
	}
	resp := &battery.LottoResponse{}
	err = test_op(req, resp, API_URI_LOTTO_OP, 1000*1000)
	*lottoid = resp.Stuff.GetLottoid()
	return
}

func test_specific_sys_lotto_request_commit(uid string, lottoid *uint64) (err error) {
	cmd := battery.LottoCmd_Lotto_Commit
	req := &battery.LottoRequest{
		Uid:           proto.String(uid),
		Cmd:           &cmd,
		Lottoid:       proto.Uint64(*lottoid),
		Parentlottoid: proto.Uint64(*lottoid),
	}
	resp := &battery.LottoResponse{}
	err = test_op(req, resp, API_URI_LOTTO_OP, 1000*1000)
	return
}

func test_pprof_sys_lotto(done, errSum *uint64, id, sum uint64) (err error) {
	beginId, uid, lottoid := id*sum, "", uint64(0)
	c := make(chan bool, DefConfig.FlowControl)
	for i := uint64(0); i < sum; i++ {
		j := i + beginId
		c <- true

		go func() {

			sid := fmt.Sprintf("sina_weibo_%d", j)
			uid, err = test_specific_login(sid,
				fmt.Sprintf("name_%d", j),
				fmt.Sprintf("devId_%d", j),
				fmt.Sprintf("iconUrl_%d", j),
				int32(CLIENT_VERSION),
				battery.ID_SOURCE_SRC_SINA_WEIBO)
			if err != xyerror.ErrOK {
				*errSum++
				<-c
				return
			} else {
				*done++
			}

			err = test_sys_lotto_request_initial(uid, &lottoid)
			if err != xyerror.ErrOK {
				*errSum++
				<-c
				return
			} else {
				*done++
			}

			err = test_specific_sys_lotto_request_commit(uid, &lottoid)
			if err != xyerror.ErrOK {
				*errSum++
				<-c
				return
			} else {
				*done++
			}
			<-c
		}()
	}

	return
}

func test_new_game() (err error) {
	err = test_specific_new_game("1431330472754676933059")
	return
}

func test_specific_new_game(uid string) (err error) {
	req := &battery.NewGameRequest{
		Uid:          proto.String(uid),
		Type:         battery.GameType_GameType_MainLine.Enum(),
		CheckPointId: proto.Uint32(1),
	}
	resp := &battery.NewGameResponse{}
	err = test_op(req, resp, API_URI_NEWGAME, 1000*1000)
	return
}

func test_specific_new_game2(uid string) (gameId string, err error) {
	req := &battery.NewGameRequest{
		Uid:          proto.String(uid),
		Type:         battery.GameType_GameType_MainLine.Enum(),
		CheckPointId: proto.Uint32(3),
	}
	resp := &battery.NewGameResponse{}
	err = test_op(req, resp, API_URI_NEWGAME, 1000*1000)
	if err == xyerror.ErrOK {
		gameId = resp.GetGameId()
	}

	return
}

func test_pprof_new_game(done, errSum *uint64, id, sum uint64) (err error) {
	err = test_pprof_op(test_specific_new_game, done, errSum, id, sum)
	return
}

func test_game_result2() (err error) {
	err = test_specific_game_result2("1421043649700167835847", "1422414311315848796847")
	return
}

func test_specific_game_result2(uid, gameId string) (err error) {
	gameResult := &battery.GameResult{
		Quotas: make([]*battery.Quota, 0),
	}

	gameResult.Quotas = append(gameResult.Quotas,
		&battery.Quota{
			Id:    battery.QuotaEnum_Quota_Score.Enum(),
			Value: proto.Uint64(1000),
		},

		&battery.Quota{
			Id:    battery.QuotaEnum_Quota_Charge.Enum(),
			Value: proto.Uint64(10),
		},

		&battery.Quota{
			Id:    battery.QuotaEnum_Quota_Coin.Enum(),
			Value: proto.Uint64(10),
		},

		&battery.Quota{
			Id:    battery.QuotaEnum_Quota_AllCheckPointStar.Enum(),
			Value: proto.Uint64(1),
		},
		&battery.Quota{
			Id:    battery.QuotaEnum_Quota_Convey.Enum(),
			Value: proto.Uint64(10),
		},
		&battery.Quota{
			Id:    battery.QuotaEnum_Quota_Stamina.Enum(),
			Value: proto.Uint64(2),
		})

	collections := []uint32{3}

	pickUps := []*battery.PropItem{&battery.PropItem{
		Type:   battery.PropType_PROP_JIGSAW.Enum(),
		Amount: proto.Uint32(2),
	}}

	req := &battery.GameResultCommitRequest{
		Uid:          proto.String(uid),
		Type:         battery.GameType_GameType_MainLine.Enum(),
		GameId:       proto.String(gameId),
		Duration:     proto.Int64(100),
		GameResult:   gameResult,
		RoleId:       proto.Uint64(150030004),
		Collections:  collections,
		Pickups:      pickUps,
		CheckPointId: proto.Uint32(3),
		IsFinish:     proto.Bool(true),
	}

	resp := &battery.GameResultCommitResponse{}

	err = test_op(req, resp, API_URI_ADD_GAMEDATA2, 5*1000*1000)

	return
}

func test_pprof_game(done, errSum *uint64, id, sum uint64) (err error) {
	//beginId, uid, gameId := id*sum, "", ""
	c := make(chan bool, DefConfig.FlowControl)
	for i := uint64(0); i < sum; i++ {
		c <- true
		go func() {
			test_pprof_sub_game(done, errSum, i, id, sum)
			<-c
		}()
		//sid := fmt.Sprintf("sina_weibo_%d", j)
		//uid, err = test_specific_login(sid, fmt.Sprintf("name_%d", j), fmt.Sprintf("devId_%d", j), fmt.Sprintf("iconUrl_%d", j), int32(CLIENT_VERSION))
		//if err != xyerror.ErrOK {
		//	*errSum++
		//	continue
		//} else {
		//	*done++
		//}

		//gameId, err = test_specific_new_game2(uid)
		//if err != xyerror.ErrOK {
		//	*errSum++
		//	continue
		//} else {
		//	*done++
		//}

		//err = test_specific_game_result2(uid, gameId)
		//if err != xyerror.ErrOK {
		//	*errSum++
		//	continue
		//} else {
		//	*done++
		//}

	}

	return
}

func test_pprof_sub_game(done, errSum *uint64, j, id, sum uint64) (err error) {
	j += id * sum
	sid := fmt.Sprintf("sina_weibo_%d", j)
	var uid string
	uid, err = test_specific_login(sid,
		fmt.Sprintf("name_%d", j),
		fmt.Sprintf("devId_%d", j),
		fmt.Sprintf("iconUrl_%d", j),
		int32(CLIENT_VERSION),
		battery.ID_SOURCE_SRC_SINA_WEIBO)
	if err != xyerror.ErrOK {
		*errSum++
		return
	} else {
		*done++
	}

	var gameId string
	gameId, err = test_specific_new_game2(uid)
	if err != xyerror.ErrOK {
		*errSum++
		return
	} else {
		*done++
	}

	err = test_specific_game_result2(uid, gameId)
	if err != xyerror.ErrOK {
		*errSum++
		return
	} else {
		*done++
	}
	return
}

func test_query_user_mission() (err error) {
	err = test_query_specific_user_mission("1425352623264100099847")
	return
}

func test_query_specific_user_mission(uid string) (err error) {
	req := &battery.QueryUserMissionRequest{
		Uid:               proto.String(uid),
		Types:             []battery.MissionType{battery.MissionType_MissionType_Daily, battery.MissionType_MissionType_Study, battery.MissionType_MissionType_MainLine},
		MissionCount:      proto.Uint32(3),
		DailyMissionCount: proto.Uint32(1),
	}

	resp := &battery.QueryUserMissionResponse{}

	err = test_op(req, resp, API_URI_QUERY_USER_MISSION, 1000*1000)

	return
}

func test_pprof_query_user_mission(done, errSum *uint64, id, sum uint64) (err error) {
	err = test_pprof_op(test_query_specific_user_mission, done, errSum, id, sum)
	return
}

func test_confirm_user_mission() (err error) {
	err = test_confirm_specific_user_mission("1425352623264100099847", 76340294420791296)
	return
}

func test_confirm_specific_user_mission(uid string, mid uint64) (err error) {
	req := &battery.ConfirmUserMissionRequest{
		Uid:  proto.String(uid),
		Type: battery.MissionType_MissionType_MainLine.Enum(),
		Mid:  proto.Uint64(mid),
	}

	resp := &battery.ConfirmUserMissionResponse{}

	err = test_op(req, resp, API_URI_CONFIRM_USER_MISSION, 1000*1000)

	return
}

func test_pprof_confirm_user_mission(done, errSum *uint64, id, sum uint64) (err error) {
	//beginId, uid := id*sum, ""
	c := make(chan bool, DefConfig.FlowControl)
	for i := uint64(0); i < sum; i++ {
		c <- true
		go func() {
			test_pprof_sub_confirm_user_mission(done, errSum, i, id, sum)
			<-c
		}()
		//j := i + beginId
		//sid := fmt.Sprintf("sina_weibo_%d", j)
		//uid, err = test_specific_login(sid, fmt.Sprintf("name_%d", j), fmt.Sprintf("devId_%d", j), fmt.Sprintf("iconUrl_%d", j), int32(CLIENT_VERSION))
		//if err != xyerror.ErrOK {
		//	*errSum++
		//	continue
		//} else {
		//	*done++
		//}

		//var missions []battery.UserMission
		//missions, err = test_query_specific_user_mission_result(uid)
		//if err != xyerror.ErrOK {
		//	*errSum++
		//	continue
		//} else {
		//	*done++
		//	xylog.DebugNoId("done missions : %v", missions)
		//	for _, m := range missions {
		//		err = test_confirm_specific_user_mission(uid, m.GetMid())
		//		if err != xyerror.ErrOK {
		//			*errSum++
		//			continue
		//		} else {
		//			*done++
		//		}
		//	}
		//}

	}

	return
}

func test_pprof_sub_confirm_user_mission(done, errSum *uint64, i, id, sum uint64) (err error) {
	j := i + id*sum
	sid := fmt.Sprintf("sina_weibo_%d", j)
	var uid string
	uid, err = test_specific_login(sid,
		fmt.Sprintf("name_%d", j),
		fmt.Sprintf("devId_%d", j),
		fmt.Sprintf("iconUrl_%d", j),
		int32(CLIENT_VERSION),
		battery.ID_SOURCE_SRC_SINA_WEIBO)
	if err != xyerror.ErrOK {
		xylog.ErrorNoId("%v", err)
		*errSum++
		return
	} else {
		*done++
	}

	//var missions []battery.UserMission
	//missions, err = test_query_specific_user_mission_result(uid)
	_, err = test_query_specific_user_mission_result(uid)
	if err != xyerror.ErrOK {
		xylog.ErrorNoId("%v", err)
		*errSum++
		return
	} /*else {
		*done++
		xylog.DebugNoId("done missions : %v", missions)
		for _, m := range missions {
			err = test_confirm_specific_user_mission(uid, m.GetMid())
			if err != xyerror.ErrOK {
				xylog.ErrorNoId("%v", err)
				*errSum++
				continue
			} else {
				*done++
			}
		}
	}*/

	return
}

func test_query_specific_user_mission_result(uid string) (missions []battery.UserMission, err error) {
	req := &battery.QueryUserMissionRequest{
		Uid:               proto.String(uid),
		Types:             []battery.MissionType{battery.MissionType_MissionType_Daily, battery.MissionType_MissionType_Study, battery.MissionType_MissionType_MainLine},
		MissionCount:      proto.Uint32(3),
		DailyMissionCount: proto.Uint32(1),
	}

	resp := &battery.QueryUserMissionResponse{}

	err = test_op(req, resp, API_URI_QUERY_USER_MISSION, 1000*1000)

	for _, entry := range resp.GetEntrys() {
		for _, m := range entry.GetDoneNotCollect() {
			missions = append(missions, *m)
		}
	}

	return
}

func test_query_specific_user_signin_activity(uid string) (err error) {
	req := &battery.QuerySignInRequest{
		Uid: proto.String(uid),
	}

	resp := &battery.QuerySignInResponse{}

	err = test_op(req, resp, API_URI_SIGNIN_ACTIVITY_QUERY, 1000*1000)

	return
}

func test_query_user_signin_activity() (err error) {

	err = test_query_specific_user_signin_activity("1425353626119417911485")

	return
}

func test_pprof_query_user_signin_activity(done, errSum *uint64, id, sum uint64) (err error) {
	err = test_pprof_op(test_query_specific_user_signin_activity, done, errSum, id, sum)
	return
}

func test_specific_user_signin(uid string) (err error) {
	req := &battery.SignInRequest{
		Uid: proto.String(uid),
		Id:  proto.Uint64(1),
	}

	resp := &battery.SignInResponse{}

	err = test_op(req, resp, API_URI_SIGNIN, 1000*1000)

	return
}

func test_user_signin() (err error) {
	err = test_specific_user_signin("1425353626119417911485")
	return
}

func test_pprof_user_signin(done, errSum *uint64, id, sum uint64) (err error) {
	err = test_pprof_op(test_specific_user_signin, done, errSum, id, sum)
	return
}

func test_specific_query_checkpoints(uid string) (err error) {
	req := &battery.QueryUserCheckPointsRequest{
		Uid:     proto.String(uid),
		BeginId: proto.Uint32(0),
		EndId:   proto.Uint32(30),
	}

	resp := &battery.QueryUserCheckPointsResponse{}

	err = test_op(req, resp, API_URI_CHECKPOINT_QUERY_RANGE, 20*1000*1000)

	return
}

func test_query_checkpoints() (err error) {
	err = test_specific_query_checkpoints("1427266457542935645408")
	return
}

func test_pprof_query_checkpoints(done, errSum *uint64, id, sum uint64) (err error) {
	return test_pprof_op(test_specific_query_checkpoints, done, errSum, id, sum)
}

func test_specific_query_checkpoint_friend_rank(uid string) (err error) {
	return test_specific_checkpoint_rank(uid, battery.CheckPointRankType_CheckPointRankType_Friend)
}

func test_query_checkpoint_friend_rank() (err error) {
	return test_specific_query_checkpoint_friend_rank("1427266457542935645408")
}

func test_pprof_query_checkpoint_friend_rank(done, errSum *uint64, id, sum uint64) (err error) {
	return test_pprof_op(test_specific_query_checkpoint_friend_rank, done, errSum, id, sum)
}

func test_specific_query_checkpoint_global_rank(uid string) (err error) {
	return test_specific_checkpoint_rank(uid, battery.CheckPointRankType_CheckPointRankType_Global)
}

func test_specific_checkpoint_rank(uid string, rankType battery.CheckPointRankType) (err error) {
	//sids :=
	req := &battery.QueryUserCheckPointDetailRequest{
		Uid:          proto.String(uid),
		CheckPointId: proto.Uint32(3),
		Sids:         []string{"sina_weibo_0", "sina_weibo_1", "sina_weibo_2"}, //friendrank 使用
		RankType:     rankType.Enum(),
		Source:       battery.ID_SOURCE_SRC_SINA_WEIBO.Enum(), //friendrank 使用
	}

	resp := &battery.QueryUserCheckPointDetailResponse{}

	err = test_op(req, resp, API_URI_CHECKPOINT_QUERY_DETAIL, 1000*1000)
	return
}

func test_query_checkpoint_global_rank() (err error) {
	return test_specific_query_checkpoint_global_rank("1427266457542935645408")
}

func test_pprof_query_checkpoint_global_rank(done, errSum *uint64, id, sum uint64) (err error) {
	return test_pprof_op(test_specific_query_checkpoint_global_rank, done, errSum, id, sum)
}

func test_commit_checkpoint() (err error) {
	req := &battery.CommitCheckPointRequest{
		Uid:          proto.String("1411365678693098337130"),
		CheckPointId: proto.Uint32(3),
		GameId:       proto.String("test_game_id"),
		Score:        proto.Uint64(102),
		Charge:       proto.Uint64(102),
	}

	resp := &battery.CommitCheckPointResponse{}

	err = test_op(req, resp, API_URI_CHECKPOINT_COMMIT, 1000*1000)

	return
}

func test_query_specific_user_wallet(uid string) (err error) {
	req := &battery.QueryWalletRequest{
		Uid: proto.String(uid),
	}

	resp := &battery.QueryWalletResponse{}

	err = test_op(req, resp, API_URI_WALLET_QUERY, REQUEST_TIMEOUT)
	return
}

func test_query_user_wallet() (err error) {
	return test_query_specific_user_wallet("1438769158813213582318")
}

func test_pprof_query_user_wallet(done, errSum *uint64, id, sum uint64) (err error) {
	err = test_pprof_op(test_query_specific_user_wallet, done, errSum, id, sum)
	return
}

func test_query_beforegame_goods() (err error) {
	req := &battery.BeforeGamePropRequest{
		Uid: proto.String("1411365678693098337130"),
		Cmd: battery.BeforeGamePropCmd_BeforeGamePropCmd_Query.Enum(),
	}

	resp := &battery.BeforeGamePropResponse{}

	err = test_op(req, resp, API_URI_BEFOREGAME_OP, REQUEST_TIMEOUT)

	return
}

func test_query_specific_beforegame_goods(uid string) (err error) {
	req := &battery.BeforeGamePropRequest{
		Uid: proto.String(uid),
		Cmd: battery.BeforeGamePropCmd_BeforeGamePropCmd_Query.Enum(),
	}

	resp := &battery.BeforeGamePropResponse{}

	err = test_op(req, resp, API_URI_BEFOREGAME_OP, 1000*1000)

	return
}

func test_pprof_query_beforegame_goods(done, errSum *uint64, id, sum uint64) (err error) {
	beginId := id * sum
	c := make(chan bool, DefConfig.FlowControl)
	for i := uint64(0); i < sum; i++ {
		var uid string
		j := i + beginId
		c <- true
		go func() {
			uid, err = test_specific_login(fmt.Sprintf("sina_weibo_%d", j),
				fmt.Sprintf("name_%d", j),
				fmt.Sprintf("devId_%d", j),
				fmt.Sprintf("iconUrl_%d", j),
				int32(CLIENT_VERSION),
				battery.ID_SOURCE_SRC_SINA_WEIBO)
			if err != nil {
				*errSum++
				<-c
				return
			}
			*done++ //login done

			err = test_query_specific_beforegame_goods(uid)
			if err != nil {
				*errSum++
				<-c
				return
			}

			*done++ //buy beforegame done
			<-c
		}()
	}

	return
}

func test_buy_beforegame_goods() (err error) {
	req := &battery.BeforeGamePropRequest{
		Uid:    proto.String("1416474160659577522511"),
		Cmd:    battery.BeforeGamePropCmd_BeforeGamePropCmd_Buy.Enum(),
		ItemId: proto.Uint64(160010000),
	}

	resp := &battery.BeforeGamePropResponse{}

	err = test_op(req, resp, API_URI_BEFOREGAME_OP, 1000*1000)

	return
}

func test_specific_buy_beforegame_goods(uid string) (err error) {

	goodids := []uint64{160010000, 160020000, 160030000, 160050000}
	goodid := goodids[rand.Intn(len(goodids))]

	req := &battery.BeforeGamePropRequest{
		Uid:    proto.String(uid),
		Cmd:    battery.BeforeGamePropCmd_BeforeGamePropCmd_Buy.Enum(),
		ItemId: proto.Uint64(goodid),
	}

	resp := &battery.BeforeGamePropResponse{}

	err = test_op(req, resp, API_URI_BEFOREGAME_OP, 1000*1000)

	return
}

func test_pprof_buy_beforegame_goods(done, errSum *uint64, id, sum uint64) (err error) {
	beginId := id * sum
	c := make(chan bool, DefConfig.FlowControl)
	for i := uint64(0); i < sum; i++ {
		var uid string
		j := beginId + i
		c <- true
		go func() {

			uid, err = test_specific_login(fmt.Sprintf("sina_weibo_%d", j),
				fmt.Sprintf("name_%d", j),
				fmt.Sprintf("devId_%d", j),
				fmt.Sprintf("iconUrl_%d", j),
				int32(CLIENT_VERSION),
				battery.ID_SOURCE_SRC_SINA_WEIBO)
			if err != nil {
				*errSum++
				<-c
				return
			}

			err = test_specific_buy_beforegame_goods(uid)
			if err != nil {
				<-c
				return
			}
			*done++
			<-c
		}()
	}

	return
}

func test_use_beforegame_goods() (err error) {
	req := &battery.BeforeGamePropRequest{
		Uid:    proto.String("1411365678693098337130"),
		Cmd:    battery.BeforeGamePropCmd_BeforeGamePropCmd_Use.Enum(),
		ItemId: proto.Uint64(160080000),
	}

	resp := &battery.BeforeGamePropResponse{}

	err = test_op(req, resp, API_URI_BEFOREGAME_OP, 1000*1000)

	return
}

func test_query_specific_user_role_info(uid string) (err error) {
	req := &battery.RoleInfoListRequest{
		Uid: proto.String(uid),
		Cmd: battery.RoleInfoListCmd_RoleInfoListCmd_RoleList.Enum(),
	}

	resp := &battery.RoleInfoListResponse{}

	err = test_op(req, resp, API_URI_ROLE_INFO, 1000*1000)

	return
}

func test_query_user_role_info() (err error) {
	return test_query_specific_user_role_info("1425353626119417911485")
}

func test_pprof_user_role_info(done, errSum *uint64, id, sum uint64) (err error) {
	err = test_pprof_op(test_query_specific_user_role_info, done, errSum, id, sum)
	return
}

func test_query_friend_role_info() (err error) {

	req := &battery.RoleInfoListRequest{
		Uid:        proto.String("1425353626119417911485"),
		Cmd:        battery.RoleInfoListCmd_RoleInfoListCmd_FriendRoleList.Enum(),
		FriendSids: []string{"sina_weibo_1", "sina_weibo_2", "sina_weibo_3"},
		Source:     battery.ID_SOURCE_SRC_SINA_WEIBO.Enum(),
	}

	resp := &battery.RoleInfoListResponse{}

	err = test_op(req, resp, API_URI_ROLE_INFO, 10*1000*1000)

	return
}

func test_set_user_select_role_id() (err error) {
	req := &battery.RoleInfoListRequest{
		Uid:    proto.String("1425353626119417911485"),
		Cmd:    battery.RoleInfoListCmd_RoleInfoListCmd_Select.Enum(),
		RoleId: proto.Uint64(150040000),
	}

	resp := &battery.RoleInfoListResponse{}

	err = test_op(req, resp, API_URI_ROLE_INFO, 1000*1000)

	return
}

func test_upgrade_user_role() (err error) {
	req := &battery.RoleInfoListRequest{
		Uid:    proto.String("1416638813923186911847"),
		Cmd:    battery.RoleInfoListCmd_RoleInfoListCmd_Upgrade.Enum(),
		RoleId: proto.Uint64(150040001),
	}

	resp := &battery.RoleInfoListResponse{}

	err = test_op(req, resp, API_URI_ROLE_INFO, 1000*1000)

	return
}

func test_query_user_jigsaw() (err error) {
	req := &battery.JigsawRequest{
		Uid: proto.String("1425353626119417911485"),
		Cmd: battery.JigsawCmd_JigsawCmd_Query.Enum(),
	}

	resp := &battery.JigsawResponse{}

	err = test_op(req, resp, API_URI_JIGSAW, 1000*1000)

	return
}

func test_buy_user_jigsaw() (err error) {
	req := &battery.JigsawRequest{
		Uid:    proto.String("1425352623264100099847"),
		Cmd:    battery.JigsawCmd_JigsawCmd_Buy.Enum(),
		ItemId: proto.Uint64(120010000),
	}

	resp := &battery.JigsawResponse{}

	err = test_op(req, resp, API_URI_JIGSAW, 1000*1000)

	return
}

func test_query_user_friend_mails() (err error) {
	req := &battery.FriendMailListRequest{
		Uid:    proto.String("1434510198689404183089"),
		Cmd:    battery.FriendMailCmd_FriendMailCmd_MailList.Enum(),
		Source: battery.ID_SOURCE_SRC_SINA_WEIBO.Enum(),
	}

	resp := &battery.FriendMailListResponse{}

	err = test_op(req, resp, API_URI_FRIEND_MAIL, 1000*1000)

	return
}

func test_specific_query_user_friend_mails(uid string) (err error) {
	req := &battery.FriendMailListRequest{
		Uid:    proto.String(uid),
		Cmd:    battery.FriendMailCmd_FriendMailCmd_MailList.Enum(),
		Source: battery.ID_SOURCE_SRC_SINA_WEIBO.Enum(),
	}

	resp := &battery.FriendMailListResponse{}

	err = test_op(req, resp, API_URI_FRIEND_MAIL, 1000*1000)

	return
}

func test_pprof_query_friend_from_friendship(done, errSum *uint64, id, sum uint64) (err error) {
	c := make(chan bool, DefConfig.FlowControl)
	for i := uint64(0); i < sum; i++ {
		c <- true
		go func() {
			test_pprof_sub_query_friend_from_friendship(done, errSum, i, id, sum)
			<-c
		}()
		//j := i + beginId
		//sid, uid := fmt.Sprintf("sina_weibo_%d", i), ""
		//uid, err = test_specific_login(sid, fmt.Sprintf("name_%d", j), fmt.Sprintf("devId_%d", j), fmt.Sprintf("iconUrl_%d", j), int32(CLIENT_VERSION))
		//if err != nil {
		//	*errSum++
		//	break
		//}

		//*done++

		//err = test_specific_query_user_friend_mails(uid)
		//if err != nil {
		//	*errSum++
		//	return
		//}

		//*done++
	}
	return
}

func test_pprof_sub_query_friend_from_friendship(done, errSum *uint64, i, id, sum uint64) (err error) {
	j := i + id*sum
	sid, uid := fmt.Sprintf("sina_weibo_%d", j), ""
	uid, err = test_specific_login(sid,
		fmt.Sprintf("name_%d", j),
		fmt.Sprintf("devId_%d", j),
		fmt.Sprintf("iconUrl_%d", j),
		int32(CLIENT_VERSION),
		battery.ID_SOURCE_SRC_SINA_WEIBO)
	if err != nil {
		xylog.ErrorNoId("%v", err)
		*errSum++
		return
	}

	*done++

	err = test_specific_query_user_friend_mails(uid)
	if err != nil {
		xylog.ErrorNoId("%v", err)
		*errSum++
		return
	}

	*done++
	return
}

func test_give_friend_from_friendship() (err error) {
	req := &battery.FriendMailListRequest{
		Uid:       proto.String("1435754950518083130089"),
		Cmd:       battery.FriendMailCmd_FriendMailCmd_StaminaGive.Enum(),
		FriendSid: proto.String("sina_weibo_0"),
		Source:    battery.ID_SOURCE_SRC_SINA_WEIBO.Enum(),
	}

	resp := &battery.FriendMailListResponse{}

	err = test_op(req, resp, API_URI_FRIEND_MAIL, 1000*1000)

	return
}

func test_specific_give_friend_from_friendship(uid, friendsid string) (err error) {
	req := &battery.FriendMailListRequest{
		Uid:       proto.String(uid),
		Cmd:       battery.FriendMailCmd_FriendMailCmd_StaminaGive.Enum(),
		FriendSid: proto.String(friendsid),
		Source:    battery.ID_SOURCE_SRC_SINA_WEIBO.Enum(),
	}

	resp := &battery.FriendMailListResponse{}

	err = test_op(req, resp, API_URI_FRIEND_MAIL, 1000*1000)

	return
}

func test_pprof_give_friend_from_friendship(done, errSum *uint64, id, sum uint64) (err error) {
	c := make(chan bool, DefConfig.FlowControl)
	for i := uint64(0); i < sum; i++ {
		c <- true
		go func() {
			test_pprof_sub_give_friend_from_friendship(done, errSum, i, id, sum)
			<-c
		}()
		//j := i + beginId
		//sid, uid := fmt.Sprintf("sina_weibo_%d", j), ""
		//uid, err = test_specific_login(sid, fmt.Sprintf("name_%d", j), fmt.Sprintf("devId_%d", j), fmt.Sprintf("iconUrl_%d", j), int32(CLIENT_VERSION))
		//if err != nil {
		//	*errSum++
		//	break
		//}

		//*done++

		//err = test_specific_give_friend_from_friendship(uid, fmt.Sprintf("sina_weibo_%d", i+1))
		//if err != nil {
		//	*errSum++
		//	return
		//}

		//*done++
	}
	return
}

func test_pprof_sub_give_friend_from_friendship(done, errSum *uint64, i, id, sum uint64) (err error) {
	j := i + id*sum
	sid, uid := fmt.Sprintf("sina_weibo_%d", j), ""
	uid, err = test_specific_login(sid,
		fmt.Sprintf("name_%d", j),
		fmt.Sprintf("devId_%d", j),
		fmt.Sprintf("iconUrl_%d", j),
		int32(CLIENT_VERSION),
		battery.ID_SOURCE_SRC_SINA_WEIBO)
	if err != nil {
		xylog.ErrorNoId("%v", err)
		*errSum++
		return
	}

	*done++

	err = test_specific_give_friend_from_friendship(uid, fmt.Sprintf("sina_weibo_%d", i+1))
	if err != nil {
		xylog.ErrorNoId("%v", err)
		*errSum++
		return
	}

	*done++

	return
}

func test_apply_friend_from_friendship() (err error) {
	req := &battery.FriendMailListRequest{
		Uid:       proto.String("1432296717359433712757"),
		Cmd:       battery.FriendMailCmd_FriendMailCmd_StaminaGetApply.Enum(),
		FriendSid: proto.String("sina_weibo_0"),
		Source:    battery.ID_SOURCE_SRC_SINA_WEIBO.Enum(),
	}

	resp := &battery.FriendMailListResponse{}

	err = test_op(req, resp, API_URI_FRIEND_MAIL, 1000*1000)

	return
}

func test_specific_apply_friend_from_friendship(uid, friendsid string) (err error) {
	req := &battery.FriendMailListRequest{
		Uid:       proto.String(uid),
		Cmd:       battery.FriendMailCmd_FriendMailCmd_StaminaGetApply.Enum(),
		FriendSid: proto.String(friendsid),
		Source:    battery.ID_SOURCE_SRC_SINA_WEIBO.Enum(),
	}

	resp := &battery.FriendMailListResponse{}

	err = test_op(req, resp, API_URI_FRIEND_MAIL, 1000*1000)

	return
}

func test_pprof_apply_friend_from_friendship(done, errSum *uint64, id, sum uint64) (err error) {
	c := make(chan bool, DefConfig.FlowControl)
	for i := uint64(0); i < sum; i++ {
		c <- true
		go func() {
			test_pprof_sub_apply_friend_from_friendship(done, errSum, i, id, sum)
			<-c
		}()
		//j := i + beginId
		//sid, uid := fmt.Sprintf("sina_weibo_%d", j), ""
		//uid, err = test_specific_login(sid, fmt.Sprintf("name_%d", j), fmt.Sprintf("devId_%d", j), fmt.Sprintf("iconUrl_%d", j), int32(CLIENT_VERSION))
		//if err != nil {
		//	*errSum++
		//	break
		//}

		//*done++

		//err = test_specific_apply_friend_from_friendship(uid, fmt.Sprintf("sina_weibo_%d", i+1))
		//if err != nil {
		//	*errSum++
		//	return
		//}

		//*done++
	}
	return
}

func test_pprof_sub_apply_friend_from_friendship(done, errSum *uint64, i, id, sum uint64) (err error) {
	j := i + id*sum
	sid, uid := fmt.Sprintf("sina_weibo_%d", j), ""
	uid, err = test_specific_login(sid,
		fmt.Sprintf("name_%d", j),
		fmt.Sprintf("devId_%d", j),
		fmt.Sprintf("iconUrl_%d", j),
		int32(CLIENT_VERSION),
		battery.ID_SOURCE_SRC_SINA_WEIBO)
	if err != nil {
		xylog.ErrorNoId("%v", err)
		*errSum++
		return
	}

	*done++

	err = test_specific_apply_friend_from_friendship(uid, fmt.Sprintf("sina_weibo_%d", i+1))
	if err != nil {
		xylog.ErrorNoId("%v", err)
		*errSum++
		return
	}

	*done++

	return
}

func test_confirm_user_friendmails() (err error) {
	req := &battery.FriendMailListRequest{
		Uid: proto.String("1421047571687749097495"),
		Cmd: battery.FriendMailCmd_FriendMailCmd_StaminatGetAll.Enum(),
	}

	resp := &battery.FriendMailListResponse{}

	err = test_op(req, resp, API_URI_FRIEND_MAIL, 1000*1000)

	return
}

func test_specific_confirm_user_friendmails(uid string) (err error) {
	req := &battery.FriendMailListRequest{
		Uid: proto.String(uid),
		Cmd: battery.FriendMailCmd_FriendMailCmd_StaminatGetAll.Enum(),
	}

	resp := &battery.FriendMailListResponse{}

	err = test_op(req, resp, API_URI_FRIEND_MAIL, 1000*1000)

	return
}

func test_pprof_confirm_friend_from_friendship(done, errSum *uint64, id, sum uint64) (err error) {
	beginId := id * sum
	for i := uint64(0); i < sum; i++ {
		sid, uid := fmt.Sprintf("sina_weibo_%d", i+beginId), ""
		uid, err = test_specific_login(sid,
			fmt.Sprintf("name_%d", i),
			fmt.Sprintf("devId_%d", i),
			fmt.Sprintf("iconUrl_%d", i),
			int32(20000),
			battery.ID_SOURCE_SRC_SINA_WEIBO)
		if err != nil {
			*errSum++
			continue
		}

		*done++

		err = test_specific_confirm_user_friendmails(uid)
		if err != nil {
			*errSum++
			continue
		}

		*done++
	}
	return
}

func test_query_user_sysmail() (err error) {
	err = test_query_specific_user_sysmail("1442828654322275370081")
	return
}

func test_query_specific_user_sysmail(uid string) (err error) {

	if "" == uid {
		return xyerror.ErrBadInputData
	}

	req := &battery.SystemMailListRequest{
		Uid: proto.String(uid),
		Cmd: battery.SystemMailCmd_SystemMailCmd_Maillist.Enum(),
	}

	resp := &battery.SystemMailListResponse{}

	err = test_op(req, resp, API_URI_SYS_MAIL, 1000*1000)

	return
}

func test_pprof_query_user_sysmail(done, errSum *uint64, id, sum uint64) (err error) {
	beginId := id * sum
	for i := uint64(0); i < sum; i++ {
		var uid string
		uid, err = test_specific_login(fmt.Sprintf("sina_weibo_%d", i+beginId),
			fmt.Sprintf("name_%d", i+beginId),
			fmt.Sprintf("devId_%d", i+beginId),
			fmt.Sprintf("iconUrl_%d", i+beginId),
			int32(20000),
			battery.ID_SOURCE_SRC_SINA_WEIBO)
		if err != nil {
			*errSum++
			continue
		}
		*done++
		err = test_query_specific_user_sysmail(uid)
		if err != nil {
			*errSum++
			return
		}
		*done++
	}

	return
}

func test_confirm_user_sysmail() (err error) {
	err = test_confirm_specific_user_sysmail("1442645456130270798081", 900001)
	return
}

func test_confirm_specific_user_sysmail(uid string, mailId int32) (err error) {
	req := &battery.SystemMailListRequest{
		Uid:    proto.String(uid),
		Cmd:    battery.SystemMailCmd_SystemMailCmd_GiftGet.Enum(),
		MailId: proto.Int32(mailId),
	}

	resp := &battery.SystemMailListResponse{}

	err = test_op(req, resp, API_URI_SYS_MAIL, 1000*1000)

	return
}

func test_pprof_confirm_user_sysmail(done, errSum *uint64, id, sum uint64) (err error) {
	beginId := id * sum
	for i := uint64(0); i < sum; i++ {
		var uid string
		uid, err = test_specific_login(fmt.Sprintf("sina_weibo_%d", i+beginId),
			fmt.Sprintf("name_%d", i+beginId),
			fmt.Sprintf("devId_%d", i+beginId),
			fmt.Sprintf("iconUrl_%d", i+beginId),
			int32(20000),
			battery.ID_SOURCE_SRC_SINA_WEIBO)
		if err != nil {
			continue
		}

		err = test_confirm_specific_user_sysmail(uid, 111005)
		if err != nil {
			continue
		}
	}

	return
}

func test_read_user_sysmail() (err error) {
	err = test_read_specific_user_sysmail("", 50003)
	return
}

func test_read_specific_user_sysmail(uid string, mailId int32) (err error) {
	req := &battery.SystemMailListRequest{
		Uid:    proto.String(uid),
		Cmd:    battery.SystemMailCmd_SystemMailCmd_MailRead.Enum(),
		MailId: proto.Int32(mailId),
	}

	resp := &battery.SystemMailListResponse{}

	err = test_op(req, resp, API_URI_SYS_MAIL, 1000*1000)

	return
}

func test_pprof_read_user_sysmail(done, errSum *uint64, id, sum uint64) (err error) {
	beginId := id * sum
	for i := uint64(0); i < sum; i++ {
		var uid string
		uid, err = test_specific_login(fmt.Sprintf("sina_weibo_%d", i+beginId),
			fmt.Sprintf("name_%d", i+beginId),
			fmt.Sprintf("devId_%d", i+beginId),
			fmt.Sprintf("iconUrl_%d", i+beginId),
			int32(20000),
			battery.ID_SOURCE_SRC_SINA_WEIBO)
		if err != nil {
			continue
		}

		err = test_read_specific_user_sysmail(uid, 50003)
		if err != nil {
			continue
		}
	}

	return
}

func test_query_user_runes() (err error) {
	req := &battery.RuneRequest{
		Uid: proto.String("1445257345918880074081"),
		Cmd: battery.RuneCmd_RuneCmd_Query.Enum(),
	}

	resp := &battery.RuneResponse{}

	err = test_op(req, resp, API_URI_RUNE, 1000*1000)

	return
}

func test_buy_user_rune() (err error) {
	req := &battery.RuneRequest{
		Uid:    proto.String("1444890726896479762300"),
		Cmd:    battery.RuneCmd_RuneCmd_Buy.Enum(),
		Itemid: proto.Uint64(110010001),
	}

	resp := &battery.RuneResponse{}

	err = test_op(req, resp, API_URI_RUNE, 1000*1000)

	return
}

func test_buy_role() (err error) {
	req := &battery.RoleInfoListRequest{
		Uid:    proto.String("1420602243680111557511"),
		Cmd:    battery.RoleInfoListCmd_RoleInfoListCmd_Upgrade.Enum(),
		RoleId: proto.Uint64(150020001),
	}

	resp := &battery.RoleInfoListResponse{}

	err = test_op(req, resp, API_URI_ROLE_INFO, REQUEST_TIMEOUT)

	return
}

func test_memcache_get() (err error) {
	return test_specific_memcache_get("1434030861182640487081")
}

func test_pprof_memcache_get(done, errSum *uint64, id, sum uint64) (err error) {
	err = test_pprof_op(test_specific_memcache_get, done, errSum, id, sum)
	return
}

func test_specific_memcache_get(uid string) (err error) {
	req := &battery.MemCacheRequest{
		Uid:          proto.String(uid),
		OpType:       battery.MemCacheOperationType_MemCacheOperationType_GET.Enum(),
		OpKey:        battery.MemCacheEnum_MemCacheEnum_UserGuideProcess.Enum(),
		PlatformType: battery.PLATFORM_TYPE_PLATFORM_TYPE_IOS.Enum(),
	}

	resp := &battery.MemCacheResponse{}

	err = test_op(req, resp, API_URI_MEMCACHE, REQUEST_TIMEOUT)

	return
}

func test_memcaches_get() (err error) {
	req := &battery.MemCacheRequest{
		Uid:    proto.String("1434030861182640487081"),
		OpType: battery.MemCacheOperationType_MemCacheOperationType_GET.Enum(),
	}

	req.Units = make([]*battery.MemCacheUnit, 0)
	req.Units = append(req.Units, &battery.MemCacheUnit{
		Key: battery.MemCacheEnum_MemCacheEnum_UserGuideProcess.Enum(),
	})
	req.Units = append(req.Units, &battery.MemCacheUnit{
		Key: battery.MemCacheEnum_MemCacheEnum_CheckPointUnLock.Enum(),
	})

	resp := &battery.MemCacheResponse{}

	err = test_op(req, resp, API_URI_MEMCACHE, REQUEST_TIMEOUT)

	return
}

func test_memcache_set() (err error) {
	return test_specific_memcache_set("1435979229390511153047")
}

func test_specific_memcache_set(uid string) (err error) {

	//xylog.DebugNoId("[%s] test_specific_memcache_set", uid)

	req := &battery.MemCacheRequest{
		Uid:     proto.String(uid),
		OpType:  battery.MemCacheOperationType_MemCacheOperationType_SET.Enum(),
		OpKey:   battery.MemCacheEnum_MemCacheEnum_UserGuideProcess.Enum(),
		OpValue: proto.String("test_value1"),
	}

	resp := &battery.MemCacheResponse{}

	err = test_op(req, resp, API_URI_MEMCACHE, REQUEST_TIMEOUT)

	return
}

func test_pprof_memcache_set(done, errSum *uint64, id, sum uint64) (err error) {
	err = test_pprof_op(test_specific_memcache_set, done, errSum, id, sum)
	return
}

func test_memcaches_set() (err error) {
	req := &battery.MemCacheRequest{
		Uid:    proto.String("1431330472754676933059"),
		OpType: battery.MemCacheOperationType_MemCacheOperationType_SET.Enum(),
	}

	req.Units = make([]*battery.MemCacheUnit, 0)
	req.Units = append(req.Units, &battery.MemCacheUnit{
		Key:   battery.MemCacheEnum_MemCacheEnum_UserGuideProcess.Enum(),
		Value: proto.String("test_value1"),
	})

	req.Units = append(req.Units, &battery.MemCacheUnit{
		Key:   battery.MemCacheEnum_MemCacheEnum_CheckPointUnLock.Enum(),
		Value: proto.String("test_value2"),
	})

	resp := &battery.MemCacheResponse{}

	err = test_op(req, resp, API_URI_MEMCACHE, REQUEST_TIMEOUT)

	return
}

func test_send_alert() (err error) {
	//err = xybusiness.SendAlert("123456", "test subject", "test content")
	return
}

func test_maintenance_prop() (err error) {

	var (
		api_url string = "http://" + DefConfig.Url + API_URI_MAINTENANCE + "/prop?amount=1&identity=AA10017&platform=1&optype=1&propid=146030000&proptype=1&ts=1431004245&sig=5437603d7655e9d2548e1d5bff0d921f&source=6"
		//api_url   = "http://" + DefConfig.Url + API_URI_MAINTENANCE + "/prop?identity=105641609835577344&platform=1&optype=4&source=6"
		timeBegin = time.Now()
		response  []byte
	)

	request := httplib.Get(api_url)
	request.SetTimeout(60*time.Second, REQUEST_TIMEOUT*time.Microsecond)

	response, err = request.Bytes()
	xylog.DebugNoId("error : %v, response : %s", err, response)
	if err != nil {
		xylog.ErrorNoId("Error Send: %s", err.Error())
		return
	}

	xylog.DebugNoId("test_maintenance_prop cost %d ms", time.Since(timeBegin).Nanoseconds()/int64(time.Millisecond))

	return
}

func test_apn_notification() (err error) {
	//var client apns.Client
	//client, err = apns.NewClientWithFiles("gateway.push.apple.com:2195", "aps_pro_com.737.batteryrun.cn.pem", "aps_pro_com.737.batteryrun.cn.pem")
	//if err != nil {
	//	xylog.ErrorNoId("can't create new client %v", err)
	//	return
	//}

	//go func() {
	//	for f := range client.FailedNotifs {
	//		fmt.Println("notification ", f.Notif.ID, "failed with", f.Err.Error())
	//	}
	//}()

	//payload := apns.NewPayload()
	//payload.APS.Alert.Title = "fuck you"
	//payload.APS.Alert.Body = "来阿里阿里"
	//badge := 10
	//payload.APS.Badge = &badge
	//payload.APS.Sound = "bingbong.aiff"
	//payload.APS.ContentAvailable = 1
	//payload.SetCustomValue("link", "zombo://dot/com")
	//payload.SetCustomValue("game", map[string]int{"score": 234})

	////pn := apns.NewPushNotification()
	//pn := apns.NewNotification()
	//pn.Payload = payload
	//pn.DeviceToken = "a73cbdea269eec3ce6b52e4c33d88c087052b403b32006670cd68577dcee90d3"
	//pn.Priority = apns.PriorityImmediate
	//pn.Identifier = 12312
	//pn.ID = "user_id:timestamp"

	////for i := 0; i < 3; i++ {
	//err = client.Send(pn)
	//xylog.DebugNoId("client.Send result : %v", err)
	//time.Sleep(time.Second * 10)
	////}
	////pn.AddPayload(payload)
	////for i := 0; i < 10000; i++ {
	////go func() {
	////	begin := time.Now()
	////	resp := client.Send(pn)
	////	fmt.Printf("cost %dms\n", time.Since(begin)/time.Millisecond)
	////	alert, _ := pn.PayloadString()
	////	fmt.Println("  Alert:", alert)
	////	fmt.Println("Success:", resp.Success)
	////	fmt.Println("  Error:", resp.Error)
	////}()
	////}

	return
}

func test_apn_notification2apns() (err error) {

	////连接到apns nats
	//natsService := xynatsservice.NewNatsService("Battery apns nats", DefConfig.ApnsNatsUrl)
	//err = natsService.Start()
	//if err != xyerror.ErrOK {
	//	xylog.ErrorNoId("natsService(%s).Start failed : %v", DefConfig.ApnsNatsUrl, err)
	//	return
	//}
	//defer natsService.Stop()

	//xyapn.InitNatsService(natsService)

	//notifiction := &xyapn.APNNotification{
	//	Cmd:         xyapn.APNNotificationCMD_Notification.Enum(),
	//	DeviceToken: proto.String("a73cbdea269eec3ce6b52e4c33d88c087052b403b32006670cd68577dcee90d3"),
	//	Title:       proto.String("打老王"),
	//	Content:     proto.String("老王来打我啊！！！~~~ 田叔"),
	//	Badge:       proto.Int32(1),
	//}

	//err = xyapn.Send(notifiction)
	//if err != xyerror.ErrOK {
	//	xylog.ErrorNoId("xyapn.Send failed : %v", err)
	//	return
	//}

	//xylog.DebugNoId("xyapn.Send(%v) succeed!", notifiction)

	return
}

func test_apn_enabledevicetoken2apns() (err error) {

	////连接到apns nats
	//natsService := xynatsservice.NewNatsService("Battery apns nats", DefConfig.ApnsNatsUrl)
	//err = natsService.Start()
	//if err != xyerror.ErrOK {
	//	xylog.ErrorNoId("natsService(%s).Start failed : %v", DefConfig.ApnsNatsUrl, err)
	//	return
	//}
	//defer natsService.Stop()

	//xyapn.InitNatsService(natsService)

	//notifiction := &xyapn.APNNotification{
	//	Cmd:         xyapn.APNNotificationCMD_EnableDeviceToken.Enum(),
	//	DeviceToken: proto.String("cb7a0d0b5cd15a4b3128e3c015d92c47f72661493fae12e6cd03082623992213"),
	//}

	//err = xyapn.Send(notifiction)
	//if err != xyerror.ErrOK {
	//	xylog.ErrorNoId("xyapn.Send failed : %v", err)
	//	return
	//}

	//xylog.DebugNoId("xyapn.Send(%v) succeed!", notifiction)

	return
}

func test_apn_disabledevicetoken2apns() (err error) {

	////连接到apns nats
	//natsService := xynatsservice.NewNatsService("Battery apns nats", DefConfig.ApnsNatsUrl)
	//err = natsService.Start()
	//if err != xyerror.ErrOK {
	//	xylog.ErrorNoId("natsService(%s).Start failed : %v", DefConfig.ApnsNatsUrl, err)
	//	return
	//}
	//defer natsService.Stop()

	//xyapn.InitNatsService(natsService)

	//notifiction := &xyapn.APNNotification{
	//	Cmd:         xyapn.APNNotificationCMD_DisableDeviceToken.Enum(),
	//	DeviceToken: proto.String("cb7a0d0b5cd15a4b3128e3c015d92c47f72661493fae12e6cd03082623992213"),
	//}

	//err = xyapn.Send(notifiction)
	//if err != xyerror.ErrOK {
	//	xylog.ErrorNoId("xyapn.Send failed : %v", err)
	//	return
	//}

	//xylog.DebugNoId("xyapn.Send(%v) succeed!", notifiction)

	return
}

func test_timer() (err error) {
	option := xytimer.TimerOption{
		Type: xytimer.TIMER_TYPE_FIXED,
	}
	option.Moments = []xytimer.TimerMoment{{Hour: 16, Minute: 27, Second: 0}, {Hour: 17, Minute: 10, Second: 00}}
	xytimer.InitTimer(option, func() {
		fmt.Println("[Fixed] Why U call me bitch?!")
	})

	option.Type = xytimer.TIMER_TYPE_INTERVAL
	option.Interval = 10
	xytimer.InitTimer(option, func() {
		fmt.Println("[Interval] Why U call me bitch?! ")
	})

	return
}

func test_advertisement() (err error) {
	//单个广告
	//多个广告，选择广告是否按照权重来
	//没有广告时，能否正常加载并返回结果
	req := &battery.AdvertisementRequest{
		Uid:          proto.String("1439973434308868422540"),
		PlatformType: battery.PLATFORM_TYPE_PLATFORM_TYPE_IOS.Enum(),
	}

	//req.AdvertisementSpaceId = make([]uint32, 0)
	//req.AdvertisementSpaceId = append(req.AdvertisementSpaceId, 1)
	req.AdvertisementSpaceId = []uint32{uint32(1), uint32(2)}

	resp := &battery.AdvertisementResponse{}

	err = test_op(req, resp, API_URI_ADVERTISEMENT, REQUEST_TIMEOUT)

	return

}

func test_umeng_apn() (err error) {
	//apnsserver.NewXYAPI().NotifyUMeng
	return
}

//内购统计上报
func test_iapstatistic() (err error) {
	urlStr := "http://gateway.dc.737.com/index.php"
	keys, values := make([]string, 0), make(map[string]string, 0)

	//values := url.Values{}
	//values.Set("app_id", "10004")
	values["app_id"] = "10004"
	keys = append(keys, "app_id")
	//values.Set("bis", "order")
	values["bis"] = "order"
	keys = append(keys, "bis")
	//values.Set("ac", "report")
	values["ac"] = "report"
	keys = append(keys, "ac")
	//values.Set("quantity", "1")
	values["quantity"] = "1"
	keys = append(keys, "quantity")
	//values.Set("amount", "1")
	values["pid"] = "76"
	keys = append(keys, "pid")
	//values.Set("username", "阿狗")
	values["username"] = "阿狗"
	keys = append(keys, "username")
	//values.Set("rname", "chongchong")
	values["rname"] = "冲冲"
	keys = append(keys, "rname")
	//values.Set("level", "20")
	values["level"] = "20"
	keys = append(keys, "level")
	values["sandbox"] = "1"
	keys = append(keys, "sandbox")
	//values.Set("oid", "123456")
	values["oid"] = "123460"
	keys = append(keys, "oid")
	//values.Set("device_id", "a73cbdea269eec3ce6b52e4c33d88c087052b403b32006670cd68577dcee90d3")
	values["device_id"] = "a73cbdea269eec3ce6b52e4c33d88c087052b403b32006670cd68577dcee90d3"
	keys = append(keys, "device_id")

	sort.Strings(keys)

	var encodeSrc string
	for i, k := range keys {
		if i == 0 {
			encodeSrc += k + "=" + values[k]
		} else {
			encodeSrc += "&" + k + "=" + values[k]
		}
	}

	urlEncode := url.QueryEscape(encodeSrc)
	urlEncode += "&a0e90ad18ce588000586ab6becb43923"

	//urlEncode = "ac%3Dreport%26amount%3D1%26app_id%3D10004%26bis%3Dorder%26device_id%3Da73cbdea269eec3ce6b52e4c33d88c087052b403b32006670cd68577dcee90d3%26level%3D20%26oid%3D123456%26quantity%3D1%26rname%3Dchongchong%26uname%3Dricktian&a0e90ad18ce588000586ab6becb43923"

	//urlEncode = "ac%3Dreport%26amount%3D1%26app_id%3D10004%26bis%3Dorder%26device_id%3Da73cbdea269eec3ce6b52e4c33d88c087052b403b32006670cd68577dcee90d3%26level%3D20%26oid%3D123456%26quantity%3D1%26rname%3Dchongchong%26uname%3D%E9%98%BF%E7%8B%97&a0e90ad18ce588000586ab6becb43923"

	xylog.DebugNoId("\nencodeSrc : %v\nurlencode : %v\n", encodeSrc, urlEncode)

	md5Ctx := md5.New()
	md5Ctx.Write([]byte(urlEncode))
	cipherStr := md5Ctx.Sum(nil)

	sign := strings.ToLower(hex.EncodeToString(cipherStr))
	encodeSrc += "&sign=" + sign

	type Response struct {
		Code int    `json:"code"`
		Desc string `json:"desc"`
		Data string `json:"data"`
	}

	//resp, err := http.PostForm(urlStr, values)
	resp, err := http.Get(fmt.Sprintf("%s?%v", urlStr, encodeSrc))
	//xylog.DebugNoId("encode : %v", values.Encode())
	xylog.DebugNoId("err : %v, resp : %v", err, resp)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	respData := &Response{}
	err = json.Unmarshal(body, respData)
	if err != nil {
		xylog.DebugNoId("failed. %v", err)
		return
	} else {
		if respData.Code != 0 {
			xylog.DebugNoId("failed.")
		} else {
			xylog.DebugNoId("succeed.")
		}
	}

	xylog.DebugNoId("resp body : %v", string(body))

	return
}

func test_transaction_iapstatistic() (err error) {

	transactionbusiness.DefConfigCache.Configs().IsProduction = false
	transactionbusiness.NewXYAPI().SendIapStatistic("1435390486436053904081", "Buy_Goods_0704", "10000")
	return
}
