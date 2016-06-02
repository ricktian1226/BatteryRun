// cache_tip
// 广告信息的缓存管理器定义
package xybusinesscache

import (
	//"time"
	//"math/rand"
	//"sort"
	"errors"
	"fmt"

	"guanghuan.com/xiaoyao/common/cache"
	xylog "guanghuan.com/xiaoyao/common/log"
	//"guanghuan.com/xiaoyao/common/util"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
)

// MAPTIPDETAIL 提示信息明细 tip identity -> tip detail
type MAPTIPDETAIL map[battery.TIP_IDENTITY]*battery.DBTip

func (m *MAPTIPDETAIL) Print() {
	for k, v := range *m {
		xylog.DebugNoId("[%d] %v", k, v)
	}
}

func (m *MAPTIPDETAIL) Clear() {
	*m = make(MAPTIPDETAIL, 0)
}

type MAPTIPS map[battery.LANGUAGE_TYPE]*MAPTIPDETAIL

func (m *MAPTIPS) Print() {
	xylog.DebugNoId("============== MAPTIPS begin ==============")
	for k, v := range *m {
		xylog.DebugNoId("%v", k)
		v.Print()
	}
	xylog.DebugNoId("============== MAPTIPS end ==============")
}

func (m *MAPTIPS) Clear() {
	*m = make(MAPTIPS, 0)
}

type TipCache struct {
	M MAPTIPS
}

func (c *TipCache) Init() {
	c.M = make(MAPTIPS, 0)
}

// 提示信息管理器
type TipManager struct {
	cache [2]TipCache
	xycache.CacheBase
}

func NewTipManager() (m *TipManager) {
	m = &TipManager{}
	for i := 0; i < 2; i++ {
		m.cache[i].Init()
	}

	return
}

var DefTipManager = NewTipManager()

// InitWhileStart 进程启动时初始化函数
func (m *TipManager) InitWhileStart() (err error) {

	//加载资源配置信息
	err = DefTipManager.Reload()
	if err != xyerror.ErrOK {
		xylog.ErrorNoId("Tip ResLoad failed : %v ", err)
		return
	}

	return
}

//MajorTips 返回主缓存
func (m *TipManager) MajorTips() *MAPTIPS {
	return &(m.cache[m.Major()].M)
}

//SecondaryTips 返回备缓存
func (m *TipManager) SecondaryTips() *MAPTIPS {
	return &(m.cache[m.Secondary()].M)
}

//重载提示信息资源
func (m *TipManager) Reload() (err error) {
	//从数据库中加载所有的提示信息
	//从数据库读取
	tips := make([]*battery.DBTip, 0)
	err = DefCacheDB.LoadTips(&tips)
	if err != xyerror.ErrOK && err != xyerror.ErrNotFound {
		xylog.ErrorNoId("LoadTips from db failed : %v", err)
		return
	}

	xylog.DebugNoId("tips : %v", tips)

	secondaryTips := m.SecondaryTips()
	secondaryTips.Clear()
	for _, tip := range tips {
		if mapDetail, ok := (*secondaryTips)[tip.GetLanguage()]; ok {
			(*mapDetail)[tip.GetId()] = tip
		} else { //如果没有level 2的列表，先创建一个
			mapTmp := make(MAPTIPDETAIL, 0)
			(*secondaryTips)[tip.GetLanguage()] = &mapTmp
			(*(*secondaryTips)[tip.GetLanguage()])[tip.GetId()] = tip
		}
	}

	m.Switch()

	m.MajorTips().Print()

	return xyerror.ErrOK

}

//Tip 返回指定的tip信息指针
// language battery.LANGUAGE_TYPE 语言
// id battery.TIP_IDENTITY
func (m *TipManager) Tip(language battery.LANGUAGE_TYPE, id battery.TIP_IDENTITY) (tipDesc string, err error) {
	var tip *battery.DBTip
	tip, err = nil, nil
	majorTips := m.MajorTips()
	if detail, ok := (*majorTips)[language]; ok {
		if tip, ok = (*detail)[id]; ok {
			tipDesc = tip.GetContent()
			return
		} else {
			err = errors.New(fmt.Sprintf("search tips for %v, %v failed", language, id))
			return
		}
	} else {
		err = errors.New(fmt.Sprintf("search tips for %v failed", language))
		return
	}
	return
}
