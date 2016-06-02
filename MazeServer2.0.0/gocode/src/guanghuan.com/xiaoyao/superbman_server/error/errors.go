package xyerror

import (
    "code.google.com/p/goprotobuf/proto"
    "encoding/xml"
    "errors"
    "fmt"
    battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
)

//返回默认的错误结构
func DefaultError() *battery.Error {
    return &battery.Error{Code: battery.ErrorCode_NoError.Enum()}
}

func ConstructError(errCode battery.ErrorCode) *battery.Error {
    return &battery.Error{Code: errCode.Enum()}
}

func ErrorStructFunc(errCode battery.ErrorCode) *battery.Error {
    return &battery.Error{Code: errCode.Enum(), Desc: proto.String(errCode.String())}
}

var (
    Resp_NoError                = ErrorStructFunc(battery.ErrorCode_NoError)
    Resp_UnknowError            = ErrorStructFunc(battery.ErrorCode_UnknowError)
    Resp_BadInputData           = ErrorStructFunc(battery.ErrorCode_BadInputData)
    Resp_ServerError            = ErrorStructFunc(battery.ErrorCode_ServerError)
    Resp_DBError                = ErrorStructFunc(battery.ErrorCode_DBError)
    Resp_ServiceNotAvaliable    = ErrorStructFunc(battery.ErrorCode_ServiceNotAvaliable)
    Resp_ClientNotSupport       = ErrorStructFunc(battery.ErrorCode_ClientVersionNotSupport)
    Resp_SendToTransactionError = ErrorStructFunc(battery.ErrorCode_SendToTransactionError)
    Resp_TimeLimitError         = ErrorStructFunc(battery.ErrorCode_TimeLimitError)

    Resp_GetAccountByUidError = ErrorStructFunc(battery.ErrorCode_GetAccountByUidError)
    Resp_UpdateAccountError   = ErrorStructFunc(battery.ErrorCode_UpdateAccountError)
    Resp_NotEnoughCurrency    = ErrorStructFunc(battery.ErrorCode_NotEnoughCurrency)
    Resp_GetUidError          = ErrorStructFunc(battery.ErrorCode_GetUidError)

    Resp_QueryStaminaError  = ErrorStructFunc(battery.ErrorCode_QueryStaminaError)
    Resp_UpdateStaminaError = ErrorStructFunc(battery.ErrorCode_UpdateStaminaError)
    Resp_NotEnoughStamina   = ErrorStructFunc(battery.ErrorCode_NotEnoughStamina)

    Resp_NotEnoughDiamond = ErrorStructFunc(battery.ErrorCode_NotEnoughDiamond)

    Resp_AddNewGameError        = ErrorStructFunc(battery.ErrorCode_AddNewGameError)
    Resp_GameNotExistError      = ErrorStructFunc(battery.ErrorCode_GameNotExistError)
    Resp_GameResultInvalidError = ErrorStructFunc(battery.ErrorCode_GameResultInvalidError)
    Resp_UpdateGameError        = ErrorStructFunc(battery.ErrorCode_UpdateGameError)

    Resp_QueryGoodsError            = ErrorStructFunc(battery.ErrorCode_QueryGoodsError)
    Resp_BuyGoodExceedLimit         = ErrorStructFunc(battery.ErrorCode_BuyGoodExceedLimit)
    Resp_BuyGoodInvalidGame         = ErrorStructFunc(battery.ErrorCode_BuyGoodInvalidGame)
    Resp_BuyGoodsError              = ErrorStructFunc(battery.ErrorCode_BuyGoodsError)
    Resp_BuyGoodOverAmountPerGame   = ErrorStructFunc(battery.ErrorCode_BuyGoodOverAmountPerGame)
    Resp_BuyGoodOverAmountPerUser   = ErrorStructFunc(battery.ErrorCode_BuyGoodOverAmountPerUser)
    Resp_BuyGoodOverAmountPerDay    = ErrorStructFunc(battery.ErrorCode_BuyGoodOverAmountPerDay)
    Resp_BuyGoodNewReceiptError     = ErrorStructFunc(battery.ErrorCode_BuyGoodNewReceiptError)
    Resp_BuyGoodConsumMoneyError    = ErrorStructFunc(battery.ErrorCode_BuyGoodConsumMoneyError)
    Resp_BuyGoodUpdateUserDataError = ErrorStructFunc(battery.ErrorCode_BuyGoodUpdateUserDataError)

    Resp_IapGoodNotFound = ErrorStructFunc(battery.ErrorCode_IapGoodNotFound)

    Resp_ResOpError      = ErrorStructFunc(battery.ErrorCode_ResOpError)
    Resp_ResUnkownOpType = ErrorStructFunc(battery.ErrorCode_ResUnkownOpType)

    Resp_QueryUserMissionError                   = ErrorStructFunc(battery.ErrorCode_QueryUserMissionError)
    Resp_AddUserMissionError                     = ErrorStructFunc(battery.ErrorCode_AddUserMissionError)
    Resp_QueryMissionError                       = ErrorStructFunc(battery.ErrorCode_QueryMissionError)
    Resp_QueryUserSignInRecordError              = ErrorStructFunc(battery.ErrorCode_QueryUserSignInRecordError)
    Resp_AddUserSignInRecordError                = ErrorStructFunc(battery.ErrorCode_AddUserSignInRecordError)
    Resp_GetSignInItemsError                     = ErrorStructFunc(battery.ErrorCode_GetSignInItemsError)
    Resp_QuerySignInActivitysFromDBError         = ErrorStructFunc(battery.ErrorCode_QuerySignInActivitysFromDBError)
    Resp_QuerySignInActivitysFromCacheError      = ErrorStructFunc(battery.ErrorCode_QuerySignInActivitysFromCacheError)
    Resp_SignInActivitysStateError               = ErrorStructFunc(battery.ErrorCode_SignInActivitysStateError)
    Resp_QueryLottoTransactionError              = ErrorStructFunc(battery.ErrorCode_QueryLottoTransactionError)
    Resp_PushLottoTransactionStateError          = ErrorStructFunc(battery.ErrorCode_PushLottoTransactionStateError)
    Resp_QueryAfterGameQuotaId2StagesError       = ErrorStructFunc(battery.ErrorCode_QueryAfterGameQuotaId2StagesError)
    Resp_QuerySysLottoWeightError                = ErrorStructFunc(battery.ErrorCode_QuerySysLottoWeightError)
    Resp_QueryAfterGameLottoWeightError          = ErrorStructFunc(battery.ErrorCode_QueryAfterGameLottoWeightError)
    Resp_QueryAfterGameStageError                = ErrorStructFunc(battery.ErrorCode_QueryAfterGameStageError)
    Resp_UnkownLottoTypeError                    = ErrorStructFunc(battery.ErrorCode_UnkownLottoTypeError)
    Resp_GetSelectedSlotError                    = ErrorStructFunc(battery.ErrorCode_GetSelectedSlotError)
    Resp_CheckAfterGameLottoTransactionsError    = ErrorStructFunc(battery.ErrorCode_CheckAfterGameLottoTransactionsError)
    Resp_AfterGameNotEnoughDeleteSlotChanceError = ErrorStructFunc(battery.ErrorCode_AfterGameNotEnoughDeleteSlotChanceError)
    Resp_QuerySysLottoInfoError                  = ErrorStructFunc(battery.ErrorCode_QuerySysLottoInfoError)
    Resp_ShareActivityFromCacheError             = ErrorStructFunc(battery.ErrorCode_QueryShareActivityError)

    Resp_QueryPropsFromDBError              = ErrorStructFunc(battery.ErrorCode_QueryPropsFromDBError)
    Resp_QueryPropsFromCacheError           = ErrorStructFunc(battery.ErrorCode_QueryPropsFromCacheError)
    Resp_PropTypeInvalidError               = ErrorStructFunc(battery.ErrorCode_PropTypeInvalidError)
    Resp_QueryNewAccountPropsFromCacheError = ErrorStructFunc(battery.ErrorCode_QueryNewAccountPropsFromCacheError)

    Resp_QueryUserMaxCheckPointError        = ErrorStructFunc(battery.ErrorCode_QueryUserMaxCheckPointError)
    Resp_QueryUserCheckPointError           = ErrorStructFunc(battery.ErrorCode_QueryUserCheckPointError)
    Resp_QueryUserCheckPointFriendRankError = ErrorStructFunc(battery.ErrorCode_QueryUserCheckPointFriendRankError)
    Resp_QueryUserCheckPointGlobalRankError = ErrorStructFunc(battery.ErrorCode_QueryUserCheckPointGlobalRankError)
    Resp_CommitUserCheckPointDetailError    = ErrorStructFunc(battery.ErrorCode_CommitUserCheckPointDetailError)
    Resp_QueryUserFriendsUidError           = ErrorStructFunc(battery.ErrorCode_QueryUserFriendsUidError)

    Resp_QueryRuneConfigsFromDBError = ErrorStructFunc(battery.ErrorCode_QueryRuneConfigsFromDBError)
    Resp_RuneConfigsFromCacheError   = ErrorStructFunc(battery.ErrorCode_QueryRuneConfigsFromCacheError)

    Resp_QueryMailConfigsFromDBError    = ErrorStructFunc(battery.ErrorCode_QueryMailConfigsFromDBError)
    Resp_QueryMailConfigsFromCacheError = ErrorStructFunc(battery.ErrorCode_QueryMailConfigsFromCacheError)

    Resp_QueryJigsawConfigsFromDBError = ErrorStructFunc(battery.ErrorCode_QueryJigsawConfigsFromDBError)
    Resp_JigsawConfigsFromCacheError   = ErrorStructFunc(battery.ErrorCode_QueryJigsawConfigsFromCacheError)

    Resp_QueryBGPropConfigsFromDBError = ErrorStructFunc(battery.ErrorCode_QueryBGPropConfigsFromDBError)
    Resp_BGPropConfigsFromCacheError   = ErrorStructFunc(battery.ErrorCode_QueryBGPropConfigsFromCacheError)
    Resp_DecreaseConsumableError       = ErrorStructFunc(battery.ErrorCode_DecreaseConsumableError)

    Resp_QueryRoleLevelBonusFromDBError    = ErrorStructFunc(battery.ErrorCode_QueryRoleLevelBonusFromDBError)
    Resp_QueryRoleInfoConfigFromDBError    = ErrorStructFunc(battery.ErrorCode_QueryRoleInfoConfigFromDBError)
    Resp_QueryRoleInfoConfigFromCacheError = ErrorStructFunc(battery.ErrorCode_QueryRoleInfoConfigFromCacheError)
    Resp_QueryUserRoleInfoError            = ErrorStructFunc(battery.ErrorCode_QueryUserRoleInfoError)
    Resp_UpgradeUserRoleError              = ErrorStructFunc(battery.ErrorCode_UpgradeUserRoleError)
)

var (
    ErrOK                  error = nil // 没有错误
    ErrUnknowError               = errors.New(Resp_UnknowError.GetDesc())
    ErrBadInputData              = errors.New(Resp_BadInputData.GetDesc())
    ErrServerError               = errors.New(Resp_ServerError.GetDesc())
    ErrDBError                   = errors.New(Resp_DBError.GetDesc())
    ErrServiceNotAvaliable       = errors.New(Resp_ServiceNotAvaliable.GetDesc())
    ErrClientNotSupport          = errors.New(Resp_ClientNotSupport.GetDesc())
    ErrSendToTransaction         = errors.New(Resp_SendToTransactionError.GetDesc())
    ErrTimeLimitError            = errors.New(Resp_TimeLimitError.GetDesc())

    ErrGetAccountByUidError = errors.New(Resp_GetAccountByUidError.GetDesc())
    ErrUpdateAccountError   = errors.New(Resp_UpdateAccountError.GetDesc())
    ErrGetUidError          = errors.New(Resp_GetUidError.GetDesc())

    ErrNotEnoughCurrency = errors.New(Resp_NotEnoughCurrency.GetDesc())

    ErrUpdateStaminaError = errors.New(Resp_UpdateStaminaError.GetDesc())
    ErrNotEnoughStamina   = errors.New(Resp_NotEnoughStamina.GetDesc())
    ErrQueryStaminaError  = errors.New(Resp_QueryStaminaError.GetDesc())

    ErrNotEnoughDiamond = errors.New(Resp_NotEnoughDiamond.GetDesc())

    ErrAddNewGameError        = errors.New(Resp_AddNewGameError.GetDesc())
    ErrGameNotExistError      = errors.New(Resp_GameNotExistError.GetDesc())
    ErrGameResultInvalidError = errors.New(Resp_GameResultInvalidError.GetDesc())
    ErrUpdateGameError        = errors.New(Resp_UpdateGameError.GetDesc())

    ErrQueryGoodsError            = errors.New(Resp_QueryGoodsError.GetDesc())
    ErrBuyGoodExceedLimit         = errors.New(Resp_BuyGoodExceedLimit.GetDesc())
    ErrBuyGoodInvalidGame         = errors.New(Resp_BuyGoodInvalidGame.GetDesc())
    ErrBuyGoodsError              = errors.New(Resp_BuyGoodsError.GetDesc())
    ErrBuyGoodOverAmountPerGame   = errors.New(Resp_BuyGoodOverAmountPerGame.GetDesc())
    ErrBuyGoodOverAmountPerUser   = errors.New(Resp_BuyGoodOverAmountPerUser.GetDesc())
    ErrBuyGoodNewReceiptError     = errors.New(Resp_BuyGoodNewReceiptError.GetDesc())
    ErrBuyGoodConsumMoneyError    = errors.New(Resp_BuyGoodConsumMoneyError.GetDesc())
    ErrBuyGoodUpdateUserDataError = errors.New(Resp_BuyGoodUpdateUserDataError.GetDesc())

    ErrIapGoodNotFound = errors.New(Resp_IapGoodNotFound.GetDesc())

    ErrResOpError      = errors.New(Resp_ResOpError.GetDesc())
    ErrResUnkownOpType = errors.New(Resp_ResUnkownOpType.GetDesc())

    ErrQueryUserMissionError              = errors.New(Resp_QueryUserMissionError.GetDesc())
    ErrAddUserMissionError                = errors.New(Resp_AddUserMissionError.GetDesc())
    ErrQueryMissionError                  = errors.New(Resp_QueryMissionError.GetDesc())
    ErrQueryUserSignInRecordError         = errors.New(Resp_QueryUserSignInRecordError.GetDesc())
    ErrAddUserSignInRecordError           = errors.New(Resp_AddUserSignInRecordError.GetDesc())
    ErrGetSignInItemsError                = errors.New(Resp_GetSignInItemsError.GetDesc())
    ErrQuerySignInActivitysFromDBError    = errors.New(Resp_QuerySignInActivitysFromDBError.GetDesc())
    ErrQuerySignInActivitysFromCacheError = errors.New(Resp_QuerySignInActivitysFromCacheError.GetDesc())
    ErrSignInActivitysStateError          = errors.New(Resp_SignInActivitysStateError.GetDesc())
    ErrQueryShareActivityFromCacheError   = errors.New(Resp_ShareActivityFromCacheError.GetDesc())

    ErrQueryLottoTransactionError              = errors.New(Resp_QueryLottoTransactionError.GetDesc())
    ErrPushLottoTransactionStateError          = errors.New(Resp_PushLottoTransactionStateError.GetDesc())
    ErrQueryAfterGameQuotaId2StagesError       = errors.New(Resp_QueryAfterGameQuotaId2StagesError.GetDesc())
    ErrQuerySysLottoWeightError                = errors.New(Resp_QuerySysLottoWeightError.GetDesc())
    ErrQueryAfterGameLottoWeightError          = errors.New(Resp_QueryAfterGameLottoWeightError.GetDesc())
    ErrQueryAfterGameStageError                = errors.New(Resp_QueryAfterGameStageError.GetDesc())
    ErrUnkownLottoTypeError                    = errors.New(Resp_UnkownLottoTypeError.GetDesc())
    ErrGetSelectedSlotError                    = errors.New(Resp_GetSelectedSlotError.GetDesc())
    ErrCheckAfterGameLottoTransactionsError    = errors.New(Resp_CheckAfterGameLottoTransactionsError.GetDesc())
    ErrAfterGameNotEnoughDeleteSlotChanceError = errors.New(Resp_AfterGameNotEnoughDeleteSlotChanceError.GetDesc())
    ErrQuerySysLottoInfoError                  = errors.New(Resp_QuerySysLottoInfoError.GetDesc())

    ErrQueryPropsFromDBError              = errors.New(Resp_QueryPropsFromDBError.GetDesc())
    ErrQueryPropsFromCacheError           = errors.New(Resp_QueryPropsFromCacheError.GetDesc())
    ErrPropTypeInvalidError               = errors.New(Resp_PropTypeInvalidError.GetDesc())
    ErrQueryNewAccountPropsFromCacheError = errors.New(Resp_QueryNewAccountPropsFromCacheError.GetDesc())

    ErrQueryUserMaxCheckPointError = errors.New(Resp_QueryUserMaxCheckPointError.GetDesc())
    ErrQueryUserFriendsUidError    = errors.New(Resp_QueryUserFriendsUidError.GetDesc())

    ErrRuneConfigsFromCacheError      = errors.New(Resp_RuneConfigsFromCacheError.GetDesc())
    ErrQueryMailConfigsFromCacheError = errors.New(Resp_QueryMailConfigsFromCacheError.GetDesc())

    ErrJigsawConfigsFromCacheError = errors.New(Resp_JigsawConfigsFromCacheError.GetDesc())

    ErrBGPropConfigsFromCacheError = errors.New(Resp_BGPropConfigsFromCacheError.GetDesc())
    ErrDecreaseConsumableError     = errors.New(Resp_DecreaseConsumableError.GetDesc())

    ErrQueryRoleLevelBonusFromDBError    = errors.New(Resp_QueryRoleLevelBonusFromDBError.GetDesc())
    ErrQueryRoleInfoConfigFromDBError    = errors.New(Resp_QueryRoleInfoConfigFromDBError.GetDesc())
    ErrQueryRoleInfoConfigFromCacheError = errors.New(Resp_QueryRoleInfoConfigFromCacheError.GetDesc())
    ErrQueryUserRoleInfoError            = errors.New(Resp_QueryUserRoleInfoError.GetDesc())
    ErrUpgradeUserRoleError              = errors.New(Resp_UpgradeUserRoleError.GetDesc())
)

const (
    RecNotFoundMsg = "not found"
)

var (
    DBErrOK                     error = ErrOK                                                // 操作完成
    DBErrFailedDueToServerError       = errors.New("Operation failed due to db error")       // 服务器原因导致操作失败
    DBErrFailedDueToClientError       = errors.New("Operation failed due to bad input data") // 调用端原因导致操作失败，比如输入参数错误，查询条件不对，等等
    ErrNoRecord                       = errors.New(RecNotFoundMsg)
    ErrNotFound                       = ErrNoRecord
    Err_nil                           = errors.New("Data == nil")
    ErrDataInvalid                    = errors.New("Data Invalid")
    ErrUidNotFound                    = errors.New("Uid no found")
    ErrNotSupport                     = errors.New("Not support")

    ErrReadIniFile = errors.New("Read ini file failed")
)

var (
    ErrNoEnoughMoney = errors.New("No enough money")
)

// The serializable Error structure.
type Error struct {
    XMLName xml.Name `json:"-" xml:"error"`
    Code    int      `json:"code" xml:"code,attr"`
    Message string   `json:"message" xml:"message"`
}

func (e *Error) Error() string {
    return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// NewError creates an error instance with the specified code and message.
func NewError(code int, msg string) *Error {
    return &Error{
        Code:    code,
        Message: msg,
    }
}
func DBError(err error) error {
    if err != nil && err.Error() == RecNotFoundMsg {
        err = ErrNotFound
    }
    return err
}
