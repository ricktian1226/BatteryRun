package xyfmt

import (
	"fmt"
	battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	xyutil "guanghuan.com/xiaoyao/superbman_server/common/util"
	xydb "guanghuan.com/xiaoyao/superbman_server/db/v2"
	//	"time"
)

// log
func FormatGift(gift battery.Gift) (str string) {
	var (
		src = gift.GetSource()
		//		gift_sid = gift.GetGiftId()
		from_sid = gift.GetFromId()
		to_sid   = gift.GetToId()
	)
	//	gift_sid, _ := xydb.SharedInstance().GetSidByGid(gift.GetGiftId(), src)
	from_sid, _ = xydb.SharedInstance().GetSidByGid(from_sid, src)
	to_sid, _ = xydb.SharedInstance().GetSidByGid(to_sid, src)

	str = fmt.Sprintf(`Gift %s (%s)
		From    : %s (%s)
		To      : %s (%s)
		IDSource: %s
		OpType  : %s
		Resource: %s
		Amount  : %d
		Rrc_id  : %s
		Note    : %s
		Confirm : %v
		Create by: %s
		Create on: %d (%s)
		Expire on: %d (%s)
		Update on: %d (%s) `,
		gift.GetGiftId(), src.String(),
		gift.GetFromId(), from_sid,
		gift.GetToId(), to_sid,
		gift.GetSource().String(),
		gift.GetOpType().String(),
		gift.GetResource().String(),
		gift.GetAmount(),
		gift.GetResourceId(),
		gift.GetNote(),
		gift.GetIsConfirmed(),
		gift.GetCreatorId(),
		gift.GetCreateDate(), xyutil.ToStrTime(gift.GetCreateDate()),
		gift.GetExpiredDate(), xyutil.ToStrTime(gift.GetExpiredDate()),
		gift.GetLastUpdateTime(), xyutil.ToStrTime(gift.GetLastUpdateTime()))

	return
}

func FormatAccount(acc battery.UserData) (str string) {
	str = fmt.Sprintf(`Uid: %s
	name		: %s
	DeviceId	: %s
	Stamina		: %d
	Diamond		: %d
	Gifts		: %d
	---
	Total.Score	: %d
	Total.Gold	: %d
	Total.Gscore	: %d
	Total.Distance	: %d
	Total.TCharge	: %d
	Total.CCharge	: %d
	---
	Best.Score	: %d
	Best.Gold	: %d
	Best.Gscore	: %d
	Best.Distance	: %d
	Best.TCharge	: %d
	Best.CCharge	: %d`,
		acc.GetUid(),
		acc.GetName(),
		acc.GetDeviceId(),
		acc.GetWallet().GetStamina(),
		acc.GetWallet().GetDiamond(),
		acc.GetTotalGiftCount(),

		acc.GetTotal().GetScore(),
		acc.GetTotal().GetGold(),
		acc.GetTotal().GetGoldScore(),
		acc.GetTotal().GetDistance(),
		acc.GetTotal().GetTotalCharge(),
		acc.GetTotal().GetComboCharge(),

		acc.GetBest().GetScore(),
		acc.GetBest().GetGold(),
		acc.GetBest().GetGoldScore(),
		acc.GetBest().GetDistance(),
		acc.GetBest().GetTotalCharge(),
		acc.GetBest().GetComboCharge())
	return
}
