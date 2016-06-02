package main

import (
	"flag"
	"log"
)

const charge_factor = 1.5
const duration_factor = 0.8

//const gold_score_factor = 1.6
const gold_score_factor = 1.1
const distance_score_factor = 1.1

// 数据合理性检查
func ValidateGameResult(start_time int64, duration int64, distance int64, total_charge int64, combo_charge int64, gold int64, gold_score int64, score int64, kills int64) (isValid bool) {
	//var (
	//	distance     = result.GetDistance()
	//	kills        = GetTotalKilled(result)
	//	total_charge = result.GetTotalCharge()
	//	combo_charge = result.GetComboCharge()
	//	gold         = result.GetGold()
	//	gold_score   = result.GetGoldScore()
	//	score        = result.GetScore()
	//)
	var (
		invalid int
	)
	if distance > MaxDistance() {
		//		log.Printf("Invalid Distance: %d > %d", distance, MaxDistance())
		invalid++
	}

	if duration < MinDurationSec(distance) {
		//		log.Printf("Invalid Duration: %d < %d", duration, MinDurationSec(distance))
		invalid++
	}
	// combo的初值为1
	if combo_charge > total_charge+1 {
		//		log.Printf("Invalid ComboCharge: %d > %d", combo_charge, total_charge)
		invalid++
	}
	if int64(total_charge) > MaxCharge(distance) {
		//		log.Printf("Invalid TotalCharge: %d > %d", total_charge, MaxCharge(distance))
		invalid++
	}
	if gold > MaxGold(distance) {
		//		log.Printf("Invalid Gold: %d > %d", gold, MaxGold(distance))
		invalid++
	}
	if gold_score > MaxGoldScore(distance) {
		//		log.Printf("Invalid GoldScore: %d > %d", gold_score, MaxGoldScore(distance))
		invalid++
	}
	if score > MaxTotalScore(distance, kills) {
		//		log.Printf("Invalid Score: %d > %d", score, MaxTotalScore(distance, kills))
		invalid++
	}

	isValid = (invalid <= 0)

	//if !isValid {
	//	log.Printf("result: %s", result.String())
	//}

	return
}
func MaxCharge(distance int64) int64 {
	return int64(float64(distance) * 0.032 * 1.5)
}

var distance_multipler []int64 = []int64{
	0:  0,
	1:  28,
	2:  56,
	3:  84,
	4:  112,
	5:  140,
	6:  168,
	7:  196,
	8:  224,
	9:  252,
	10: 280,
	11: 308,
}

func MaxDistance() int64 {
	return int64(float64(10000) * 1.1)
}

func MaxSpeed() int64 {
	return 20
}

func MinDurationSec(distance int64) int64 {
	return int64((float64(distance) / float64(MaxSpeed())) * duration_factor)
}

func MaxGold(distance int64) int64 {
	return int64(float64(distance)*1.5) + 200
}

func DistanceMultipler(distance int64) int64 {
	idx := int64(distance/1000 + 1)
	if idx > int64(len(distance_multipler)) {
		idx = int64(len(distance_multipler)) - 1
		//		log.Printf("distance may be out of range")
	}
	return distance_multipler[idx]
}

func MaxGoldScore(distance int64) int64 {
	return int64(float64(distance*DistanceMultipler(distance)) * gold_score_factor)
}

func MaxDistanceScore(distance int64) int64 {
	return int64(float64(distance*DistanceMultipler(distance)) * distance_score_factor)
}

// 杀怪分数
func MonsterKillScore(kills int64) int64 {
	return 1000 * kills
}

func MaxTotalScore(distance int64, monster_killed int64) int64 {
	return MaxGoldScore(distance) + MaxDistanceScore(distance) + MonsterKillScore(monster_killed)
}

func main() {
	var (
		distance     = flag.Int64("distance", 0, "distance")
		kills        = flag.Int64("kills", 0, "kills")
		score        = flag.Int64("score", 0, "score")
		start_time   = flag.Int64("start_time", 0, "start_time")
		duration     = flag.Int64("duration", 0, "duration")
		total_charge = flag.Int64("total_charge", 0, "total_charge")
		combo_charge = flag.Int64("combo_charge", 0, "combo_charge")
		gold         = flag.Int64("gold", 0, "gold")
		gold_score   = flag.Int64("gold_score", 0, "gold_score")
	)

	flag.Parse()
	log.Printf(`Result:
	distance   = %d
	kills      = %d
	start_time = %d
	duration   = %d  (%d)
	gold       = %d  (%d)
	gold_score = %d  (%d)
	ttl_charge = %d  (%d)
	cmb_charge = %d  (%d)
	fnl_score  = %d  (%d)
	
	Validation result: %t
	`, *distance,
		*kills,
		*start_time,
		*duration, MinDurationSec(*distance),
		*gold, MaxGold(*distance),
		*gold_score, MaxGoldScore(*distance),
		*total_charge, MaxCharge(*distance),
		*combo_charge, MaxCharge(*distance)+1,
		*score, MaxTotalScore(*distance, *kills),
		ValidateGameResult(*start_time, *duration, *distance, *total_charge, *combo_charge, *gold, *gold_score, *score, *kills))
}
