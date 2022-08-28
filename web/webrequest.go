package web

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	httpEndpointGetShipDetails       = "https://api.spacetraders.io/my/ships/%s?token=%s"
	httpEndpointGetFlightPlanDetails = "https://api.spacetraders.io/my/flight-plans/%s?token=%s"
	httpEndpointPostFlightPlanNew    = "https://api.spacetraders.io/my/flight-plans?token=%s"
	httpEndpointGetMarketplaceInfo   = "https://api.spacetraders.io/locations/%s/marketplace?token=%s"
	httpEndpointPostSellOrderNew     = "https://api.spacetraders.io/my/sell-orders?token=%s"
	httpEndpointPostBuyOrderNew      = "https://api.spacetraders.io/my/purchase-orders?token=%s"

	maxRetriesTimeout = 5
	waitTimeout       = time.Duration(10)
)

// WebProxy is an implementation of web.Proxy.
type WebProxy struct {
	id      string
	token   string
	baseUrl string
}

// NewWebProxy creates a new instance of WebProxy.
//
// token is provided by the game when you claim your username.
func NewWebProxy(id, token string) Proxy {
	return &WebProxy{
		id:      id,
		token:   token,
		baseUrl: ""}
}

// GetShipInfo collects information about specific ship.
//
// https://api.spacetraders.io/#api-ships-GetShip
func (wp *WebProxy) GetShipInfo() ([]byte, error) {
	return wp.get(fmt.Sprintf(httpEndpointGetShipDetails, wp.id, wp.token))
}

// GetMarketplaceProducts gathers information about products available to trade in the planet where the ship is.
//
// https://api.spacetraders.io/#api-locations-GetMarketplace
func (wp *WebProxy) GetMarketplaceProducts(location string) ([]byte, error) {
	return wp.get(fmt.Sprintf(httpEndpointGetMarketplaceInfo, location, wp.token))
}

// SetNewFlightPlan sends to game a new destination where the ships needs to fly to.
//
// https://api.spacetraders.io/#api-flight_plans-NewFlightPlan
func (wp *WebProxy) SetNewFlightPlan(destination string) ([]byte, error) {
	return wp.post(
		fmt.Sprintf(httpEndpointPostFlightPlanNew, wp.token),
		bytes.NewReader(
			[]byte(fmt.Sprintf("{\"shipId\": \"%s\", \"destination\": \"%s\"}", wp.id, destination))))
}

// GetFlightPlan retrieves information about current flight plan for specific ship, if any.
//
// https://api.spacetraders.io/#api-flight_plans-GetFlightPlan
func (wp *WebProxy) GetFlightPlan(planId string) ([]byte, error) {
	return wp.get(fmt.Sprintf(httpEndpointGetFlightPlanDetails, planId, wp.token))
}

// BuyGood sends to game a purchase order.
//
// https://api.spacetraders.io/#api-purchase_orders-NewPurchaseOrder
func (wp *WebProxy) BuyGood(good string, quantity int) ([]byte, error) {
	return wp.post(
		fmt.Sprintf(httpEndpointPostBuyOrderNew, wp.token),
		bytes.NewReader(
			[]byte(
				fmt.Sprintf("{\"shipId\": \"%s\", \"good\": \"%s\", \"quantity\": %d}", wp.id, good, quantity))))
}

// SellGood sends to game a sell order.
//
// https://api.spacetraders.io/#api-sell_orders-NewSellOrder
func (wp *WebProxy) SellGood(good string, quantity int) ([]byte, error) {
	return wp.post(
		fmt.Sprintf(httpEndpointPostSellOrderNew, wp.token),
		bytes.NewReader(
			[]byte(
				fmt.Sprintf("{\"shipId\": \"%s\", \"good\": \"%s\", \"quantity\": %d}", wp.id, good, quantity))))
}

// get is a generic GET request, used by the other methods.
func (wp *WebProxy) get(uri string) ([]byte, error) {
	count := 0
	for count < maxRetriesTimeout {
		response, err := http.Get(uri)
		if err == nil {
			defer response.Body.Close()
			return io.ReadAll(response.Body)
		}

		errUrl := err.(*url.Error)
		if errUrl.Timeout() {
			time.Sleep(waitTimeout * time.Second)
		} else {
			return []byte{}, err
		}
	}
	return []byte{}, fmt.Errorf("reached unexpected piece of code in get")
}

// post is a generic POST request, used by the other methods.
func (wp *WebProxy) post(uri string, data io.Reader) ([]byte, error) {
	count := 0
	for count < maxRetriesTimeout {
		response, err := http.Post(uri, "application/json", data)
		if err == nil {
			defer response.Body.Close()
			return io.ReadAll(response.Body)
		}

		errUrl := err.(*url.Error)
		if errUrl.Timeout() {
			time.Sleep(waitTimeout * time.Second)
		} else {
			return []byte{}, err
		}
	}
	return []byte{}, fmt.Errorf("reached unexpected piece of code in post")
}
