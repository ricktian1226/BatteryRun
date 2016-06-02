// cache_lotto
// 抽奖相关信息的缓存管理器定义
package xybusinesscache

import (
	"guanghuan.com/xiaoyao/common/cache"
	xylog "guanghuan.com/xiaoyao/common/log"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyerror "guanghuan.com/xiaoyao/superbman_server/error"
	"sort"
	//"time"
)

//抽奖配置项
type LottoConfig struct {
	LottoSlotCount     uint32 //抽奖格子数
	LottoInitUserValue int32  //默认权值
	LottoCostPerTime   int32  //抽奖消耗
	LottoDeduct        int32  //抽水比率
	SysLottoFreeCount  int32  //免费次数初值
}

//------------------ 奖池缓存定义 begin ------------------
//每个格子奖池，权重值对应的礼包id
type MAPWeight2Prop map[uint32]uint64

type SlotMapStruct struct {
	M MAPWeight2Prop
	S sort.IntSlice //排序的权重key值
}

func NewSlotMapStruct() *SlotMapStruct {
	return &SlotMapStruct{
		M: make(MAPWeight2Prop, 0),
		S: make(sort.IntSlice, 0),
	}
}

func (sms *SlotMapStruct) Clear() {
	sms.M = make(MAPWeight2Prop, 0)
	sms.S = make(sort.IntSlice, 0)
}

//奖池信息，系统抽奖
type SlotMap map[uint32]*SlotMapStruct

func (s *SlotMap) Clear() {
	*s = make(SlotMap, 0)
}

func (s *SlotMap) Print() {
	for slotid, m := range *s {
		xylog.DebugNoId("--Slot(%d)--", slotid)
		for weight, propid := range m.M {
			xylog.DebugNoId("prop(%d) weight(%d)", propid, weight)
		}
	}
}

func NewSlotMap() (s *SlotMap) {
	*s = make(SlotMap, 0)
	return
}

//游戏后抽奖是相同的
type AfterGameStage2SlotMap map[uint32]*SlotMap

func (a *AfterGameStage2SlotMap) Clear() {
	*a = make(AfterGameStage2SlotMap, 0)
}

func NewAfterGameStage2SlotMap() (afterGameStage2SlotMap *AfterGameStage2SlotMap) {
	*afterGameStage2SlotMap = make(AfterGameStage2SlotMap, 0)
	return
}

//------------------ 奖池缓存定义 end ------------------

//----------- 游戏 指标->阶段 映射信息 begin -------------
//指标值对应的阶段编号
type MAPQuotaValue2Stage map[uint64]uint32

//指标标识对应的阶段信息
type AfterGameMAPQuotaId2Stage map[battery.QuotaEnum]*StageMapStruct

func (a *AfterGameMAPQuotaId2Stage) Clear() {
	*a = make(AfterGameMAPQuotaId2Stage, 0)
}

type StageMapStruct struct {
	M MAPQuotaValue2Stage //阶段列表
	S sort.IntSlice       //排序的阶段key值(注意是由大到小进行排序)
}

func NewStageMapStruct() *StageMapStruct {
	return &StageMapStruct{
		M: make(MAPQuotaValue2Stage, 0),
		S: make(sort.IntSlice, 0),
	}
}

func (sms *StageMapStruct) Clear() {
	sms.M = make(MAPQuotaValue2Stage, 0)
	sms.S = make(sort.IntSlice, 0)
}

//----------- 游戏 指标->阶段 映射信息 end -------------

//------------------ 权重方案缓存定义 begin ------------------
//权重->格子号 映射关系
type MAPWeight2Slot map[uint32]uint32

type MAPWeightStruct struct {
	M MAPWeight2Slot
	S sort.IntSlice
}

func NewMAPWeightStruct() *MAPWeightStruct {
	return &MAPWeightStruct{
		M: make(MAPWeight2Slot, 0),
		S: make(sort.IntSlice, 0),
	}
}

func (mws *MAPWeightStruct) Clear() {
	mws.M = make(MAPWeight2Slot, 0)
	mws.S = make(sort.IntSlice, 0)
}

//系统抽奖 内部价值->权重列表
type SysWeightMap map[int32]*MAPWeightStruct
type SysWeightStruct struct {
	M SysWeightMap
	S sort.IntSlice
}

func NewSysWeightStruct() *SysWeightStruct {
	return &SysWeightStruct{
		M: make(SysWeightMap, 0),
		S: make(sort.IntSlice, 0),
	}
}

func (s *SysWeightStruct) Clear() {
	s.M = make(SysWeightMap, 0)
	s.S = make(sort.IntSlice, 0)
}

func (s *SysWeightStruct) Print() {
	for value, w := range s.M {
		xylog.DebugNoId("--value(%d)--", value)
		for weight, slotid := range w.M {
			xylog.DebugNoId("slotid(%d) %d", slotid, weight)
		}
	}
}

//系统抽奖 抽奖编号->权重列表
type SysSerialNum2SlotsMap map[int32]*battery.LottoSerialNumSlot

type SysSerialNum2SlotsStruct struct {
	M SysSerialNum2SlotsMap
}

func NewSysSerialNum2SlotsStruct() *SysSerialNum2SlotsStruct {
	return &SysSerialNum2SlotsStruct{
		M: make(SysSerialNum2SlotsMap, 0),
	}
}

func (s *SysSerialNum2SlotsStruct) Clear() {
	s.M = make(SysSerialNum2SlotsMap, 0)
}

func (s *SysSerialNum2SlotsStruct) Print() {
	for num, slots := range s.M {
		xylog.DebugNoId("--num(%d)--", num)
		xylog.DebugNoId("%v", slots)
	}
}

//游戏后抽奖 阶段号->权重列表
type AfterGameWeightMap map[uint32]*MAPWeightStruct

//游戏后抽奖 格子号->权重值
type MAPSlot2Weight map[uint32]uint32
type AfterGameWeightOriginalMap map[uint32]*MAPSlot2Weight

func (a *AfterGameWeightOriginalMap) Clear() {
	*a = make(AfterGameWeightOriginalMap, 0)
}

func (a *AfterGameWeightMap) Clear() {
	*a = make(AfterGameWeightMap, 0)
}

const WeightKeyMAX int32 = 1<<31 - 1 //权重最大值为2^32-1

//------------------ 权重方案缓存定义 end ------------------

type LottoCache struct {
	//系统抽奖
	sysSlotMap               SlotMap                   //奖池
	sysWeightStruct          *SysWeightStruct          //内部价值对应权重列表
	sysSerialNum2SlotsStruct *SysSerialNum2SlotsStruct //序号对应抽奖列表

	//游戏后抽奖
	afterGameStage2SlotMap     AfterGameStage2SlotMap     //奖池
	afterGameWeightMap         AfterGameWeightMap         //阶段对应的权重列表
	afterGameWeightOriginalMap AfterGameWeightOriginalMap //游戏后抽奖的原始权重表，用于删除格子后权重的重算
	afterGameMAPQuotaId2Stage  AfterGameMAPQuotaId2Stage  //游戏指标对应的权重列表
}

func (l *LottoCache) Clear() {
	l.sysSlotMap.Clear()
	l.sysWeightStruct.Clear()
	l.sysSerialNum2SlotsStruct.Clear()
	l.afterGameMAPQuotaId2Stage = make(AfterGameMAPQuotaId2Stage, 0)
	l.afterGameStage2SlotMap.Clear()
	l.afterGameWeightMap.Clear()
	l.afterGameWeightOriginalMap.Clear()
}

func NewLottoCache() *LottoCache {
	return &LottoCache{
		sysSlotMap:                 make(SlotMap, 0),
		sysWeightStruct:            NewSysWeightStruct(),
		sysSerialNum2SlotsStruct:   NewSysSerialNum2SlotsStruct(),
		afterGameStage2SlotMap:     make(AfterGameStage2SlotMap, 0),
		afterGameWeightMap:         make(AfterGameWeightMap, 0),
		afterGameWeightOriginalMap: make(AfterGameWeightOriginalMap, 0),
		afterGameMAPQuotaId2Stage:  make(AfterGameMAPQuotaId2Stage, 0),
	}
}

type LottoCacheManager struct {
	config LottoConfig
	cache  [2]*LottoCache
	xycache.CacheBase
}

//抽奖缓存管理器
var DefLottoCacheManager = NewLottoCacheManager()

func NewLottoCacheManager() *LottoCacheManager {
	return &LottoCacheManager{}
}

//进程启动时的初始化
func (pm *LottoCacheManager) InitWhileStart(config LottoConfig) (failReason battery.ErrorCode, err error) {
	//初始化数据库操作指针
	pm.Init()
	//加载资源配置信息
	failReason, err = DefLottoCacheManager.ReLoad(config)
	if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
		xylog.ErrorNoId("LottoResLoad failed : %v ", err)
		//os.Exit(-1) //加载失败，进程直接退出
	}
	return
}

func (pm *LottoCacheManager) Init() {
	pm.cache[0] = NewLottoCache()
	pm.cache[1] = NewLottoCache()
}

func (pm *LottoCacheManager) ReLoad(config LottoConfig) (failReason battery.ErrorCode, err error) {
	//begin := time.Now()
	//defer xylog.Debug("LottoCacheManager.Reload cost %d us", time.Since(begin)/time.Microsecond)

	pm.config = config
	failReason, err = pm.Load()
	return
}

func (pm *LottoCacheManager) SysSlots() *SlotMap {
	return &(pm.cache[pm.Major()].sysSlotMap)
}

func (pm *LottoCacheManager) SecondarySysSlots() *SlotMap {
	return &(pm.cache[pm.Secondary()].sysSlotMap)
}

func (pm *LottoCacheManager) AfterGameSlots() *AfterGameStage2SlotMap {
	return &(pm.cache[pm.Major()].afterGameStage2SlotMap)
}

func (pm *LottoCacheManager) SecondaryAfterGameStage2Slots() *AfterGameStage2SlotMap {
	return &(pm.cache[pm.Secondary()].afterGameStage2SlotMap)
}

func (pm *LottoCacheManager) SysWeights() *SysWeightStruct {
	return pm.cache[pm.Major()].sysWeightStruct
}

func (pm *LottoCacheManager) SecondarySysWeights() *SysWeightStruct {
	return pm.cache[pm.Secondary()].sysWeightStruct
}

func (pm *LottoCacheManager) SysSerialNumSlots() *SysSerialNum2SlotsStruct {
	return pm.cache[pm.Major()].sysSerialNum2SlotsStruct
}

func (pm *LottoCacheManager) SpecificSysSerialNumSlots(num int32) (battery.ErrorCode, *battery.LottoSerialNumSlot) {
	if s, ok := pm.cache[pm.Major()].sysSerialNum2SlotsStruct.M[num]; ok {
		return battery.ErrorCode_NoError, s
	} else { //没有找到，返回错误
		return battery.ErrorCode_GetSysSpecificLottoSerialNumSlotsError, nil
	}
}

func (pm *LottoCacheManager) SecondarySysSerialNumSlots() *SysSerialNum2SlotsStruct {
	return pm.cache[pm.Secondary()].sysSerialNum2SlotsStruct
}

func (pm *LottoCacheManager) AfterGameWeights() *AfterGameWeightMap {
	return &(pm.cache[pm.Major()].afterGameWeightMap)
}

func (pm *LottoCacheManager) SecondaryAfterGameWeights() *AfterGameWeightMap {
	return &(pm.cache[pm.Secondary()].afterGameWeightMap)
}

func (pm *LottoCacheManager) AfterGameWeightsOriginal() *AfterGameWeightOriginalMap {
	return &(pm.cache[pm.Major()].afterGameWeightOriginalMap)
}

func (pm *LottoCacheManager) SecondaryAfterGameWeightsOriginal() *AfterGameWeightOriginalMap {
	return &(pm.cache[pm.Secondary()].afterGameWeightOriginalMap)
}

func (pm *LottoCacheManager) AfterGameQuotaId2Stages() *AfterGameMAPQuotaId2Stage {
	return &(pm.cache[pm.Major()].afterGameMAPQuotaId2Stage)
}

func (pm *LottoCacheManager) SecondaryAfterGameQuotaId2Stages() *AfterGameMAPQuotaId2Stage {
	return &(pm.cache[pm.Secondary()].afterGameMAPQuotaId2Stage)
}

func (pm *LottoCacheManager) Load() (failReason battery.ErrorCode, err error) {

	failReason, err = pm.loadSlots()
	if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
		return
	}

	failReason, err = pm.loadSysWeights()
	if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
		return
	}

	failReason, err = pm.loadSysSerialNumSlots()
	if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
		return
	}

	//failReason, err = pm.loadAfterGameWeights()
	//if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
	//	return
	//}

	//failReason, err = pm.loadAfterGameQuotaId2Stages()
	//if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
	//	return
	//}

	pm.switchCache()

	return
}

func (pm *LottoCacheManager) loadSlots() (failReason battery.ErrorCode, err error) {
	slotitems := make([]battery.LottoSlotItem, 0)
	err = DefCacheDB.LoadSlotItems(&slotitems)
	if err != nil || len(slotitems) <= 0 {
		failReason = xyerror.LOTTO_DB_ERR_LOAD_SLOTITEMS
		return
	}
	//xylog.Debug("slotitems : %d %v", len(slotitems), slotitems)

	sysSlotMap := pm.SecondarySysSlots()
	afterGameStage2SlotMap := pm.SecondaryAfterGameStage2Slots()

	sysSlotMap.Clear()
	afterGameStage2SlotMap.Clear()

	//step 1.加载slotitems信息
	sysTmpWeight := make(map[uint32]uint32, pm.config.LottoSlotCount) //slotid -> weight
	var aftergameTmpWeight map[uint32]map[uint32]uint32               //stage -> slotid -> weight

	for _, s := range slotitems {
		dAType := s.GetDatype()
		propid := s.GetPropid()
		weight := s.GetWeight()
		slotid := s.GetSlotid()
		stage := s.GetStage()

		if weight == 0 {
			continue
		}

		switch dAType {
		case battery.DrawAwardType_DrawAwardType_System:
			sysTmpWeight[slotid] += weight

			var slotMapStruct *SlotMapStruct
			if _, ok := (*sysSlotMap)[slotid]; !ok {
				slotMapStruct = NewSlotMapStruct()
				(*sysSlotMap)[slotid] = slotMapStruct
			}

			slotMapStruct = (*sysSlotMap)[slotid]
			slotMapStruct.M[sysTmpWeight[slotid]] = propid
			slotMapStruct.S = append(slotMapStruct.S, int(sysTmpWeight[slotid]))
		case battery.DrawAwardType_DrawAwardType_GameFinish:
			if _, ok := aftergameTmpWeight[stage]; !ok {
				aftergameTmpWeight[stage] = make(map[uint32]uint32, pm.config.LottoSlotCount)
			}

			aftergameTmpWeight[stage][slotid] += weight

			var slotMapStruct *SlotMapStruct
			if _, ok := (*afterGameStage2SlotMap)[stage]; !ok {
				(*afterGameStage2SlotMap)[stage] = NewSlotMap()
				(*(*afterGameStage2SlotMap)[stage])[slotid] = NewSlotMapStruct()
			} else if _, ok := (*(*afterGameStage2SlotMap)[stage])[slotid]; !ok {
				(*(*afterGameStage2SlotMap)[stage])[slotid] = NewSlotMapStruct()
			}

			slotMapStruct = (*((*afterGameStage2SlotMap)[stage]))[slotid]
			slotMapStruct.M[aftergameTmpWeight[stage][slotid]] = propid
			slotMapStruct.S = append(slotMapStruct.S, int(aftergameTmpWeight[stage][slotid]))
		default:
			xylog.ErrorNoId("worry DrawAwardType[%d]", dAType)
		}
	}

	for _, slot := range *sysSlotMap {
		sort.Sort(slot.S)
	}

	for _, stage := range *afterGameStage2SlotMap {
		for _, slot := range *stage {
			sort.Sort(slot.S)
		}
	}
	//xylog.Debug(" sys slot items: %v", *sysSlotMap)
	//sysSlotMap.Print()
	//xylog.Debug(" aftergame slot items : %v", *afterGameStage2SlotMap)

	return
}

func (pm *LottoCacheManager) loadSysWeights() (failReason battery.ErrorCode, err error) {
	weights := make([]battery.LottoWeight, 0)
	err = DefCacheDB.LoadWeights(&weights)
	if err != nil || len(weights) <= 0 {
		failReason = xyerror.Resp_QuerySysLottoWeightError.GetCode()
		return
	}

	//xylog.Debug("weights : %v", weights)

	sysWeights := pm.SecondarySysWeights()

	sysWeights.Clear()

	for _, w := range weights {
		key := w.GetEndvalue()
		if w.Endvalue == nil {
			key = WeightKeyMAX
		}

		mws, _ := pm.ParseWeight(w.GetWeightlist())
		(sysWeights.M)[key] = mws
		sysWeights.S = append(sysWeights.S, int(key))
	}

	sort.Sort(sysWeights.S)

	//sysWeights.Print()

	return
}

func (pm *LottoCacheManager) loadSysSerialNumSlots() (failReason battery.ErrorCode, err error) {
	items := make([]battery.LottoSerialNumSlot, 0)
	err = DefCacheDB.LoadSerialNumSlots(&items)
	if err != nil || len(items) <= 0 {
		failReason = battery.ErrorCode_QuerySysLottoSerialNumSlotsError
		return
	}

	//xylog.Debug("weights : %v", weights)

	sysWeights := pm.SecondarySysSerialNumSlots()

	sysWeights.Clear()

	for _, item := range items {
		itemTmp := item
		num := itemTmp.GetSerialNum()
		if num <= 0 {
			xylog.ErrorNoId("Error SysSerialNum2Slot %v", itemTmp)
			continue
		}

		(sysWeights.M)[num] = &itemTmp
	}

	//sort.Sort(sysWeights.S)

	sysWeights.Print()

	return
}

//func (pm *LottoCacheManager) loadAfterGameWeights() (failReason battery.ErrorCode, err error) {
//	weights := make([]battery.LottoAfterGameWeight, 0)
//	err = DefCacheDB.LoadAfterGameWeights(&weights)
//	if err != nil || len(weights) <= 0 {
//		failReason = xyerror.Resp_QueryAfterGameLottoWeightError.GetCode()
//		return
//	}

//	afterGameWeights := pm.SecondaryAfterGameWeights()
//	afterGameWeights.Clear()
//	afterGameWeightOriginal := pm.SecondaryAfterGameWeightsOriginal()
//	afterGameWeightOriginal.Clear()

//	for _, w := range weights {
//		key := w.GetStage()
//		mws, msw := pm.ParseWeight(w.GetWeightlist())
//		(*afterGameWeights)[uint32(key)] = mws
//		(*afterGameWeightOriginal)[uint32(key)] = msw
//	}

//	xylog.Debug("sysWeights : %v", afterGameWeights)

//	return
//}

func (pm *LottoCacheManager) ParseWeight(weightlist []uint32) (mws *MAPWeightStruct, msw *MAPSlot2Weight) {
	msw = new(MAPSlot2Weight)
	*msw = make(MAPSlot2Weight, 0)
	mws = new(MAPWeightStruct)
	mws.M = make(MAPWeight2Slot, 0)
	mws.S = make(sort.IntSlice, pm.config.LottoSlotCount)
	var weightTmp uint32
	for i, v := range weightlist {

		if v == 0 {
			continue
		}

		(*msw)[uint32(i)] = v

		weightTmp += v
		mws.M[weightTmp] = uint32(i)
		mws.S[i] = int(weightTmp)
	}

	if len(mws.S) > 0 {
		sort.Sort(mws.S)
	}

	return
}

//func (pm *LottoCacheManager) loadAfterGameQuotaId2Stages() (failReason battery.ErrorCode, err error) {

//	stages := make([]battery.LottoStageItem, 0)
//	err = DefCacheDB.LoadStages(&stages)
//	if err != nil || len(stages) <= 0 {
//		failReason = xyerror.Resp_QueryAfterGameQuotaId2StagesError.GetCode()
//		return
//	}

//	afterGameMAPQuotaId2Stage := pm.SecondaryAfterGameQuotaId2Stages()
//	afterGameMAPQuotaId2Stage.Clear()

//	var stageMapStruct *StageMapStruct
//	for _, stage := range stages {
//		quotaId := stage.GetQuotaId()
//		quotaValue := stage.GetQuotaValue()
//		stageId := stage.GetStage()

//		if _, ok := (*afterGameMAPQuotaId2Stage)[quotaId]; !ok {
//			(*afterGameMAPQuotaId2Stage)[quotaId] = NewStageMapStruct()
//		}

//		stageMapStruct = (*afterGameMAPQuotaId2Stage)[quotaId]

//		stageMapStruct.M[quotaValue] = stageId

//		stageMapStruct.S = append(stageMapStruct.S, int(quotaValue))
//	}

//	//排一下序
//	for _, S := range *afterGameMAPQuotaId2Stage {
//		sort.Sort(S.S)
//	}

//	xylog.Debug("afterGameMAPQuotaId2Stage : %v", afterGameMAPQuotaId2Stage)

//	return
//}

func (pm *LottoCacheManager) switchCache() (fail_reason int32, err error) {
	pm.Switch()
	//xylog.Debug("now lotto cache switch to %d", pm.Major())
	return
}
