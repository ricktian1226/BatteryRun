// business
package main

import (
	"testing"

	proto "code.google.com/p/goprotobuf/proto"

	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
)

func TestGetIds(t *testing.T) {
	//
	ids := getIds("1;2;3;4;")
	if len(ids) != 4 || ids[1] != 2 {
		t.Fatalf("result should be [1,2,3,4] but %v", ids)
	}
}

func TestGetIds0(t *testing.T) {
	//
	ids := getIds("a;2;;4;")
	if len(ids) != 2 || ids[1] != 4 {
		t.Fatalf("result should be [2,4] but %v", ids)
	}
}

func TestGetMallItemFlag(t *testing.T) {
	//
	mapId2Flag := MAPId2Flag{
		1: &battery.MallItemFlag{
			Id:   proto.Uint64(1),
			Flag: proto.Int64(2),
		},
	}

	ids := []uint64{1}

	setMallItemFlag(mapId2Flag, ids, battery.MALL_ITEM_FLAG_MALL_ITEM_FLAG_DISPLAY)

	if mapId2Flag[1].GetFlag() != 3 {
		t.Fatalf("flag should be 3 but %d", mapId2Flag[1].GetFlag())
	}
}

func TestGetMallItemFlag0(t *testing.T) {
	//
	mapId2Flag := MAPId2Flag{
		1: &battery.MallItemFlag{
			Id:   proto.Uint64(1),
			Flag: proto.Int64(0),
		},
	}

	ids := []uint64{1}

	setMallItemFlag(mapId2Flag, ids, battery.MALL_ITEM_FLAG_MALL_ITEM_FLAG_GIFT)

	if mapId2Flag[1].GetFlag() != 2 {
		t.Fatalf("flag should be 2 but %d", mapId2Flag[1].GetFlag())
	}
}
