package component

// Trade represents a trading order (buy or sell).
type Trade struct {
	Credits int        `yaml:"credits"`
	Order   TradeOrder `yaml:"order"`
	Ship    Ship       `yaml:"ship"`
	Error   Error      `yaml:"error"`
}

// TradingOrder contains the details (to buy or to sell) of a good in the marketplace.
type TradeOrder struct {
	Good         string `yaml:"good"`
	PricePerUnit int    `yaml:"pricePerUnit"`
	Quantity     int    `yaml:"quantity"`
	Total        int    `yaml:"total"`
}

// Marketplace contains the list of products that can be traded.
type Marketplace struct {
	Products []Product `yaml:"marketplace"`
	Error    Error     `yaml:"error"`
}

// Product contains the details about a tradeable product.
type Product struct {
	PricePerUnit         int    `yaml:"pricePerUnit"`
	PurchasePricePerUnit int    `yaml:"purchasePricePerUnit"`
	QuantityAvailable    int    `yaml:"quantityAvailable"`
	SellPricePerUnit     int    `yaml:"sellPricePerUnit"`
	Spread               int    `yaml:"spread"`
	Symbol               string `yaml:"symbol"`
	VolumePerUnit        int    `yaml:"volumePerUnit"`
}
