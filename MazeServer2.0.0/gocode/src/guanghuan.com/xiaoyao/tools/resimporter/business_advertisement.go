// business_advertisement
package main

import (
	proto "code.google.com/p/goprotobuf/proto"
	xylog "guanghuan.com/xiaoyao/common/log"
	//"guanghuan.com/xiaoyao/common/util"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	"strconv"
	"strings"
)

//广告信息
const (
	AdvertisementId          = iota //广告标识
	AdvertisementViewUrl            //广告曝光链接
	AdvertisementMaterialUrl        //广告素材链接
	AdvertisementClickUrl           //广告点击链接
)

//广告位信息
const (
	AdvertisementSpaceId     = iota //广告位标识
	AdvertisementSpaceItems         //广告列表
	AdvertisementSpaceEnable        //是否播放广告
	AdvertisementSpaceFlag          //播放标识
)

//解析广告信息
type Advertisement struct{}

func (a *Advertisement) Analyse(lineNum int, line string, ops *[]*battery.ResOpItem, optype *battery.RES_OP_TYPE) {
	subs := strings.Split(line, SEP)
	if 0 >= len(subs) {
		return
	}
	advertisement := new(battery.Advertisement)

	for i, v := range subs {
		switch i {
		case AdvertisementId:
			convertAtoui(v, &advertisement.Id)
		case AdvertisementViewUrl:
			advertisement.ViewUrl = proto.String(v)
		case AdvertisementMaterialUrl:
			advertisement.MaterialUrl = proto.String(v)
		case AdvertisementClickUrl:
			advertisement.ClickUrl = proto.String(v)
		}
	}
	op := &battery.ResOpItem{
		Optype: optype,
	}

	op.Advertisement = advertisement
	*ops = append(*ops, op)

}

//解析广告位信息
type AdvertisementSpace struct{}

func (a *AdvertisementSpace) Analyse(lineNum int, line string, ops *[]*battery.ResOpItem, optype *battery.RES_OP_TYPE) {
	subs := strings.Split(line, SEP)
	if 0 >= len(subs) {
		return
	}
	advertisementSpace := &battery.AdvertisementSpace{}

	for i, v := range subs {
		switch i {
		case AdvertisementSpaceId:
			convertAtoui(v, &advertisementSpace.Id)
		case AdvertisementSpaceItems:
			a.AnalyseItem(v, &advertisementSpace.Items)
		case AdvertisementSpaceEnable:
			enable, _ := strconv.ParseBool(v)
			advertisementSpace.Enable = proto.Bool(enable)
		case AdvertisementSpaceFlag:
			a.AnalyseFlag(v, &advertisementSpace.Flags)
		}
	}
	op := &battery.ResOpItem{
		Optype: optype,
	}

	op.AdvertisementSpace = advertisementSpace
	*ops = append(*ops, op)

}

// 解析item字段
func (a *AdvertisementSpace) AnalyseItem(v string, items *[]*battery.AdvertisementWeight) {

	subs := strings.Split(v, SEP1)
	if 0 >= len(subs) {
		xylog.ErrorNoId("AdvertisementSpace.AnalyseItem(%v) failed.", v)
		return
	}

	for _, sub := range subs {
		values := strings.Split(sub, SEP2)
		if 2 != len(values) {
			xylog.ErrorNoId("AdvertisementSpace.AnalyseSubItem(%v) failed.", sub)
			return
		}

		advertisementWeight := &battery.AdvertisementWeight{}
		convertAtoui(values[0], &advertisementWeight.Id)
		convertAtoui(values[1], &advertisementWeight.Weight)

		*items = append(*items, advertisementWeight)
	}
}

func (a *AdvertisementSpace) AnalyseFlag(v string, flags *[]*battery.AdvertisementFlag) {
	subs := strings.Split(v, SEP1)
	if 0 >= len(subs) {
		xylog.ErrorNoId("AdvertisementSpace.AnalyseFlag(%v) failed.", v)
		return
	}

	for _, sub := range subs {
		values := strings.Split(sub, SEP2)
		if 2 != len(values) {
			xylog.ErrorNoId("AdvertisementSpace.AnalyseSubFlag(%v) failed.", sub)
			return
		}

		advertisementFlag := &battery.AdvertisementFlag{}
		flag := new(int32)
		convertAtoi(values[0], &flag)
		advertisementFlag.Flag = battery.ADVERTISEMENT_FLAG(*flag).Enum()
		convertAtoi(values[1], &advertisementFlag.Value)

		*flags = append(*flags, advertisementFlag)
	}
}
