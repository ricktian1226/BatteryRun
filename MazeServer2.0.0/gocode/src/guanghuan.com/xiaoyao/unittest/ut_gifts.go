package main

import (
	proto "code.google.com/p/goprotobuf/proto"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"log"
)

func unit_confirm_all_gift(fbid string) {
	uid := test_login(fbid, "")
	if uid != "" {
		//		test_query_gift(uid)
		test_confirm_all_gifts(uid)
		test_query_gift(uid)
	}
}
func unit_confirm_gifts_one_by_one(fbid string) {
	uid := test_login(fbid, "")
	if uid != "" {
		gifts := test_query_gift(uid)
		for _, gift := range gifts {
			if gift.GetOpType() == battery.Gift_Op_Ask {
				test_confirm_one_gift(uid, gift.GetGiftId(), battery.Gift_Op_Approve)
			} else if gift.GetOpType() == battery.Gift_Op_Give {
				test_confirm_one_gift(uid, gift.GetGiftId(), battery.Gift_Op_Accept)
			}
		}
	}
}

func unit_new_gifts(fbid string, friend_fbid string, isAsk bool) (gift_id string) {
	uid := test_login(fbid, "")
	if uid != "" {
		//		if gift_fbid == "" {
		//			gift_fbid = new_gift_id(fbid, friend_fbid, isAsk)
		//		}
		gift_id = test_new_gift(uid, friend_fbid, isAsk)
		test_query_gift(uid)
	}
	return
}

func unit_approve_gift(fbid string, gift_fbid string) (isSuccess bool) {
	uid := test_login(fbid, "")
	if uid != "" {
		test_one_gift(uid, gift_fbid, battery.Gift_Op_Approve)
	}
	return
}
func unit_accept_gift(fbid string, gift_fbid string) (isSuccess bool) {
	uid := test_login(fbid, "")
	if uid != "" {
		test_one_gift(uid, gift_fbid, battery.Gift_Op_Accept)
	}
	return
}

func unit_query_gift(fbid string) {
	uid := test_login(fbid, "")
	if uid != "" {
		test_query_gift(uid)
	}
}

func test_query_gift(uid string) (gifts []*battery.Gift) {
	log.Printf("==== test_query_gift start")
	defer log.Printf("==== test_query_gift done")

	req := &battery.QueryStaminaGiftRequest{}
	req.Uid = proto.String(uid)

	resp := &battery.QueryStaminaGiftResponse{}
	unit_test("query gift", "/v1/gift/query/315", req, resp)

	gifts = resp.GetGifts()
	log.Printf("gifts count: %d", len(gifts))
	for i, gift := range gifts {
		//		gifts_id[i] = gift.GetGiftId()
		log.Printf("#%02d.gift[%s]: [%s]->[%s] type:%s", i, gift.GetGiftId(), gift.GetFromId(),
			gift.GetToId(), gift.GetOpType().String())
	}
	return
}

func test_confirm_all_gifts(uid string) {
	log.Printf("==== test_confirm_all_gifts start")
	defer log.Printf("==== test_confirm_all_gifts done")

	req := &battery.StaminaGiftRequest{}
	req.Uid = proto.String(uid)
	req.IsToAll = proto.Bool(true)

	resp := &battery.StaminaGiftResponse{}
	unit_test("confirm all gift", "/v1/gift/op/315", req, resp)
}

func test_confirm_one_gift(uid string, gift_fbid string, op_code battery.Gift_OpType) {
	log.Printf("==== test_confirm_one_gift start")
	defer log.Printf("==== test_confirm_one_gift done")
	log.Printf("op: %s", op_code.String())
	req := &battery.StaminaGiftRequest{
		Uid:     proto.String(uid),
		IsToAll: proto.Bool(false),
		Requests: []*battery.StaminaGiftRequest_SingleRequest{
			&battery.StaminaGiftRequest_SingleRequest{
				GiftSid:   proto.String(gift_fbid),
				Source:    battery.ID_SOURCE_SRC_SINA_WEIBO.Enum(),
				FriendSid: proto.String(""),
			},
		},
		OpType: op_code.Enum(),
	}

	resp := &battery.StaminaGiftResponse{}

	unit_test("confirm one gift", "/v1/gift/op/315", req, resp)
}

func test_new_gift(uid string, friend_fbid string, isAsk bool) (gift_sid string) {
	log.Printf("==== test_new_gift start")
	defer log.Printf("==== test_new_gift done")
	var op_code = battery.Gift_Op_Ask.Enum()
	if !isAsk {
		op_code = battery.Gift_Op_Give.Enum()
	}

	req := &battery.StaminaGiftRequest{
		Uid:     proto.String(uid),
		IsToAll: proto.Bool(false),
		Requests: []*battery.StaminaGiftRequest_SingleRequest{
			&battery.StaminaGiftRequest_SingleRequest{
				GiftSid:   proto.String(""),
				Source:    battery.ID_SOURCE_SRC_SINA_WEIBO.Enum(),
				FriendSid: proto.String(friend_fbid),
			},
		},
		OpType: op_code.Enum(),
	}

	resp := &battery.StaminaGiftResponse{}
	unit_test("new gift", "/v1/gift/op/315", req, resp)

	if len(resp.GiftStatus) > 0 {
		for _, status := range resp.GetGiftStatus() {
			if status.GetIsSuccess() {
				log.Printf("new gift: %s created", status.GetGiftSid())
				gift_sid = status.GetGiftSid()
			} else {
				log.Printf("new gift: %s failed, reason: %d", status.GetGiftSid(), status.GetFailReason())
			}
		}
	}
	return
}
func test_one_gift(uid string, gift_fbid string, op_code battery.Gift_OpType) (isSuccess bool, err error) {

	req := &battery.StaminaGiftRequest{
		Uid:     proto.String(uid),
		IsToAll: proto.Bool(false),
		Requests: []*battery.StaminaGiftRequest_SingleRequest{
			&battery.StaminaGiftRequest_SingleRequest{
				GiftSid:   proto.String(gift_fbid),
				Source:    battery.ID_SOURCE_SRC_SINA_WEIBO.Enum(),
				FriendSid: proto.String(""),
			},
		},
		OpType: op_code.Enum(),
	}

	resp := &battery.StaminaGiftResponse{}
	unit_test("confirm gift", "/v1/gift/op/315", req, resp)

	if len(resp.GiftStatus) > 0 {
		for _, status := range resp.GetGiftStatus() {
			isSuccess = status.GetIsSuccess()
			if status.GetIsSuccess() {
				log.Printf("confirm gift: %s success", status.GetGiftSid())
			} else {
				log.Printf("confirm gift: %s failed, reason: %d", status.GetGiftSid(), status.GetFailReason())
			}
		}
	}
	return
}
func test_approve_gift(uid string, gift_fbid string) {
	log.Printf("==== test_approve_gift start")
	defer log.Printf("==== test_approve_gift done")
	var op_code = battery.Gift_Op_Approve

	test_one_gift(uid, gift_fbid, op_code)
	//	return
}
func test_accept_gift(uid string, gift_fbid string) {
	log.Printf("==== test_accept_gift start")
	defer log.Printf("==== test_accept_gift done")
	var op_code = battery.Gift_Op_Accept

	test_one_gift(uid, gift_fbid, op_code)
}

func process_gift_status(statuses []*battery.StaminaGiftResponse_GiftStatus) {
	if len(statuses) > 0 {
		for i, status := range statuses {
			log.Printf("#%02d gift: %s, success: %d", i, status.GetGiftSid(), status.GetIsSuccess())
		}
	}
}

/*
func new_gift_id(from string, to string, isAsk bool) (name string) {
	if isAsk {
		name = fmt.Sprintf("gift_ask_%s_%s_%d", from, to, xyutil.CurTimeNs())
	} else {
		name = fmt.Sprintf("gift_give_%s_%s_%d", from, to, xyutil.CurTimeNs())
	}
	return
}
*/
