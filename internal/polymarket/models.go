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

type Activity struct {
	ID          string  `json:"id"`
	User        string  `json:"user"`        // может быть proxyWallet
	Type        string  `json:"type"`        // TRADE / MERGE / REDEEM / SPLIT / REWARD ...
	MarketTitle string  `json:"marketTitle"` // часто так
	Title       string  `json:"title"`       // иногда так
	Slug        string  `json:"slug"`
	ConditionID string  `json:"conditionId"`
	Side        string  `json:"side"`    // BUY / SELL
	Size        float64 `json:"sizeUsd"` // бывает sizeUsd / usdAmount
	CreatedAt   string  `json:"createdAt"`
}

type UserPosition struct {
	ProxyWallet        string  `json:"proxyWallet"`
	Asset              string  `json:"asset"`
	ConditionID        string  `json:"conditionId"`
	Size               float64 `json:"size"`               // сколько у пользователя штук
	AvgPrice           float64 `json:"avgPrice"`           // средняя цена покупки
	InitialValue       float64 `json:"initialValue"`       // сколько стоило при покупке
	CurrentValue       float64 `json:"currentValue"`       // сколько стоит сейчас в $
	CashPnL            float64 `json:"cashPnl"`            // прибыль/убыток в $
	PercentPnL         float64 `json:"percentPnl"`         // прибыль/убыток в %
	TotalBought        float64 `json:"totalBought"`        // всего куплено на $
	RealizedPnL        float64 `json:"realizedPnl"`        // уже зафиксированная прибыль
	PercentRealizedPnL float64 `json:"percentRealizedPnl"` // %
	CurPrice           float64 `json:"curPrice"`           // 🔥 текущая цена этого исхода
	Redeemable         bool    `json:"redeemable"`
	Mergeable          bool    `json:"mergeable"`
	Title              string  `json:"title"` // человеческое имя маркета
	Slug               string  `json:"slug"`
	Icon               string  `json:"icon"`
	EventSlug          string  `json:"eventSlug"`
	Outcome            string  `json:"outcome"` // "Yes" / "No" / другой вариант
	OutcomeIndex       int     `json:"outcomeIndex"`
	OppositeOutcome    string  `json:"oppositeOutcome"`
	OppositeAsset      string  `json:"oppositeAsset"`
	EndDate            string  `json:"endDate"`
	NegativeRisk       bool    `json:"negativeRisk"`
}

// ответ для /value
type UserValue struct {
	User  string  `json:"user"`
	Value float64 `json:"value"`
}
