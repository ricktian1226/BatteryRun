package main

import (
	proto "code.google.com/p/goprotobuf/proto"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"log"
)

// http 查询
/*
func test_confirm_one_gift(uid string, gift_fbid string, op_code battery.Gift_OpType) {
	log.Printf("==== test_confirm_one_gift start")
	defer log.Printf("==== test_confirm_one_gift done")
	log.Printf("op: %s", op_code.String())
	req := &battery.StaminaGiftRequest{}
	req.Uid = proto.String(uid)
	req.GiftSid = proto.String(gift_fbid)
	req.Source = battery.ID_SOURCE_FACEBOOK.Enum()
	req.OpType = op_code.Enum()

	resp := &battery.StaminaGiftResponse{}

	unit_test("confirm one gift", "/v1/gift/op/315", req, resp)
}
*/
func ut_query_goods(fbid string) {
	uid := test_login(fbid, "")
	if uid != "" {
		HttpQueryGoods(uid, 0)
	}
}
func ut_buy_goods(fbid string, goods_id string) {
	uid := test_login(fbid, "")
	if uid != "" {
		HttpBuyGoods(uid, goods_id)
	}
}

func HttpQueryGoods(uid string, cat_id int32) (err error) {
	log.Printf("==== HttpQueryGoods start")
	defer log.Printf("==== HttpQueryGoods done")
	/*
		req := &battery.QueryGoodsRequest{
			Uid:      proto.String(uid),
			Category: proto.Int32(cat_id),
		}
		resp := &battery.QueryGoodsResponse{}

		err = unit_test("query goods", "/v1/goods/query/315", req, resp)

		if err == nil {
			goods_list := resp.GetGoodsList()
			for i, goods := range goods_list {
				units := goods.GetAllGoods()
				log.Printf("#%d.[%s] contains %d units", i+1, goods.GetId(), len(units))
				for j, unit := range units {
					log.Printf("\t%d. %s x %d = %s x %d", j+1,
						unit.GetResource().String(), unit.GetAmount(),
						goods.GetPrice().GetCurrency().String(), goods.GetPrice().GetAmount())
				}
			}
		}
	*/
	return
}

func HttpBuyGoods(uid string, goods_id string) (err error) {
	log.Printf("==== HttpBuyGoods start")
	defer log.Printf("==== HttpBuyGoods done")
	log.Printf("uid(%s) buy goods %s", uid, goods_id)
	req := &battery.BuyGoodsRequest{
		Uid: proto.String(uid),
		//GoodsId: proto.String(goods_id),
	}

	resp := &battery.BuyGoodsResponse{}

	err = unit_test("HttpBuyGoods", "/v1/goods/buy/315", req, resp)
	if err == nil {
		log.Printf("receipt: %s", resp.GetReceiptId())
	}
	return
}
