// main
package main

import (
	proto "code.google.com/p/goprotobuf/proto"
	"flag"
	"fmt"
	httplib "github.com/astaxie/beego/httplib"
	xyencoder "guanghuan.com/xiaoyao/common/encoding"
	xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	crypto "guanghuan.com/xiaoyao/superbman_server/crypto"
	"time"
)

var HTTP_URL string = "192.168.21.91:10003"

var timewait time.Duration = 10000 * 1000

func test_op(req proto.Message, resp proto.Message, s string, rwtimeout time.Duration) (err error) {
	var (
		data    []byte
		api_url string = "http://" + HTTP_URL + s
	)
	data, err = xyencoder.PbEncode(req)
	if err != nil {
		xylog.Error("Error Encoding: %s", err.Error())
		return err
	}
	xylog.Debug("reqdata after PbEncode: [%d][%v]\n", len(data), data)

	data, err = crypto.Encrypt(data)
	if err != nil {
		xylog.Error("Error Encrypt: %s", err.Error())
		return err
	}
	//xylog.Debug("reqdata after Encrypt: [%d][%v]\n", len(data), data)

	fmt.Println(api_url)
	request := httplib.Post(api_url)
	request.Body(data)
	request.SetTimeout(60*time.Second, rwtimeout*time.Microsecond)

	data, err = request.Bytes()
	if err != nil {
		xylog.Error("Error Send: %s", err.Error())
		return err
	}
	//xylog.Debug("respdata after Send: [%d][%v]\n", len(data), data)

	data, err = crypto.Decrypt(data)
	if err != nil {
		xylog.Error("Error Decrypt: %s", err.Error())
		return err
	}
	//xylog.Debug("respdata after Decrypt: [%d][%v]\n", len(data), data)

	err = xyencoder.PbDecode(data, resp)
	if err != nil {
		xylog.Error("Error Decoding: %s", err.Error())
		return err
	}

	if err == nil {
		xylog.Debug("response : %v", resp)
	}

	return err
}

func main() {

	flag.StringVar(&HTTP_URL, "url", HTTP_URL, "gateway server url")
	xylog.ProcessCmdAndApply()

	var err error
	var uid string

	//--------------------roleinfolist_op--------------------
	/*
		reqRole := &battery.RoleInfoListRequest{}
		uid = "1411475570219295000099"
		reqRole.Uid = proto.String(uid)
		var cmdRole battery.RoleInfoListCmd
		var roleID uint64
		reqRole.Cmd = &cmdRole
		reqRole.Roleid = &roleID

		//test 获取角色列表
		cmdRole = battery.RoleInfoListCmd_RoleInfoListCmd_Rolelist

		//test 选中某个角色
		cmdRole = battery.RoleInfoListCmd_RoleInfoListCmd_Select
		roleID = 150020000

		//test  升级某个角色
		cmdRole = battery.RoleInfoListCmd_RoleInfoListCmd_Upgrade
		roleID = 150050003

		xylog.Debug("reqT : %v", reqRole)

		respRole := &battery.RoleInfoListResponse{}

		err = test_op(reqRole, respRole, "/v2/roleinfolist/roleinfolist_op", timewait)

		if err != nil {
			xylog.Error("======== roleinfolist_config import failed =========")
		} else {
			xylog.Debug("======== roleinfolist_config import succeed =========")
		}
	*/
	//-------------------system_mail_op--------------------
	/*
		reqSMail := &battery.SystemMailListRequest{}
		uid = "1411475570219295000099"
		reqSMail.Uid = proto.String(uid)

		var cmdMail battery.SystemMailCmd
		var mailID int32
		reqSMail.Cmd = &cmdMail
		reqSMail.MailID = &mailID

		//test 获取邮件列表
		cmdMail = battery.SystemMailCmd_SystemMailCmd_Maillist

		//cmdMail = battery.SystemMailCmd_SystemMailCmd_MailRead
		//mailID = 10001

		//test 领取礼包
		//cmdMail = battery.SystemMailCmd_SystemMailCmd_GiftGet
		//mailID = 10002

		xylog.Debug("reqSMail : %v", reqSMail)

		respSMail := &battery.SystemMailListResponse{}

		err = test_op(reqSMail, respSMail, "/v2/systemmail/systemmail_op", timewait)

		if err != nil {
			xylog.Error("======== system_mail_op failed =========")
		} else {
			xylog.Debug("========system_mail_op succeed =========")
		}*/

	//-------------------friend_mail_op--------------------
	/**/
	reqFriendMail := &battery.FriendMailListRequest{}
	//uid = "1414395373585955463059"
	uid = "1414403548023489495847"
	reqFriendMail.Uid = proto.String(uid)

	var cmdMail battery.FriendMailCmd
	reqFriendMail.Cmd = &cmdMail
	//reqFriendMail.FriendSid = proto.String("12580")
	//reqFriendMail.FriendSid = proto.String("10000")
	var source battery.ID_SOURCE = battery.ID_SOURCE_SRC_SINA_WEIBO
	reqFriendMail.Source = &source

	//reqFriendMail.CreateTime = proto.Int64(1411480222)

	//test 获取邮件列表
	cmdMail = battery.FriendMailCmd_FriendMailCmd_MailList
	//cmdMail = battery.FriendMailCmd_FriendMailCmd_StaminaGive //赠送
	//cmdMail = battery.FriendMailCmd_FriendMailCmd_StaminaGetApply //赠送申请
	//cmdMail = battery.FriendMailCmd_FriendMailCmd_MailStaminaGive //体力赠送（邮件里）
	//cmdMail = battery.FriendMailCmd_FriendMailCmd_StaminatGet    //领取
	cmdMail = battery.FriendMailCmd_FriendMailCmd_StaminatGetAll //领取全部

	respFriendMail := &battery.FriendMailListResponse{}

	err = test_op(reqFriendMail, respFriendMail, "/v2/friendmail/friendmail_op", timewait)

	xylog.Debug("respFriendMail : %v", respFriendMail)
	if err != nil {
		xylog.Error("======== friend_mail_op failed =========")
	} else {
		xylog.Debug("========friend_mail_op succeed =========")
	}

	//-------------------jigsaw_op--------------------
	/*
		reqJigsaw := &battery.JigsawRequest{}
		uid = "1411475570219295000099"
		reqJigsaw.Uid = proto.String(uid)

		var cmdJigsaw battery.JigsawCmd
		reqJigsaw.Cmd = &cmdJigsaw
		reqJigsaw.Itemid = proto.Uint64(120010001)

		//test
		cmdJigsaw = battery.JigsawCmd_JigsawCmd_Query //赠送
		cmdJigsaw = battery.JigsawCmd_JigsawCmd_Buy   //赠送申请

		respJigsaw := &battery.JigsawResponse{}

		err = test_op(reqJigsaw, respJigsaw, "/v2/jigsaw/jigsaw_op", timewait)

		//reqJigsaw.Itemid = proto.Uint64(120010002)
		//cmdJigsaw = battery.JigsawCmd_JigsawCmd_Query //赠送
		//cmdJigsaw = battery.JigsawCmd_JigsawCmd_Buy   //赠送申请
		//err = test_op(reqJigsaw, respJigsaw, "/v2/jigsaw/jigsaw_op", timewait)

		//reqJigsaw.Itemid = proto.Uint64(120010003)
		//cmdJigsaw = battery.JigsawCmd_JigsawCmd_Query //赠送
		//cmdJigsaw = battery.JigsawCmd_JigsawCmd_Buy   //赠送申请
		//err = test_op(reqJigsaw, respJigsaw, "/v2/jigsaw/jigsaw_op", timewait)

		//reqJigsaw.Itemid = proto.Uint64(120010004)
		//cmdJigsaw = battery.JigsawCmd_JigsawCmd_Query //赠送
		//cmdJigsaw = battery.JigsawCmd_JigsawCmd_Buy   //赠送申请
		//err = test_op(reqJigsaw, respJigsaw, "/v2/jigsaw/jigsaw_op", timewait)

		if err != nil {
			xylog.Error("======== jigsaw_op failed =========")
		} else {
			xylog.Debug("========jigsaw_op succeed =========")
		}
	*/
	time.Sleep(50 * time.Second)
}
