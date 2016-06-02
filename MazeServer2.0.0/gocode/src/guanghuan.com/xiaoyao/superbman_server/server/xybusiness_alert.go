// xybusiness
package xybusiness

import (
	proto "code.google.com/p/goprotobuf/proto"
	"fmt"
	beegoconf "github.com/astaxie/beego/config"
	"guanghuan.com/xiaoyao/common/log"
	"guanghuan.com/xiaoyao/common/mail"
	xynatsservice "guanghuan.com/xiaoyao/common/service/nats"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	"strings"
)

const (
	INI_CONFIG_ITEM_SMTP_HOST   = "Smtp::host"
	INI_CONFIG_ITEM_SMTP_PORT   = "Smtp::port"
	INI_CONFIG_ITEM_SMTP_USER   = "Smtp::user"
	INI_CONFIG_ITEM_SMTP_PASSWD = "Smtp::passwd"
	INI_CONFIG_ITEM_SMTP_FROM   = "Smtp::from"
	INI_CONFIG_ITEM_SMTP_TO     = "Smtp::to"
)

const ALERT_API = "alert"

const (
	ALERT_EVENT_SERVER_START = "Server.Start"
)

//告警消息结构体
// subect string 告警标题
// content string 告警内容
type AlertStruct struct {
	Uid, Subect, Content string
}

//根据ini配置文件初始化mail配置项
// config *beegoconf.ConfigContainer ini配置项容器对象
func Init(config *beegoconf.ConfigContainer) error {

	if !checkIniItem(INI_CONFIG_ITEM_SMTP_HOST, config, &(xymail.DefSmtpConfig.Host)) {
		return xyerror.ErrBadInputData
	}

	if !checkIniItem(INI_CONFIG_ITEM_SMTP_PORT, config, &(xymail.DefSmtpConfig.Port)) {
		return xyerror.ErrBadInputData
	}

	if !checkIniItem(INI_CONFIG_ITEM_SMTP_USER, config, &(xymail.DefSmtpConfig.User)) {
		return xyerror.ErrBadInputData
	}

	if !checkIniItem(INI_CONFIG_ITEM_SMTP_PASSWD, config, &(xymail.DefSmtpConfig.Passwd)) {
		return xyerror.ErrBadInputData
	}

	if !checkIniItem(INI_CONFIG_ITEM_SMTP_FROM, config, &(xymail.DefSmtpConfig.From)) {
		return xyerror.ErrBadInputData
	}

	var tmp string
	if !checkIniItem(INI_CONFIG_ITEM_SMTP_TO, config, &(tmp)) {
		return xyerror.ErrBadInputData
	}
	xymail.DefSmtpConfig.To = strings.Split(tmp, ",")

	return xyerror.ErrOK
}

func InitAlertNats(alertNatsService *xynatsservice.NatsService) {
	natservice = alertNatsService
}

//校验配置项，如果配置项不存在或者为""，则校验失败
// item string 配置项名称
// config *beegoconf.ConfigContainer 配置项管理器
//return:
// value *string 配置值
func checkIniItem(item string, config *beegoconf.ConfigContainer, value *string) (valid bool) {

	tmp := (*config).String(item)
	if tmp == "" {
		xylog.Error("%s invalid.", item)
		valid = false
	}
	*value = tmp

	return true
}

var (
	natservice *xynatsservice.NatsService //nats服务指针
)

//发送业务告警邮件
// subject string 邮件标题
// message string 邮件正文
// isHtml bool 邮件正文是否使用html格式
func SendMail(subject, message, node string, isHtml bool) error {
	subject = string("[BusinessAlert!!!] ") + subject

	message = string("===== This mail is sent by system, please don't reply =====\n") +
		fmt.Sprintf("from %s:\n", node) +
		message +
		string("\n==============\n")

	return xymail.Send(subject, message, isHtml)
}

//业务发送告警
// uid string 玩家id
// subject string 邮件标题
// message string 邮件正文
func SendAlert(uid, subject, content, node string) (err error) {
	var (
		alert = &battery.BusinessAlert{
			Uid:     &uid,
			Subject: &subject,
			Content: &content,
			Node:    &node,
		}
		data []byte
	)

	data, err = MarshalAlert(alert)
	if err != xyerror.ErrOK {
		xylog.Error(uid, "MarshalAlert failed : %v", err)
		return
	}

	natservice.Publish(ALERT_API, data)

	return
}

//告警信息序列化
// alert *battery.BusinessAlert 告警信息结构体
//return:
// data []byte 序列化后的字节序列
func MarshalAlert(alert *battery.BusinessAlert) (data []byte, err error) {
	data, err = proto.Marshal(alert)
	return
}

//告警信息反序列化
// alert *battery.BusinessAlert 保存告警信息的结构体指针
// data []byte 字节序列
func UnMarshalAlert(alert *battery.BusinessAlert, data []byte) (err error) {
	err = proto.Unmarshal(data, alert)
	return
}
