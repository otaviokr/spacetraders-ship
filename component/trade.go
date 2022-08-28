package component

// Trade represents a trading order (buy or sell).
type Trade struct {
	Credits int        `json:"credits"`
	Order   TradeOrder `json:"order"`
	Ship    Ship       `json:"ship"`
	Error   Error      `yaml:"error"`
}

// TradingOrder contains the details (to buy or to sell) of a good in the marketplace.
type TradeOrder struct {
	Good         string `json:"good"`
	PricePerUnit int    `json:"pricePerUnit"`
	Quantity     int    `json:"quantity"`
	Total        int    `json:"total"`
}

// Marketplace contains the list of products that can be traded.
type Marketplace struct {
	Products []Product `json:"marketplace"`
	Error    Error     `yaml:"error"`
}

// Product contains the details about a tradeable product.
type Product struct {
	PricePerUnit         int    `json:"pricePerUnit"`
	PurchasePricePerUnit int    `json:"purchasePricePerUnit"`
	QuantityAvailable    int    `json:"quantityAvailable"`
	SellPricePerUnit     int    `json:"sellPricePerUnit"`
	Spread               int    `json:"spread"`
	Symbol               string `json:"symbol"`
	VolumePerUnit        int    `json:"volumePerUnit"`
}
