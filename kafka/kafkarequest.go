package kafka

import (
	"fmt"
	"log"
	"time"
)

const (
	httpEndpointGetShipDetails       = "GetShipDetails"
	httpEndpointGetMarketplaceInfo   = "GetMarketplaceInfo"
	httpEndpointPostFlightPlanNew    = "PostFlightPlanNew"
	httpEndpointGetFlightPlanDetails = "GetFlightPlanDetails"
	httpEndpointPostBuyOrderNew      = "PostBuyOrderNew"
	httpEndpointPostSellOrderNew     = "PostSellOrderNew"
)

// GetShipInfo collects information about specific ship.
//
// https://api.spacetraders.io/#api-ships-GetShip
func (kp *KafkaProxy) GetShipInfo() ([]byte, error) {
	// return wp.get(fmt.Sprintf(httpEndpointGetShipDetails, wp.id, wp.token))
	msg := []byte{}
	err := kp.Write(fmt.Sprintf("{\"id\": \"%s\", \"action\": \"%s\"}", kp.id, httpEndpointGetShipDetails))
	if err != nil {
		log.Println("Error sending request to kafka:", err)
		return msg, err
	}

	for len(msg) < 1 {
		log.Println("Waiting for response from GetShipInfo...")
		time.Sleep(1 * time.Second)
		msg = kp.Read()
	}

	return msg, nil
}

// GetMarketplaceProducts gathers information about products available to trade in the planet where the ship is.
//
// https://api.spacetraders.io/#api-locations-GetMarketplace
func (kp *KafkaProxy) GetMarketplaceProducts(location string) ([]byte, error) {
	// return wp.get(fmt.Sprintf(httpEndpointGetMarketplaceInfo, location, wp.token))
	kp.Write(fmt.Sprintf("{\"action\": \"%s\", \"id\": \"%s\"}", httpEndpointGetMarketplaceInfo, kp.id))

	msg := []byte{}
	for len(msg) < 1 {
		time.Sleep(1 * time.Second)
		msg = kp.Read()
	}

	return msg, nil
}

// SetNewFlightPlan sends to game a new destination where the ships needs to fly to.
//
// https://api.spacetraders.io/#api-flight_plans-NewFlightPlan
func (kp *KafkaProxy) SetNewFlightPlan(destination string) ([]byte, error) {
	// return wp.post(
	// 	fmt.Sprintf(httpEndpointPostFlightPlanNew, wp.token),
	// 	bytes.NewReader(
	// 		[]byte(fmt.Sprintf("{\"shipId\": \"%s\", \"destination\": \"%s\"}", wp.id, destination))))
	kp.Write(fmt.Sprintf("{\"action\": \"%s\",\"shipId\": \"%s\",\"destination\":\"%s\"}", httpEndpointPostFlightPlanNew, kp.id, destination))

	msg := []byte{}
	for len(msg) < 1 {
		time.Sleep(1 * time.Second)
		msg = kp.Read()
	}

	return msg, nil
}

// GetFlightPlan retrieves information about current flight plan for specific ship, if any.
//
// https://api.spacetraders.io/#api-flight_plans-GetFlightPlan
func (kp *KafkaProxy) GetFlightPlan(planId string) ([]byte, error) {
	// return wp.get(fmt.Sprintf(httpEndpointGetFlightPlanDetails, planId, wp.token))
	kp.Write(fmt.Sprintf("{\"action\": \"%s\",\"planId\": \"%s\"}", httpEndpointGetFlightPlanDetails, planId))

	msg := []byte{}
	for len(msg) < 1 {
		time.Sleep(1 * time.Second)
		msg = kp.Read()
	}

	return msg, nil
}

// BuyGood sends to game a purchase order.
//
// https://api.spacetraders.io/#api-purchase_orders-NewPurchaseOrder
func (kp *KafkaProxy) BuyGood(good string, quantity int) ([]byte, error) {
	// return wp.post(
	// 	fmt.Sprintf(httpEndpointPostBuyOrderNew, wp.token),
	// 	bytes.NewReader(
	// 		[]byte(
	// 			fmt.Sprintf("{\"shipId\": \"%s\", \"good\": \"%s\", \"quantity\": %d}", wp.id, good, quantity))))
	kp.Write(fmt.Sprintf(
		"{\"action\": \"%s\",\"shipId\": \"%s\",\"good\": \"%s\",\"quantity\": %d}",
		httpEndpointPostBuyOrderNew, kp.id, good, quantity))

	msg := []byte{}
	for len(msg) < 1 {
		time.Sleep(1 * time.Second)
		msg = kp.Read()
	}

	return msg, nil
}

// SellGood sends to game a sell order.
//
// https://api.spacetraders.io/#api-sell_orders-NewSellOrder
func (kp *KafkaProxy) SellGood(good string, quantity int) ([]byte, error) {
	// return wp.post(
	// 	fmt.Sprintf(httpEndpointPostSellOrderNew, wp.token),
	// 	bytes.NewReader(
	// 		[]byte(
	// 			fmt.Sprintf("{\"shipId\": \"%s\", \"good\": \"%s\", \"quantity\": %d}", wp.id, good, quantity))))
	kp.Write(fmt.Sprintf(
		"{\"action\": \"%s\",\"shipId\": \"%s\",\"good\": \"%s\",\"quantity\": %d}",
		httpEndpointPostSellOrderNew, kp.id, good, quantity))

	msg := []byte{}
	for len(msg) < 1 {
		time.Sleep(1 * time.Second)
		msg = kp.Read()
	}

	return msg, nil
}
