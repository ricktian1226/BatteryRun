package batteryapi

//import (
//	//proto "code.google.com/p/goprotobuf/proto"
//	//"encoding/binary"
//	//apns "github.com/timehop/apns"
//	//batterydb "guanghuan.com/xiaoyao/battery_apns_server/db"
//	"guanghuan.com/xiaoyao/common/apn"
//	//xyconf "guanghuan.com/xiaoyao/common/conf"
//	//"guanghuan.com/xiaoyao/common/db"
//	"guanghuan.com/xiaoyao/common/log"
//	//battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
//	//"guanghuan.com/xiaoyao/superbman_server/error"
//	"guanghuan.com/xiaoyao/superbman_server/server"
//	"sync"
//	"time"
//)

////使能devicetoken
////把devicetoken从黑名单中删除
//func (api *XYAPI) EnableDeviceToken(uid, dt string) {
//	if dt == "" {
//		xylog.Error("[%s] EnableDeviceToken devicetoken is nil", uid)
//		return
//	}

//	DefDeviceTokenBlackListManager.Remove(dt)
//	err := api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_DEVICETOKENBLACKLIST).RemoveDeviceTokenFromBlackList(dt)
//	if err != nil {
//		xylog.Error("[%s] RemoveDeviceToken(%s)FromBlackList failed : %v", uid, dt, err)
//	}
//}

////去使能devicetoken
//func (api *XYAPI) DisableDeviceToken(uid, dt string) {
//	if dt == "" {
//		xylog.Error("[%s] DisableDeviceToken devicetoken is nil", uid)
//		return
//	}

//	unit := &DeviceTokenBlackUnit{
//		DeviceToken: dt,
//		Timestamp:   time.Now(),
//	}

//	xylog.Debug("[%s] DisableDeviceToken devicetoken(%s) ", uid, dt)

//	if DefDeviceTokenBlackListManager.Add(unit) { //现将device token加到内存黑名单，如果内存黑名单修改了再将变化修改到数据库中
//		err := api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_DEVICETOKENBLACKLIST).UpsertDeviceTokenToBlackList(unit.DeviceToken, unit.Timestamp.Unix())
//		if err != nil {
//			xylog.ErrorNoId("UpsertDeviceToken(%s)ToBlackList failed : %v", unit, err)
//		} else {
//			xylog.DebugNoId("UpsertDeviceToken(%s)ToBlackList succeed : %v", unit, err)
//		}
//	}
//}

////加载device token黑名单
//func LoadDeviceTokenBlackList() (err error) {
//	var tokens []*xyapn.BlackDeviceToken
//	tokens, err = NewXYAPI().GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_DEVICETOKENBLACKLIST).QueryDeviceTokenBlackList()

//	if len(tokens) > 0 {
//		DefDeviceTokenBlackListManager.AddTokens(tokens)
//	}

//	return
//}

//type DeviceTokenBlackUnit struct {
//	DeviceToken string    //devicetoken
//	Timestamp   time.Time //时间戳
//}

//type DeviceTokenBlackListItem struct {
//	Timestamp time.Time //时间戳
//}

//type MapDeviceTokenBlackList map[string]*DeviceTokenBlackListItem

////devicetoke黑名单管理器
//type DeviceTokenBlackListManager struct {
//	items MapDeviceTokenBlackList //devicetoken黑名单，在黑名单中的，将不进行消息推送
//	mutex sync.RWMutex            //读写锁
//}

//func NewDeviceTokenBlackListManager() *DeviceTokenBlackListManager {
//	return &DeviceTokenBlackListManager{
//		items: make(MapDeviceTokenBlackList, 0),
//	}
//}

//var DefDeviceTokenBlackListManager = NewDeviceTokenBlackListManager()

//func (m *DeviceTokenBlackListManager) Lock() {
//	m.mutex.Lock()
//}

//func (m *DeviceTokenBlackListManager) Unlock() {
//	m.mutex.Unlock()
//}

//func (m *DeviceTokenBlackListManager) RLock() {
//	m.mutex.RLock()
//}

//func (m *DeviceTokenBlackListManager) RUnlock() {
//	m.mutex.RUnlock()
//}

////增加黑名单device token
//func (m *DeviceTokenBlackListManager) Add(item *DeviceTokenBlackUnit) (change bool) {
//	m.mutex.Lock()
//	defer m.mutex.Unlock()

//	change = false
//	if itemTmp, ok := m.items[item.DeviceToken]; ok {
//		if itemTmp.Timestamp.Before(item.Timestamp) {
//			itemTmp.Timestamp = item.Timestamp
//			change = true
//		}
//	} else {
//		m.items[item.DeviceToken] = &DeviceTokenBlackListItem{
//			Timestamp: item.Timestamp,
//		}
//		change = true
//	}
//	return
//}

////批量增加黑名单device token
//func (m *DeviceTokenBlackListManager) AddTokens(items []*xyapn.BlackDeviceToken) {
//	m.mutex.Lock()
//	defer m.mutex.Unlock()

//	m.items = make(MapDeviceTokenBlackList, 0) //先将列表清空

//	for _, item := range items {
//		m.items[item.GetToken()] = &DeviceTokenBlackListItem{
//			Timestamp: time.Unix(item.GetTimestamp(), 0),
//		}
//	}
//}

////删除黑名单device token
//func (m *DeviceTokenBlackListManager) Remove(dt string) {
//	m.mutex.Lock()
//	defer m.mutex.Unlock()

//	delete(m.items, dt)
//}

////判断device token是否在黑名单中
//func (m *DeviceTokenBlackListManager) BeWanted(dt string) bool {
//	m.mutex.RLock()
//	defer m.mutex.RUnlock()

//	if _, ok := m.items[dt]; ok {
//		return true
//	}

//	return false
//}

////打印当前的device token黑名单信息
//func (m *DeviceTokenBlackListManager) String() (str string) {
//	//m.mutex.RLock()
//	//defer m.mutex.RUnlock()

//	for dt, _ := range m.items {
//		str += dt + "\n"
//	}

//	return
//}
