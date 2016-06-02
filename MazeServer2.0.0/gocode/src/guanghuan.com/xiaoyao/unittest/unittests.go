package main

import (
	"bytes"
	proto "code.google.com/p/goprotobuf/proto"
	//	"flag"
	"fmt"
	xyutil "guanghuan.com/xiaoyao/common/util"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

//var server_url = "http://localhost:12356"
//var server_url = "http://localhost:8888"

//var server_url = "http://117.25.150.86:12356"

var server_url = "http://192.168.1.205:10003"

func main() {
	log.Printf("Server: %s", server_url)
	rand.Seed(time.Now().UnixNano())

	var err error
	test_count := 1
	delay := 10
	args := os.Args[1:]

	test := "all"
	fbid := ""
	if len(args) >= 2 {
		test = args[0]
		fbid = args[1]
		args = args[2:]
	} else {
		log.Printf("usage: <test> <fbid> ...")
		return
	}
	log.Printf("test bed: [%s], fbid: [%s], args:%v", test, fbid, args)

	//	for i := 0; i < 1; i++ {
	//		fbid := fmt.Sprintf("fb%03d", i)
	//		//		test_login(fbid)
	//		unit_test_gift(fbid)
	//	}

	// 体力请求
	if test == "gift/new" {
		if len(args) < 2 {
			log.Printf("gift/new <your_fbid> <friend_fbid> <ask|give>")
		} else {
			friend_fbid := args[0]
			//			gift_fbid := args[1]
			isAsk := (args[1] == "ask")
			unit_new_gifts(fbid, friend_fbid, isAsk)
		}
	}
	if test == "gift/query" {
		unit_query_gift(fbid)
	}

	if test == "gift/confirm/all" {
		unit_confirm_all_gift(fbid)
	}
	/*
		if test == "gift/confirm/one" {
			if len(args) < 2 {
				log.Printf("gif/confirm/one <your fbid> <approve|accept> <gift1_fbid>, <gift2_fbid> ...")
			} else {
				isApprove := (args[0] == "approve")
				args = args[1:]
				unit_approve_gift(fbid, gift_fbid)
			}
		}
	*/
	if test == "gift/ask" {
		if len(args) < 1 {
			log.Printf("gift/ask <your_fbid> <friend_fbid>")
		} else {
			friend_fbid := args[0]
			//			gift_fbid := args[1]
			unit_new_gifts(fbid, friend_fbid, true)
		}
	}
	if test == "gift/give" {
		if len(args) < 1 {
			log.Printf("gift/ask <your_fbid> <friend_fbid>")
		} else {
			friend_fbid := args[0]
			//			gift_fbid := args[1]
			unit_new_gifts(fbid, friend_fbid, false)
		}
	}
	if test == "gift/approve" {
		if len(args) < 1 {
			log.Printf("gift/approve <your_fbid> <gift_id>")
		} else {
			gift_fbid := args[0]
			unit_approve_gift(fbid, gift_fbid)
		}
	}
	if test == "gift/accept" {
		if len(args) < 1 {
			log.Printf("gift/accept <your_fbid> <gift_id>")
		} else {
			gift_fbid := args[0]
			unit_accept_gift(fbid, gift_fbid)
		}
	}

	// 商品
	if test == "goods/query" {
		ut_query_goods(fbid)
	}

	if test == "goods/buy" {
		if len(args) < 1 {
			log.Printf("goods/buy <your_fbid> <goods_id>")
		} else {
			goods_id := args[0]
			ut_buy_goods(fbid, goods_id)
		}
	}

	if test == "game" {

		if len(args) > 1 {
			test_count, err = strconv.Atoi(args[0])
			if err != nil {
				test_count = 1
			}
		}
		if len(args) > 2 {
			delay, err = strconv.Atoi(args[1])
			if err != nil {
				delay = 10
			}
		}

		log.Printf("Testing %s, times: %d, delay: %d", server_url, test_count, delay)

		Test_Game(fbid, test_count, delay)
	}

	if test == "normal" {
		if len(args) > 1 {
			test_count, err = strconv.Atoi(args[0])
			if err != nil {
				test_count = 1
			}
		}
		if len(args) > 2 {
			delay, err = strconv.Atoi(args[1])
			if err != nil {
				delay = 10
			}
		}

		log.Printf("Testing %s, times: %d, delay: %d", server_url, test_count, delay)
		fbid := "fb000"
		Test_Normal(fbid, delay)
	}

	if test == "user/new/rand" {
		var count int = 1
		if len(args) > 0 {
			count, _ = strconv.Atoi(args[0])
		}
		if count > 0 {
			for i := 0; i < count; i++ {
				fbid := fmt.Sprintf("gc%03d%d", xyutil.CurTimeSec(), i)
				test_login(fbid, fbid)
			}
		}
	}

	if test == "user/new" {
		var (
			from  int
			count int
		)
		if len(args) < 2 {
			log.Printf("user/new <prefix> <from> <count>")
		} else {
			prefix := fbid
			from, _ = strconv.Atoi(args[0])
			count, _ = strconv.Atoi(args[1])

			if count > 0 {
				for i := 0; i < count; i++ {
					login_id := fmt.Sprintf("%s%03d", prefix, i+from)
					//		name := fmt.Sprintf("%s-name-%03d", prefix, i+from)
					test_login(login_id, login_id)
				}
			}
		}
	}
	if test == "user/login" {
		test_login(fbid, "")
	}
	if test == "user/query" {
		if len(args) < 1 {
			log.Printf("user/query <fbid>")
		} else {
			//	fbid := args[0]
			return
		}
	}
	if test == "echo" {
		test_echo()
	}
	if test == "friend" {
		friends := args
		if len(friends) <= 0 {
			log.Printf("friend <your id> <friend 1> <friend 2> ...")
			friends = make([]string, 3)
			friends[0] = "100007944793956"
			friends[1] = "100001062771429"
			friends[2] = "100008018798533"
		}
		uid := test_login(fbid, "")
		if uid != "" {
			log.Printf("friends: %v", friends)
			test_friend_data(uid, friends)
		}
	}
}

func handleResp(resp *http.Response, obj_out proto.Message) (err error) {
	data_out, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Printf("Response Error: %s", err.Error())
	} else {
		err = proto.Unmarshal(data_out, obj_out)
		if err != nil {
			log.Printf("Response Unmarshal Error:%s", err.Error())
		}
	}
	return err
}

func unit_test(name string, uri string, req proto.Message, resp proto.Message) error {
	log.Printf("==== Unit test (%s)[%s] start ", name, uri)
	defer log.Printf("==== Unit test (%s) done", name)

	var data_in []byte
	data_in, err := proto.Marshal(req)
	if err != nil {
		log.Printf("Error:%s", err.Error())
	}
	ret_resp, err := http.Post(server_url+uri, "", bytes.NewReader(data_in))
	if err != nil {
		log.Printf("Error: %s", err.Error())
	} else {
		//		resp := &battery.QueryStaminaGiftResponse{}
		err = handleResp(ret_resp, resp)
		if err != nil {
			log.Printf("Error: %s", err.Error())
		} else {
			log.Printf("resp: %v", resp)
		}
	}
	return err
}

func test_echo() {
	req := &battery.Request{
		Data: proto.String(""),
	}
	resp := &battery.Response{}
	unit_test("echo", "/echo", req, resp)

}
