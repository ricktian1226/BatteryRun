// xyapi_prop
package batteryapi

import (
    "code.google.com/p/goprotobuf/proto"
    "gopkg.in/mgo.v2"
    xylog "guanghuan.com/xiaoyao/common/log"
    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    "guanghuan.com/xiaoyao/superbman_server/cache"
    "guanghuan.com/xiaoyao/superbman_server/error"
    xymoney "guanghuan.com/xiaoyao/superbman_server/money"
    "guanghuan.com/xiaoyao/superbman_server/server"
)

//运营操作接口——发放物品
func (api *XYAPI) OperationMaintenanceProp(req *battery.MaintenancePropRequest, resp *battery.MaintenancePropResponse) (err error) {

    var (
        uid             = req.GetUid()
        propType        = req.GetPropType()
        propId          = req.GetPropId()
        maintenanceType = req.GetMaintenanceType()
        platform        = req.GetPlatformType()
        amount          = req.GetAmount()
    )

    //获取请求的终端平台类型
    api.SetDB(platform)

    //初始化返回值
    resp.Error = xyerror.Resp_NoError

    xylog.Debug(uid, "OperationMaintenanceProp : %v", req)

    switch maintenanceType {
    case battery.MAINTENANCE_TYPE_MAINTENANCE_TYPE_PROP_ADD, battery.MAINTENANCE_TYPE_MAINTENANCE_TYPE_CDKEY_EXCHANGE:
        //api.maintenanceAddProp(uid, propId, propType, amount, resp.Error)
        api.maintenanceAddSystemmail(uid, propId, propType, amount, SYSMAIL_TYPE_DYNAMIC*SYSMAIL_SEGMENT, resp.Error)

    //case battery.MAINTENANCE_TYPE_MAINTENANCE_TYPE_PROP_DEL:
    //do nothing，暂时没有做删除物品的接口，如果有需要再开发
    default: //do nothing
    }

    return
}

//判断是否发放物品是否需要刷新玩家
// 根据物品类型判断是否需要刷新玩家账户
func (api *XYAPI) maintenancePropNeedUpdateAccount(propType battery.PropType) bool {
    return propType == battery.PropType_PROP_COIN ||
        propType == battery.PropType_PROP_CHIP ||
        propType == battery.PropType_PROP_BADGE ||
        propType == battery.PropType_PROP_DIAMOND ||
        propType == battery.PropType_PROP_STAMINA ||
        propType == battery.PropType_PROP_RUNE ||
        propType == battery.PropType_PROP_JIGSAW ||
        propType == battery.PropType_PROP_PACKAGE ||
        propType == battery.PropType_PROP_ROLE
}

// 运营增加道具（直接增加）
// uid string 玩家标识
// propId uint64 道具id
// propType battery.PropType 道具类型
// amount int32 道具数目
// errStruct *battery.Error 错误信息
func (api *XYAPI) maintenanceAddProp(uid string,
    propId uint64,
    propType battery.PropType,
    amount int32,
    errStruct *battery.Error) {

    var (
        err             error
        accountWithFlag *AccountWithFlag
    )

    if api.maintenancePropNeedUpdateAccount(propType) {
        account := &battery.DBAccount{}
        err = api.GetDBAccountDirect(uid, account, mgo.Strong)
        if err != xyerror.ErrOK {
            xylog.Error(uid, "GetDBAccountDirect failed : %v", err)
            errStruct.Code = battery.ErrorCode_GetAccountByUidError.Enum()
            return
        }

        accountWithFlag = &AccountWithFlag{
            account: account,
            bChange: false,
        }
    }
    propItem := &battery.PropItem{
        Id:     proto.Uint64(propId),
        Type:   propType.Enum(),
        Amount: proto.Uint32(uint32(amount)),
    }

    err = api.GainProp(uid, accountWithFlag, propItem, ACCOUNT_UPDATE_NO_DELAY, battery.MoneySubType_gain)
    if err != xyerror.ErrOK {
        xylog.Error(uid, "GainProp failed : %v", err)
        errStruct.Code = battery.ErrorCode_GainPropsError.Enum()
        return
    }
}

// 运营增加道具（系统邮件增加）
// uid string 玩家标识
// propId uint64 道具id
// propType battery.PropType 道具类型
// amount int32 道具数目
// platform battery.PLATFORM_TYPE 平台类型
// errStruct *battery.Error 错误信息
func (api *XYAPI) maintenanceAddSystemmail(uid string,
    propId uint64,
    propType battery.PropType,
    amount int32,
    systemBaseMailId int32,
    errStruct *battery.Error) {

    err := api.addMaintenanceSysmail(uid, propId, propType, amount, systemBaseMailId)
    if err != xyerror.ErrOK {
        errStruct.Code = battery.ErrorCode_AddMaintenanceSysmailError.Enum()
        return
    }
}

//查询道具详细信息
// propid uint64 道具id
func (api *XYAPI) PropDetail(propid uint64) (detail *xybusinesscache.PropStruct, err error) {
    detail = xybusinesscache.DefPropCacheManager.Prop(propid)
    if nil == detail { //没找到
        err = xyerror.ErrQueryPropsFromCacheError
    }

    return
}

//分解道具
// uid string 玩家id
// propItem *battery.PropItem 道具详细信息
// accountWithFlag *AccountWithFlag 玩家账户信息临时变量
// delay bool 是否延迟刷新（延迟刷新是为性能考虑，将多次刷新玩家账户信息的操作合并成一个）
func (api *XYAPI) ResolveProp(uid string, accountWithFlag *AccountWithFlag, propItem *battery.PropItem, delay bool) (err error) {

    id := propItem.GetId()
    amount := propItem.GetAmount()
    //查找道具详细信息
    detail, err := api.PropDetail(id)
    if err != nil {
        xylog.Error(uid, "get PropDetail for %d failed :　%v", id, err)
        return
    }

    // 获取道具的分解价值，分发到用户钱包
    moneys := detail.ResolveValue
    if len(moneys) > 0 {
        // 如果玩家已经拥有了分解加成道具，则算上加成
        var value int32
        if api.IsRuneValid(uid, RUNE_RESOLVE_ADDITIONAL) {
            value = api.RuneConfigValue(RUNE_RESOLVE_ADDITIONAL)
        }

        // 加成后的权值
        //factor := RUNE_BASE_FACTOR + accountWithFlag.account.GetResolveAddtional()
        factor := RUNE_BASE_FACTOR + value

        xylog.Debug(uid, "ResolveProp (%d*%d) to %v", id, amount, moneys)

        if delay { // 延迟刷新玩家账户数据
            if accountWithFlag != nil {
                for _, money := range moneys {
                    xymoney.Add(uid, api.GetXYDB(xybusiness.BUSINESS_COLLECTION_INDEX_ACCOUNT), money.GetAmount()*uint32(factor)/100, money.GetType(), battery.MoneySubType_gain, accountWithFlag.account, ACCOUNT_UPDATE_DELAY)
                }

                accountWithFlag.SetChange()
            } else { //accountWithFlag == nil...do you  want to crash me?! SOB!
                err = xyerror.ErrBadInputData
                return
            }
        } else { //即刻刷新玩家账户数据
            xymoney.AddMultiple(uid, api.GetXYDB(xybusiness.BUSINESS_COLLECTION_INDEX_ACCOUNT), moneys, amount*uint32(factor))
        }

    }
    return
}

//延后发放道具接口，避免多次刷新，针对代币道具
// uid string 玩家id
// propItems []*battery.PropItem 道具列表
func (api *XYAPI) GainProps(uid string, accountWithFlag *AccountWithFlag, propItems []*battery.PropItem, delay bool, moneySubType battery.MoneySubType) (err error) {
    for _, propItem := range propItems {
        err = api.GainProp(uid, accountWithFlag, propItem, delay, moneySubType)
        if err != xyerror.ErrOK {
            return
        }
    }

    return
}

//玩家获取收集物
// pickUps []*battery.PropItem 收集物信息
func (api *XYAPI) GainPickUps(uid string, checkPointId uint32, pickUps []*battery.PropItem, accountWithFlag *AccountWithFlag) {
    //moneyPropItems := make([]*battery.PropItem, 0)   //代币类型
    //noMoneyPropItems := make([]*battery.PropItem, 0) //非代币类型
    propItems := make([]*battery.PropItem, 0)
    for _, pickUp := range pickUps { //按照记忆点id，道具类型获取收集的实际道具
        propType := pickUp.GetType()
        switch propType {
        case battery.PropType_PROP_COIN:
            fallthrough
        case battery.PropType_PROP_CHIP:
            fallthrough
        case battery.PropType_PROP_BADGE:
            fallthrough
        case battery.PropType_PROP_DIAMOND:
            fallthrough
        case battery.PropType_PROP_STAMINA:
            //代币类型都可以直接发放，不需要id
            propItem := pickUp
            propItems = append(propItems, propItem)
        default: //非代币类型需要随机出一个具体的道具项
            amount := pickUp.GetAmount()
            for ; amount > 0; amount-- {
                result, propItem := api.getCollectProp(uid, checkPointId, propType)
                if result {
                    propItems = append(propItems, propItem)
                }
            }
        }
    }

    if len(propItems) > 0 {
        xylog.Debug(uid, "pickUp moneyPropItems : %v", propItems)
        api.GainProps(uid, accountWithFlag, propItems, ACCOUNT_UPDATE_DELAY, battery.MoneySubType_gain)
    }
}

//获取收集物道具，根据记忆点id、道具类型获取实际收集的道具
// uid string 玩家id
// checkPointId uint32 记忆点id
// propType  battery.PropType 道具类型
//return:
// result bool 查询结果
// propItem *battery.PropItem 实际收集的道具信息
func (api *XYAPI) getCollectProp(uid string, checkPointId uint32, propType battery.PropType) (result bool, propItem *battery.PropItem) {
    propId := xybusinesscache.DefPickUpCacheManager.PickUp(checkPointId, propType)
    if xybusinesscache.INVALID_PROPID == propId {
        result = false
    } else {
        result = true
        propItem = &battery.PropItem{
            Id:     &propId,
            Amount: proto.Uint32(1),
            Type:   &propType,
        }
    }

    xylog.Debug(uid, "GainCollectProp checkPointId(%d) propType(%v) result(%t) : %d", checkPointId, propType, result, propId)

    return
}

//玩家获取单个道具，GainProps函数调用
// uid string 玩家id
// propItem *battery.PropItem 道具信息
func (api *XYAPI) GainProp(uid string, accountWithFlag *AccountWithFlag, propItem *battery.PropItem, delay bool, moneySubType battery.MoneySubType) (err error) {
    id := propItem.GetId()
    ptype := propItem.GetType()
    amount := propItem.GetAmount()

    switch ptype {
    case battery.PropType_PROP_COIN:
        err = api.updateUserDataWithCoin(uid, accountWithFlag, amount, delay, moneySubType)
    case battery.PropType_PROP_CHIP:
        err = api.updateUserDataWithChip(uid, accountWithFlag, amount, delay, moneySubType)
    case battery.PropType_PROP_BADGE:
        err = api.updateUserDataWithBadge(uid, accountWithFlag, amount, delay, moneySubType)
    case battery.PropType_PROP_DIAMOND:
        err = api.updateUserDataWithDiamond(uid, accountWithFlag, amount, delay, moneySubType)
    case battery.PropType_PROP_STAMINA:
        err = api.updateUserDataWithStamina(uid, accountWithFlag, amount, delay)
    case battery.PropType_PROP_RUNE:
        err = api.updateUserDataWithRune(uid, accountWithFlag, propItem, delay)
    case battery.PropType_PROP_ROLE:
        err = api.updateUserDataWithRole(uid, accountWithFlag, id, delay, UNLOCK_ROLE_NEED_ADD_JIGSAW)
    case battery.PropType_PROP_JIGSAW:
        err = api.updateUserDataWithJigsaw(uid, accountWithFlag, propItem, delay)
    case battery.PropType_PROP_PACKAGE:
        err = api.updateUserDataWithPackage(uid, accountWithFlag, propItem, delay)
    case battery.PropType_PROP_CONSUM:
        err = api.updateUserDataWithConsumable(uid, id, amount)
    case battery.PropType_PROP_LOTTO_TICKET:
        err = api.updateUserDataWithLottoTicket(uid, amount)
    default:
        xylog.Error(uid, "unkown PropType %v ", ptype)
        return
    }

    return
}

//玩家是否拥有道具
// uid string 玩家id
// propItem *battery.PropItem 道具信息
func (api *XYAPI) OwnProp(uid string, propItem *battery.PropItem) (isOwn bool) {
    id := propItem.GetId()
    ptype := propItem.GetType()
    amount := propItem.GetAmount()
    isOwn = false
    switch ptype {
    case battery.PropType_PROP_RUNE:
        isOwn = api.IsRuneValid(uid, id)
    case battery.PropType_PROP_ROLE:
        isOwn = api.IsRoleExisting(uid, id)
    case battery.PropType_PROP_JIGSAW:
        isOwn, _ = api.isJigsawExisting(uid, id)
    case battery.PropType_PROP_CONSUM:
        isOwn = api.IsConsumableExisting(uid, id, amount)
    default:
        xylog.Error(uid, " unkown PropType %v ", ptype)
        return
    }

    return
}

//更新礼包信息到玩家
// uid string 玩家id
//accountWithFlag *AccountWithFlag 玩家账户信息
// propItem *battery.PropItem 礼包信息
// delay bool 是否延迟发放
func (api *XYAPI) updateUserDataWithPackage(uid string, accountWithFlag *AccountWithFlag, propItem *battery.PropItem, delay bool) (err error) {
    //查询礼包子项数据
    var (
        id         = propItem.GetId()
        propType   = propItem.GetType()
        propStruct = xybusinesscache.DefPropCacheManager.Prop(id)
    )
    if nil == propStruct {
        xylog.Error(uid, "get propStruct of %d", id)
        err = xyerror.ErrQueryPropsFromCacheError
        return
    }

    //校验一把类型，如果类型不匹配，直接返回
    if propType != propItem.GetType() {
        xylog.Error(uid, "invalid propType %v in cache and request %v", propItem.GetType(), propType)
        err = xyerror.ErrPropTypeInvalidError
        return
    }

    //按照礼包里的道具子项，个个添加
    for _, item := range propStruct.Items {
        err = api.GainProp(uid, accountWithFlag, item, delay, battery.MoneySubType_gain)
        if err != xyerror.ErrOK {
            xylog.Error(uid, "GainProp %v failed : %v", item, err)
            continue
        }
    }

    return
}

//登录礼包发放
// uid string 玩家标识
// units []xycache.DispenseUnit 待发放的礼包信息
func (api *XYAPI) newAccountGainProps(uid string, source battery.ID_SOURCE, errStruct *battery.Error) {
    units := xybusinesscache.DefNewAccountPropManager.Units(source)
    if nil == units {
        xylog.DebugNoId(uid, "DefNewAccountPropManager.Units is nil for %v", source)

        return
    }

    //遍历整个待分发的礼包列表，进行礼包内容的发放
    for _, unit := range units {
        switch unit.Type {
        case battery.DISPENSE_TYPE_DISPENSE_TYPE_DIRECT: //直接发放
            for _, item := range unit.Items {
                api.maintenanceAddProp(
                    uid,
                    item.GetId(),
                    item.GetType(),
                    int32(item.GetAmount()),
                    errStruct)
            }

        case battery.DISPENSE_TYPE_DISPENSE_TYPE_SYSTEMMAIL: //通过系统邮件发放
            for _, item := range unit.Items {
                api.maintenanceAddSystemmail(
                    uid,
                    item.GetId(),
                    item.GetType(),
                    int32(item.GetAmount()),
                    unit.MailId,
                    errStruct)
            }
        }
    }

    return
}
