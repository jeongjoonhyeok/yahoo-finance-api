package yahoofinanceapi

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
)

/*
 * Stock Quote API Module
 *
 * 이 파일은 Yahoo Finance API를 통해 일반 주식의 실시간 quote 정보를 가져오는 기능을 제공합니다.
 * 옵션이 아닌 일반 주식의 시가총액, 주가, 거래량 등의 정보를 조회합니다.
 *
 * 주요 기능:
 * - 개별 주식의 실시간 quote 정보 조회
 * - 여러 주식의 quote 정보 일괄 조회
 * - 시가총액 및 주요 재무 지표 제공
 */

// QuoteResponse는 Yahoo Finance quote API의 응답을 파싱하기 위한 구조체입니다.
type QuoteResponse struct {
	QuoteResponse QuoteResponseData `json:"quoteResponse"`
}

// QuoteResponseData는 quote 응답의 메인 데이터를 담는 구조체입니다.
type QuoteResponseData struct {
	Result []StockQuote `json:"result"`
	Error  interface{}  `json:"error"`
}

// StockQuote는 일반 주식의 실시간 quote 정보를 담는 구조체입니다.
type StockQuote struct {
	Language                          string  `json:"language"`
	Region                            string  `json:"region"`
	QuoteType                         string  `json:"quoteType"`
	TypeDisp                          string  `json:"typeDisp"`
	QuoteSourceName                   string  `json:"quoteSourceName"`
	Triggerable                       bool    `json:"triggerable"`
	CustomPriceAlertConfidence        string  `json:"customPriceAlertConfidence"`
	Currency                          string  `json:"currency"`
	MarketState                       string  `json:"marketState"`
	RegularMarketChangePercent        float64 `json:"regularMarketChangePercent"`
	RegularMarketPrice                float64 `json:"regularMarketPrice"`
	Exchange                          string  `json:"exchange"`
	ShortName                         string  `json:"shortName"`
	LongName                          string  `json:"longName"`
	MessageBoardId                    string  `json:"messageBoardId"`
	ExchangeTimezoneName              string  `json:"exchangeTimezoneName"`
	ExchangeTimezoneShortName         string  `json:"exchangeTimezoneShortName"`
	GmtOffSetMilliseconds             int64   `json:"gmtOffSetMilliseconds"`
	Market                            string  `json:"market"`
	EsgPopulated                      bool    `json:"esgPopulated"`
	HasPrePostMarketData              bool    `json:"hasPrePostMarketData"`
	FirstTradeDateMilliseconds        int64   `json:"firstTradeDateMilliseconds"`
	PriceHint                         int     `json:"priceHint"`
	PostMarketChangePercent           float64 `json:"postMarketChangePercent"`
	PostMarketTime                    int64   `json:"postMarketTime"`
	PostMarketPrice                   float64 `json:"postMarketPrice"`
	PostMarketChange                  float64 `json:"postMarketChange"`
	RegularMarketChange               float64 `json:"regularMarketChange"`
	RegularMarketTime                 int64   `json:"regularMarketTime"`
	RegularMarketDayHigh              float64 `json:"regularMarketDayHigh"`
	RegularMarketDayRange             string  `json:"regularMarketDayRange"`
	RegularMarketDayLow               float64 `json:"regularMarketDayLow"`
	RegularMarketVolume               int64   `json:"regularMarketVolume"`
	RegularMarketPreviousClose        float64 `json:"regularMarketPreviousClose"`
	Bid                               float64 `json:"bid"`
	Ask                               float64 `json:"ask"`
	BidSize                           int     `json:"bidSize"`
	AskSize                           int     `json:"askSize"`
	FullExchangeName                  string  `json:"fullExchangeName"`
	FinancialCurrency                 string  `json:"financialCurrency"`
	RegularMarketOpen                 float64 `json:"regularMarketOpen"`
	AverageDailyVolume3Month          int64   `json:"averageDailyVolume3Month"`
	AverageDailyVolume10Day           int64   `json:"averageDailyVolume10Day"`
	FiftyTwoWeekLowChange             float64 `json:"fiftyTwoWeekLowChange"`
	FiftyTwoWeekLowChangePercent      float64 `json:"fiftyTwoWeekLowChangePercent"`
	FiftyTwoWeekRange                 string  `json:"fiftyTwoWeekRange"`
	FiftyTwoWeekHighChange            float64 `json:"fiftyTwoWeekHighChange"`
	FiftyTwoWeekHighChangePercent     float64 `json:"fiftyTwoWeekHighChangePercent"`
	FiftyTwoWeekLow                   float64 `json:"fiftyTwoWeekLow"`
	FiftyTwoWeekHigh                  float64 `json:"fiftyTwoWeekHigh"`
	FiftyTwoWeekChangePercent         float64 `json:"fiftyTwoWeekChangePercent"`
	DividendDate                      int64   `json:"dividendDate"`
	EarningsTimestamp                 int64   `json:"earningsTimestamp"`
	EarningsTimestampStart            int64   `json:"earningsTimestampStart"`
	EarningsTimestampEnd              int64   `json:"earningsTimestampEnd"`
	EarningsCallTimestampStart        int64   `json:"earningsCallTimestampStart"`
	EarningsCallTimestampEnd          int64   `json:"earningsCallTimestampEnd"`
	IsEarningsDateEstimate            bool    `json:"isEarningsDateEstimate"`
	TrailingAnnualDividendRate        float64 `json:"trailingAnnualDividendRate"`
	TrailingPE                        float64 `json:"trailingPE"`
	DividendRate                      float64 `json:"dividendRate"`
	TrailingAnnualDividendYield       float64 `json:"trailingAnnualDividendYield"`
	DividendYield                     float64 `json:"dividendYield"`
	EpsTrailingTwelveMonths           float64 `json:"epsTrailingTwelveMonths"`
	EpsForward                        float64 `json:"epsForward"`
	EpsCurrentYear                    float64 `json:"epsCurrentYear"`
	PriceEpsCurrentYear               float64 `json:"priceEpsCurrentYear"`
	SharesOutstanding                 int64   `json:"sharesOutstanding"`
	BookValue                         float64 `json:"bookValue"`
	FiftyDayAverage                   float64 `json:"fiftyDayAverage"`
	FiftyDayAverageChange             float64 `json:"fiftyDayAverageChange"`
	FiftyDayAverageChangePercent      float64 `json:"fiftyDayAverageChangePercent"`
	TwoHundredDayAverage              float64 `json:"twoHundredDayAverage"`
	TwoHundredDayAverageChange        float64 `json:"twoHundredDayAverageChange"`
	TwoHundredDayAverageChangePercent float64 `json:"twoHundredDayAverageChangePercent"`
	MarketCap                         int64   `json:"marketCap"`
	ForwardPE                         float64 `json:"forwardPE"`
	PriceToBook                       float64 `json:"priceToBook"`
	SourceInterval                    int     `json:"sourceInterval"`
	ExchangeDataDelayedBy             int     `json:"exchangeDataDelayedBy"`
	AverageAnalystRating              string  `json:"averageAnalystRating"`
	Tradeable                         bool    `json:"tradeable"`
	CryptoTradeable                   bool    `json:"cryptoTradeable"`
	DisplayName                       string  `json:"displayName"`
	Symbol                            string  `json:"symbol"`
}

// MarketCapInfo는 사용자에게 제공되는 정제된 시가총액 정보 구조체입니다.
type MarketCapInfo struct {
	Symbol             string  `json:"symbol"`
	CompanyName        string  `json:"companyName"`
	MarketCap          int64   `json:"marketCap"`
	MarketCapFormatted string  `json:"marketCapFormatted"`
	SharesOutstanding  int64   `json:"sharesOutstanding"`
	CurrentPrice       float64 `json:"currentPrice"`
	Currency           string  `json:"currency"`
	Exchange           string  `json:"exchange"`
	MarketState        string  `json:"marketState"`
	PriceChange        float64 `json:"priceChange"`
	PriceChangePercent float64 `json:"priceChangePercent"`
}

// Quote는 주식 quote 관련 API 호출을 담당하는 구조체입니다.
type Quote struct {
	client *Client
}

// NewQuote는 새로운 Quote 인스턴스를 생성합니다.
func NewQuote() *Quote {
	return &Quote{client: GetClient()}
}

// GetQuote는 단일 심볼의 실시간 quote 정보를 조회합니다.
//
// 매개변수:
// - symbol: 조회할 주식의 심볼 (예: "AAPL", "TSLA")
//
// 반환값:
// - StockQuote: 주식의 실시간 quote 정보가 담긴 구조체
// - error: 조회 중 발생한 오류
func (q *Quote) GetQuote(symbol string) (StockQuote, error) {
	endpoint := fmt.Sprintf("%s/v7/finance/quote", BASE_URL)
	params := url.Values{}
	params.Add("symbols", symbol)

	resp, err := q.client.Get(endpoint, params)
	if err != nil {
		slog.Error("Failed to get quote data", "symbol", symbol, "err", err)
		return StockQuote{}, err
	}
	defer resp.Body.Close()

	var quoteResponse QuoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&quoteResponse); err != nil {
		slog.Error("Failed to decode quote JSON response", "err", err)
		return StockQuote{}, err
	}

	if len(quoteResponse.QuoteResponse.Result) == 0 {
		return StockQuote{}, fmt.Errorf("no data found for symbol: %s", symbol)
	}

	return quoteResponse.QuoteResponse.Result[0], nil
}

// GetMultipleQuotes는 여러 심볼의 실시간 quote 정보를 일괄 조회합니다.
//
// 매개변수:
// - symbols: 조회할 주식 심볼들의 슬라이스 (예: []string{"AAPL", "GOOGL", "MSFT"})
//
// 반환값:
// - []StockQuote: 각 심볼의 quote 정보가 담긴 구조체 슬라이스
// - error: 조회 중 발생한 오류
func (q *Quote) GetMultipleQuotes(symbols []string) ([]StockQuote, error) {
	if len(symbols) == 0 {
		return []StockQuote{}, fmt.Errorf("no symbols provided")
	}

	endpoint := fmt.Sprintf("%s/v7/finance/quote", BASE_URL)
	params := url.Values{}

	// 여러 심볼을 콤마로 구분하여 전달
	symbolsStr := strings.Join(symbols, ",")
	params.Add("symbols", symbolsStr)

	resp, err := q.client.Get(endpoint, params)
	if err != nil {
		slog.Error("Failed to get multiple quote data", "symbols", symbols, "err", err)
		return []StockQuote{}, err
	}
	defer resp.Body.Close()

	var quoteResponse QuoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&quoteResponse); err != nil {
		slog.Error("Failed to decode multiple quote JSON response", "err", err)
		return []StockQuote{}, err
	}

	return quoteResponse.QuoteResponse.Result, nil
}

// GetMarketCap은 단일 심볼의 시가총액 정보를 조회합니다.
//
// 매개변수:
// - symbol: 조회할 주식의 심볼 (예: "AAPL", "TSLA")
//
// 반환값:
// - MarketCapInfo: 시가총액 정보가 담긴 구조체
// - error: 조회 중 발생한 오류
func (q *Quote) GetMarketCap(symbol string) (MarketCapInfo, error) {
	quote, err := q.GetQuote(symbol)
	if err != nil {
		return MarketCapInfo{}, err
	}

	return q.transformQuoteToMarketCapInfo(quote), nil
}

// GetMultipleMarketCaps는 여러 심볼의 시가총액 정보를 일괄 조회합니다.
//
// 매개변수:
// - symbols: 조회할 주식 심볼들의 슬라이스 (예: []string{"AAPL", "GOOGL", "MSFT"})
//
// 반환값:
// - []MarketCapInfo: 각 심볼의 시가총액 정보가 담긴 구조체 슬라이스
// - error: 조회 중 발생한 오류
func (q *Quote) GetMultipleMarketCaps(symbols []string) ([]MarketCapInfo, error) {
	quotes, err := q.GetMultipleQuotes(symbols)
	if err != nil {
		return []MarketCapInfo{}, err
	}

	var results []MarketCapInfo
	for _, quote := range quotes {
		results = append(results, q.transformQuoteToMarketCapInfo(quote))
	}

	return results, nil
}

// transformQuoteToMarketCapInfo는 StockQuote를 MarketCapInfo로 변환합니다.
//
// 매개변수:
// - quote: Yahoo Finance API에서 받은 StockQuote 데이터
//
// 반환값:
// - MarketCapInfo: 정제된 시가총액 정보
func (q *Quote) transformQuoteToMarketCapInfo(quote StockQuote) MarketCapInfo {
	companyName := quote.LongName
	if companyName == "" {
		companyName = quote.ShortName
	}

	return MarketCapInfo{
		Symbol:             quote.Symbol,
		CompanyName:        companyName,
		MarketCap:          quote.MarketCap,
		MarketCapFormatted: q.formatMarketCap(quote.MarketCap),
		SharesOutstanding:  quote.SharesOutstanding,
		CurrentPrice:       quote.RegularMarketPrice,
		Currency:           quote.Currency,
		Exchange:           quote.Exchange,
		MarketState:        quote.MarketState,
		PriceChange:        quote.RegularMarketChange,
		PriceChangePercent: quote.RegularMarketChangePercent,
	}
}

// formatMarketCap은 시가총액을 읽기 쉬운 형태로 포맷팅합니다.
//
// 매개변수:
// - marketCap: 시가총액 값 (정수)
//
// 반환값:
// - string: 포맷팅된 시가총액 문자열 (예: "2.5T", "150.2B", "5.8M")
func (q *Quote) formatMarketCap(marketCap int64) string {
	if marketCap == 0 {
		return "N/A"
	}

	marketCapFloat := float64(marketCap)

	if marketCapFloat >= 1e12 {
		return fmt.Sprintf("%.1fT", marketCapFloat/1e12)
	} else if marketCapFloat >= 1e9 {
		return fmt.Sprintf("%.1fB", marketCapFloat/1e9)
	} else if marketCapFloat >= 1e6 {
		return fmt.Sprintf("%.1fM", marketCapFloat/1e6)
	} else if marketCapFloat >= 1e3 {
		return fmt.Sprintf("%.1fK", marketCapFloat/1e3)
	}

	return fmt.Sprintf("%.0f", marketCapFloat)
}

// GetMultipleMarketCapsGlobal은 패키지 레벨에서 여러 심볼의 시가총액을 일괄 조회하는 편의 함수입니다.
//
// 매개변수:
// - symbols: 조회할 주식 심볼들의 슬라이스 (예: []string{"AAPL", "GOOGL", "MSFT"})
//
// 반환값:
// - []MarketCapInfo: 각 심볼의 시가총액 정보가 담긴 구조체 슬라이스
// - error: 조회 중 발생한 오류
func GetMultipleMarketCaps(symbols []string) ([]MarketCapInfo, error) {
	quote := NewQuote()
	return quote.GetMultipleMarketCaps(symbols)
}

// GetMarketCapGlobal은 패키지 레벨에서 단일 심볼의 시가총액을 조회하는 편의 함수입니다.
//
// 매개변수:
// - symbol: 조회할 주식의 심볼 (예: "AAPL", "TSLA")
//
// 반환값:
// - MarketCapInfo: 시가총액 정보가 담긴 구조체
// - error: 조회 중 발생한 오류
func GetMarketCap(symbol string) (MarketCapInfo, error) {
	quote := NewQuote()
	return quote.GetMarketCap(symbol)
}
