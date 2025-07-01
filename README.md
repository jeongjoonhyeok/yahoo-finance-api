# yahoo-finance-api

## Motivation

- I used to write Python programs and use Yahoo Finance data. The [yfinance](https://github.com/ranaroussi/yfinance) library is an awesome library which I enjoyed a lot.
- Could not find similar packages in Go.
- Learn Go
- Able to use this package on my other Go projects


```
func main() {
	// 1. 단일 심볼 시가총액 조회 테스트
	fmt.Println("=== 단일 심볼 시가총액 조회 ===")
	info, err := yahoofinanceapi.GetMarketCap("AAPL")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Symbol: %s\n", info.Symbol)
	fmt.Printf("Company: %s\n", info.CompanyName)
	fmt.Printf("MarketCap: %s\n", info.MarketCapFormatted)
	fmt.Printf("Price: $%.2f\n", info.CurrentPrice)

	// 2. 여러 심볼 시가총액 일괄 조회 테스트
	fmt.Println("\n=== 여러 심볼 시가총액 일괄 조회 ===")
	symbols := []string{"AAPL", "GOOGL", "MSFT"}
	infos, err := yahoofinanceapi.GetMultipleMarketCaps(symbols)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, info := range infos {
		fmt.Printf("%s: %s (Price: $%.2f)\n", info.Symbol, info.MarketCapFormatted, info.CurrentPrice)
	}

	// 3. Ticker를 통한 조회 (lazy initialization 테스트)
	fmt.Println("\n=== Ticker를 통한 조회 ===")
	ticker := yahoofinanceapi.NewTicker("TSLA")
	marketCapInfo, err := ticker.MarketCap()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Ticker %s: %s\n", marketCapInfo.Symbol, marketCapInfo.MarketCapFormatted)
}
```