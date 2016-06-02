// xymail
//后台邮件公共类，可以为各个服务调用
package xymail

import (
	"fmt"
	//"guanghuan.com/xiaoyao/common/log"
	"net/smtp"
	"strings"
)

//smtp服务器信息
type SmtpConfig struct {
	User   string   //smtp用户
	Passwd string   //smtp用户密码
	Host   string   //smtp服务器地址,e.g "smtp.163.com"
	Port   string   //smtp服务器端口，e.g "25"
	From   string   //邮件发送者
	To     []string //邮件接受者列表
}

//全局的smtp配置
var DefSmtpConfig SmtpConfig

//邮件发送接口定义
// subject string 邮件标题
// message string 邮件内容
// from string 邮件发送者
// to []string 邮件接受者，可以是多个
// smtpConfig SmtpConfig smtp服务器配置信息
// isHtml bool 邮件内容是否是html格式文本
func Send(subject, message string /* from string, to []string, , smtpConfig SmtpConfig*/, isHtml bool) error {

	auth := smtp.PlainAuth(
		"",
		DefSmtpConfig.User,
		DefSmtpConfig.Passwd,
		DefSmtpConfig.Host,
	)
	contentType := "text/plain"
	if isHtml {
		contentType = "text/html"
	}
	msg := fmt.Sprintf("To: %s\r\nFrom: %s\r\nSubject: %s\r\nContent-Type: %s; charset=UTF-8\r\n\r\n%s", strings.Join(DefSmtpConfig.To, ";"), DefSmtpConfig.From, subject, contentType, message)
	err := smtp.SendMail(
		DefSmtpConfig.Host+":"+DefSmtpConfig.Port,
		auth,
		DefSmtpConfig.User, //注意这里的from必须和鉴权的user一致
		DefSmtpConfig.To,
		[]byte(msg),
	)

	return err
}
