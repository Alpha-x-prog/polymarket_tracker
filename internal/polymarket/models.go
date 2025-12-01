package polymarket

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// type MarketsResponse struct {
// 	Markets []Market `json:"markets"`
// }

// ===== activity =====
type Activity struct {
	ID          string  `json:"id"`
	User        string  `json:"user"`
	Type        string  `json:"type"`
	MarketTitle string  `json:"marketTitle"`
	Title       string  `json:"title"`
	Slug        string  `json:"slug"`
	ConditionID string  `json:"conditionId"`
	Side        string  `json:"side"`
	SizeUSD     float64 `json:"sizeUsd"`
	USDValue    float64 `json:"usdAmount"`
	CreatedAt   string  `json:"createdAt"`
}

// ===== open positions =====
type UserPosition struct {
	ProxyWallet        string  `json:"proxyWallet"`
	Asset              string  `json:"asset"`
	ConditionID        string  `json:"conditionId"`
	Size               float64 `json:"size"`
	AvgPrice           float64 `json:"avgPrice"`
	InitialValue       float64 `json:"initialValue"`
	CurrentValue       float64 `json:"currentValue"`
	CashPnL            float64 `json:"cashPnl"`
	PercentPnL         float64 `json:"percentPnl"`
	TotalBought        float64 `json:"totalBought"`
	RealizedPnL        float64 `json:"realizedPnl"`
	PercentRealizedPnL float64 `json:"percentRealizedPnl"`
	CurPrice           float64 `json:"curPrice"`
	Redeemable         bool    `json:"redeemable"`
	Mergeable          bool    `json:"mergeable"`
	Title              string  `json:"title"`
	Slug               string  `json:"slug"`
	Icon               string  `json:"icon"`
	EventSlug          string  `json:"eventSlug"`
	Outcome            string  `json:"outcome"`
	OutcomeIndex       int     `json:"outcomeIndex"`
	OppositeOutcome    string  `json:"oppositeOutcome"`
	OppositeAsset      string  `json:"oppositeAsset"`
	EndDate            string  `json:"endDate"`
	NegativeRisk       bool    `json:"negativeRisk"`
}

// ===== closed positions (добавили) =====
type ClosedPosition struct {
	User               string  `json:"user"`
	ConditionID        string  `json:"conditionId"`
	Title              string  `json:"title"`
	Slug               string  `json:"slug"`
	Outcome            string  `json:"outcome"`
	Size               float64 `json:"size"`
	AvgPrice           float64 `json:"avgPrice"`
	RealizedPnL        float64 `json:"realizedPnl"`
	PercentRealizedPnL float64 `json:"percentRealizedPnl"`
	ClosedAt           string  `json:"closedAt"`
}

// ===== /value =====
type UserValue struct {
	User  string  `json:"user"`
	Value float64 `json:"value"`
}

// ===== /traded =====
type UserTraded struct {
	User   string `json:"user"`
	Traded int    `json:"traded"`
}

// Базовый Outcome (используется когда приходят полноформатные объекты)
type Outcome struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// Гибкий срез исходов. Умеет парсить:
// - []Outcome
// - []string
// - "Yes"/"No" (одна строка)
// - null
type FlexibleOutcomes []Outcome

func (fo *FlexibleOutcomes) UnmarshalJSON(b []byte) error {
	// []Outcome
	var detailed []Outcome
	if err := json.Unmarshal(b, &detailed); err == nil {
		*fo = detailed
		return nil
	}

	// []string
	var names []string
	if err := json.Unmarshal(b, &names); err == nil {
		out := make([]Outcome, 0, len(names))
		for _, n := range names {
			out = append(out, Outcome{Name: n})
		}
		*fo = out
		return nil
	}

	// "Yes"/"No" (одна строка)
	var single string
	if err := json.Unmarshal(b, &single); err == nil {
		*fo = []Outcome{{Name: single}}
		return nil
	}

	// null
	if string(b) == "null" {
		*fo = nil
		return nil
	}

	return fmt.Errorf("unknown outcomes format: %s", string(b))
}

type FlexiblePrices []float64

func (fp *FlexiblePrices) UnmarshalJSON(b []byte) error {
	// 1) []float64
	var f []float64
	if err := json.Unmarshal(b, &f); err == nil {
		*fp = f
		return nil
	}
	// 2) []string
	var ss []string
	if err := json.Unmarshal(b, &ss); err == nil {
		out := make([]float64, 0, len(ss))
		for _, s := range ss {
			if x, err := strconv.ParseFloat(strings.TrimSpace(s), 64); err == nil {
				out = append(out, x)
			}
		}
		*fp = out
		return nil
	}
	// 3) "0.63,0.37" (одна строка)
	var one string
	if err := json.Unmarshal(b, &one); err == nil {
		parts := strings.Split(one, ",")
		out := make([]float64, 0, len(parts))
		for _, p := range parts {
			if x, err := strconv.ParseFloat(strings.TrimSpace(p), 64); err == nil {
				out = append(out, x)
			}
		}
		*fp = out
		return nil
	}
	// 4) null
	if string(b) == "null" {
		*fp = nil
		return nil
	}
	return fmt.Errorf("unknown outcomePrices format: %s", string(b))
}

// Удобный хелпер: вернуть только имена исходов
func (fo FlexibleOutcomes) Names() []string {
	if len(fo) == 0 {
		return nil
	}
	names := make([]string, 0, len(fo))
	for _, o := range fo {
		names = append(names, o.Name)
	}
	return names
}

// -------- ТВОИ МОДЕЛИ С УЧЁТОМ ГИБКИХ OUTCOMES --------

type MarketsResponse struct {
	Markets []Market `json:"markets"`
}

type MarketCategory struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Slug  string `json:"slug"`
}

type Market struct {
	ID          string `json:"id"`
	ConditionID string `json:"conditionId"`

	Question string `json:"question"`
	Slug     string `json:"slug"`

	Category   string           `json:"category"`
	Categories []MarketCategory `json:"categories"`

	VolumeNum    float64 `json:"volumeNum"`
	Volume24hr   float64 `json:"volume24hr"`
	LiquidityNum float64 `json:"liquidityNum"`

	Spread  float64 `json:"spread"`
	BestBid float64 `json:"bestBid"`
	BestAsk float64 `json:"bestAsk"`

	Description      string `json:"description"`
	ResolutionSource string `json:"resolutionSource"`

	Outcomes      FlexibleOutcomes `json:"outcomes"`
	OutcomePrices FlexiblePrices   `json:"outcomePrices"`
}
