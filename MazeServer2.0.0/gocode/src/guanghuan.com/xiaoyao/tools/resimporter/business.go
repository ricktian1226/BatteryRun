// business
package main

import (
    "bufio"
    proto "code.google.com/p/goprotobuf/proto"
    //"fmt"
    "errors"
    xylog "guanghuan.com/xiaoyao/common/log"
    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    xyerror "guanghuan.com/xiaoyao/superbman_server/error"
    "os"
    "strconv"
    "strings"
    "time"
)

var (
    DEF_BEGINTIME = int64(0)
    DEF_ENDTIME   = int64(1<<63 - 1)
    ErrResError   = errors.New("ErrResError")
    TimeFormat    = "2006/1/2 15:04"
    DateFormat    = "2006/1/2"
    NullChar      = ""  //空，无内容
    NullChar1     = "#" //部分配置项无内容的时候会配置"#"，如：道具子项和道具分解价值
    NullChar2     = "0" //部分配置项无内容的时候会配置"0"，如：道具子项和道具分解价值
    NullChar3     = "$" //价格配置项无内容的时候会配置"$"
)

const (
    FILE_NAME_PROPS                     = "道具配置信息.txt"
    FILE_NAME_SLOTITEMS                 = "奖池配置信息.txt"
    FILE_NAME_WEIGHTS                   = "系统抽奖格子权重信息.txt"
    FILE_NAME_STAGES                    = "游戏阶段对应信息.txt"
    FILE_NAME_GOODS                     = "商品配置信息.txt"
    FILE_NAME_MISSION_CONFIG            = "任务配置信息.txt"
    FILE_NAME_SIGNIN_ACTIVITY_CONFIG    = "签到活动.txt"
    FILE_NAME_SIGNIN_AWARD_CONFIG       = "签到活动奖励.txt"
    FILE_NAME_SYSTEM_PROPS              = "系统道具数值配置.txt"
    FILE_NAME_BEFOREGAME_PROPS          = "游戏前随机商品配置信息.txt"
    FILE_NAME_ROLE_CONFIG               = "角色信息.txt"
    FILE_NAME_ROLE_LEVEL_GOODS          = "角色升级信息.txt"
    FILE_NAME_MAIL_CONFIG               = "系统基础邮件信息.txt"
    FILE_NAME_JIGSAW_CONFIG             = "拼图配置信息.txt"
    FILE_NAME_PICKUP_CONFIG             = "收集物配置信息.txt"
    FILE_NAME_ANNOUNCEMENT_CONFIG       = "公告.txt"
    FILE_NAME_LOTTO_SPECIAL_CONFIG      = "特殊抽奖配置信息.txt"
    FILE_NAME_ADVERTISEMENT_CONFIG      = "广告配置信息.txt"
    FILE_NAME_ADVERTISEMENTSPACE_CONFIG = "广告位配置信息.txt"
    FILE_NAME_TIP_CONFIG                = "提示配置信息.txt"
    FILE_NAME_NEWACCOUNTPROP_CONFIG     = "首次登录礼包配置信息.txt"
    FILE_NAME_SHARE_ACTIVITY_CONFIG     = "分享活动配置.txt"
    FILE_NAME_SHARE_AWARDS_CONFIG       = "分享奖励配置.txt"
    FILE_NAME_CHECKPOINT_UNLOCK_CONFIG  = "关卡解锁商品配置.txt"
)

const (
    SEP  = "\t"
    SEP1 = ";"
    SEP2 = ":"
)

func SetResFlag(flag *int32, operationResType battery.OPERATION_RES_TYPE) {

    if operationResType == battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_UNKOWN {
        return
    }

    *flag = (*flag) | (1 << (uint32(operationResType)))
    //xylog.Debug("*flag = (*flag) | (1 << (uint32(operationResType) - 1)) : *flag = (*flag) | (1 << (%d - 1))", uint32(operationResType))
}

func FlagMask(resType battery.OPERATION_RES_TYPE) int32 {
    return int32(1 << uint32(resType))
}

func FlagString(flag int32) (str string) {

    if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_PROP)) > 0 {
        str += "道具配置信息,"
    }
    if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_LOTTO_SLOTITEM)) > 0 {
        str += "奖池配置信息,"
    }
    if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_LOTTO_WEIGHT)) > 0 {
        str += "系统抽奖格子权重信息,"
    }
    if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_LOTTO_STAGE)) > 0 {
        str += "游戏阶段对应信息,"
    }
    if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_MALLITEM)) > 0 {
        str += "商品配置信息,"
    }
    if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_MISSION)) > 0 {
        str += "任务配置信息,"
    }
    if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_SIGNIN_ACTIVITY)) > 0 {
        str += "签到活动配置信息,"
    }
    if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_SIGNIN_ITEM)) > 0 {
        str += "签到活动奖励配置信息,"
    }
    if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_RUNE_VALUE)) > 0 {
        str += "系统道具数值配置,"
    }
    if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_BEFOREGAME_RANDOM_WEIGHT)) > 0 {
        str += "游戏前随机商品配置信息,"
    }
    if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_ROLE)) > 0 {
        str += "角色信息,"
    }
    if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_ROLE_LEVEL)) > 0 {
        str += "角色升级信息,"
    }
    if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_SYS_MAIL)) > 0 {
        str += "系统基础邮件信息,"
    }
    if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_JIASAW)) > 0 {
        str += "拼图配置信息,"
    }
    if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_PICKUP_WEIGHT)) > 0 {
        str += "收集物配置信息,"
    }
    if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_ANNOUNCEMENT)) > 0 {
        str += "公告,"
    }
    if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_ADVERTISEMENT)) > 0 {
        str += "广告,"
    }
    if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_ADVERTISEMENT_SPACE)) > 0 {
        str += "广告位,"
    }
    if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_TIP)) > 0 {
        str += "提示,"
    }
    if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_SHARE_ACTIVITY)) > 0 {
        str += "分享活动配置"
    }
    if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_SHARE_AWARDS)) > 0 {
        str += "分享奖励配置"
    }

    return
}

func GetProperty(ops *[]*battery.ResOpItem) (flag int32) {

    var err error

    optype := battery.RES_OP_TYPE_OP_ADD
    {
        m := &Props{}
        err = GetOps(FILE_NAME_PROPS, ops, &optype, m)
        if err == xyerror.ErrOK {
            SetResFlag(&flag, battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_PROP)
        } else {
            xylog.ErrorNoId("err : %v", err)
        }
    }

    {
        m := &LottoSlotItems{}
        err = GetOps(FILE_NAME_SLOTITEMS, ops, &optype, m)
        if err == xyerror.ErrOK {
            SetResFlag(&flag, battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_LOTTO_SLOTITEM)
        } else {
            xylog.ErrorNoId("err : %v", err)
        }
    }

    {
        m := &LottoWeight{}
        err = GetOps(FILE_NAME_WEIGHTS, ops, &optype, m)
        if err == xyerror.ErrOK {
            SetResFlag(&flag, battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_LOTTO_WEIGHT)
        } else {
            xylog.ErrorNoId("err : %v", err)
        }
    }

    //游戏后抽奖暂时不实现，屏蔽先
    //{
    //	m := &LottoStage{}
    //	err = GetOps(FILE_NAME_STAGES, ops, &optype, m)
    //	if err == xyerror.ErrOK {
    //		SetResFlag(&flag, battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_LOTTO_STAGE)
    //	} else {
    //		xylog.Error("err : %v", err)
    //	}
    //}

    {
        m := &Goods{}
        err = GetOps(FILE_NAME_GOODS, ops, &optype, m)
        if err == xyerror.ErrOK {
            SetResFlag(&flag, battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_MALLITEM)
        } else {
            xylog.ErrorNoId("err : %v", err)
        }
    }

    {
        m := &MissionConfig{}
        err = GetOps(FILE_NAME_MISSION_CONFIG, ops, &optype, m)
        if err == xyerror.ErrOK {
            SetResFlag(&flag, battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_MISSION)
        } else {
            xylog.ErrorNoId("err : %v", err)
        }
    }

    {
        m := &SignInActivity{}
        err = GetOps(FILE_NAME_SIGNIN_ACTIVITY_CONFIG, ops, &optype, m)
        if err == xyerror.ErrOK {
            SetResFlag(&flag, battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_SIGNIN_ACTIVITY)
        } else {
            xylog.ErrorNoId("err : %v", err)
        }
    }

    {
        m := &SignInItem{}
        err = GetOps(FILE_NAME_SIGNIN_AWARD_CONFIG, ops, &optype, m)
        if err == xyerror.ErrOK {
            SetResFlag(&flag, battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_SIGNIN_ITEM)
        } else {
            xylog.ErrorNoId("err : %v", err)
        }
    }

    {
        m := &SystemProps{}
        err = GetOps(FILE_NAME_SYSTEM_PROPS, ops, &optype, m)
        if err == xyerror.ErrOK {
            SetResFlag(&flag, battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_RUNE_VALUE)
        } else {
            xylog.ErrorNoId("err : %v", err)
        }
    }

    {
        m := &BeforeGameProps{}
        err = GetOps(FILE_NAME_BEFOREGAME_PROPS, ops, &optype, m)
        if err == xyerror.ErrOK {
            SetResFlag(&flag, battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_BEFOREGAME_RANDOM_WEIGHT)
        } else {
            xylog.ErrorNoId("err : %v", err)
        }
    }

    {
        m := &RoleInfos{}
        err = GetOps(FILE_NAME_ROLE_CONFIG, ops, &optype, m)
        if err == xyerror.ErrOK {
            SetResFlag(&flag, battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_ROLE)
        } else {
            xylog.ErrorNoId("err : %v", err)
        }
    }

    {
        m := &RoleLevelMallItems{}
        err = GetOps(FILE_NAME_ROLE_LEVEL_GOODS, ops, &optype, m)
        if err == xyerror.ErrOK {
            SetResFlag(&flag, battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_ROLE_LEVEL)
        } else {
            xylog.ErrorNoId("err : %v", err)
        }
    }

    {
        m := &MailConfigs{}
        err = GetOps(FILE_NAME_MAIL_CONFIG, ops, &optype, m)
        if err == xyerror.ErrOK {
            SetResFlag(&flag, battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_SYS_MAIL)
        } else {
            xylog.ErrorNoId("err : %v", err)
        }
    }

    {
        m := &JigsawConfigs{}
        err = GetOps(FILE_NAME_JIGSAW_CONFIG, ops, &optype, m)
        if err == xyerror.ErrOK {
            SetResFlag(&flag, battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_JIASAW)
        } else {
            xylog.ErrorNoId("err : %v", err)
        }
    }

    {
        m := &PickUps{}
        err = GetOps(FILE_NAME_PICKUP_CONFIG, ops, &optype, m)
        if err == xyerror.ErrOK {
            SetResFlag(&flag, battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_PICKUP_WEIGHT)
        } else {
            xylog.ErrorNoId("err : %v", err)
        }
    }

    {
        m := &Announcement{}
        err = GetOps(FILE_NAME_ANNOUNCEMENT_CONFIG, ops, &optype, m)
        if err == xyerror.ErrOK {
            SetResFlag(&flag, battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_ANNOUNCEMENT)
        } else {
            xylog.ErrorNoId("err : %v", err)
        }
    }

    {
        m := &LottoSerialNumSlot{}
        err = GetOps(FILE_NAME_LOTTO_SPECIAL_CONFIG, ops, &optype, m)
        if err == xyerror.ErrOK {
            SetResFlag(&flag, battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_LOTTO_SERIALNUM_SLOT)
        } else {
            xylog.ErrorNoId("err : %v", err)
        }
    }

    {
        m := &Advertisement{}
        err = GetOps(FILE_NAME_ADVERTISEMENT_CONFIG, ops, &optype, m)
        if err == xyerror.ErrOK {
            SetResFlag(&flag, battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_ADVERTISEMENT)
        } else {
            xylog.ErrorNoId("err : %v", err)
        }
    }

    {
        m := &AdvertisementSpace{}
        err = GetOps(FILE_NAME_ADVERTISEMENTSPACE_CONFIG, ops, &optype, m)
        if err == xyerror.ErrOK {
            SetResFlag(&flag, battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_ADVERTISEMENT_SPACE)
        } else {
            xylog.ErrorNoId("err : %v", err)
        }
    }

    {
        m := &Tip{}
        err = GetOps(FILE_NAME_TIP_CONFIG, ops, &optype, m)
        if err == xyerror.ErrOK {
            SetResFlag(&flag, battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_TIP)
        } else {
            xylog.ErrorNoId("err : %v", err)
        }
    }

    {
        m := &NewAccountProp{}
        err = GetOps(FILE_NAME_NEWACCOUNTPROP_CONFIG, ops, &optype, m)
        if err == xyerror.ErrOK {
            SetResFlag(&flag, battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_NEWACCOUNTPROP)
        } else {
            xylog.ErrorNoId("err : %v", err)
        }
    }

    {
        m := &ShareActivity{}
        err = GetOps(FILE_NAME_SHARE_ACTIVITY_CONFIG, ops, &optype, m)
        if err == xyerror.ErrOK {
            SetResFlag(&flag, battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_SHARE_ACTIVITY)
        } else {
            xylog.ErrorNoId("err : %v", err)
        }
    }
    {
        m := &ShareAwrads{}
        err = GetOps(FILE_NAME_SHARE_AWARDS_CONFIG, ops, &optype, m)
        if err == xyerror.ErrOK {
            SetResFlag(&flag, battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_SHARE_AWARDS)
        } else {
            xylog.ErrorNoId("err : %v", err)
        }
    }
    {
        m := &UnlockGoodsConfig{}
        err = GetOps(FILE_NAME_CHECKPOINT_UNLOCK_CONFIG, ops, &optype, m)
        if err == xyerror.ErrOK {
            SetResFlag(&flag, battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_CHECKPOINT_UNLOCK)
        } else {
            xylog.ErrorNoId("err :%v", err)
        }
    }
    xylog.DebugNoId("Ops res type : %s ", FlagString(flag))

    return
}

type OpsInterface interface {
    Analyse(lineNum int, line string, ops *[]*battery.ResOpItem, optype *battery.RES_OP_TYPE)
}

func GetOps(filename string, ops *[]*battery.ResOpItem, optype *battery.RES_OP_TYPE, opinterface OpsInterface) (err error) {
    xylog.DebugNoId("GetOps %s", filename)
    var file *os.File
    file, err = os.Open(filename)
    if err != nil {
        xylog.ErrorNoId("Open file %s failed", filename)
        return
    }
    defer file.Close()
    rb := bufio.NewReader(file)
    skipLine := 2
    lineNum := 0 //记录行号，方便查看错误日志时能快速定位到行
    for {
        lineNum++
        skipLine--
        var line string
        line, err = rb.ReadString('\n')
        if err != nil {
            err = xyerror.ErrOK
            return
        }
        //xylog.Debug("line(len(%d)) : [%v] ", len(line), line)

        //跳过第一行
        //xylog.Debug("skipLine : [%v] ", skipLine)
        if skipLine > 0 {
            continue
        }

        //xylog.Debug("line(len(%d)) : [%v] ", len(line), line)
        line = strings.TrimRight(line, "\r\n")
        //xylog.Debug("line(len(%d)) : [%v] ", len(line), line)

        if len(line) <= 0 {
            continue
        }

        //xylog.Debug("line(len(%d)) : [%v] ", len(line), line)
        if line[0] == '#' || len(line) <= 0 { //跳过注释和空行
            continue
        }

        opinterface.Analyse(lineNum, line, ops, optype)
    }
    return
}

func convertParseBool(a string, i **bool) (err error) {
    if len(a) > 0 {
        var value bool
        value, err = strconv.ParseBool(a)
        if err == nil {
            *i = proto.Bool(value)
        }
    }

    return
}

func convertAtoi(a string, i **int32) (err error) {
    if len(a) > 0 {
        var value int
        value, err = strconv.Atoi(a)
        if err == nil {
            *i = proto.Int32(int32(value))
        }
    }

    return
}

func convertAtoui(a string, i **uint32) (err error) {
    if len(a) > 0 {
        var value int
        value, err = strconv.Atoi(a)
        if err == nil {
            *i = proto.Uint32(uint32(value))
        }
    }

    return
}

func convertAtoui64(a string, i **uint64) (err error) {
    if len(a) > 0 {
        var value int
        value, err = strconv.Atoi(a)
        if err == nil {
            *i = proto.Uint64(uint64(value))
        }
    }

    return
}

func convertAtoi64(a string, i **int64) (err error) {
    if len(a) > 0 {
        var value int
        value, err = strconv.Atoi(a)
        if err == nil {
            *i = proto.Int64(int64(value))
        }
    }

    return
}

//解析时间类型(分钟)的参数
func convertAtoTimestamp(timeStr string) (timestamp int64, err error) {

    if len(timeStr) <= 0 {
        timestamp = 0
        return
    }

    var timeTmp time.Time
    timeTmp, err = time.Parse(TimeFormat, timeStr)
    if err != nil {
        //err = errors.New(fmt.Sprintf("convertStrToTimestamp %s to timestamp failed : %v ", timeStr, err))
        timestamp = 0
        return
    }

    timestamp = timeTmp.Unix()
    return
}

//解析时间类型(天)的参数
func convertAtoTimestampDate(timeStr string) (timestamp int64, err error) {

    //如果是空的，未设置则返回0
    if len(timeStr) <= 0 {
        timestamp = 0
        return
    }
    var timeTmp time.Time
    timeTmp, err = time.Parse(DateFormat, timeStr)
    if err != nil {
        //err = errors.New(fmt.Sprintf("convertAtoTimestampDate %s to timestamp failed : %v ", timeStr, err))
        timestamp = 0
        return
    }

    timestamp = timeTmp.Unix()
    return
}

//解析单个PropItem类型的参数
func getPropItem(propItem *battery.PropItem, propStr string) bool {
    propPropertys := strings.Split(propStr, SEP2)
    if len(propPropertys) != 3 {
        return false
    }

    for i, _ := range propPropertys {
        switch i {
        case 0:
            convertAtoui64(propPropertys[i], &(propItem.Id))
        case 1:
            convertAtoui(propPropertys[i], &(propItem.Amount))
        case 2:
            propType := new(int32)
            convertAtoi(propPropertys[i], &propType)
            propItem.Type = battery.PropType(*propType).Enum()
        }

    }
    xylog.DebugNoId("propItem:%v", propItem)
    if !checkProp(propItem) {
        return false
    }

    return true
}

func getIds(strIds string) (ids []uint64) {
    subStrIds := strings.Split(strIds, SEP1)
    if len(subStrIds) < 1 {
        return
    }

    for _, v := range subStrIds {
        if v != "" {
            id := new(uint64)
            err := convertAtoui64(v, &id)
            if err == nil {
                ids = append(ids, *id)
            }
        }
    }

    return
}

// 设置商品子项对应的标识位
func setMallItemFlag(mapId2Flag MAPId2Flag, ids []uint64, flagBit battery.MALL_ITEM_FLAG) {
    for _, id := range ids {
        if flag, ok := mapId2Flag[id]; ok {
            flag.Flag = proto.Int64(flag.GetFlag() | int64(1<<uint(flagBit)))
        }
    }
}

//校验道具配置信息的id和type是否匹配
func checkProp(propItem *battery.PropItem) bool {
    id := propItem.GetId()
    propType := uint64(propItem.GetType())
    if id/10000000 != propType {
        xylog.DebugNoId("id :%v,prpoType:%v", id, propType)
        return false
    } else {
        return true
    }
}

//解析多个PropItem类型的参数
func getPropItems(lineNum int, propStrs string, items *[]*battery.PropItem) {
    //xylog.Debug("%d : propStrs [%s]", lineNum, propStrs)
    propStrs = strings.TrimSpace(propStrs)
    if propStrs == NullChar1 || propStrs == NullChar {
        return
    }

    props := strings.Split(propStrs, SEP1)
    for _, prop := range props {
        if len(prop) <= 0 {
            continue
        }
        propItem := &battery.PropItem{}
        if getPropItem(propItem, prop) {
            *items = append(*items, propItem)
        } else {
            xylog.ErrorNoId("line %d : propItem Error : prop(%s) props(%s) propStrs(%s)", lineNum, prop, props, propStrs)
        }
    }
    xylog.DebugNoId("propitem :%v", items)
}

//解析MoneyItem类型的参数
func getMoneys(lineNum int, moneyString string, moneys *[]*battery.MoneyItem) {

    if moneyString == NullChar1 || moneyString == NullChar2 || moneyString == NullChar3 {
        return
    }

    if nil == *moneys {
        *moneys = make([]*battery.MoneyItem, 0)
    }

    moneyStrs := strings.Split(moneyString, SEP1)

    for _, moneyStr := range moneyStrs {

        if len(moneyStr) <= 0 {
            continue
        }

        mi := &battery.MoneyItem{}
        subs := strings.Split(moneyStr, SEP2)
        if len(subs) < 2 {
            xylog.ErrorNoId("line %d : Wrong money format %s", lineNum, moneyStr)
            break
        }
        mt := new(int32)
        amount := new(uint32)
        convertAtoi(subs[0], &mt)
        convertAtoui(subs[1], &amount)
        moneyType := battery.MoneyType(*mt)
        mi.Type = &moneyType
        mi.Amount = proto.Uint32(*amount)

        *moneys = append(*moneys, mi)
    }

}
