package main

import (
    proto "code.google.com/p/goprotobuf/proto"
    business "guanghuan.com/xiaoyao/battery_transaction_server/business"
    xynatsservice "guanghuan.com/xiaoyao/common/service/nats"
    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
    "guanghuan.com/xiaoyao/superbman_server/server"
)

var DefMsgHandlerMap xynatsservice.MsgCodeHandlerMap = make(xynatsservice.MsgCodeHandlerMap, 20)

func OperationIapValidate(req proto.Message, resp proto.Message) error {
    return business.NewXYAPI().OperationIapValidate(req.(*battery.OrderVerifyRequest), resp.(*battery.OrderVerifyResponse))
}

func OperationLottoRequest(req proto.Message, resp proto.Message) error {
    return business.NewXYAPI().OperationLottoRequest(req.(*battery.LottoRequest), resp.(*battery.LottoResponse))
}

func OperationLogin(req proto.Message, resp proto.Message) error {
    return business.NewXYAPI().OperationLogin(req.(*battery.LoginRequest), resp.(*battery.LoginResponse), "")
}

func OperationNewGame(req proto.Message, resp proto.Message) error {
    return business.NewXYAPI().OperationNewGame(req.(*battery.NewGameRequest), resp.(*battery.NewGameResponse))
}

func OperationGameResultCommit(req proto.Message, resp proto.Message) error {
    return business.NewXYAPI().OperationGameResultCommit(req.(*battery.GameResultCommitRequest), resp.(*battery.GameResultCommitResponse))
}

func OperationQueryWallet(req proto.Message, resp proto.Message) error {
    return business.NewXYAPI().OperationQueryWallet(req.(*battery.QueryWalletRequest), resp.(*battery.QueryWalletResponse))
}

func OperationQueryGoods(req proto.Message, resp proto.Message) error {
    return business.NewXYAPI().OperationQueryGoods(req.(*battery.QueryGoodsRequest), resp.(*battery.QueryGoodsResponse))
}

func OperationBuyGoods(req proto.Message, resp proto.Message) error {
    return business.NewXYAPI().OperationBuyGoods(req.(*battery.BuyGoodsRequest), resp.(*battery.BuyGoodsResponse))
}

func OperationQuerySignIn(req proto.Message, resp proto.Message) error {
    return business.NewXYAPI().OperationQuerySignIn(req.(*battery.QuerySignInRequest), resp.(*battery.QuerySignInResponse))
}

// 添加新版本签到消息接口
func OperationQuerySignInNew(req proto.Message, resp proto.Message) error {
    return business.NewXYAPI().OperationQuerySignIn2(req.(*battery.NewQuerySignInRequest), resp.(*battery.NewQuerySignInResponse))
}
func OperationSignIn(req proto.Message, resp proto.Message) error {
    return business.NewXYAPI().OperationSignIn(req.(*battery.SignInRequest), resp.(*battery.SignInResponse))
}

func OperationQueryUserMission(req proto.Message, resp proto.Message) error {
    return business.NewXYAPI().OperationQueryUserMission(req.(*battery.QueryUserMissionRequest), resp.(*battery.QueryUserMissionResponse))
}

func OperationConfirmUserMission(req proto.Message, resp proto.Message) error {
    return business.NewXYAPI().OperationConfirmUserMission(req.(*battery.ConfirmUserMissionRequest), resp.(*battery.ConfirmUserMissionResponse))
}

func OperationRoleInfoList(req proto.Message, resp proto.Message) error {
    return business.NewXYAPI().OperationRoleInfoList(req.(*battery.RoleInfoListRequest), resp.(*battery.RoleInfoListResponse))
}
func OperationFriendMailInfoList(req proto.Message, resp proto.Message) error {
    return business.NewXYAPI().OperationFriendMailInfoList(req.(*battery.FriendMailListRequest), resp.(*battery.FriendMailListResponse))
}

func OperationSystemMailInfoList(req proto.Message, resp proto.Message) error {
    return business.NewXYAPI().OperationSystemMailInfoList(req.(*battery.SystemMailListRequest), resp.(*battery.SystemMailListResponse))
}

func OperationAnnouncement(req proto.Message, resp proto.Message) error {
    return business.NewXYAPI().OperationAnnouncement(req.(*battery.AnnouncementRequest), resp.(*battery.AnnouncementResponse))
}

func OperationJigsaw(req proto.Message, resp proto.Message) error {
    return business.NewXYAPI().OperationJigsaw(req.(*battery.JigsawRequest), resp.(*battery.JigsawResponse))
}

func OperationRune(req proto.Message, resp proto.Message) error {
    return business.NewXYAPI().OperationRune(req.(*battery.RuneRequest), resp.(*battery.RuneResponse))
}

func OperationBeforeGameProp(req proto.Message, resp proto.Message) error {
    return business.NewXYAPI().OperationBeforeGameProp(req.(*battery.BeforeGamePropRequest), resp.(*battery.BeforeGamePropResponse))
}

func OperationQueryUserCheckPoints(req proto.Message, resp proto.Message) error {
    return business.NewXYAPI().OperationQueryUserCheckPoints(req.(*battery.QueryUserCheckPointsRequest), resp.(*battery.QueryUserCheckPointsResponse))
}

func OperationQueryUserCheckPointDetail(req proto.Message, resp proto.Message) error {
    return business.NewXYAPI().OperationQueryCheckPointDetail(req.(*battery.QueryUserCheckPointDetailRequest), resp.(*battery.QueryUserCheckPointDetailResponse))
}

// func OperationCommitCheckPoint(req proto.Message, resp proto.Message) error {
//     return business.NewXYAPI().OperationCommitCheckPoint(req.(*battery.CommitCheckPointRequest), resp.(*battery.CommitCheckPointResponse))
// }

func OperationCheckPointUnlock(req proto.Message, resp proto.Message) error {
    return business.NewXYAPI().OperationCheckPointUnlock(req.(*battery.CheckPointUnlockRequest), resp.(*battery.CheckPointUnlockResponse))
}

func OperationMemCache(req proto.Message, resp proto.Message) error {
    return business.NewXYAPI().OperationMemCache(req.(*battery.MemCacheRequest), resp.(*battery.MemCacheResponse))
}

func OperationMaintenanceProp(req proto.Message, resp proto.Message) error {
    return business.NewXYAPI().OperationMaintenanceProp(req.(*battery.MaintenancePropRequest), resp.(*battery.MaintenancePropResponse))
}

func OperationBind(req proto.Message, resp proto.Message) error {
    return business.NewXYAPI().OperationBind(req.(*battery.BindRequest), resp.(*battery.BindResponse))
}

func OperationSharedQuery(req proto.Message, resp proto.Message) error {
    return business.NewXYAPI().OperationSharedQuery(req.(*battery.SharedQueryRequest), resp.(*battery.SharedQueryResponse))
}

func OperationSharedRequest(req proto.Message, resp proto.Message) error {
    return business.NewXYAPI().OperationSharedRequest(req.(*battery.SharedRequest), resp.(*battery.SharedResponse))
}

func OperationSDKOrderOp(req proto.Message, resp proto.Message) error {
    return business.NewXYAPI().OperationSDKOrderOp(req.(*battery.SDKOrderOperationRequest), resp.(*battery.SDKOrderOperationResponse))
}

func OperationSDKOrderQuery(req proto.Message, resp proto.Message) error {
    return business.NewXYAPI().OperationSDKOrderQuery(req.(*battery.SDKOrderRequest), resp.(*battery.SDKOrderResponse))
}

func OperationSDKAddOrder(req proto.Message, resp proto.Message) error {
    return business.NewXYAPI().OperationSDKAddOrder(req.(*battery.SDKAddOrderRequest), resp.(*battery.SDKAddOrderResponse))
}

func OperationGetGlobalRankList(req proto.Message, resp proto.Message) error {
    return business.NewXYAPI().OperationGetGlobalRankList(req.(*battery.QueryGlobalRankRequest), resp.(*battery.QueryGlobalRankResponse))
}
func OperationCreatName(req proto.Message, resp proto.Message) error {
    return business.NewXYAPI().OperationCreatName(req.(*battery.CreatNameRequest), resp.(*battery.CreatNameResponse))
}
func initMsgHandlerMap() {
    DefMsgHandlerMap.AddHandler(xybusiness.BusinessCode_Login, battery.LoginRequest{}, battery.LoginResponse{}, OperationLogin)
    DefMsgHandlerMap.AddHandler(xybusiness.BusinessCode_NewGame, battery.NewGameRequest{}, battery.NewGameResponse{}, OperationNewGame)
    DefMsgHandlerMap.AddHandler(xybusiness.BusinessCode_GameResult2, battery.GameResultCommitRequest{}, battery.GameResultCommitResponse{}, OperationGameResultCommit)
    DefMsgHandlerMap.AddHandler(xybusiness.BusinessCode_QueryWallet, battery.QueryWalletRequest{}, battery.QueryWalletResponse{}, OperationQueryWallet)
    DefMsgHandlerMap.AddHandler(xybusiness.BusinessCode_GoodsBuy, battery.BuyGoodsRequest{}, battery.BuyGoodsResponse{}, OperationBuyGoods)
    DefMsgHandlerMap.AddHandler(xybusiness.BusinessCode_QueryGoods, battery.QueryGoodsRequest{}, battery.QueryGoodsResponse{}, OperationQueryGoods)
    DefMsgHandlerMap.AddHandler(xybusiness.BusinessCode_IapValidate, battery.OrderVerifyRequest{}, battery.OrderVerifyResponse{}, OperationIapValidate)
    DefMsgHandlerMap.AddHandler(xybusiness.BusinessCode_RoleInfoList, battery.RoleInfoListRequest{}, battery.RoleInfoListResponse{}, OperationRoleInfoList)
    DefMsgHandlerMap.AddHandler(xybusiness.BusinessCode_SystemMailInfoList, battery.SystemMailListRequest{}, battery.SystemMailListResponse{}, OperationSystemMailInfoList)
    DefMsgHandlerMap.AddHandler(xybusiness.BusinessCode_FriendMailInfoList, battery.FriendMailListRequest{}, battery.FriendMailListResponse{}, OperationFriendMailInfoList)
    DefMsgHandlerMap.AddHandler(xybusiness.BusinessCode_Announcement, battery.AnnouncementRequest{}, battery.AnnouncementResponse{}, OperationAnnouncement)
    DefMsgHandlerMap.AddHandler(xybusiness.BusinessCode_Lotto, battery.LottoRequest{}, battery.LottoResponse{}, OperationLottoRequest)
    DefMsgHandlerMap.AddHandler(xybusiness.BusinessCode_QuerySignIn, battery.QuerySignInRequest{}, battery.QuerySignInResponse{}, OperationQuerySignIn)
    DefMsgHandlerMap.AddHandler(xybusiness.BusinessCode_SignIn, battery.SignInRequest{}, battery.SignInResponse{}, OperationSignIn)
    DefMsgHandlerMap.AddHandler(xybusiness.BusinessCode_Jigsaw, battery.JigsawRequest{}, battery.JigsawResponse{}, OperationJigsaw)
    DefMsgHandlerMap.AddHandler(xybusiness.BusinessCode_Rune, battery.RuneRequest{}, battery.RuneResponse{}, OperationRune)
    DefMsgHandlerMap.AddHandler(xybusiness.BusinessCode_BeforeGameProp, battery.BeforeGamePropRequest{}, battery.BeforeGamePropResponse{}, OperationBeforeGameProp)
    DefMsgHandlerMap.AddHandler(xybusiness.BusinessCode_QueryUserMission, battery.QueryUserMissionRequest{}, battery.QueryUserMissionResponse{}, OperationQueryUserMission)
    DefMsgHandlerMap.AddHandler(xybusiness.BusinessCode_ConfirmUserMission, battery.ConfirmUserMissionRequest{}, battery.ConfirmUserMissionResponse{}, OperationConfirmUserMission)
    DefMsgHandlerMap.AddHandler(xybusiness.BusinessCode_QueryUserCheckPoints, battery.QueryUserCheckPointsRequest{}, battery.QueryUserCheckPointsResponse{}, OperationQueryUserCheckPoints)
    DefMsgHandlerMap.AddHandler(xybusiness.BusinessCode_QueryUserCheckPointDetail, battery.QueryUserCheckPointDetailRequest{}, battery.QueryUserCheckPointDetailResponse{}, OperationQueryUserCheckPointDetail)
    // DefMsgHandlerMap.AddHandler(xybusiness.BusinessCode_CommitCheckPoint, battery.CommitCheckPointRequest{}, battery.CommitCheckPointResponse{}, OperationCommitCheckPoint)
    DefMsgHandlerMap.AddHandler(xybusiness.BusinessCode_CheckPointUnlock, battery.CheckPointUnlockRequest{}, battery.CheckPointUnlockResponse{}, OperationCheckPointUnlock)
    DefMsgHandlerMap.AddHandler(xybusiness.BusinessCode_MemCache, battery.MemCacheRequest{}, battery.MemCacheResponse{}, OperationMemCache)
    DefMsgHandlerMap.AddHandler(xybusiness.BusinessCode_MaintenanceProp, battery.MaintenancePropRequest{}, battery.MaintenancePropResponse{}, OperationMaintenanceProp)
    DefMsgHandlerMap.AddHandler(xybusiness.BusinessCode_Bind, battery.BindRequest{}, battery.BindResponse{}, OperationBind)

    DefMsgHandlerMap.AddHandler(xybusiness.BusinessCode_QuerySignIn2, battery.NewQuerySignInRequest{}, battery.NewQuerySignInResponse{}, OperationQuerySignInNew)

    DefMsgHandlerMap.AddHandler(xybusiness.BusinessCode_ShareQuery, battery.SharedQueryRequest{}, battery.SharedQueryResponse{}, OperationSharedQuery)
    DefMsgHandlerMap.AddHandler(xybusiness.BusinessCode_ShareRequest, battery.SharedRequest{}, battery.SharedResponse{}, OperationSharedRequest)
    DefMsgHandlerMap.AddHandler(xybusiness.BusinessCode_SDKOrderOp, battery.SDKOrderOperationRequest{}, battery.SDKOrderOperationResponse{}, OperationSDKOrderOp)
    DefMsgHandlerMap.AddHandler(xybusiness.BusinessCode_SDKOrderQuery, battery.SDKOrderRequest{}, battery.SDKOrderResponse{}, OperationSDKOrderQuery)
    DefMsgHandlerMap.AddHandler(xybusiness.BusinessCode_SDKAddOrder, battery.SDKAddOrderRequest{}, battery.SDKAddOrderResponse{}, OperationSDKAddOrder)

    DefMsgHandlerMap.AddHandler(xybusiness.BusinessCode_GlobalRankList, battery.QueryGlobalRankRequest{}, battery.QueryGlobalRankResponse{}, OperationGetGlobalRankList)
    DefMsgHandlerMap.AddHandler(xybusiness.BusinessCode_CreatName, battery.CreatNameRequest{}, battery.CreatNameResponse{}, OperationCreatName)
}
