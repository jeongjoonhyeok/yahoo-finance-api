package yahoofinanceapi

import (
	"fmt"
	"log"
)

func TestVolume() {
	fmt.Println("=== Volume 테스트 프로그램 ===")

	// AAPL 티커 생성
	ticker := NewTicker("AAPL")

	// Premarket 데이터 포함하여 조회
	query := HistoryQuery{
		Range:    "1d",
		Interval: "1m",
		Prepost:  true,
	}

	fmt.Printf("Symbol: AAPL\n")
	fmt.Printf("Range: %s\n", query.Range)
	fmt.Printf("Interval: %s\n", query.Interval)
	fmt.Printf("Prepost: %t\n\n", query.Prepost)

	// 데이터 조회
	history, err := ticker.History(query)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("총 데이터 포인트: %d\n\n", len(history))

	// 처음 10개 데이터 포인트 확인
	count := 0
	for timeKey, data := range history {
		fmt.Printf("%s: Open=%.2f, High=%.2f, Low=%.2f, Close=%.2f, Volume=%d\n",
			timeKey, data.Open, data.High, data.Low, data.Close, data.Volume)
		count++
		if count >= 10 {
			break
		}
	}

	// Volume이 0이 아닌 데이터 개수 확인
	nonZeroVolume := 0
	zeroVolume := 0
	for _, data := range history {
		if data.Volume > 0 {
			nonZeroVolume++
		} else {
			zeroVolume++
		}
	}

	fmt.Printf("\n=== Volume 통계 ===\n")
	fmt.Printf("Volume > 0인 데이터: %d\n", nonZeroVolume)
	fmt.Printf("Volume = 0인 데이터: %d\n", zeroVolume)
	fmt.Printf("전체 데이터: %d\n", len(history))

	// Volume이 0이 아닌 첫 5개 데이터 표시
	fmt.Printf("\n=== Volume > 0인 데이터 샘플 ===\n")
	count = 0
	for timeKey, data := range history {
		if data.Volume > 0 && count < 5 {
			fmt.Printf("%s: Volume=%d, Price=%.2f\n", timeKey, data.Volume, data.Close)
			count++
		}
	}
}
