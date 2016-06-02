// xyapi_res
package batteryapi

import (
    "errors"
    "fmt"
    xylog "guanghuan.com/xiaoyao/common/log"
    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    xyerror "guanghuan.com/xiaoyao/superbman_server/error"
    "guanghuan.com/xiaoyao/superbman_server/server"
)

func FlagMask(resType battery.OPERATION_RES_TYPE) int32 {
    return int32(1 << uint32(resType))
}

//资源配置信息接口消息
func (api *XYAPI) OperationResRequest(req *battery.ResRequest, resp *battery.ResResponse) (err error) {
    xylog.DebugNoId("[ResRequest] start")
    defer xylog.DebugNoId("[ResRequest] done")

    var (
        failReason = xyerror.Resp_NoError.GetCode()
    )

    ops := req.GetOps()

    xylog.DebugNoId("[ResRequest] remove old : %t, flag : %d, ops : %v", req.GetRemoveOld(), req.GetFlag(), req.GetOps())

    //删除老数据
    if req.GetRemoveOld() {

        flag := req.GetFlag()
        if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_PROP)) > 0 {
            api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_PROP_CONFIG).RemoveAllProps()
            xylog.DebugNoId("OperationResRequest RemoveAllProps")
        }

        if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_LOTTO_SLOTITEM)) > 0 {
            api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_LOTTOSLOTITEM_CONFIG).RemoveAllSlotItems()
            xylog.DebugNoId("OperationResRequest RemoveAllLottoSlotItems")
        }
        if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_LOTTO_WEIGHT)) > 0 {
            api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_LOTTOWEIGHT_CONFIG).RemoveAllWeights()
            xylog.DebugNoId("OperationResRequest RemoveAllLottoWeights")
        }
        if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_LOTTO_SERIALNUM_SLOT)) > 0 {
            api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_LOTTOSERIALNUMSLOT_CONFIG).RemoveAllSerialNumSlots()
            xylog.DebugNoId("OperationResRequest RemoveAllSerialNumSlots")
        }
        //if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_LOTTO_STAGE)) > 0 {
        //	api.GetDB().RemoveAllStages()
        //	xylog.Debug("OperationResRequest RemoveAllLottoStages")
        //}
        if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_MALLITEM)) > 0 {
            api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_GOODS_CONFIG).RemoveAllMallItems()
            xylog.DebugNoId("OperationResRequest RemoveAllGoods")
        }
        if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_MISSION)) > 0 {
            api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_MISSION_CONFIG).RemoveAllMissionItems()
            xylog.DebugNoId("OperationResRequest RemoveAllMissionItems")
        }
        if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_SIGNIN_ACTIVITY)) > 0 {
            api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_SIGNINACTIVITY_CONFIG).RemoveAllSignInActivitys()
            xylog.DebugNoId("OperationResRequest RemoveAllSignInActivitys")
        }
        if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_SIGNIN_ITEM)) > 0 {
            api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_SIGNAWARD_CONFIG).RemoveAllSignInItems()
            xylog.DebugNoId("OperationResRequest RemoveAllSignInItems")
        }
        if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_RUNE_VALUE)) > 0 {
            api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_RUNE_CONFIG).RemoveAllRuneConfig()
            xylog.DebugNoId("OperationResRequest RemoveAllRuneConfig")
        }
        if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_BEFOREGAME_RANDOM_WEIGHT)) > 0 {
            api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_BEFOREGAMERANDOM_CONFIG).RemoveAllBGPropConfig()
            xylog.DebugNoId("OperationResRequest RemoveAllBGPropConfig")
        }
        if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_ROLE)) > 0 {
            api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_ROLEINFO_CONFIG).RemoveAllRoleInfoConfig()
            xylog.DebugNoId("OperationResRequest RemoveAllRoleInfoConfig")
        }
        if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_ROLE_LEVEL)) > 0 {
            api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_ROLELEVELBONUS_CONFIG).RemoveAllRoleLevelBonus()
            xylog.DebugNoId("OperationResRequest RemoveAllRoleLevelBonus")
        }
        if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_SYS_MAIL)) > 0 {
            api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_MAILINFO_CONFIG).RemoveAllMaillnfoConfig()
            xylog.DebugNoId("OperationResRequest RemoveAllMaillnfoConfig")
        }
        if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_JIASAW)) > 0 {
            api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_JIGSAW_CONFIG).RemoveAllJigsawConfig()
            xylog.DebugNoId("OperationResRequest RemoveAllJigsawConfig")
        }
        if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_PICKUP_WEIGHT)) > 0 {
            api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_PICKUPWEIGHT_CONFIG).RemoveAllPickUpWeightItems()
            xylog.DebugNoId("OperationResRequest RemoveAllPickUpWeightItems")
        }
        if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_ANNOUNCEMENT)) > 0 {
            api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_ANNOUNCEMENT_CONFIG).RemoveAllAnnouncementConfig()
            xylog.DebugNoId("OperationResRequest RemoveAllAnnouncementConfig")
        }
        if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_ADVERTISEMENT)) > 0 {
            api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_ADVERTISEMENT_CONFIG).RemoveAllAdvertisementConfig()
            xylog.DebugNoId("OperationResRequest RemoveAllAdvertisementConfig")
        }
        if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_ADVERTISEMENT_SPACE)) > 0 {
            api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_ADVERTISEMENT_SPACE_CONFIG).RemoveAllAdvertisementSpaceConfig()
            xylog.DebugNoId("OperationResRequest RemoveAllAdvertisementSpaceConfig")
        }
        if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_TIP)) > 0 {
            api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_TIP_CONFIG).RemoveAllTipConfig()
            xylog.DebugNoId("OperationResRequest RemoveAllTipConfig")
        }
        if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_NEWACCOUNTPROP)) > 0 {
            api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_TIP_CONFIG).RemoveAllNewAccountPropConfig()
            xylog.DebugNoId("OperationResRequest RemoveAllNewAccountPropConfig")
        }
        if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_SHARE_ACTIVITY)) > 0 {
            api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_SHARE_ACTIVITY).RemoveAllShareActivityConfig()
            xylog.DebugNoId("OperationResRequest RemoveAllShareAccountConfig")
        }
        if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_SHARE_AWARDS)) > 0 {
            api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_SHARE_AWARDS).RemoveAllShareAwardsConfig()
            xylog.DebugNoId("OperationResRequest RemoveAllShareAccountConfig")
        }
        if (flag & FlagMask(battery.OPERATION_RES_TYPE_OPERATION_RES_TYPE_CHECKPOINT_UNLOCK)) > 0 {
            api.GetCommonDB(xybusiness.Business_COMMON_COLLECTION_INDEX_CHECKPOINTUNLOCK_GOODS_CONFIG).RemoveAllUnlockGoodsConfig()
            xylog.DebugNoId("OperationResRequest RemoveAllUnlockGoodsConfig")
        }
        xylog.DebugNoId("RemoveOld true")

    }

    if nil != ops {
        xylog.DebugNoId("[ResRequest] ops : %v ", ops)
        for _, op := range ops {
            xylog.DebugNoId("[ResRequest] op : %v ", *op)

            //道具操作
            if op.Prop != nil {
                failReason, err = api.opProp(op.GetOptype(), op.Prop)
                if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
                    xylog.ErrorNoId("[ResRequest] opProp failed :%v", err)
                }
                xylog.DebugNoId("[ResRequest] opProp : %v ", *op)
            }

            //格子子项操作
            if op.Slotitem != nil {
                failReason, err = api.opSlotItem(op.GetOptype(), op.Slotitem)
                if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
                    xylog.ErrorNoId("[ResRequest] opSlotItem failed :%v", err)
                }
                //xylog.Debug("[ResRequest] opSlotItem : %v ", *op)
            }

            //权重操作
            if op.Weight != nil {
                failReason, err = api.opWeight(op.GetOptype(), op.Weight)
                if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
                    xylog.ErrorNoId("[ResRequest] opWeight failed :%v", err)
                }
                //xylog.Debug("[ResRequest] opweight : %v ", *op)
            }

            //权重操作
            if op.LottoSerialNumSlot != nil {
                failReason, err = api.opLottoSerialNumSlot(op.GetOptype(), op.LottoSerialNumSlot)
                if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
                    xylog.ErrorNoId("[ResRequest] opLottoSerialNumSlot failed :%v", err)
                }
                xylog.DebugNoId("[ResRequest] opLottoSerialNumSlot : %v ", op.LottoSerialNumSlot)
            }

            //if op.Stage != nil {
            //	failReason, err = api.opStage(op.GetOptype(), op.Stage)
            //	if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
            //		xylog.Error("[ResRequest] opStage failed :%v", err)
            //	}
            //	//xylog.Debug("[ResRequest] opstage : %v ", *op)
            //}

            if op.Mallitem != nil {
                failReason, err = api.opGoods(op.GetOptype(), op.Mallitem)
                if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
                    xylog.ErrorNoId("[ResRequest] opGoods failed :%v", err)
                }
                //xylog.Debug("[ResRequest] opGoods : %v ", *op)
            }

            if op.MissionItem != nil {
                failReason, err = api.opMissions(op.GetOptype(), op.MissionItem)
                if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
                    xylog.ErrorNoId("[ResRequest] opMissions failed :%v", err)
                }
                //xylog.Debug("[ResRequest] opMissions : %v ", *op)
            }

            if op.Activity != nil {
                failReason, err = api.opSignInActivity(op.GetOptype(), op.Activity)
                if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
                    xylog.ErrorNoId("[ResRequest] opSignInActivity failed :%v", err)
                }
                //xylog.Debug("[ResRequest] opSignInActivity : %v ", *op)
            }

            if op.SigninItem != nil {
                failReason, err = api.opSignInItem(op.GetOptype(), op.SigninItem)
                if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
                    xylog.ErrorNoId("[ResRequest] opSignInActivity failed :%v", err)
                }
                //xylog.Debug("[ResRequest] opSignInActivity : %v ", *op)
            }
            if op.RuneConfig != nil {
                failReason, err = api.opRuneConfig(op.GetOptype(), op.RuneConfig)
                if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
                    xylog.ErrorNoId("[ResRequest] opRuneConfig failed :%v", err)
                }
                //xylog.Debug("[ResRequest] opRuneConfig : %v ", *op)
            }
            if op.BeforeGameRandomGoodWeight != nil {
                failReason, err = api.opBeforeGameRandomWeightConfig(op.GetOptype(), op.BeforeGameRandomGoodWeight)
                if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
                    xylog.ErrorNoId("[ResRequest] opBeforGamePropConfig failed :%v", err)
                }
                //xylog.Debug("[ResRequest] opBeforGamePropConfig : %v ", *op)
            }
            if op.RoleInfoConfig != nil {
                failReason, err = api.opRoleInfoConfig(op.GetOptype(), op.RoleInfoConfig)
                if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
                    xylog.ErrorNoId("[ResRequest] opRoleInfoConfig failed :%v", err)
                }
                //xylog.Debug("[ResRequest] opRoleInfoConfig : %v ", *op)
            }
            if op.JigsawConfig != nil {
                failReason, err = api.opJigsawConfig(op.GetOptype(), op.JigsawConfig)
                if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
                    xylog.ErrorNoId("[ResRequest] opJigsawConfig failed :%v", err)
                }
                //xylog.Debug("[ResRequest] opBeforGamePropConfig : %v ", *op)
            }
            if op.MailConfig != nil {
                failReason, err = api.opMailConfig(op.GetOptype(), op.MailConfig)
                if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
                    xylog.ErrorNoId("[ResRequest] opMailConfig failed :%v", err)
                }
                //xylog.Debug("[ResRequest] opMailConfig : %v ", *op)
            }
            if op.PickUp != nil {
                failReason, err = api.opPickUp(op.GetOptype(), op.PickUp)
                if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
                    xylog.ErrorNoId("[ResRequest] opMailConfig failed :%v", err)
                }
                //xylog.Debug("[ResRequest] opPickUp : %v ", *op)
            }

            if op.RoleLevelBonus != nil {
                failReason, err = api.opRoleLevelBonus(op.GetOptype(), op.RoleLevelBonus)
                if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
                    xylog.ErrorNoId("[ResRequest] opRoleLevelBonus failed :%v", err)
                }
                //xylog.Debug("[ResRequest] opRoleLevelBonus : %v ", *op)
            }

            if op.AnnouncementItem != nil {
                failReason, err = api.opAnnouncement(op.GetOptype(), op.AnnouncementItem)
                if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
                    xylog.ErrorNoId("[ResRequest] opAnnouncement failed :%v", err)
                }
                //xylog.Debug("[ResRequest] opAnnouncement : %v ", *op)
            }

            if op.Advertisement != nil {
                failReason, err = api.opAdvertisement(op.GetOptype(), op.Advertisement)
                if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
                    xylog.ErrorNoId("[ResRequest] opAdvertisement failed :%v", err)
                }
            }

            if op.AdvertisementSpace != nil {
                failReason, err = api.opAdvertisementSpace(op.GetOptype(), op.AdvertisementSpace)
                if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
                    xylog.ErrorNoId("[ResRequest] opAdvertisementSpace failed :%v", err)
                }
            }

            if op.Tip != nil {
                failReason, err = api.opTip(op.GetOptype(), op.Tip)
                if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
                    xylog.ErrorNoId("[ResRequest] opTip failed :%v", err)
                }
            }

            if op.NewAccountProp != nil {
                failReason, err = api.opNewAccountProp(op.GetOptype(), op.NewAccountProp)
                if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
                    xylog.ErrorNoId("[ResRequest] opNewAccountProp failed :%v", err)
                }
            }

            if op.ShareActivity != nil {
                failReason, err = api.opShareActivity(op.GetOptype(), op.ShareActivity)
                if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
                    xylog.ErrorNoId("[ResRequest] opShareGift failed :%v", err)
                }
            }
            if op.ShareAwards != nil {
                failReason, err = api.opShareAwards(op.GetOptype(), op.ShareAwards)
                if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
                    xylog.ErrorNoId("[ResRequest] opShareGift failed :%v", err)
                }
            }
            if op.UnlockGoodsConfig != nil {
                xylog.DebugNoId("add unlock goods config: %v,%v", op.UnlockGoodsConfig.GetCheckPointId(), op.UnlockGoodsConfig.GetGoodsId())
                failReason, err = api.opUnlockGoodsConfig(op.GetOptype(), op.UnlockGoodsConfig)
                if failReason != xyerror.Resp_NoError.GetCode() || err != nil {
                    xylog.ErrorNoId("[ResRequest] opUnlockGoodsConfig failed : %v", err)
                }
            }
        }
    }

    return
}

//道具配置信息操作接口
func (api *XYAPI) opProp(optype battery.RES_OP_TYPE, prop *battery.Prop) (failReason battery.ErrorCode, err error) {
    switch optype {
    case battery.RES_OP_TYPE_OP_ADD:
        err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_PROP_CONFIG).UpsertProp(prop)
    case battery.RES_OP_TYPE_OP_DEL:
        err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_PROP_CONFIG).DelProp(prop)
    case battery.RES_OP_TYPE_OP_MOD:
        err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_PROP_CONFIG).ModProp(prop)
    default:
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
        err = errors.New(fmt.Sprintf("res op unkown op type %v", optype))
        return
    }

    if err != nil {
        failReason = xyerror.Resp_ResOpError.GetCode()
    }

    return
}

//抽奖槽位信息操作接口
func (api *XYAPI) opSlotItem(optype battery.RES_OP_TYPE, si *battery.LottoSlotItem) (failReason battery.ErrorCode, err error) {
    switch optype {
    case battery.RES_OP_TYPE_OP_ADD:
        err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_LOTTOSLOTITEM_CONFIG).UpsertSlotItem(si)
    case battery.RES_OP_TYPE_OP_DEL:
        err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_LOTTOSLOTITEM_CONFIG).DelSlotItem(si)
    case battery.RES_OP_TYPE_OP_MOD:
        err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_LOTTOSLOTITEM_CONFIG).ModSlotItem(si)
    default:
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
        err = errors.New(fmt.Sprintf("lotto res op unkown op type %v", optype))
        return
    }

    if err != nil {
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
    }

    return
}

//增加系统抽奖格子权重信息
func (api *XYAPI) opWeight(optype battery.RES_OP_TYPE, w *battery.LottoWeight) (failReason battery.ErrorCode, err error) {
    switch optype {
    case battery.RES_OP_TYPE_OP_ADD:
        err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_LOTTOWEIGHT_CONFIG).UpsertWeight(w)
    case battery.RES_OP_TYPE_OP_DEL:
        err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_LOTTOWEIGHT_CONFIG).DelWeight(w)
    case battery.RES_OP_TYPE_OP_MOD:
        err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_LOTTOWEIGHT_CONFIG).ModWeight(w)
    default:
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
        err = errors.New(fmt.Sprintf("res op unkown op type %v", optype))
        return
    }

    if err != nil {
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
    }

    return
}

//增加系统抽奖格子权重信息
func (api *XYAPI) opLottoSerialNumSlot(optype battery.RES_OP_TYPE, l *battery.LottoSerialNumSlot) (failReason battery.ErrorCode, err error) {
    switch optype {
    case battery.RES_OP_TYPE_OP_ADD:
        err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_LOTTOSERIALNUMSLOT_CONFIG).AddSerialNumSlot(l)
    default:
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
        err = errors.New(fmt.Sprintf("res op unkown op type %v", optype))
        return
    }

    if err != nil {
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
    }

    return
}

//游戏后抽奖阶段信息
//func (api *XYAPI) opStage(optype battery.RES_OP_TYPE, s *battery.LottoStageItem) (failReason battery.ErrorCode, err error) {
//	switch optype {
//	case battery.RES_OP_TYPE_OP_ADD:
//		err = api.GetDB().UpsertStage(s)
//	default:
//		failReason = xyerror.Resp_ResUnkownOpType.GetCode()
//		err = errors.New(fmt.Sprintf("res op unkown op type %v", optype))
//		return
//	}

//	if err != nil {
//		failReason = xyerror.Resp_ResUnkownOpType.GetCode()
//	}

//	return
//}

//商品配置信息
func (api *XYAPI) opGoods(optype battery.RES_OP_TYPE, good *battery.DBMallItem) (failReason battery.ErrorCode, err error) {
    switch optype {
    case battery.RES_OP_TYPE_OP_ADD:
        err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_GOODS_CONFIG).UpsertMallItem(good)
    default:
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
        err = errors.New(fmt.Sprintf("res op unkown op type %v", optype))
        return
    }

    if err != nil {
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
    }

    return
}

//任务配置信息
func (api *XYAPI) opMissions(optype battery.RES_OP_TYPE, mission *battery.MissionItem) (failReason battery.ErrorCode, err error) {
    switch optype {
    case battery.RES_OP_TYPE_OP_ADD:
        err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_MISSION_CONFIG).UpsertMissionItem(mission)
    default:
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
        err = xyerror.ErrResUnkownOpType
        return
    }

    if err != nil {
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
    }

    return
}

//签到活动配置信息
func (api *XYAPI) opSignInActivity(optype battery.RES_OP_TYPE, activity *battery.DBSignInActivity) (failReason battery.ErrorCode, err error) {
    switch optype {
    case battery.RES_OP_TYPE_OP_ADD:
        err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_SIGNINACTIVITY_CONFIG).UpsertSignInActivity(activity)
    default:
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
        err = xyerror.ErrResUnkownOpType
        return
    }

    if err != nil {
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
    }

    return
}

//签到奖励配置信息
func (api *XYAPI) opSignInItem(optype battery.RES_OP_TYPE, item *battery.DBSignInItem) (failReason battery.ErrorCode, err error) {
    switch optype {
    case battery.RES_OP_TYPE_OP_ADD:
        err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_SIGNAWARD_CONFIG).UpsertSignInItem(item)
    default:
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
        err = xyerror.ErrResUnkownOpType
        return
    }

    if err != nil {
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
    }

    return
}

//系统道具配置
func (api *XYAPI) opRuneConfig(optype battery.RES_OP_TYPE, item *battery.RuneConfig) (failReason battery.ErrorCode, err error) {
    switch optype {
    case battery.RES_OP_TYPE_OP_ADD:
        err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_RUNE_CONFIG).UpsertRuneConfig(item)
    default:
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
        err = xyerror.ErrResUnkownOpType
        return
    }

    if err != nil {
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
    }

    return
}

//赛前道具配置
func (api *XYAPI) opBeforeGameRandomWeightConfig(optype battery.RES_OP_TYPE, item *battery.DBBeforeGameRandomGoodWeight) (failReason battery.ErrorCode, err error) {
    switch optype {
    case battery.RES_OP_TYPE_OP_ADD:
        err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_BEFOREGAMERANDOM_CONFIG).UpsertBeforeGameRandomGoodsConfig(item)
    default:
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
        err = xyerror.ErrResUnkownOpType
        return
    }

    if err != nil {
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
    }

    return
}

//角色信息列表
func (api *XYAPI) opRoleInfoConfig(optype battery.RES_OP_TYPE, item *battery.DBRoleInfoConfig) (failReason battery.ErrorCode, err error) {
    switch optype {
    case battery.RES_OP_TYPE_OP_ADD:
        err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_ROLEINFO_CONFIG).UpsertRoleConfig(*item)
    default:
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
        err = xyerror.ErrResUnkownOpType
        return
    }

    if err != nil {
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
    }

    return
}

//邮件新消息列表
func (api *XYAPI) opMailConfig(optype battery.RES_OP_TYPE, item *battery.SystemMailInfoConfig) (failReason battery.ErrorCode, err error) {
    switch optype {
    case battery.RES_OP_TYPE_OP_ADD:
        err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_MAILINFO_CONFIG).UpsertMailInfoConfig(*item)
    default:
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
        err = xyerror.ErrResUnkownOpType
        return
    }

    if err != nil {
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
    }

    return
}

//拼图信息列表
func (api *XYAPI) opJigsawConfig(optype battery.RES_OP_TYPE, item *battery.JigsawConfig) (failReason battery.ErrorCode, err error) {
    switch optype {
    case battery.RES_OP_TYPE_OP_ADD:
        err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_JIGSAW_CONFIG).UpsertJigsawConfig(*item)
    default:
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
        err = xyerror.ErrResUnkownOpType
        return
    }

    if err != nil {
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
    }

    return
}

//收集物信息列表
func (api *XYAPI) opPickUp(optype battery.RES_OP_TYPE, item *battery.DBPickUpItem) (failReason battery.ErrorCode, err error) {
    switch optype {
    case battery.RES_OP_TYPE_OP_ADD:
        err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_PICKUPWEIGHT_CONFIG).AddPickUp(item)
    default:
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
        err = xyerror.ErrResUnkownOpType
        return
    }

    if err != nil {
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
    }

    return
}

//角色加成信息
func (api *XYAPI) opRoleLevelBonus(optype battery.RES_OP_TYPE, item *battery.DBRoleLevelBonusItem) (failReason battery.ErrorCode, err error) {
    switch optype {
    case battery.RES_OP_TYPE_OP_ADD:
        err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_ROLELEVELBONUS_CONFIG).UpsertRoleLevelBonus(item)
    default:
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
        err = xyerror.ErrResUnkownOpType
        return
    }

    if err != nil {
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
    }

    return
}

//公告信息
func (api *XYAPI) opAnnouncement(optype battery.RES_OP_TYPE, item *battery.DBAnnouncementConfig) (failReason battery.ErrorCode, err error) {
    switch optype {
    case battery.RES_OP_TYPE_OP_ADD:
        err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_ANNOUNCEMENT_CONFIG).AddAnnouncementConfig(item)
    default:
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
        err = xyerror.ErrResUnkownOpType
        return
    }

    if err != nil {
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
    }

    return
}

//广告信息
func (api *XYAPI) opAdvertisement(optype battery.RES_OP_TYPE, item *battery.Advertisement) (failReason battery.ErrorCode, err error) {
    switch optype {
    case battery.RES_OP_TYPE_OP_ADD:
        err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_ANNOUNCEMENT_CONFIG).AddAdvertisementConfig(item)
    default:
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
        err = xyerror.ErrResUnkownOpType
        return
    }

    if err != nil {
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
    }

    return
}

func (api *XYAPI) opAdvertisementSpace(optype battery.RES_OP_TYPE, item *battery.AdvertisementSpace) (failReason battery.ErrorCode, err error) {
    switch optype {
    case battery.RES_OP_TYPE_OP_ADD:
        err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_ADVERTISEMENT_SPACE_CONFIG).AddAdvertisementSpaceConfig(item)
    default:
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
        err = xyerror.ErrResUnkownOpType
        return
    }

    if err != nil {
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
    }

    return
}

func (api *XYAPI) opTip(optype battery.RES_OP_TYPE, item *battery.DBTip) (failReason battery.ErrorCode, err error) {
    switch optype {
    case battery.RES_OP_TYPE_OP_ADD:
        err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_TIP_CONFIG).AddTipConfig(item)
    default:
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
        err = xyerror.ErrResUnkownOpType
        return
    }

    if err != nil {
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
    }

    return
}

func (api *XYAPI) opNewAccountProp(optype battery.RES_OP_TYPE, item *battery.DBNewAccountProp) (failReason battery.ErrorCode, err error) {
    switch optype {
    case battery.RES_OP_TYPE_OP_ADD:
        err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_NEWACCOUNTPROP_CONFIG).AddNewAccountPropConfig(item)
    default:
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
        err = xyerror.ErrResUnkownOpType
        return
    }

    if err != nil {
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
    }

    return
}

// 分享活动配置
func (api *XYAPI) opShareActivity(optype battery.RES_OP_TYPE, item *battery.DBShareActivity) (failReason battery.ErrorCode, err error) {
    switch optype {
    case battery.RES_OP_TYPE_OP_ADD:
        err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_SHARE_ACTIVITY).AddShareActivityConfig(item)
    default:
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
        err = xyerror.ErrResUnkownOpType
        return
    }

    if err != nil {
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
    }

    return
}

// 分享活动配置
func (api *XYAPI) opShareAwards(optype battery.RES_OP_TYPE, item *battery.DBShareAward) (failReason battery.ErrorCode, err error) {
    switch optype {
    case battery.RES_OP_TYPE_OP_ADD:
        err = api.GetCommonDB(xybusiness.BUSINESS_COMMON_COLLECTION_INDEX_SHARE_AWARDS).AddShareAwardsConfig(item)
    default:
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
        err = xyerror.ErrResUnkownOpType
        return
    }

    if err != nil {
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
    }

    return
}

// 关卡解锁商品配置
func (api *XYAPI) opUnlockGoodsConfig(optype battery.RES_OP_TYPE, item *battery.DBCheckPointUnlockGoodsConfig) (failReason battery.ErrorCode, err error) {
    switch optype {
    case battery.RES_OP_TYPE_OP_ADD:
        err = api.GetCommonDB(xybusiness.Business_COMMON_COLLECTION_INDEX_CHECKPOINTUNLOCK_GOODS_CONFIG).AddCheckPointUnlockGoodsConfig(item)
    default:
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
        err = xyerror.ErrResUnkownOpType
        return
    }

    if err != nil {
        failReason = xyerror.Resp_ResUnkownOpType.GetCode()
    }
    return
}
