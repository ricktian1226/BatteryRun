// a service wrapper to call Apple Push Notification service (apn)
package xyapnservice

import (
	xylog "guanghuan.com/xiaoyao/common/log"
	xyservice "guanghuan.com/xiaoyao/common/service"
	apns "guanghuan.com/xiaoyao/external/apns"
	"time"
)

const (
	APN_Push_Gateway     = "gateway.push.apple.com:2195"
	APN_Push_Sandbox     = "gateway.sandbox.push.apple.com:2195"
	APN_Feedback_Sandbox = "feedback.sandbox.push.apple.com:2196"
	APN_Feedback_Gateway = "feedback.push.apple.com:2196"
)
const (
	APN_Cert_Dev = "certs/aps_dev_cert.pem"
	APN_Key_Dev  = "certs/key.unencrypted.pem"
	APN_Cert_Pro = "certs/aps_pro_cert.pem"
	APN_Key_Pro  = "certs/key.unencrypted.pem"
)

type ApnService struct {
	xyservice.DefaultService
	push_gateway string
	cert_file    string
	key_file     string
	apnPnClient  *apns.Client
	async_mode   bool
	resp_handler ResponseHandler
}

type ResponseHandler func(sender string, pn_resp apns.PushNotificationResponse) (err error)

//type FeedbackHandler func(sender string, fb_resp apns.FeedbackResponse) (err error)

///////////////////////////
func NewApnService(svc_name string, pn_url string, cert_file string, key_file string) (svc *ApnService) {
	svc = &ApnService{
		DefaultService: *xyservice.NewDefaultService(svc_name),
		push_gateway:   pn_url,
		cert_file:      cert_file,
		key_file:       key_file,
		apnPnClient:    apns.NewClient(pn_url, cert_file, key_file),
		resp_handler:   DefaultResponseHandler,
		async_mode:     false,
	}
	svc.apnPnClient.AsyncRespHandler = svc.OnAsyncResponse
	return
}

func DefaultApnService(svc_name string, use_production_gateway bool, cert_file string, key_file string) (svc *ApnService) {
	var (
		pnurl string
		//		fburl string
		cert string
		key  string
	)
	if use_production_gateway {
		pnurl = APN_Push_Gateway
		//		fburl = APN_Feedback_Gateway
	} else {
		pnurl = APN_Push_Sandbox
		//		fburl = APN_Feedback_Sandbox
	}

	if cert_file != "" {
		cert = cert_file
	} else {
		if use_production_gateway {
			cert = APN_Cert_Pro
		} else {
			cert = APN_Cert_Dev
		}
	}

	if key_file != "" {
		key = key_file
	} else {
		if use_production_gateway {
			key = APN_Key_Pro
		} else {
			key = APN_Key_Dev
		}
	}

	svc = NewApnService(svc_name, pnurl, cert, key)

	return
}

func DefaultResponseHandler(sender string, pn_resp apns.PushNotificationResponse) (err error) {
	xylog.DebugNoId("Default handler: [%s] got response: %s", sender, pn_resp.ToString())
	return
}

//func DefaultRespTimeoutHandler() (err error) {
//	return
//}
//func DefaultFeedbackHandler(sender string, resp apns.FeedbackResponse) (err error) {
//	return
//}

func (svc *ApnService) SetAsyncMode(async bool, h ResponseHandler) {
	svc.apnPnClient.SetAsyncMode(async)
	if h == nil {
		svc.resp_handler = DefaultResponseHandler
	} else {
		svc.resp_handler = h
	}

}

func (svc *ApnService) SetTimeout(t time.Duration) {
	//	svc.resp_timeout = t
	svc.apnPnClient.SetTimeout(t)
	return
}

//func (svc *ApnService) SetResponseHanlder(h ResponseHandler) {
//	svc.resp_handler = h
//}

func (svc *ApnService) OnAsyncResponse(resp apns.PushNotificationResponse) (err error) {
	if !resp.Success {
		svc.Reset()
	}
	return svc.resp_handler(svc.Name(), resp)
}

//func (svc *ApnService) SetResponseTimeoutHandler(h TimeoutHandler) {
//}

//func (svc *ApnService) AsyncSend(pn apns.PushNotification) (err error) {
//	var (
//		resp apns.PushNotificationResponse
//	)
//	resp, err = svc.ApnPnClient.Send(pn)
//	if err != nil || !resp.Success {
//		xylog.Error("Error sending push notification: %s", err.Error())
//		svc.ApnPnClient.Close()
//	}
//	err = svc.resp_handler(svc.Name(), resp)
//	return
//}

func (svc *ApnService) Send(pn apns.PushNotification) (resp apns.PushNotificationResponse, err error) {
	resp, err = svc.apnPnClient.Send(pn)
	if err != nil {
		xylog.ErrorNoId("Error sending push notification: %s", err.Error())
	}
	if err != nil || !resp.Success {
		// 只要有错误就把连接重置
		svc.Reset()
	}
	//if !svc.async_mode {
	//	svc.resp_handler(svc.Name(), resp)
	//}
	return
}
func (svc *ApnService) Reset() {
	svc.apnPnClient.Close()
}

///////////////////////////
// XYService interface
func (svc *ApnService) Init() (err error) {
	err = svc.DefaultService.Init()
	return
}
func (svc *ApnService) Start() (err error) {
	if svc.IsRunning() {
		return
	}
	err = svc.DefaultService.Start()
	return
}
func (svc *ApnService) Stop() (err error) {
	if svc.IsRunning() {
		svc.apnPnClient.Close()
		err = svc.DefaultService.Stop()
	}
	return
}
