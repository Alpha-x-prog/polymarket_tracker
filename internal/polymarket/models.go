package polymarket

type MarketsResponse struct {
	Markets []Market `json:"markets"`
}

type Market struct {
	ID       string    `json:"conditionId"`
	Question string    `json:"question"`
	Slug     string    `json:"slug"`
	Category string    `json:"category"`
	Volume24 float64   `json:"volume24h"`
	OI       float64   `json:"openInterest"`
	Outcomes []Outcome `json:"outcomes"`
}

type Outcome struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

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
