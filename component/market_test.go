package component_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/otaviokr/spacetraders-ship/component"
	"github.com/otaviokr/spacetraders-ship/mocks"
	"go.opentelemetry.io/otel/trace"
)

// GetMarketplaceProducts will fetch the products that are available to be traded in the current marketplace.
func TestGetMarketplaceProducts(t *testing.T) {
	useCases := map[string]map[string]interface{}{
		"uc1": {
			"id":              "id0001",
			"tracer":          trace.NewNoopTracerProvider().Tracer(""),
			"location":        "Local0001",
			"detailsResponse": "{\"ship\":{\"id\":\"id0001\",\"location\":\"Local0001\",\"x\":52,\"y\":3,\"cargo\":[{\"good\":\"FUEL\",\"quantity\":14,\"totalVolume\":14}],\"spaceAvailable\":286,\"type\":\"GR-MK-II\",\"class\":\"MK-II\",\"maxCargo\":300,\"loadingSpeed\":500,\"speed\":1,\"manufacturer\":\"Gravager\",\"plating\":10,\"weapons\":5}}",
			"marketResponse":  "{\"marketplace\": [{\"pricePerUnit\": 1,\"purchasePricePerUnit\": 2,\"quantityAvailable\": 3,\"sellPricePerUnit\": 4,\"spread\": 6,\"symbol\": \"Good0001\",\"volumePerUnit\": 7}]}",
			"expectedMarketplace": &component.Marketplace{
				Products: []component.Product{
					{
						PricePerUnit:         1,
						PurchasePricePerUnit: 2,
						QuantityAvailable:    3,
						SellPricePerUnit:     4,
						Spread:               6,
						Symbol:               "Good0001",
						VolumePerUnit:        7}}},
			"expectedProducts": &map[string]component.Product{
				"Good0001": {
					PricePerUnit:         1,
					PurchasePricePerUnit: 2,
					QuantityAvailable:    3,
					SellPricePerUnit:     4,
					Spread:               6,
					Symbol:               "Good0001",
					VolumePerUnit:        7}}}}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	proxy := mocks.NewMockProxy(ctrl)

	for _, uc := range useCases {
		proxy.EXPECT().GetShipInfo().Return([]byte(fmt.Sprintf("%v", uc["detailsResponse"])), nil)
		proxy.EXPECT().GetShipInfo().Return([]byte(fmt.Sprintf("%v", uc["detailsResponse"])), nil)
		proxy.EXPECT().
			GetMarketplaceProducts(fmt.Sprintf("%v", uc["location"])).
			Return([]byte(fmt.Sprintf("%v", uc["marketResponse"])), nil)
		ship, err := component.NewShipCustomProxy(
			context.TODO(),
			trace.NewNoopTracerProvider().Tracer(""),
			proxy,
			fmt.Sprintf("%v", uc["id"]))
		if err != nil {
			t.Fail()
		}

		actualMarketplace, actualProducts, err := ship.GetMarketplaceProducts(context.TODO())
		if err != nil {
			t.Fail()
		}

		if !reflect.DeepEqual(actualMarketplace, uc["expectedMarketplace"].(*component.Marketplace)) {
			t.Fatalf("\nACTUAL: %+v\nEXPECT: %+v\n", actualMarketplace, uc["expectedMarketplace"].(*component.Marketplace))
		}

		if !reflect.DeepEqual(actualProducts, uc["expectedProducts"].(*map[string]component.Product)) {
			t.Fatalf("\n%+v\n%+v\n", actualProducts, uc["expectedProducts"].(*map[string]component.Product))
		}
	}
}

func TestDoCommerce(t *testing.T) {
	useCases := map[string]map[string]interface{}{
		"uc1": {
			"id":                 "id0001",
			"tracer":             trace.NewNoopTracerProvider().Tracer(""),
			"location":           "Local0001",
			"buy":                map[string]int{"good0001": 2},
			"sell":               map[string]int{"good0001": 2},
			"detailsResponse":    "{\"ship\":{\"id\":\"id0001\",\"location\":\"Local0001\",\"x\":52,\"y\":3,\"cargo\":[{\"good\":\"FUEL\",\"quantity\":14,\"totalVolume\":14}],\"spaceAvailable\":286,\"type\":\"GR-MK-II\",\"class\":\"MK-II\",\"maxCargo\":300,\"loadingSpeed\":500,\"speed\":1,\"manufacturer\":\"Gravager\",\"plating\":10,\"weapons\":5}}",
			"flightPlanResponse": "{\"flightPlan\": {\"arrivesAt\": \"2021-05-13T18:41:24.963Z\",\"createdAt\": \"2021-05-13T18:40:23.003Z\",\"departure\": \"OE-PM-TR\",\"destination\": \"OE-PM\",\"distance\": 1,\"fuelConsumed\": 1,\"fuelRemaining\": 18,\"id\": \"flightplanid0001\",\"shipId\": \"id0001\",\"terminatedAt\": null,\"timeRemainingInSeconds\": 1}}",
			"marketResponse":     "{\"marketplace\": [{\"pricePerUnit\": 1,\"purchasePricePerUnit\": 2,\"quantityAvailable\": 3,\"sellPricePerUnit\": 4,\"spread\": 6,\"symbol\": \"Good0001\",\"volumePerUnit\": 7}]}",
			"expectedMarketplace": &component.Marketplace{
				Products: []component.Product{
					{
						PricePerUnit:         1,
						PurchasePricePerUnit: 2,
						QuantityAvailable:    3,
						SellPricePerUnit:     4,
						Spread:               6,
						Symbol:               "Good0001",
						VolumePerUnit:        7}}},
			"expectedProducts": &map[string]component.Product{
				"Good0001": {
					PricePerUnit:         1,
					PurchasePricePerUnit: 2,
					QuantityAvailable:    3,
					SellPricePerUnit:     4,
					Spread:               6,
					Symbol:               "Good0001",
					VolumePerUnit:        7}}}}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	proxy := mocks.NewMockProxy(ctrl)

	for _, uc := range useCases {
		proxy.EXPECT().GetShipInfo().Return([]byte(fmt.Sprintf("%v", uc["detailsResponse"])), nil)
		proxy.EXPECT().GetShipInfo().Return([]byte(fmt.Sprintf("%v", uc["detailsResponse"])), nil)
		proxy.EXPECT().GetShipInfo().Return([]byte(fmt.Sprintf("%v", uc["detailsResponse"])), nil)
		proxy.EXPECT().
			GetMarketplaceProducts(fmt.Sprintf("%v", uc["location"])).
			Return(
				[]byte(fmt.Sprintf("%v", uc["marketResponse"])),
				nil)
		// proxy.EXPECT().SetNewFlightPlan(
		// 	fmt.Sprintf("%v", uc["destination"])).Return([]byte(fmt.Sprintf("%v", uc["flightPlanResponse"])),
		// 	nil)
		ship, err := component.NewShipCustomProxy(
			context.TODO(),
			trace.NewNoopTracerProvider().Tracer(""),
			proxy,
			fmt.Sprintf("%v", uc["id"]))
		if err != nil {
			t.Fail()
		}

		err = ship.DoCommerce(context.TODO(), uc["sell"].(map[string]int), uc["buy"].(map[string]int))
		if err != nil {
			t.Log(err)
			t.Fail()
		}
	}
}

func TestSellAll(t *testing.T) {
	useCases := map[string]map[string]interface{}{
		"uc1": {
			"id":                 "id0001",
			"tracer":             trace.NewNoopTracerProvider().Tracer(""),
			"location":           "Local0001",
			"buy":                map[string]int{"good0001": 2},
			"sell":               map[string]int{"good0001": 2},
			"detailsResponse":    "{\"ship\":{\"id\":\"id0001\",\"location\":\"Local0001\",\"x\":52,\"y\":3,\"cargo\":[{\"good\":\"FUEL\",\"quantity\":14,\"totalVolume\":14}],\"spaceAvailable\":286,\"type\":\"GR-MK-II\",\"class\":\"MK-II\",\"maxCargo\":300,\"loadingSpeed\":500,\"speed\":1,\"manufacturer\":\"Gravager\",\"plating\":10,\"weapons\":5}}",
			"flightPlanResponse": "{\"flightPlan\": {\"arrivesAt\": \"2021-05-13T18:41:24.963Z\",\"createdAt\": \"2021-05-13T18:40:23.003Z\",\"departure\": \"OE-PM-TR\",\"destination\": \"OE-PM\",\"distance\": 1,\"fuelConsumed\": 1,\"fuelRemaining\": 18,\"id\": \"flightplanid0001\",\"shipId\": \"id0001\",\"terminatedAt\": null,\"timeRemainingInSeconds\": 1}}",
			"marketResponse":     "{\"marketplace\": [{\"pricePerUnit\": 1,\"purchasePricePerUnit\": 2,\"quantityAvailable\": 3,\"sellPricePerUnit\": 4,\"spread\": 6,\"symbol\": \"Good0001\",\"volumePerUnit\": 7}]}",
			"marketplace": map[string]component.Product{
				"Good0001": {
					PricePerUnit:         1,
					PurchasePricePerUnit: 2,
					QuantityAvailable:    3,
					SellPricePerUnit:     4,
					Spread:               6,
					Symbol:               "Good0001",
					VolumePerUnit:        7}},
			"expectedProducts": &map[string]component.Product{
				"Good0001": {
					PricePerUnit: 1,

					PurchasePricePerUnit: 2,
					QuantityAvailable:    3,
					SellPricePerUnit:     4,
					Spread:               6,
					Symbol:               "Good0001",
					VolumePerUnit:        7}}}}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	proxy := mocks.NewMockProxy(ctrl)

	for _, uc := range useCases {
		proxy.EXPECT().GetShipInfo().Return([]byte(fmt.Sprintf("%v", uc["detailsResponse"])), nil)
		// proxy.EXPECT().GetShipInfo().Return([]byte(fmt.Sprintf("%v", uc["detailsResponse"])), nil)
		// proxy.EXPECT().GetShipInfo().Return([]byte(fmt.Sprintf("%v", uc["detailsResponse"])), nil)
		// proxy.EXPECT().
		// 	GetMarketplaceProducts(fmt.Sprintf("%v", uc["location"])).
		// 	Return(
		// 		[]byte(fmt.Sprintf("%v", uc["marketResponse"])),
		// 		nil)
		// proxy.EXPECT().SetNewFlightPlan(
		// 	fmt.Sprintf("%v", uc["destination"])).Return([]byte(fmt.Sprintf("%v", uc["flightPlanResponse"])),
		// 	nil)
		ship, err := component.NewShipCustomProxy(
			context.TODO(),
			trace.NewNoopTracerProvider().Tracer(""),
			proxy,
			fmt.Sprintf("%v", uc["id"]))
		if err != nil {
			t.Fail()
		}

		err = ship.SellAll(context.TODO(), uc["sell"].(map[string]int), uc["marketplace"].(map[string]component.Product))
		if err != nil {
			t.Log(err)
			t.Fail()
		}
	}
}

func TestBuyAll(t *testing.T) {
	useCases := map[string]map[string]interface{}{
		"uc1": {
			"id":                 "id0001",
			"tracer":             trace.NewNoopTracerProvider().Tracer(""),
			"location":           "Local0001",
			"buy":                map[string]int{"good0001": 2},
			"sell":               map[string]int{"good0001": 2},
			"detailsResponse":    "{\"ship\":{\"id\":\"id0001\",\"location\":\"Local0001\",\"x\":52,\"y\":3,\"cargo\":[{\"good\":\"FUEL\",\"quantity\":14,\"totalVolume\":14}],\"spaceAvailable\":286,\"type\":\"GR-MK-II\",\"class\":\"MK-II\",\"maxCargo\":300,\"loadingSpeed\":500,\"speed\":1,\"manufacturer\":\"Gravager\",\"plating\":10,\"weapons\":5}}",
			"flightPlanResponse": "{\"flightPlan\": {\"arrivesAt\": \"2021-05-13T18:41:24.963Z\",\"createdAt\": \"2021-05-13T18:40:23.003Z\",\"departure\": \"OE-PM-TR\",\"destination\": \"OE-PM\",\"distance\": 1,\"fuelConsumed\": 1,\"fuelRemaining\": 18,\"id\": \"flightplanid0001\",\"shipId\": \"id0001\",\"terminatedAt\": null,\"timeRemainingInSeconds\": 1}}",
			"marketResponse":     "{\"marketplace\": [{\"pricePerUnit\": 1,\"purchasePricePerUnit\": 2,\"quantityAvailable\": 3,\"sellPricePerUnit\": 4,\"spread\": 6,\"symbol\": \"Good0001\",\"volumePerUnit\": 7}]}",
			"marketplace": map[string]component.Product{
				"Good0001": {
					PricePerUnit:         1,
					PurchasePricePerUnit: 2,
					QuantityAvailable:    3,
					SellPricePerUnit:     4,
					Spread:               6,
					Symbol:               "Good0001",
					VolumePerUnit:        7}},
			"expectedProducts": &map[string]component.Product{
				"Good0001": {
					PricePerUnit: 1,

					PurchasePricePerUnit: 2,
					QuantityAvailable:    3,
					SellPricePerUnit:     4,
					Spread:               6,
					Symbol:               "Good0001",
					VolumePerUnit:        7}}}}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	proxy := mocks.NewMockProxy(ctrl)

	for _, uc := range useCases {
		proxy.EXPECT().GetShipInfo().Return([]byte(fmt.Sprintf("%v", uc["detailsResponse"])), nil)
		// proxy.EXPECT().GetShipInfo().Return([]byte(fmt.Sprintf("%v", uc["detailsResponse"])), nil)
		// proxy.EXPECT().GetShipInfo().Return([]byte(fmt.Sprintf("%v", uc["detailsResponse"])), nil)
		// proxy.EXPECT().
		// 	GetMarketplaceProducts(fmt.Sprintf("%v", uc["location"])).
		// 	Return(
		// 		[]byte(fmt.Sprintf("%v", uc["marketResponse"])),
		// 		nil)
		// proxy.EXPECT().SetNewFlightPlan(
		// 	fmt.Sprintf("%v", uc["destination"])).Return([]byte(fmt.Sprintf("%v", uc["flightPlanResponse"])),
		// 	nil)
		ship, err := component.NewShipCustomProxy(
			context.TODO(),
			trace.NewNoopTracerProvider().Tracer(""),
			proxy,
			fmt.Sprintf("%v", uc["id"]))
		if err != nil {
			t.Fail()
		}

		err = ship.BuyAll(context.TODO(), uc["sell"].(map[string]int), uc["marketplace"].(map[string]component.Product))
		if err != nil {
			t.Log(err)
			t.Fail()
		}
	}
}

func TestSell(t *testing.T) {
	useCases := map[string]map[string]interface{}{
		"uc1": {
			"id":                 "id0001",
			"tracer":             trace.NewNoopTracerProvider().Tracer(""),
			"good":               "Good0001",
			"quantity":           2,
			"detailsResponse":    "{\"ship\":{\"id\":\"id0001\",\"flightPlanId\":\"flightPlanId0001\",\"cargo\":[{\"good\":\"FUEL\",\"quantity\":14,\"totalVolume\":14}],\"spaceAvailable\":286,\"type\":\"GR-MK-II\",\"class\":\"MK-II\",\"maxCargo\":300,\"loadingSpeed\":500,\"speed\":1,\"manufacturer\":\"Gravager\",\"plating\":10,\"weapons\":5}}",
			"flightPlanResponse": "{\"flightPlan\": {\"arrivesAt\": \"2021-05-13T18:41:24.963Z\",\"createdAt\": \"2021-05-13T18:40:23.003Z\",\"departure\": \"OE-PM-TR\",\"destination\": \"OE-PM\",\"distance\": 1,\"fuelConsumed\": 1,\"fuelRemaining\": 18,\"id\": \"flightplanid0001\",\"shipId\": \"id0001\",\"terminatedAt\": null,\"timeRemainingInSeconds\": 1}}",
			"sellGoodResponse":   "{\"credits\": 123, \"order\": {\"good\": \"Good0001\",\"pricePerUnit\": 12,\"quantity\": 2,\"total\": 5}}",
			"expectedSell": &component.Trade{
				Credits: 123,
				Order: component.TradeOrder{
					Good:         "Good0001",
					PricePerUnit: 12,
					Quantity:     2,
					Total:        5},
				Ship: component.Ship{},
				Error: component.Error{
					Code:    0,
					Message: ""}}},
		"uc2": {
			"id":                 "id0002",
			"tracer":             trace.NewNoopTracerProvider().Tracer(""),
			"good":               "Good0002",
			"quantity":           3,
			"detailsResponse":    "{\"ship\":{\"id\":\"id0002\",\"flightPlanId\":\"flightPlanId0002\",\"cargo\":[{\"good\":\"FUEL\",\"quantity\":41,\"totalVolume\":55}],\"spaceAvailable\":286,\"type\":\"GR-MK-II\",\"class\":\"MK-II\",\"maxCargo\":300,\"loadingSpeed\":500,\"speed\":1,\"manufacturer\":\"Gravager\",\"plating\":10,\"weapons\":5}}",
			"flightPlanResponse": "{\"flightPlan\": {\"arrivesAt\": \"2022-06-13T19:44:24.963Z\",\"createdAt\": \"2022-06-13T19:44:23.003Z\",\"departure\": \"OE-PM-TR\",\"destination\": \"OE-PM\",\"distance\": 1,\"fuelConsumed\": 1,\"fuelRemaining\": 18,\"id\": \"flightplanid0001\",\"shipId\": \"id0002\",\"terminatedAt\": null,\"timeRemainingInSeconds\": 1}}",
			"sellGoodResponse":   "{\"credits\": 456, \"order\": {\"good\": \"Good0002\",\"pricePerUnit\": 9,\"quantity\": 3,\"total\": 7}}",
			"expectedSell": &component.Trade{
				Credits: 456,
				Order: component.TradeOrder{
					Good:         "Good0002",
					PricePerUnit: 9,
					Quantity:     3,
					Total:        7},
				Ship: component.Ship{},
				Error: component.Error{
					Code:    0,
					Message: ""}}}}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	proxy := mocks.NewMockProxy(ctrl)

	for _, uc := range useCases {
		proxy.EXPECT().GetShipInfo().Return([]byte(fmt.Sprintf("%v", uc["detailsResponse"])), nil)
		// proxy.EXPECT().GetShipInfo().Return([]byte(fmt.Sprintf("%v", uc["detailsResponse"])), nil)
		proxy.EXPECT().SellGood(fmt.Sprintf("%v", uc["good"]), uc["quantity"].(int)).
			Return([]byte(fmt.Sprintf("%v", uc["sellGoodResponse"])), nil)
		ship, err := component.NewShipCustomProxy(
			context.TODO(),
			trace.NewNoopTracerProvider().Tracer(""),
			proxy,
			fmt.Sprintf("%v", uc["id"]))
		if err != nil {
			t.Fail()
		}

		actualSell, err := ship.Sell(context.TODO(), fmt.Sprintf("%v", uc["good"]), uc["quantity"].(int))
		if err != nil {
			t.Log(err)
			t.Fail()
		}

		if actualSell == nil {
			t.Fatal("actualSell is nil")
		} else if !reflect.DeepEqual(actualSell, uc["expectedSell"].(*component.Trade)) {
			t.Fatalf("\n%+v\n%+v\n", actualSell, uc["expectedSell"].(*component.Trade))
		}
	}
}

func TestBuy(t *testing.T) {
	useCases := map[string]map[string]interface{}{
		"uc1": {
			"id":                 "id0001",
			"tracer":             trace.NewNoopTracerProvider().Tracer(""),
			"good":               "Good0001",
			"quantity":           2,
			"detailsResponse":    "{\"ship\":{\"id\":\"id0001\",\"flightPlanId\":\"flightPlanId0001\",\"cargo\":[{\"good\":\"FUEL\",\"quantity\":14,\"totalVolume\":14}],\"spaceAvailable\":286,\"type\":\"GR-MK-II\",\"class\":\"MK-II\",\"maxCargo\":300,\"loadingSpeed\":500,\"speed\":1,\"manufacturer\":\"Gravager\",\"plating\":10,\"weapons\":5}}",
			"flightPlanResponse": "{\"flightPlan\": {\"arrivesAt\": \"2021-05-13T18:41:24.963Z\",\"createdAt\": \"2021-05-13T18:40:23.003Z\",\"departure\": \"OE-PM-TR\",\"destination\": \"OE-PM\",\"distance\": 1,\"fuelConsumed\": 1,\"fuelRemaining\": 18,\"id\": \"flightplanid0001\",\"shipId\": \"id0001\",\"terminatedAt\": null,\"timeRemainingInSeconds\": 1}}",
			"buyGoodResponse":    "{\"credits\": 123, \"order\": {\"good\": \"Good0001\",\"pricePerUnit\": 12,\"quantity\": 2,\"total\": 5}}",
			"expectedBuy": &component.Trade{
				Credits: 123,
				Order: component.TradeOrder{
					Good:         "Good0001",
					PricePerUnit: 12,
					Quantity:     2,
					Total:        5},
				Ship: component.Ship{},
				Error: component.Error{
					Code:    0,
					Message: ""}}},
		"uc2": {
			"id":                 "id0002",
			"tracer":             trace.NewNoopTracerProvider().Tracer(""),
			"good":               "Good0002",
			"quantity":           3,
			"detailsResponse":    "{\"ship\":{\"id\":\"id0002\",\"flightPlanId\":\"flightPlanId0002\",\"cargo\":[{\"good\":\"FUEL\",\"quantity\":41,\"totalVolume\":55}],\"spaceAvailable\":286,\"type\":\"GR-MK-II\",\"class\":\"MK-II\",\"maxCargo\":300,\"loadingSpeed\":500,\"speed\":1,\"manufacturer\":\"Gravager\",\"plating\":10,\"weapons\":5}}",
			"flightPlanResponse": "{\"flightPlan\": {\"arrivesAt\": \"2022-06-13T19:44:24.963Z\",\"createdAt\": \"2022-06-13T19:44:23.003Z\",\"departure\": \"OE-PM-TR\",\"destination\": \"OE-PM\",\"distance\": 1,\"fuelConsumed\": 1,\"fuelRemaining\": 18,\"id\": \"flightplanid0001\",\"shipId\": \"id0002\",\"terminatedAt\": null,\"timeRemainingInSeconds\": 1}}",
			"buyGoodResponse":    "{\"credits\": 456, \"order\": {\"good\": \"Good0002\",\"pricePerUnit\": 9,\"quantity\": 3,\"total\": 7}}",
			"expectedBuy": &component.Trade{
				Credits: 456,
				Order: component.TradeOrder{
					Good:         "Good0002",
					PricePerUnit: 9,
					Quantity:     3,
					Total:        7},
				Ship: component.Ship{},
				Error: component.Error{
					Code:    0,
					Message: ""}}}}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	proxy := mocks.NewMockProxy(ctrl)

	for _, uc := range useCases {
		proxy.EXPECT().GetShipInfo().Return([]byte(fmt.Sprintf("%v", uc["detailsResponse"])), nil)
		// proxy.EXPECT().GetShipInfo().Return([]byte(fmt.Sprintf("%v", uc["detailsResponse"])), nil)
		proxy.EXPECT().BuyGood(fmt.Sprintf("%v", uc["good"]), uc["quantity"].(int)).
			Return([]byte(fmt.Sprintf("%v", uc["buyGoodResponse"])), nil)
		ship, err := component.NewShipCustomProxy(
			context.TODO(),
			trace.NewNoopTracerProvider().Tracer(""),
			proxy,
			fmt.Sprintf("%v", uc["id"]))
		if err != nil {
			t.Fail()
		}

		actualBuy, err := ship.Buy(context.TODO(), fmt.Sprintf("%v", uc["good"]), uc["quantity"].(int))
		if err != nil {
			t.Log(err)
			t.Fail()
		}

		if actualBuy == nil {
			t.Fatal("actualBuy is nil")
		} else if !reflect.DeepEqual(actualBuy, uc["expectedBuy"].(*component.Trade)) {
			t.Fatalf("\n%+v\n%+v\n", actualBuy, uc["expectedBuy"].(*component.Trade))
		}
	}
}

func TestForceBuyFuel(t *testing.T) {
	useCases := map[string]map[string]interface{}{
		"uc1": {
			"id":                 "id0001",
			"tracer":             trace.NewNoopTracerProvider().Tracer(""),
			"good":               "FUEL",
			"quantity":           2,
			"detailsResponse":    "{\"ship\":{\"id\":\"id0001\",\"flightPlanId\":\"flightPlanId0001\",\"cargo\":[{\"good\":\"FUEL\",\"quantity\":14,\"totalVolume\":14}],\"spaceAvailable\":286,\"type\":\"GR-MK-II\",\"class\":\"MK-II\",\"maxCargo\":300,\"loadingSpeed\":500,\"speed\":1,\"manufacturer\":\"Gravager\",\"plating\":10,\"weapons\":5}}",
			"flightPlanResponse": "{\"flightPlan\": {\"arrivesAt\": \"2021-05-13T18:41:24.963Z\",\"createdAt\": \"2021-05-13T18:40:23.003Z\",\"departure\": \"OE-PM-TR\",\"destination\": \"OE-PM\",\"distance\": 1,\"fuelConsumed\": 1,\"fuelRemaining\": 18,\"id\": \"flightplanid0001\",\"shipId\": \"id0001\",\"terminatedAt\": null,\"timeRemainingInSeconds\": 1}}",
			"buyGoodResponse":    "{\"credits\": 123, \"order\": {\"good\": \"FUEL\",\"pricePerUnit\": 12,\"quantity\": 2,\"total\": 5}}",
			"expectedbuy": &component.Trade{
				Credits: 123,
				Order: component.TradeOrder{
					Good:         "FUEL",
					PricePerUnit: 12,
					Quantity:     2,
					Total:        5},
				Ship: component.Ship{},
				Error: component.Error{
					Code:    0,
					Message: ""}}}}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	proxy := mocks.NewMockProxy(ctrl)

	for _, uc := range useCases {
		proxy.EXPECT().GetShipInfo().Return([]byte(fmt.Sprintf("%v", uc["detailsResponse"])), nil)
		proxy.EXPECT().BuyGood(fmt.Sprintf("%v", uc["good"]), uc["quantity"].(int)).
			Return([]byte(fmt.Sprintf("%v", uc["buyGoodResponse"])), nil)
		ship, err := component.NewShipCustomProxy(
			context.TODO(),
			trace.NewNoopTracerProvider().Tracer(""),
			proxy,
			fmt.Sprintf("%v", uc["id"]))
		if err != nil {
			t.Fail()
		}

		err = ship.ForceBuyFuel(context.TODO(), uc["quantity"].(int))
		if err != nil {
			t.Log(err)
			t.Fail()
		}

		// if actualSell == nil {
		// 	t.Fatal("actualSell is nil")
		// } else if !reflect.DeepEqual(actualSell, uc["expectedSell"].(*component.Trade)) {
		// 	t.Fatalf("\n%+v\n%+v\n", actualSell, uc["expectedSell"].(*component.Trade))
		// }
	}
}

func TestForceBuyFuelNoSpace(t *testing.T) {
	useCases := map[string]map[string]interface{}{
		"uc2": {
			"id":                  "id0002",
			"tracer":              trace.NewNoopTracerProvider().Tracer(""),
			"good":                "FUEL",
			"quantity":            2,
			"location":            "testlocation",
			"detailsResponse":     "{\"ship\":{\"id\":\"id0001\",\"flightPlanId\":\"flightPlanId0001\",\"location\":\"testlocation\",\"cargo\":[{\"good\":\"FUEL\",\"quantity\":1,\"totalVolume\":1},{\"good\":\"Good001\",\"quantity\":50,\"totalVolume\":50}],\"spaceAvailable\":2,\"type\":\"GR-MK-II\",\"class\":\"MK-II\",\"maxCargo\":300,\"loadingSpeed\":500,\"speed\":1,\"manufacturer\":\"Gravager\",\"plating\":10,\"weapons\":5}}",
			"flightPlanResponse":  "{\"flightPlan\": {\"arrivesAt\": \"2021-05-13T18:41:24.963Z\",\"createdAt\": \"2021-05-13T18:40:23.003Z\",\"departure\": \"OE-PM-TR\",\"destination\": \"OE-PM\",\"distance\": 1,\"fuelConsumed\": 1,\"fuelRemaining\": 18,\"id\": \"flightplanid0001\",\"shipId\": \"id0001\",\"terminatedAt\": null,\"timeRemainingInSeconds\": 1}}",
			"buyGoodResponse":     "{\"credits\": 123, \"order\": {\"good\": \"FUEL\",\"pricePerUnit\": 12,\"quantity\": 2,\"total\": 5}}",
			"marketplaceResponse": "{\"marketplace\":[{\"symbol\":\"Good001\",\"pricePerUnit\":1,\"purchasePricePerUnit\":1,\"sellPricePerUnit\":1,\"volumePerUnit\":1}]}",
			"sellGoodResponse":    "{\"credits\": 456, \"order\": {\"good\": \"Good0002\",\"pricePerUnit\": 9,\"quantity\": 3,\"total\": 7}}",
			"products": &map[string]component.Product{
				"FUEL": {
					Symbol:               "FUEL",
					PricePerUnit:         1,
					QuantityAvailable:    20,
					SellPricePerUnit:     5,
					VolumePerUnit:        1,
					PurchasePricePerUnit: 2,
					Spread:               1}},
			"marketplace": &component.Marketplace{
				Products: []component.Product{},
				Error: component.Error{
					Code:    0,
					Message: ""}},
			"expectedbuy": &component.Trade{
				Credits: 123,
				Order: component.TradeOrder{
					Good:         "FUEL",
					PricePerUnit: 12,
					Quantity:     2,
					Total:        5},
				Ship: component.Ship{},
				Error: component.Error{
					Code:    0,
					Message: ""}}}}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	proxy := mocks.NewMockProxy(ctrl)

	for _, uc := range useCases {
		proxy.EXPECT().GetShipInfo().Return([]byte(fmt.Sprintf("%v", uc["detailsResponse"])), nil)
		proxy.EXPECT().GetShipInfo().Return([]byte(fmt.Sprintf("%v", uc["detailsResponse"])), nil)
		proxy.EXPECT().GetMarketplaceProducts(uc["location"]).
			Return([]byte(fmt.Sprintf("%v", uc["marketplaceResponse"])), nil)
		proxy.EXPECT().SellGood(fmt.Sprintf("%v", "Good001"), uc["quantity"].(int)).
			Return([]byte(fmt.Sprintf("%v", uc["sellGoodResponse"])), nil)
		proxy.EXPECT().BuyGood(fmt.Sprintf("%v", uc["good"]), uc["quantity"].(int)).
			Return([]byte(fmt.Sprintf("%v", uc["buyGoodResponse"])), nil)
		ship, err := component.NewShipCustomProxy(
			context.TODO(),
			trace.NewNoopTracerProvider().Tracer(""),
			proxy,
			fmt.Sprintf("%v", uc["id"]))
		if err != nil {
			t.Fail()
		}

		err = ship.ForceBuyFuel(context.TODO(), uc["quantity"].(int))
		if err != nil {
			t.Log(err)
			t.Fail()
		}

		// if actualSell == nil {
		// 	t.Fatal("actualSell is nil")
		// } else if !reflect.DeepEqual(actualSell, uc["expectedSell"].(*component.Trade)) {
		// 	t.Fatalf("\n%+v\n%+v\n", actualSell, uc["expectedSell"].(*component.Trade))
		// }
	}
}
