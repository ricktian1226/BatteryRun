package main

//import (
//	proto "code.google.com/p/goprotobuf/proto"
//	business "guanghuan.com/xiaoyao/battery_app_server/business"
//	xynatsservice "guanghuan.com/xiaoyao/common/service/nats"
//	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
//)

//var DefMsgHandlerMap xynatsservice.MsgHandlerMap = make(xynatsservice.MsgHandlerMap, 20)

//func OperationLogin(req proto.Message, resp proto.Message) error {
//	return business.NewXYAPI().OperationLogin(req.(*battery.LoginRequest), resp.(*battery.LoginResponse), "")
//}
//func OperationGetFriendData(req proto.Message, resp proto.Message) error {
//	return business.NewXYAPI().OperationGetFriendData(req.(*battery.QueryFriendsDataRequest), resp.(*battery.QueryFriendsDataResponse))
//}
//func OperationNewGame(req proto.Message, resp proto.Message) error {
//	return business.NewXYAPI().OperationNewGame(req.(*battery.NewGameRequest), resp.(*battery.NewGameResponse))
//}
//func OperationAddGameData(req proto.Message, resp proto.Message) error {
//	return business.NewXYAPI().OperationAddGameData(req.(*battery.GameDataRequest), resp.(*battery.GameDataResponse))
//}
//func OperationGameResultCommit(req proto.Message, resp proto.Message) error {
//	return business.NewXYAPI().OperationGameResultCommit(req.(*battery.GameResultCommitRequest), resp.(*battery.GameResultCommitResponse))
//}

////func OperationQueryStamina(req proto.Message, resp proto.Message) error {
////	return api_service.Api.OperationQueryStamina(req.(*battery.QueryStaminaRequest), resp.(*battery.QueryStaminaResponse))
////}

//func OperationQueryWallet(req proto.Message, resp proto.Message) error {
//	return business.NewXYAPI().OperationQueryWallet(req.(*battery.QueryWalletRequest), resp.(*battery.QueryWalletResponse))
//}

////func OperationQueryGift(req proto.Message, resp proto.Message) error {
////	return api_service.Api.OperationQueryGift(req.(*battery.QueryStaminaGiftRequest), resp.(*battery.QueryStaminaGiftResponse))
////}
//func OperationStaminaGiftOp(req proto.Message, resp proto.Message) error {
//	return business.NewXYAPI().OperationStaminaGiftOp(req.(*battery.StaminaGiftRequest), resp.(*battery.StaminaGiftResponse))
//}
//func OperationQueryGoods(req proto.Message, resp proto.Message) error {
//	return business.NewXYAPI().OperationQueryGoods(req.(*battery.QueryGoodsRequest), resp.(*battery.QueryGoodsResponse))
//}
//func OperationBuyGoods(req proto.Message, resp proto.Message) error {
//	return business.NewXYAPI().OperationBuyGoods(req.(*battery.BuyGoodsRequest), resp.(*battery.BuyGoodsResponse))
//}

////func OperationIapOrder(req proto.Message, resp proto.Message) error {
////	return api_service.Api.OperationIapOrder(req.(*battery.OrderNumRequest), resp.(*battery.OrderNumResponse))
////}

//func OperationIapValidate(req proto.Message, resp proto.Message) error {
//	return business.NewXYAPI().OperationIapValidate(req.(*battery.OrderVerifyRequest), resp.(*battery.OrderVerifyResponse))
//}
//func OperationSubmitDeviceId(req proto.Message, resp proto.Message) error {
//	return business.NewXYAPI().OperationSubmitDeviceId(req.(*battery.DeviceIdSubmitRequest), resp.(*battery.DeviceIdSubmitResponse))
//}
//func OperationAnnouncement(req proto.Message, resp proto.Message) error {
//	return business.NewXYAPI().OperationAnnouncement(req.(*battery.AnnouncementRequest), resp.(*battery.AnnouncementResponse))
//}

//func OperationResRequest(req proto.Message, resp proto.Message) error {
//	return business.NewXYAPI().OperationResRequest(req.(*battery.ResRequest), resp.(*battery.ResResponse))
//}

//func OperationLottoRequest(req proto.Message, resp proto.Message) error {
//	return business.NewXYAPI().OperationLottoRequest(req.(*battery.LottoRequest), resp.(*battery.LottoResponse))
//}

//func OperationSystemMailListRequest(req proto.Message, resp proto.Message) error {
//	return business.NewXYAPI().OperationSystemMailListRequest(req.(*battery.SystemMailListRequest), resp.(*battery.SystemMailListResponse))
//}

//func OperationFriendMailListRequest(req proto.Message, resp proto.Message) error {
//	return business.NewXYAPI().OperationFriendMailListRequest(req.(*battery.FriendMailListRequest), resp.(*battery.FriendMailListResponse))
//}

//func OperationRoleInfoListRequest(req proto.Message, resp proto.Message) error {
//	return business.NewXYAPI().OperationRoleInfoListRequest(req.(*battery.RoleInfoListRequest), resp.(*battery.RoleInfoListResponse))
//}

//func OperationQuerySignIn(req proto.Message, resp proto.Message) error {
//	return business.NewXYAPI().OperationQuerySignIn(req.(*battery.QuerySignInRequest), resp.(*battery.QuerySignInResponse))
//}

//func OperationSignIn(req proto.Message, resp proto.Message) error {
//	return business.NewXYAPI().OperationSignIn(req.(*battery.SignInRequest), resp.(*battery.SignInResponse))
//}

//func OperationQueryUserMission(req proto.Message, resp proto.Message) error {
//	return business.NewXYAPI().OperationQueryUserMission(req.(*battery.QueryUserMissionRequest), resp.(*battery.QueryUserMissionResponse))
//}

//func OperationConfirmUserMission(req proto.Message, resp proto.Message) error {
//	return business.NewXYAPI().OperationConfirmUserMission(req.(*battery.ConfirmUserMissionRequest), resp.(*battery.ConfirmUserMissionResponse))
//}

//func OperationRuneRequest(req proto.Message, resp proto.Message) error {
//	return business.NewXYAPI().OperationRuneRequest(req.(*battery.RuneRequest), resp.(*battery.RuneResponse))
//}

//func OperationJigsawRequest(req proto.Message, resp proto.Message) error {
//	return business.NewXYAPI().OperationJigsawRequest(req.(*battery.JigsawRequest), resp.(*battery.JigsawResponse))
//}

//func OperationBeforeGameProp(req proto.Message, resp proto.Message) error {
//	return business.NewXYAPI().OperationBeforeGameProp(req.(*battery.BeforeGamePropRequest), resp.(*battery.BeforeGamePropResponse))
//}

//func OperationQueryPropRes(req proto.Message, resp proto.Message) error {
//	return business.NewXYAPI().OperationQueryPropRes(req.(*battery.QueryPropResRequest), resp.(*battery.QueryPropResResponse))
//}

//func OperationQueryUserCheckPoints(req proto.Message, resp proto.Message) error {
//	return business.NewXYAPI().OperationQueryUserCheckPoints(req.(*battery.QueryUserCheckPointsRequest), resp.(*battery.QueryUserCheckPointsResponse))
//}

//func OperationQueryUserCheckPointDetail(req proto.Message, resp proto.Message) error {
//	return business.NewXYAPI().OperationQueryCheckPointDetail(req.(*battery.QueryUserCheckPointDetailRequest), resp.(*battery.QueryUserCheckPointDetailResponse))
//}

//func OperationCommitCheckPoint(req proto.Message, resp proto.Message) error {
//	return business.NewXYAPI().OperationCommitCheckPoint(req.(*battery.CommitCheckPointRequest), resp.(*battery.CommitCheckPointResponse))
//}

//func OperationMemCache(req proto.Message, resp proto.Message) error {
//	return business.NewXYAPI().OperationMemCache(req.(*battery.MemCacheRequest), resp.(*battery.MemCacheResponse))
//}

//func initMsgHandlerMap() {
//	DefMsgHandlerMap.AddHandler("login", battery.LoginRequest{}, battery.LoginResponse{}, OperationLogin)
//	DefMsgHandlerMap.AddHandler("frienddata", battery.QueryFriendsDataRequest{}, battery.QueryFriendsDataResponse{}, OperationGetFriendData)

//	DefMsgHandlerMap.AddHandler("newgame", battery.NewGameRequest{}, battery.NewGameResponse{}, OperationNewGame)
//	DefMsgHandlerMap.AddHandler("gameresult", battery.GameDataRequest{}, battery.GameDataResponse{}, OperationAddGameData)
//	DefMsgHandlerMap.AddHandler("gameresult2", battery.GameResultCommitRequest{}, battery.GameResultCommitResponse{}, OperationGameResultCommit)

//	DefMsgHandlerMap.AddHandler("wallet_query", battery.QueryWalletRequest{}, battery.QueryWalletResponse{}, OperationQueryWallet)

//	DefMsgHandlerMap.AddHandler("gift_op", battery.StaminaGiftRequest{}, battery.StaminaGiftResponse{}, OperationStaminaGiftOp)

//	DefMsgHandlerMap.AddHandler("goods_query", battery.QueryGoodsRequest{}, battery.QueryGoodsResponse{}, OperationQueryGoods)
//	DefMsgHandlerMap.AddHandler("goods_buy", battery.BuyGoodsRequest{}, battery.BuyGoodsResponse{}, OperationBuyGoods)

//	DefMsgHandlerMap.AddHandler("iap_verify_validate", battery.OrderVerifyRequest{}, battery.OrderVerifyResponse{}, OperationIapValidate)

//	DefMsgHandlerMap.AddHandler("device_id", battery.DeviceIdSubmitRequest{}, battery.DeviceIdSubmitResponse{}, OperationSubmitDeviceId)
//	DefMsgHandlerMap.AddHandler("announcement", battery.AnnouncementRequest{}, battery.AnnouncementResponse{}, OperationAnnouncement)

//	DefMsgHandlerMap.AddHandler("lotto_op", battery.LottoRequest{}, battery.LottoResponse{}, OperationLottoRequest)

//	DefMsgHandlerMap.AddHandler("systemmail_op", battery.SystemMailListRequest{}, battery.SystemMailListResponse{}, OperationSystemMailListRequest)
//	DefMsgHandlerMap.AddHandler("friendmail_op", battery.FriendMailListRequest{}, battery.FriendMailListResponse{}, OperationFriendMailListRequest)

//	DefMsgHandlerMap.AddHandler("roleinfolist_op", battery.RoleInfoListRequest{}, battery.RoleInfoListResponse{}, OperationRoleInfoListRequest)

//	DefMsgHandlerMap.AddHandler("jigsaw_op", battery.JigsawRequest{}, battery.JigsawResponse{}, OperationJigsawRequest)

//	DefMsgHandlerMap.AddHandler("signin_query", battery.QuerySignInRequest{}, battery.QuerySignInResponse{}, OperationQuerySignIn)
//	DefMsgHandlerMap.AddHandler("signin_sign", battery.SignInRequest{}, battery.SignInResponse{}, OperationSignIn)

//	DefMsgHandlerMap.AddHandler("usermission_query", battery.QueryUserMissionRequest{}, battery.QueryUserMissionResponse{}, OperationQueryUserMission)
//	DefMsgHandlerMap.AddHandler("usermission_confirm", battery.ConfirmUserMissionRequest{}, battery.ConfirmUserMissionResponse{}, OperationConfirmUserMission)

//	DefMsgHandlerMap.AddHandler("checkpoint_query_range", battery.QueryUserCheckPointsRequest{}, battery.QueryUserCheckPointsResponse{}, OperationQueryUserCheckPoints)
//	DefMsgHandlerMap.AddHandler("checkpoint_query_detail", battery.QueryUserCheckPointDetailRequest{}, battery.QueryUserCheckPointDetailResponse{}, OperationQueryUserCheckPointDetail)
//	DefMsgHandlerMap.AddHandler("checkpoint_commit", battery.CommitCheckPointRequest{}, battery.CommitCheckPointResponse{}, OperationCommitCheckPoint)

//	DefMsgHandlerMap.AddHandler("prop_res_query", battery.QueryPropResRequest{}, battery.QueryPropResResponse{}, OperationQueryPropRes)

//	DefMsgHandlerMap.AddHandler("rune_op", battery.RuneRequest{}, battery.RuneResponse{}, OperationRuneRequest)
//	DefMsgHandlerMap.AddHandler("beforegameprop_op", battery.BeforeGamePropRequest{}, battery.BeforeGamePropResponse{}, OperationBeforeGameProp)

//	//运营接口
//	DefMsgHandlerMap.AddHandler("res_op", battery.ResRequest{}, battery.ResResponse{}, OperationResRequest)

//	//玩家的memcache接口
//	DefMsgHandlerMap.AddHandler("memcache_op", battery.MemCacheRequest{}, battery.MemCacheResponse{}, OperationMemCache)

//}
