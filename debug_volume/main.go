package main

import (
	"fmt"
	"log"

	yahoofinanceapi "github.com/oscarli916/yahoo-finance-api"
)

func main() {
	fmt.Println("=== Volume 디버깅 프로그램 ===")

	// AAPL 티커 생성
	ticker := yahoofinanceapi.NewTicker("AAPL")

	// 1분 간격, 하루 데이터, premarket 포함
	query := yahoofinanceapi.HistoryQuery{
		Range:    "1d",
		Interval: "1m",
		Prepost:  true,
	}

	fmt.Printf("Symbol: AAPL\n")
	fmt.Printf("Range: %s\n", query.Range)
	fmt.Printf("Interval: %s\n", query.Interval)
	fmt.Printf("Prepost: %t\n\n", query.Prepost)

	// 데이터 조회 (디버깅 정보 포함)
	history, err := ticker.History(query)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("총 데이터 포인트: %d\n\n", len(history))

	// Volume 통계 분석
	nonZeroVolume := 0
	zeroVolume := 0
	totalVolume := int64(0)

	for _, data := range history {
		totalVolume += data.Volume
		if data.Volume > 0 {
			nonZeroVolume++
		} else {
			zeroVolume++
		}
	}

	fmt.Printf("=== Volume 통계 ===\n")
	fmt.Printf("Volume > 0인 데이터: %d\n", nonZeroVolume)
	fmt.Printf("Volume = 0인 데이터: %d\n", zeroVolume)
	fmt.Printf("전체 데이터: %d\n", len(history))
	fmt.Printf("총 Volume: %d\n", totalVolume)
	fmt.Printf("Volume=0 비율: %.1f%%\n\n", float64(zeroVolume)/float64(len(history))*100)

	// Volume이 0이 아닌 데이터 샘플
	fmt.Printf("=== Volume > 0인 데이터 샘플 (최대 10개) ===\n")
	count := 0
	for timeKey, data := range history {
		if data.Volume > 0 && count < 10 {
			fmt.Printf("%s: Volume=%d, Close=%.2f\n", timeKey, data.Volume, data.Close)
			count++
		}
	}

	if count == 0 {
		fmt.Println("Volume > 0인 데이터가 없습니다!")
		fmt.Println("\n=== 전체 데이터 샘플 (Volume=0 포함) ===")
		count = 0
		for timeKey, data := range history {
			fmt.Printf("%s: Volume=%d, Close=%.2f\n", timeKey, data.Volume, data.Close)
			count++
			if count >= 10 {
				break
			}
		}
	}

	fmt.Println("\n=== 분석 결과 ===")
	if zeroVolume == len(history) {
		fmt.Println("❌ 모든 데이터의 Volume이 0입니다!")
		fmt.Println("가능한 원인:")
		fmt.Println("1. Yahoo Finance API의 premarket volume 데이터 제한")
		fmt.Println("2. 해당 종목의 premarket 거래량 부족")
		fmt.Println("3. API 응답에서 volume 데이터가 누락")
	} else if zeroVolume > len(history)/2 {
		fmt.Printf("⚠️  Volume=0인 데이터가 많습니다 (%.1f%%)\n", float64(zeroVolume)/float64(len(history))*100)
		fmt.Println("Premarket 시간대에는 거래량이 적은 것이 정상입니다.")
	} else {
		fmt.Printf("✅ Volume 데이터가 정상적으로 조회되었습니다.\n")
	}
}
