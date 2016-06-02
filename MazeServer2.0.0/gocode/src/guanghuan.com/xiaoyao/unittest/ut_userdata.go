package main

import (
	"bytes"
	proto "code.google.com/p/goprotobuf/proto"
	//xylog "guanghuan.com/xiaoyao/common/log"
	xyutil "guanghuan.com/xiaoyao/common/util"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	//	xyfmt "guanghuan.com/xiaoyao/superbman_server/fmt"
	"log"
	"net/http"
	"time"
)

func Test_Normal(fbid string, delay int) {

	i := 0
	for {
		i++
		log.Printf("<<<<<<<<<<< test #%d start >>>>>>>>>>>>", i)
		uid := test_login(fbid, "")
		if uid != "" {
			test_user_data(uid)
			test_query_stamina(uid)
			/*
				if stamina > 2 {
					for j := 0; j < 3; j++ {
						game_id := test_new_game(uid)
						if game_id != "" && game_id != "0" {
							test_add_game_data(uid, game_id)
						} else {
							log.Printf("no new game created, by pass adding result test")
						}
						time.Sleep(5 * time.Second)
					}
					//				test_user_data(uid)
				}
				//			test_friend_data(uid)
			*/
		} else {
			log.Printf("login failed")
		}
		//		test_echo()
		log.Printf("<<<<<<<<<<< test #%d done  >>>>>>>>>>>>", i)
		log.Printf("sleep %d second", delay)
		time.Sleep(time.Duration(delay) * time.Second)
	}
}

func test_friend_data(uid string, friends []string) {
	log.Printf("==== test_friend_data start ")
	defer log.Printf("==== test_friend_data done ")

	req := battery.QueryFriendsDataRequest{}
	req.Uid = proto.String(uid)
	req.Sids = friends

	req.Source = battery.ID_SOURCE_SRC_SINA_WEIBO.Enum()
	//	log.Printf("request: %s", req.String())

	var data_in []byte
	data_in, err := proto.Marshal(&req)
	if err != nil {
		log.Printf("Marshal Error:%s", err.Error())
	}

	resp, err := http.Post(server_url+"/v1/friend/315", "", bytes.NewReader(data_in))
	if err != nil {
		log.Printf("Error: %s", err.Error())
	} else {
		resp_obj := &battery.QueryFriendsDataResponse{}
		err = handleResp(resp, resp_obj)
		if err != nil {
			log.Printf("Error: %s", err.Error())
		} else {
			log.Printf("response: %s", (*resp_obj).String())
		}
	}
}

func test_user_data(uid string) {
	log.Printf("==== test_user_data start ")
	defer log.Printf("==== test_user_data done ")
	req := battery.QueryUserDataRequest{}
	req.Uid = proto.String(uid)

	//	log.Printf("request: %s", req.String())

	var data_in []byte
	data_in, err := proto.Marshal(&req)
	if err != nil {
		log.Printf("Error:%s", err.Error())
	}

	resp, err := http.Post(server_url+"/v1/user/315", "", bytes.NewReader(data_in))
	if err != nil {
		log.Printf("Error: %s", err.Error())
	} else {
		resp_obj := &battery.QueryUserDataResponse{}
		err = handleResp(resp, resp_obj)
		if err != nil {
			log.Printf("Error: %s", err.Error())
		} else {
			//			log.Printf("response: %s", (*resp_obj).String())
			//if resp_obj.GetData().Wallet != nil {
			//	//				stm := resp_obj.GetData().GetWallet().GetStamina()
			//	//log.Printf("stamina = %d", stm)
			//}
			//if resp_obj.GetData().Total != nil {
			//	score := resp_obj.GetData().GetTotal().GetScore()
			//	log.Printf("score = %d", score)
			//}
			//			log.PrintAccount(*resp_obj.GetData())
			//xylog.Debug(xyfmt.FormatAccount(*resp_obj.GetData()))
		}
	}
}

func test_login(fbid string, name string) (uid string) {
	log.Printf("==== test_login start ")
	defer log.Printf("==== test_login done")

	req := battery.LoginRequest{}
	req.Token = proto.String("abin")
	req.LoginId = &battery.TPID{
		Id:     proto.String(fbid),
		Name:   proto.String(name),
		Source: battery.ID_SOURCE_SRC_SINA_WEIBO.Enum(),
	}
	req.DeviceId = &battery.TPID{
		Id: proto.String(""),
	}

	var data_in []byte
	data_in, err := proto.Marshal(&req)
	if err != nil {
		log.Printf("Error:%s", err.Error())
	}
	log.Printf("login with gcid: %s", fbid)
	resp, err := http.Post(server_url+"/v1/login/315", "", bytes.NewReader(data_in))
	if err != nil {
		log.Printf("Error: %s", err.Error())
	} else {
		resp_obj := &battery.LoginResponse{}
		err = handleResp(resp, resp_obj)
		if err != nil {
			log.Printf("Error: %s", err.Error())
		} else {
			uid = resp_obj.GetUid()
			server_time := resp_obj.GetServerTime()
			log.Printf("uid: %s", uid)
			log.Printf("server time: %s (%d)", xyutil.ToStrTime(server_time), server_time)
			//			acc := resp_obj.GetData()
			//			xyutil.PrintAccount(*acc)
			//			xylog.Debug(xyfmt.FormatAccount(*resp_obj.GetData()))
		}
	}
	return
}

func test_query_stamina(uid string) (stamina int32) {
	log.Printf("==== test_query_stamina start")
	defer log.Printf("==== test_query_stamina done")

	req := battery.QueryStaminaRequest{}
	req.Uid = proto.String(uid)
	var data_in []byte
	data_in, err := proto.Marshal(&req)
	if err != nil {
		log.Printf("Error:%s", err.Error())
	}

	resp, err := http.Post(server_url+"/v1/stamina/315", "", bytes.NewReader(data_in))
	if err != nil {
		log.Printf("Error: %s", err.Error())
	} else {
		resp_obj := &battery.QueryStaminaResponse{}
		err = handleResp(resp, resp_obj)
		if err != nil {
			log.Printf("Error: %s", err.Error())
		} else {
			stamina = resp_obj.GetStamina()
			log.Printf("Current Stamina=%d", stamina)
			//	update_time := resp_obj.GetStartTime()
			//	cur_time := xyutil.CurTimeSec()
			time_left := resp_obj.GetTimeleft()
			log.Printf("time left before update: %ds", time_left)
		}
	}
	return
}
