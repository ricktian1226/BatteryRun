package batteryapi

import (
	//"fmt"
	"testing"

	proto "code.google.com/p/goprotobuf/proto"

	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
)

//func Test_example(t *testing.T) {
//	var a = 1
//	t.Logf("here")
//	if a == 1 {
//		t.Fatalf("a == 1")
//	}
//}

// distinctDeviceToken 单元测试
func Test_distinctDeviceToken_equal(t *testing.T) {

	//相同的deviceid去重
	var accounts = []*battery.DBAccount{
		&battery.DBAccount{
			Deviceid: proto.String("abcdef"),
		},
		&battery.DBAccount{
			Deviceid: proto.String("abcdef"),
		},
	}

	accounts = NewXYAPI().distinctDeviceToken(accounts)

	if len(accounts) != 1 ||
		accounts[0].GetDeviceid() != "abcdef" {
		t.Fatalf("distinctDeviceToken 0 failed")
	}

}

// distinctDeviceToken 单元测试
func Test_distinctDeviceToken_noequal(t *testing.T) {
	//不同的deviceid保留
	accounts := []*battery.DBAccount{
		&battery.DBAccount{
			Deviceid: proto.String("abcdef"),
		},
		&battery.DBAccount{
			Deviceid: proto.String("123456"),
		},
	}

	accounts = NewXYAPI().distinctDeviceToken(accounts)

	if len(accounts) != 2 ||
		accounts[0].GetDeviceid() != "abcdef" ||
		accounts[1].GetDeviceid() != "123456" {
		t.Fatalf("distinctDeviceToken 1 failed")
	}
}

// distinctDeviceToken 单元测试
func Test_distinctDeviceToken_null(t *testing.T) {
	//不同的deviceid保留
	accounts := []*battery.DBAccount{
		&battery.DBAccount{
			Deviceid: proto.String("abcdef"),
		},
		&battery.DBAccount{
			Deviceid: proto.String(""),
		},
	}

	accounts = NewXYAPI().distinctDeviceToken(accounts)

	if len(accounts) != 1 ||
		accounts[0].GetDeviceid() != "abcdef" {
		t.Fatalf("distinctDeviceToken 1 failed")
	}
}
