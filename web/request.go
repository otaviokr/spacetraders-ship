package web

type Proxy interface {
	// GetShipInfo collects information about specific ship.
	//
	// https://api.spacetraders.io/#api-ships-GetShip
	GetShipInfo() ([]byte, error)

	// GetMarketplaceProducts gathers information about products available to trade in the planet where the ship is.
	//
	// https://api.spacetraders.io/#api-locations-GetMarketplace
	GetMarketplaceProducts(string) ([]byte, error)

	// SetNewFlightPlan sends to game a new destination where the ships needs to fly to.
	//
	// https://api.spacetraders.io/#api-flight_plans-NewFlightPlan
	SetNewFlightPlan(string) ([]byte, error)

	// GetFlightPlan retrieves information about current flight plan for specific ship, if any.
	//
	// https://api.spacetraders.io/#api-flight_plans-GetFlightPlan
	GetFlightPlan(string) ([]byte, error)

	// BuyGood sends to game a purchase order.
	//
	// https://api.spacetraders.io/#api-purchase_orders-NewPurchaseOrder
	BuyGood(string, int) ([]byte, error)

	// SellGood sends to game a sell order.
	//
	// https://api.spacetraders.io/#api-sell_orders-NewSellOrder
	SellGood(string, int) ([]byte, error)
}
