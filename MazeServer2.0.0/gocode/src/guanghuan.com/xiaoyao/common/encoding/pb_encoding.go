package xyencoder

import (
	proto "code.google.com/p/goprotobuf/proto"
	//	xylog "guanghuan.com/xiaoyao/common/log"
)

//type PbNatsEncoder struct {
//}

/*
encoder interface of NATS:
   func (ge *GobEncoder) Decode(subject string, data []byte, vPtr interface{}) (err error)
   func (ge *GobEncoder) Encode(subject string, v interface{}) ([]byte, error)
*/
/*
func (enc *PbNatsEncoder) Decode(subject string, data []byte, vPtr interface{}) (err error) {
	return
}
func (enc *PbNatsEncoder) Encode(subject string, v interface{}) (data []byte, err error) {
	return
}
*/

// 编码
func PbEncode(msg_in proto.Message) (data_out []byte, err error) {
	data_out, err = proto.Marshal(msg_in)
	//if err == nil {
	//	xylog.Debug("msg enc: %s", msg_in.String())
	//}
	return
}

// 解码
func PbDecode(data_in []byte, obj_out proto.Message) (err error) {
	err = proto.Unmarshal(data_in, obj_out)
	//if err == nil {
	//	xylog.Debug("msg dec: %s", obj_out.String())
	//}
	return
}
