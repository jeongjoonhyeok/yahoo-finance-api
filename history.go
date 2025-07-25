package yahoofinanceapi

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type YahooHistoryRespose struct {
	Chart YahooChart `json:"chart"`
}

type YahooChart struct {
	Result []YahooHistoryResult `json:"result"`
}

type YahooHistoryResult struct {
	Meta       YahooMeta      `json:"meta"`
	Timestamp  []int64        `json:"timestamp"`
	Indicators YahooIndicator `json:"indicators"`
}

type YahooMeta struct {
	Currency             string             `json:"currency"`
	Symbol               string             `json:"symbol"`
	ExchangeName         string             `json:"exchangeName"`
	FullExchangeName     string             `json:"fullExchangeName"`
	InstrumentType       string             `json:"instrumentType"`
	FirstTradeDate       int64              `json:"firstTradeDate"`
	RegularMarketTime    int64              `json:"regularMarketTime"`
	HasPrePostMarketData bool               `json:"hasPrePostMarketData"`
	GmtOffset            int                `json:"gmtoffset"`
	Timezone             string             `json:"timezone"`
	ExchangeTimezoneName string             `json:"exchangeTimezoneName"`
	RegularMarketPrice   float64            `json:"regularMarketPrice"`
	FiftyTwoWeekHigh     float64            `json:"fiftyTwoWeekHigh"`
	FiftyTwoWeekLow      float64            `json:"fiftyTwoWeekLow"`
	RegularMarketDayHigh float64            `json:"regularMarketDayHigh"`
	RegularMarketDayLow  float64            `json:"regularMarketDayLow"`
	RegularMarketVolume  int64              `json:"regularMarketVolume"`
	LongName             string             `json:"longName"`
	ShortName            string             `json:"shortName"`
	ChartPreviousClose   float64            `json:"chartPreviousClose"`
	PreviousClose        float64            `json:"previousClose"`
	Scale                int                `json:"scale"`
	PriceHint            int                `json:"priceHint"`
	CurrentTradingPeriod YahooTradingPeriod `json:"currentTradingPeriod"`
	TradingPeriods       json.RawMessage    `json:"tradingPeriods"`
	DataGranularity      string             `json:"dataGranularity"`
	Range                string             `json:"range"`
	ValidRanges          []string           `json:"validRanges"`
}

type YahooTradingPeriod struct {
	Timezone  string `json:"timezone"`
	End       int64  `json:"end"`
	Start     int64  `json:"start"`
	GmtOffset int    `json:"gmtoffset"`
}

// GetTradingPeriods는 TradingPeriods 데이터를 안전하게 파싱하여 반환합니다.
//
// 반환값:
// - [][]YahooTradingPeriod: 파싱된 거래 기간 데이터
// - error: 파싱 중 발생한 오류
func (ym *YahooMeta) GetTradingPeriods() ([][]YahooTradingPeriod, error) {
	if len(ym.TradingPeriods) == 0 {
		return nil, nil
	}

	// 먼저 배열의 배열로 파싱 시도
	var periodsArray [][]YahooTradingPeriod
	if err := json.Unmarshal(ym.TradingPeriods, &periodsArray); err == nil {
		return periodsArray, nil
	}

	// 실패하면 객체로 파싱 시도
	var periodsMap map[string]interface{}
	if err := json.Unmarshal(ym.TradingPeriods, &periodsMap); err != nil {
		// 둘 다 실패하면 빈 배열 반환
		log.Printf("Failed to parse tradingPeriods: %v", err)
		return [][]YahooTradingPeriod{}, nil
	}

	// 객체에서 유용한 데이터를 추출하여 배열로 변환
	result := [][]YahooTradingPeriod{}

	// 일반적으로 "regular" 키가 있는 경우가 많음
	if regular, ok := periodsMap["regular"]; ok {
		if regularArray, ok := regular.([]interface{}); ok {
			var periods []YahooTradingPeriod
			for _, item := range regularArray {
				if itemMap, ok := item.(map[string]interface{}); ok {
					period := YahooTradingPeriod{}
					if timezone, ok := itemMap["timezone"].(string); ok {
						period.Timezone = timezone
					}
					if end, ok := itemMap["end"].(float64); ok {
						period.End = int64(end)
					}
					if start, ok := itemMap["start"].(float64); ok {
						period.Start = int64(start)
					}
					if gmtOffset, ok := itemMap["gmtoffset"].(float64); ok {
						period.GmtOffset = int(gmtOffset)
					}
					periods = append(periods, period)
				}
			}
			if len(periods) > 0 {
				result = append(result, periods)
			}
		}
	}

	return result, nil
}

// DebugVolumeData는 volume 데이터를 디버깅하기 위한 헬퍼 함수입니다.
// Premarket/Postmarket 데이터의 volume 통계와 패턴을 분석합니다.
//
// 매개변수:
// - symbol: 디버깅할 심볼
// - data: Yahoo 응답 데이터
func (h *History) DebugVolumeData(symbol string, data YahooHistoryRespose) {
	if len(data.Chart.Result) == 0 {
		log.Printf("No chart results for %s", symbol)
		return
	}

	result := data.Chart.Result[0]
	log.Printf("\n=== Volume Analysis for %s (Prepost: %t) ===", symbol, h.query.Prepost)
	log.Printf("Total data points: %d", len(result.Timestamp))

	if len(result.Indicators.Quote) > 0 {
		quote := result.Indicators.Quote[0]

		// Volume 통계 계산
		totalVolume := int64(0)
		nonZeroVolume := 0
		premarketVolume := int64(0)
		regularVolume := int64(0)
		postmarketVolume := int64(0)

		premarketNonZero := 0
		regularNonZero := 0
		postmarketNonZero := 0

		for i, vol := range quote.Volume {
			if i < len(result.Timestamp) {
				timestamp := time.Unix(result.Timestamp[i], 0)
				hour := timestamp.Hour()

				totalVolume += vol
				if vol > 0 {
					nonZeroVolume++
				}

				// 시간대별 분류 (EST 기준)
				if hour < 9 || (hour == 9 && timestamp.Minute() < 30) {
					// Premarket: 9:30 AM EST 이전
					premarketVolume += vol
					if vol > 0 {
						premarketNonZero++
					}
				} else if hour < 16 {
					// Regular: 9:30 AM - 4:00 PM EST
					regularVolume += vol
					if vol > 0 {
						regularNonZero++
					}
				} else {
					// Postmarket: 4:00 PM EST 이후
					postmarketVolume += vol
					if vol > 0 {
						postmarketNonZero++
					}
				}
			}
		}

		log.Printf("Volume Statistics:")
		log.Printf("  Total Volume: %d", totalVolume)
		log.Printf("  Non-zero entries: %d/%d (%.1f%%)", nonZeroVolume, len(quote.Volume),
			float64(nonZeroVolume)/float64(len(quote.Volume))*100)

		if h.query.Prepost {
			log.Printf("Time-based Volume Distribution:")
			log.Printf("  Premarket:  %d volume, %d non-zero entries", premarketVolume, premarketNonZero)
			log.Printf("  Regular:    %d volume, %d non-zero entries", regularVolume, regularNonZero)
			log.Printf("  Postmarket: %d volume, %d non-zero entries", postmarketVolume, postmarketNonZero)

			// 첫 5개와 마지막 5개 volume이 0이 아닌 경우만 샘플 출력
			log.Printf("Sample volume data (non-zero only):")
			sampleCount := 0
			for i := 0; i < len(quote.Volume) && sampleCount < 5; i++ {
				if quote.Volume[i] > 0 && i < len(result.Timestamp) {
					timestamp := time.Unix(result.Timestamp[i], 0)
					log.Printf("  %s: %d", timestamp.Format("2006-01-02 15:04:05"), quote.Volume[i])
					sampleCount++
				}
			}
		}

		log.Printf("=== End Volume Analysis ===\n")
	}
}

type YahooIndicator struct {
	Quote []YahooQuote `json:"quote"`
}

type YahooQuote struct {
	Open   []float64 `json:"open"`
	High   []float64 `json:"high"`
	Low    []float64 `json:"low"`
	Close  []float64 `json:"close"`
	Volume []int64   `json:"volume"`
}

type PriceData struct {
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int64
}

type HistoryQuery struct {
	Range     string
	Interval  string
	Start     string
	End       string
	Prepost   bool
	UserAgent string
}

func (hq *HistoryQuery) SetDefault() {
	if hq.Range == "" && hq.Start == "" {
		hq.Range = "1mo"
	}
	if hq.Interval == "" {
		hq.Interval = "1d"
	}
	// Prepost는 기본적으로 false, 사용자가 명시적으로 설정하지 않는 한
	if hq.Start != "" {
		t, err := time.Parse("2006-01-02", hq.Start)
		if err != nil {
			log.Printf("Failed to parse start date: %v\n", err)
			hq.Start = "default"
		} else {
			hq.Start = fmt.Sprintf("%d", t.Unix())
		}
	}
	if hq.End == "" {
		hq.End = fmt.Sprintf("%d", time.Now().Unix())
	}
	if hq.UserAgent == "" {
		hq.UserAgent = USER_AGENTS[rand.Intn(len(USER_AGENTS))]
	}
}

type History struct {
	query  *HistoryQuery
	client *Client
}

func NewHistory() *History {
	return &History{query: &HistoryQuery{}, client: GetClient()}
}

func (h *History) SetQuery(query HistoryQuery) {
	h.query = &query
}

func (h *History) GetHistory(symbol string) (YahooHistoryRespose, error) {
	h.query.SetDefault()

	params := url.Values{}
	if h.query.Range != "" {
		params.Add("range", h.query.Range)
	}
	params.Add("interval", h.query.Interval)
	params.Add("period1", h.query.Start)
	params.Add("period2", h.query.End)

	// 사용자가 설정한 Prepost 값에 따라 includePrePost 파라미터 설정
	if h.query.Prepost {
		params.Add("includePrePost", "true")
	} else {
		params.Add("includePrePost", "false")
	}

	endpoint := fmt.Sprintf("%s/v8/finance/chart/%s", BASE_URL, symbol)
	resp, err := h.client.Get(endpoint, params)
	if err != nil {
		slog.Error("Failed to get history", "err", err)
		return YahooHistoryRespose{}, err
	}
	defer resp.Body.Close()

	var historyResponse YahooHistoryRespose
	if err := json.NewDecoder(resp.Body).Decode(&historyResponse); err != nil {
		// tradingPeriods 관련 오류인지 확인
		if strings.Contains(err.Error(), "tradingPeriods") {
			// tradingPeriods 오류는 경고로 처리하고 재시도
			slog.Warn("TradingPeriods parsing failed, retrying with custom parsing", "symbol", symbol, "err", err)

			// 응답을 다시 읽기 위해 재요청
			resp.Body.Close()
			resp, err = h.client.Get(endpoint, params)
			if err != nil {
				return YahooHistoryRespose{}, err
			}
			defer resp.Body.Close()

			// 커스텀 파싱으로 재시도
			return h.parseResponseWithFallback(resp, symbol)
		}

		slog.Error("Failed to decode history data JSON response", "err", err)
		return YahooHistoryRespose{}, err
	}

	if len(historyResponse.Chart.Result) == 0 {
		return YahooHistoryRespose{}, fmt.Errorf("no data found for symbol: %s", symbol)
	}

	// Premarket 요청 시 volume 데이터 디버깅
	if h.query.Prepost {
		h.DebugVolumeData(symbol, historyResponse)
	}

	return historyResponse, nil
}

// parseResponseWithFallback은 tradingPeriods 파싱 오류 시 대체 파싱을 수행합니다.
//
// 매개변수:
// - resp: HTTP 응답
// - symbol: 심볼명
//
// 반환값:
// - YahooHistoryRespose: 파싱된 응답 데이터
// - error: 파싱 중 발생한 오류
func (h *History) parseResponseWithFallback(resp *http.Response, symbol string) (YahooHistoryRespose, error) {
	// 수동으로 필요한 부분만 파싱
	var rawData map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&rawData); err != nil {
		slog.Error("Failed to decode raw JSON response", "err", err)
		return YahooHistoryRespose{}, err
	}

	// 차트 데이터 추출
	if chart, ok := rawData["chart"].(map[string]interface{}); ok {
		if result, ok := chart["result"].([]interface{}); ok && len(result) > 0 {
			if resultData, ok := result[0].(map[string]interface{}); ok {
				// 필수 데이터만 추출
				var response YahooHistoryRespose
				response.Chart.Result = make([]YahooHistoryResult, 1)

				// 타임스탬프 추출
				if timestamps, ok := resultData["timestamp"].([]interface{}); ok {
					for _, ts := range timestamps {
						if tsFloat, ok := ts.(float64); ok {
							response.Chart.Result[0].Timestamp = append(response.Chart.Result[0].Timestamp, int64(tsFloat))
						}
					}
				}

				// 인디케이터 데이터 추출
				if indicators, ok := resultData["indicators"].(map[string]interface{}); ok {
					if quote, ok := indicators["quote"].([]interface{}); ok && len(quote) > 0 {
						if quoteData, ok := quote[0].(map[string]interface{}); ok {
							response.Chart.Result[0].Indicators.Quote = make([]YahooQuote, 1)

							// OHLCV 데이터 추출
							extractFloatArray := func(key string) []float64 {
								var result []float64
								if arr, ok := quoteData[key].([]interface{}); ok {
									for _, v := range arr {
										if val, ok := v.(float64); ok {
											result = append(result, val)
										} else {
											result = append(result, 0.0) // null 값 처리
										}
									}
								}
								return result
							}

							extractIntArray := func(key string) []int64 {
								var result []int64
								if arr, ok := quoteData[key].([]interface{}); ok {
									for _, v := range arr {
										switch val := v.(type) {
										case float64:
											result = append(result, int64(val))
										case int64:
											result = append(result, val)
										case int:
											result = append(result, int64(val))
										case nil:
											result = append(result, 0) // null 값 처리
										default:
											result = append(result, 0) // 기타 타입 처리
										}
									}
								}
								return result
							}

							response.Chart.Result[0].Indicators.Quote[0].Open = extractFloatArray("open")
							response.Chart.Result[0].Indicators.Quote[0].High = extractFloatArray("high")
							response.Chart.Result[0].Indicators.Quote[0].Low = extractFloatArray("low")
							response.Chart.Result[0].Indicators.Quote[0].Close = extractFloatArray("close")
							response.Chart.Result[0].Indicators.Quote[0].Volume = extractIntArray("volume")
						}
					}
				}

				// 메타데이터 기본값 설정
				if meta, ok := resultData["meta"].(map[string]interface{}); ok {
					if symbolStr, ok := meta["symbol"].(string); ok {
						response.Chart.Result[0].Meta.Symbol = symbolStr
					}
					if currency, ok := meta["currency"].(string); ok {
						response.Chart.Result[0].Meta.Currency = currency
					}
					if timezone, ok := meta["timezone"].(string); ok {
						response.Chart.Result[0].Meta.Timezone = timezone
					}
					// tradingPeriods는 빈 RawMessage로 설정
					response.Chart.Result[0].Meta.TradingPeriods = json.RawMessage("{}")
				}

				return response, nil
			}
		}
	}

	return YahooHistoryRespose{}, fmt.Errorf("failed to parse response for symbol %s even with fallback method", symbol)
}

func (h *History) transformData(data YahooHistoryRespose) map[string]PriceData {
	d := make(map[string]PriceData)

	if len(data.Chart.Result) == 0 {
		return d
	}

	result := data.Chart.Result[0]
	if len(result.Indicators.Quote) == 0 {
		return d
	}

	quote := result.Indicators.Quote[0]

	for i, timestamp := range result.Timestamp {
		t := time.Unix(timestamp, 0)
		var key string
		if strings.HasSuffix(h.query.Interval, "d") || strings.HasSuffix(h.query.Interval, "wk") || strings.HasSuffix(h.query.Interval, "mo") {
			key = t.Format("2006-01-02")
		} else {
			key = t.Format("2006-01-02 15:04:05")
		}

		// 안전한 배열 접근 함수
		getFloatAt := func(arr []float64, index int) float64 {
			if index < len(arr) {
				return arr[index]
			}
			return 0.0
		}

		getVolumeAt := func(arr []int64, index int) int64 {
			if index < len(arr) {
				return arr[index]
			}
			return 0
		}

		priceData := PriceData{
			Open:   getFloatAt(quote.Open, i),
			High:   getFloatAt(quote.High, i),
			Low:    getFloatAt(quote.Low, i),
			Close:  getFloatAt(quote.Close, i),
			Volume: getVolumeAt(quote.Volume, i),
		}

		// premarket 시간대 volume 통계를 위한 카운팅 (개별 로그는 제거)
		if h.query.Prepost && priceData.Volume == 0 {
			// volume이 0인 경우의 통계는 DebugVolumeData에서 처리하므로
			// 여기서는 개별 로그를 출력하지 않음
		}

		d[key] = priceData
	}

	return d
}
