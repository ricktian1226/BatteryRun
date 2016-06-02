package batteryapi

import (
	"testing"
)

func TestSendIapStatistic(t *testing.T) {
	NewXYAPI().SendIapStatistic("1435390486436053904081", "Buy_Goods_0704", "10000")
}
