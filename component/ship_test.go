package component_test

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/otaviokr/spacetraders-ship/component"
	"github.com/otaviokr/spacetraders-ship/mocks"
	"go.opentelemetry.io/otel/trace"
)

func TestNewShip(t *testing.T) {
	useCases := map[string]map[string]interface{}{
		"uc1": {
			"id":              "id0001",
			"tracer":          trace.NewNoopTracerProvider().Tracer(""),
			"detailsResponse": "{\"ship\":{\"id\":\"id0001\",\"location\":\"XV-OS\",\"x\":52,\"y\":3,\"cargo\":[{\"good\":\"FUEL\",\"quantity\":14,\"totalVolume\":14}],\"spaceAvailable\":286,\"type\":\"GR-MK-II\",\"class\":\"MK-II\",\"maxCargo\":300,\"loadingSpeed\":500,\"speed\":1,\"manufacturer\":\"Gravager\",\"plating\":10,\"weapons\":5}}",
			"expected": &component.Ship{
				Details: component.ShipDetails{
					Id:             "id0001",
					Location:       "XV-OS",
					X:              52,
					Y:              3,
					Cargo:          []component.ShipCargo{{Good: "FUEL", Quantity: 14, TotalVolume: 14}},
					SpaceAvailable: 286,
					Type:           "GR-MK-II",
					Class:          "MK-II",
					MaxCargo:       300,
					LoadingSpeed:   500,
					Speed:          1,
					Manufacturer:   "Gravager",
					Plating:        10,
					Weapons:        5},
				Error: component.Error{Code: -1}}}}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	proxy := mocks.NewMockProxy(ctrl)

	for _, uc := range useCases {
		proxy.EXPECT().GetShipInfo().Return([]byte(fmt.Sprintf("%v", uc["detailsResponse"])), nil)
		actual, err := component.NewShipCustomProxy(
			context.TODO(),
			trace.NewNoopTracerProvider().Tracer(""),
			proxy,
			fmt.Sprintf("%v", uc["id"]))
		if err != nil {
			t.Fail()
		}

		if !reflect.DeepEqual(actual.Details, uc["expected"].(*component.Ship).Details) {
			t.Fatalf("\n%+v\n%+v\n", actual.Details, uc["expected"].(*component.Ship).Details)
		}

		if !reflect.DeepEqual(actual.Error, uc["expected"].(*component.Ship).Error) {
			t.Fatalf("\n%+v\n%+v\n", actual.Error, uc["expected"].(*component.Ship).Error)
		}
	}
}

func TestNewShipFailed(t *testing.T) {
	errorMessage := "This is error"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	proxy := mocks.NewMockProxy(ctrl)

	proxy.EXPECT().GetShipInfo().Return(nil, fmt.Errorf(errorMessage))
	_, err := component.NewShipCustomProxy(
		context.TODO(),
		trace.NewNoopTracerProvider().Tracer(""),
		proxy,
		"irrelevant_id")
	if err == nil {
		t.Fatalf("\ngot no error, but expected error: %s", err.Error())
	} else if strings.Compare(err.Error(), errorMessage) != 0 {
		t.Fatalf("\nACTUAL: %s\nEXPECT: %s\n", err.Error(), errorMessage)
	}
}

// func TestNewShipFailed(t *testing.T) {
// 	errorMessage := "This is error"
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	proxy := mocks.NewMockProxy(ctrl)

// 	for _, uc := range useCases {
// 		proxy.EXPECT().GetShipInfo().Return([]byte(fmt.Sprintf("%v", uc["detailsResponse"])), nil)
// 		proxy.EXPECT().GetShipInfo().Return([]byte(fmt.Sprintf("%v", uc["detailsResponse"])), nil)
// 		proxy.EXPECT().
// 			GetMarketplaceProducts(fmt.Sprintf("%v", uc["location"])).
// 			Return([]byte(fmt.Sprintf("%v", uc["marketResponse"])), nil)
// 		ship, err := component.NewShipCustomProxy(
// 			context.TODO(),
// 			trace.NewNoopTracerProvider().Tracer(""),
// 			proxy,
// 			fmt.Sprintf("%v", uc["id"]),
// 			fmt.Sprintf("%v", uc["token"]))
// 		if err != nil {
// 			t.Fail()
// 		}

// 		actualMarketplace, actualProducts, err := ship.GetMarketplaceProducts(context.TODO())
// 		if err != nil {
// 			t.Fail()
// 		}

// 		if !reflect.DeepEqual(actualMarketplace, uc["expectedMarketplace"].(*component.Marketplace)) {
// 			t.Fatalf("\nACTUAL: %+v\nEXPECT: %+v\n", actualMarketplace, uc["expectedMarketplace"].(*component.Marketplace))
// 		}

// 		if !reflect.DeepEqual(actualProducts, uc["expectedProducts"].(*map[string]component.Product)) {
// 			t.Fatalf("\n%+v\n%+v\n", actualProducts, uc["expectedProducts"].(*map[string]component.Product))
// 		}
// 	}
// }

func TestFly(t *testing.T) {
	useCases := map[string]map[string]interface{}{
		"uc1": {
			"id":                 "id0001",
			"tracer":             trace.NewNoopTracerProvider().Tracer(""),
			"location":           "Local0001",
			"detailsResponse":    "{\"ship\":{\"id\":\"id0001\",\"location\":\"Local0001\",\"x\":52,\"y\":3,\"cargo\":[{\"good\":\"FUEL\",\"quantity\":14,\"totalVolume\":14}],\"spaceAvailable\":286,\"type\":\"GR-MK-II\",\"class\":\"MK-II\",\"maxCargo\":300,\"loadingSpeed\":500,\"speed\":1,\"manufacturer\":\"Gravager\",\"plating\":10,\"weapons\":5}}",
			"flightPlanResponse": "{\"flightPlan\": {\"arrivesAt\": \"2021-05-13T18:41:24.963Z\",\"createdAt\": \"2021-05-13T18:40:23.003Z\",\"departure\": \"OE-PM-TR\",\"destination\": \"OE-PM\",\"distance\": 1,\"fuelConsumed\": 1,\"fuelRemaining\": 18,\"id\": \"flightplanid0001\",\"shipId\": \"id0001\",\"terminatedAt\": null,\"timeRemainingInSeconds\": 1}}",
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
		proxy.EXPECT().SetNewFlightPlan(
			fmt.Sprintf("%v", uc["destination"])).Return([]byte(fmt.Sprintf("%v", uc["flightPlanResponse"])),
			nil)
		ship, err := component.NewShipCustomProxy(
			context.TODO(),
			trace.NewNoopTracerProvider().Tracer(""),
			proxy,
			fmt.Sprintf("%v", uc["id"]))
		if err != nil {
			t.Fail()
		}

		err = ship.Fly(context.TODO(), fmt.Sprintf("%v", uc["destination"]))
		if err != nil {
			t.Log(err)
			t.Fail()
		}
	}
}

func TestNewFlightPlan(t *testing.T) {
	useCases := map[string]map[string]interface{}{
		"uc1": {
			"id":                 "id0001",
			"tracer":             trace.NewNoopTracerProvider().Tracer(""),
			"location":           "Local0001",
			"detailsResponse":    "{\"ship\":{\"id\":\"id0001\",\"location\":\"Local0001\",\"x\":52,\"y\":3,\"cargo\":[{\"good\":\"FUEL\",\"quantity\":14,\"totalVolume\":14}],\"spaceAvailable\":286,\"type\":\"GR-MK-II\",\"class\":\"MK-II\",\"maxCargo\":300,\"loadingSpeed\":500,\"speed\":1,\"manufacturer\":\"Gravager\",\"plating\":10,\"weapons\":5}}",
			"flightPlanResponse": "{\"flightPlan\": {\"arrivesAt\": \"2021-05-13T18:41:24.963Z\",\"createdAt\": \"2021-05-13T18:40:23.003Z\",\"departure\": \"OE-PM-TR\",\"destination\": \"OE-PM\",\"distance\": 1,\"fuelConsumed\": 1,\"fuelRemaining\": 18,\"id\": \"flightplanid0001\",\"shipId\": \"id0001\",\"terminatedAt\": null,\"timeRemainingInSeconds\": 1}}",
			"expectedFlightPlan": &component.FlightPlan{
				Details: component.FlightPlanDetails{
					ArrivesAt:              "2021-05-13T18:41:24.963Z",
					CreatedAt:              "2021-05-13T18:40:23.003Z",
					Departure:              "OE-PM-TR",
					Destination:            "OE-PM",
					Distance:               1,
					FuelConsumed:           1,
					FuelRemaining:          18,
					Id:                     "flightplanid0001",
					ShipId:                 "id0001",
					TerminatedAt:           "",
					TimeRemainingInSeconds: 1},
				Error: component.Error{
					Code:    0,
					Message: ""}}}}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	proxy := mocks.NewMockProxy(ctrl)

	for _, uc := range useCases {
		proxy.EXPECT().GetShipInfo().Return([]byte(fmt.Sprintf("%v", uc["detailsResponse"])), nil)
		// proxy.EXPECT().GetShipInfo().Return([]byte(fmt.Sprintf("%v", uc["detailsResponse"])), nil)
		proxy.EXPECT().SetNewFlightPlan(
			fmt.Sprintf("%v", uc["destination"])).Return([]byte(fmt.Sprintf("%v", uc["flightPlanResponse"])),
			nil)
		ship, err := component.NewShipCustomProxy(
			context.TODO(),
			trace.NewNoopTracerProvider().Tracer(""),
			proxy,
			fmt.Sprintf("%v", uc["id"]))
		if err != nil {
			t.Fail()
		}

		actualFlightPlan, err := ship.NewFlightPlan(context.TODO(), fmt.Sprintf("%v", uc["destination"]))
		if err != nil {
			t.Log(err)
			t.Fail()
		}

		if !reflect.DeepEqual(actualFlightPlan, uc["expectedFlightPlan"].(*component.FlightPlan)) {
			t.Fatalf("\n%+v\n%+v\n", actualFlightPlan, uc["expectedFlightPlan"].(*component.FlightPlan))
		}
	}
}

func TestGetFlightPlan(t *testing.T) {
	useCases := map[string]map[string]interface{}{
		"uc1": {
			"id":                 "id0001",
			"tracer":             trace.NewNoopTracerProvider().Tracer(""),
			"flightPlanId":       "flightPlanId0001",
			"detailsResponse":    "{\"ship\":{\"id\":\"id0001\",\"flightPlanId\":\"flightPlanId0001\",\"cargo\":[{\"good\":\"FUEL\",\"quantity\":14,\"totalVolume\":14}],\"spaceAvailable\":286,\"type\":\"GR-MK-II\",\"class\":\"MK-II\",\"maxCargo\":300,\"loadingSpeed\":500,\"speed\":1,\"manufacturer\":\"Gravager\",\"plating\":10,\"weapons\":5}}",
			"flightPlanResponse": "{\"flightPlan\": {\"arrivesAt\": \"2021-05-13T18:41:24.963Z\",\"createdAt\": \"2021-05-13T18:40:23.003Z\",\"departure\": \"OE-PM-TR\",\"destination\": \"OE-PM\",\"distance\": 1,\"fuelConsumed\": 1,\"fuelRemaining\": 18,\"id\": \"flightplanid0001\",\"shipId\": \"id0001\",\"terminatedAt\": null,\"timeRemainingInSeconds\": 1}}",
			"expectedFlightPlan": &component.FlightPlan{
				Details: component.FlightPlanDetails{
					ArrivesAt:              "2021-05-13T18:41:24.963Z",
					CreatedAt:              "2021-05-13T18:40:23.003Z",
					Departure:              "OE-PM-TR",
					Destination:            "OE-PM",
					Distance:               1,
					FuelConsumed:           1,
					FuelRemaining:          18,
					Id:                     "flightplanid0001",
					ShipId:                 "id0001",
					TerminatedAt:           "",
					TimeRemainingInSeconds: 1},
				Error: component.Error{
					Code:    0,
					Message: ""}}}}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	proxy := mocks.NewMockProxy(ctrl)

	for _, uc := range useCases {
		proxy.EXPECT().GetShipInfo().Return([]byte(fmt.Sprintf("%v", uc["detailsResponse"])), nil)
		proxy.EXPECT().GetShipInfo().Return([]byte(fmt.Sprintf("%v", uc["detailsResponse"])), nil)
		proxy.EXPECT().GetFlightPlan(fmt.Sprintf("%v", uc["flightPlanId"])).
			Return([]byte(fmt.Sprintf("%v", uc["flightPlanResponse"])), nil)
		ship, err := component.NewShipCustomProxy(
			context.TODO(),
			trace.NewNoopTracerProvider().Tracer(""),
			proxy,
			fmt.Sprintf("%v", uc["id"]))
		if err != nil {
			t.Fail()
		}

		actualFlightPlan, err := ship.GetFlightPlan(context.TODO())
		if err != nil {
			t.Log(err)
			t.Fail()
		}

		if !reflect.DeepEqual(actualFlightPlan, uc["expectedFlightPlan"].(*component.FlightPlan)) {
			t.Fatalf("\n%+v\n%+v\n", actualFlightPlan, uc["expectedFlightPlan"].(*component.FlightPlan))
		}
	}
}

func TestShipSell(t *testing.T) {
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

func TestShipBuy(t *testing.T) {
}

func TestShipForceBuyFuel(t *testing.T) {
}
