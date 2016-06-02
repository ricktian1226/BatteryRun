package main

import (
	"bytes"
	proto "code.google.com/p/goprotobuf/proto"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"log"
	"net/http"
	"time"
)

func Test_Game(fbid string, count int, delay int) {
	var uid = ""
	time.Sleep(5 * time.Second)
	for i := 0; i < count; i++ {
		log.Printf(">>>>>>>>> Game #%d start <<<<<<<<<<", i)
		uid = test_login(fbid, "")
		if uid != "" {
			stamina := test_query_stamina(uid)
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
				test_user_data(uid)
			}
		} else {
			log.Printf("login failed")
		}
		log.Printf(">>>>>>>>>>> Game #%d done  <<<<<<<<<", i)
		time.Sleep(time.Duration(delay) * time.Second)
	}
}

func test_new_game(uid string) (game_id string) {
	log.Printf("==== test_new_game start")
	defer log.Printf("==== test_new_game done")

	req := battery.NewGameRequest{}
	req.Uid = proto.String(uid)

	//	log.Printf("request: %s", req.String())

	var data_in []byte
	data_in, err := proto.Marshal(&req)
	if err != nil {
		log.Printf("Error:%s", err.Error())
	}

	resp, err := http.Post(server_url+"/v1/newgame/315", "", bytes.NewReader(data_in))
	if err != nil {
		log.Printf("Error: %s", err.Error())
	} else {
		resp_obj := &battery.NewGameResponse{}
		err = handleResp(resp, resp_obj)
		if err != nil {
			log.Printf("Error: %s", err.Error())
		} else {
			//			log.Printf("response: %s", (*resp_obj).String())
			game_id = resp_obj.GetGameId()
			log.Printf("game id = %s", game_id)
		}
	}
	return
}

func test_add_game_data(uid string, game_id string) {
	log.Printf("==== test_add_game start ")
	defer log.Printf("==== test_add_game done ")

	req := battery.GameDataRequest{}

	req.Uid = proto.String(uid)
	req.GameId = proto.String(game_id)
	/*	req.Duration = proto.Int64(30)
		game := &battery.GameResult{
			Score:       proto.Int64(10000),
			Gold:        proto.Int64(1000),
			GoldScore:   proto.Int64(8000),
			Distance:    proto.Int64(2000),
			TotalCharge: proto.Int32(10),
			ComboCharge: proto.Int32(5),
		}
		kills := make([]*battery.MonsterKilled, 3)
		kills[0] = &battery.MonsterKilled{Id: proto.String("m1"), Killed: proto.Int32(1)}
		kills[1] = &battery.MonsterKilled{Id: proto.String("m2"), Killed: proto.Int32(2)}
		kills[2] = &battery.MonsterKilled{Id: proto.String("m3"), Killed: proto.Int32(0)}
		game.MonsterKilled = kills

		gains := make([]*battery.ItemGained, 6)
		gains[0] = &battery.ItemGained{Id: proto.String("10000"), Amount: proto.Int32(1)}
		gains[1] = &battery.ItemGained{Id: proto.String("20000"), Amount: proto.Int32(2)}
		gains[2] = &battery.ItemGained{Id: proto.String("30000"), Amount: proto.Int32(3)}
		gains[3] = &battery.ItemGained{Id: proto.String("40000"), Amount: proto.Int32(4)}
		gains[4] = &battery.ItemGained{Id: proto.String("50000"), Amount: proto.Int32(5)}
		gains[5] = &battery.ItemGained{Id: proto.String("60000"), Amount: proto.Int32(0)}
		game.ItemGained = gains
	*/
	//req.GameResult = game
	//	log.Printf("request: %s", req.String())

	var data_in []byte
	data_in, err := proto.Marshal(&req)
	if err != nil {
		log.Printf("Error:%s", err.Error())
	}

	resp, err := http.Post(server_url+"/v1/gameresult/315", "", bytes.NewReader(data_in))
	if err != nil {
		log.Printf("Error: %s", err.Error())
	} else {
		resp_obj := &battery.GameDataResponse{}
		err = handleResp(resp, resp_obj)
		if err != nil {
			log.Printf("Error: %s", err.Error())
		} else {
			if resp_obj.Data != nil {
				//if resp_obj.GetData().Wallet != nil {
				//	stm := resp_obj.GetData().GetWallet().GetStamina()
				//	log.Printf("stamina = %d", stm)
				//}
				//if resp_obj.GetData().Total != nil {
				//	score := resp_obj.GetData().GetTotal().GetScore()
				//	log.Printf("score = %d", score)
				//}
			}
		}
	}
}
