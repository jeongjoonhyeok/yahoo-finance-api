package yahoofinanceapi

/*
 * Ticker Module
 *
 * 이 파일은 개별 주식 심볼에 대한 통합된 인터페이스를 제공합니다.
 *
 * 주요 기능:
 * - 주식 히스토리 데이터 조회
 * - 옵션 체인 정보 조회
 * - 실시간 quote 및 시가총액 정보 조회
 */

type Ticker struct {
	Symbol  string
	history *History
	option  *Option
	quote   *Quote
}

func NewTicker(symbol string) *Ticker {
	return &Ticker{Symbol: symbol}
}

// History는 주식의 과거 가격 데이터를 조회합니다.
//
// 매개변수:
// - query: 조회 조건을 담은 HistoryQuery 구조체
//
// 반환값:
// - map[string]PriceData: 날짜별 가격 데이터
// - error: 조회 중 발생한 오류
func (t *Ticker) History(query HistoryQuery) (map[string]PriceData, error) {
	if t.history == nil {
		t.history = NewHistory()
	}
	t.history.SetQuery(query)
	history, err := t.history.GetHistory(t.Symbol)
	if err != nil {
		return nil, err
	}
	return t.history.transformData(history), nil
}

// OptionChain은 주식의 옵션 체인 정보를 조회합니다.
//
// 반환값:
// - OptionData: 옵션 체인 데이터
func (t *Ticker) OptionChain() OptionData {
	if t.option == nil {
		t.option = NewOption()
	}
	optionChain := t.option.GetOptionChain(t.Symbol)
	return t.option.transformData(optionChain)
}

// OptionChainByExpiration은 특정 만료일의 옵션 체인 정보를 조회합니다.
//
// 매개변수:
// - expiration: 만료일 (형식: "2006-01-02")
//
// 반환값:
// - OptionData: 해당 만료일의 옵션 체인 데이터
func (t *Ticker) OptionChainByExpiration(expiration string) OptionData {
	if t.option == nil {
		t.option = NewOption()
	}
	optionChain := t.option.GetOptionChainByExpiration(t.Symbol, expiration)
	return t.option.transformData(optionChain)
}

// ExpirationDates는 옵션의 가능한 만료일 목록을 조회합니다.
//
// 반환값:
// - []string: 만료일 목록
func (t *Ticker) ExpirationDates() []string {
	if t.option == nil {
		t.option = NewOption()
	}
	expirationDates := t.option.GetExpirationDates(t.Symbol)
	return expirationDates
}

// Quote는 주식의 실시간 quote 정보를 조회합니다.
//
// 반환값:
// - StockQuote: 실시간 quote 정보
// - error: 조회 중 발생한 오류
func (t *Ticker) Quote() (StockQuote, error) {
	if t.quote == nil {
		t.quote = NewQuote()
	}
	return t.quote.GetQuote(t.Symbol)
}

// MarketCap은 주식의 시가총액 정보를 조회합니다.
//
// 반환값:
// - MarketCapInfo: 시가총액 관련 정보
// - error: 조회 중 발생한 오류
func (t *Ticker) MarketCap() (MarketCapInfo, error) {
	if t.quote == nil {
		t.quote = NewQuote()
	}
	return t.quote.GetMarketCap(t.Symbol)
}
