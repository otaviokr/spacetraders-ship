package component

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/otaviokr/spacetraders-ship/web"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gopkg.in/yaml.v3"
)

// Ship contains the essential information to authenticate in the game, but also to map the response from ship details.
type Ship struct {
	token    string
	tracer   trace.Tracer
	webProxy web.Proxy
	Details  ShipDetails `yaml:"ship"`
	Error    Error       `yaml:"error"`
}

// ShipDetails is the response from the Ship Detail API.
type ShipDetails struct {
	Id             string      `yaml:"id"`
	FlightPlanId   string      `yaml:"flightPlanId"`
	Location       string      `yaml:"location"`
	X              int         `yaml:"x"`
	Y              int         `yaml:"y"`
	Cargo          []ShipCargo `yaml:"cargo"`
	SpaceAvailable int         `yaml:"spaceAvailable"`
	Type           string      `yaml:"type"`
	Class          string      `yaml:"class"`
	MaxCargo       int         `yaml:"maxCargo"`
	LoadingSpeed   int         `yaml:"loadingSpeed"`
	Speed          int         `yaml:"speed"`
	Manufacturer   string      `yaml:"manufacturer"`
	Plating        int         `yaml:"plating"`
	Weapons        int         `yaml:"weapons"`
}

// ShipCargo contains the details about the products stored in the ship cargo.
type ShipCargo struct {
	Good        string `yaml:"good"`
	Quantity    int    `yaml:"quantity"`
	TotalVolume int    `yaml:"totalVolume"`
}

// NewShip creates a new instance of component.Ship.
func NewShip(ctx context.Context, tracer trace.Tracer, id, token string) (*Ship, error) {
	return NewShipCustomProxy(ctx, tracer, web.NewWebProxy(id, token), id, token)
}

// NewShipCustomProxy creates a new instance of component.Ship, using a provided custom web.WebProxy.
func NewShipCustomProxy(ctx context.Context, tracer trace.Tracer, proxy web.Proxy, id, token string) (*Ship, error) {
	shipCtx, span := tracer.Start(ctx, "Activate Ship")
	defer span.End()
	ship := Ship{
		token:    token,
		tracer:   tracer,
		webProxy: proxy,
		Details: ShipDetails{
			Id: id}}
	if err := ship.GetDetails(shipCtx); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return nil, err
	}
	return &ship, nil
}

// GetDetails will get the ship details from the game.
func (s *Ship) GetDetails(ctx context.Context) error {
	_, span := s.tracer.Start(
		ctx,
		"Get Ship Details",
		trace.WithAttributes(
			attribute.Key("ship.id").String(s.Details.Id)))
	defer span.End()

	data, err := s.webProxy.GetShipInfo()
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Key("data").String(string(data)))
		span.SetStatus(codes.Error, err.Error())
		return err
	} else {
		s.Error.Code = -1
		s.Error.Message = ""
	}

	decoder := yaml.NewDecoder(bytes.NewReader(data))
	if err = decoder.Decode(&s); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	if len(s.Error.Message) > 0 {
		// Error from the server, we should still report it.
		err = fmt.Errorf("ERROR FROM SERVER (%d): %s", s.Error.Code, s.Error.Message)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

// GetMarketplaceProducts will fetch the products that are available to be traded in the current marketplace.
func (s *Ship) GetMarketplaceProducts(ctx context.Context) (*Marketplace, *map[string]Product, error) {
	newCtx, span := s.tracer.Start(
		ctx,
		"Get Products from Marketplace",
		trace.WithAttributes(
			attribute.Key("ship.id").String(s.Details.Id),
			attribute.Key("location").String(s.Details.Location)))
	defer span.End()

	err := s.GetDetails(newCtx)
	if err != nil {
		span.RecordError(err)
	}

	data, err := s.webProxy.GetMarketplaceProducts(s.Details.Location)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, nil, err
	}

	var m Marketplace
	decoder := json.NewDecoder(bytes.NewReader(data))
	if err = decoder.Decode(&m); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, nil, err
	}

	if len(m.Error.Message) > 0 {
		// Error from the server, we should still report it.
		err = fmt.Errorf("ERROR FROM SERVER (%d): %s", s.Error.Code, s.Error.Message)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, nil, err
	}

	p := map[string]Product{}
	for _, product := range m.Products {
		p[product.Symbol] = product
	}

	return &m, &p, nil
}

// Fly will set the FlightPlan to a new destination, and wait until the flight is finished before returning from the method.
func (s *Ship) Fly(ctx context.Context, destination string) error {
	flyCtx, flySpan := s.tracer.Start(
		ctx,
		"Fly",
		trace.WithAttributes(
			attribute.Key("ship.id").String(s.Details.Id),
			attribute.Key("Destination").String(destination)))
	defer flySpan.End()

	flightPlan, err := s.NewFlightPlan(flyCtx, destination)
	if err != nil {
		flySpan.RecordError(err)
		flySpan.SetStatus(codes.Error, err.Error())
		return err
	}

	flySpan.AddEvent(
		"Flight plan ready",
		trace.WithAttributes(
			attribute.Key("flightplan.id").String(flightPlan.Details.Id),
			attribute.Key("flightplan.remaining").Int(flightPlan.Details.TimeRemainingInSeconds),
			attribute.Key("flightplan.departure").String(flightPlan.Details.Departure),
			attribute.Key("flightplan.destination").String(flightPlan.Details.Destination),
			attribute.Key("flightplan.fuel.consumed").Int(flightPlan.Details.FuelConsumed),
			attribute.Key("flightplan.distance").Int(flightPlan.Details.Distance)))
	web.FuelConsumed.WithLabelValues(s.Details.Id).Add(float64(flightPlan.Details.FuelConsumed))

	log.Printf("Flight Plan defined to %s in %ds (%+v)\n",
		flightPlan.Details.Destination,
		flightPlan.Details.TimeRemainingInSeconds,
		flightPlan.Details.ArrivesAt)

	time.Sleep(time.Duration(flightPlan.Details.TimeRemainingInSeconds+5) * time.Second)

	flySpan.AddEvent("Check flight status")
	err = s.GetDetails(flyCtx)
	if err != nil {
		flySpan.RecordError(err)
		flySpan.SetStatus(codes.Error, err.Error())
		return err
	}

	for len(s.Details.FlightPlanId) > 0 {
		flightPlan, err = s.GetFlightPlan(flyCtx)
		if err != nil {
			flySpan.RecordError(err)
			flySpan.SetStatus(codes.Error, err.Error())
			return err
		}

		flySpan.AddEvent(
			"Extending flight time",
			trace.WithAttributes(
				attribute.Key("flightplan.id").String(flightPlan.Details.Id),
				attribute.Key("flightplan.remaining").Int(flightPlan.Details.TimeRemainingInSeconds),
				attribute.Key("flightplan.destination").String(flightPlan.Details.Destination)))
		time.Sleep(time.Duration(flightPlan.Details.TimeRemainingInSeconds) * time.Second)

		flySpan.AddEvent("Update flight status")
		err = s.GetDetails(flyCtx)
		if err != nil {
			flySpan.RecordError(err)
			flySpan.SetStatus(codes.Error, err.Error())
			return err
		}
	}

	return nil
}

// DoCommerce places the buy and sell orders to the game.
func (s *Ship) DoCommerce(ctx context.Context, sell, buy map[string]int) error {
	newCtx, span := s.tracer.Start(ctx, "Commerce")
	defer span.End()

	_, products, err := s.GetMarketplaceProducts(newCtx)
	if err != nil {
		span.RecordError(err)
	}

	for _, good := range s.Details.Cargo {
		if _, ok := (*products)[good.Good]; ok {
			if sell[good.Good] == -1 {
				log.Printf("Selling the whole lot of %s: %d", good.Good, good.Quantity)
				sell[good.Good] = good.Quantity
			} else {
				log.Printf("Selling pre-defined lot of %s: %d", good.Good, sell[good.Good])
			}
		}
	}

	err = s.SellAll(newCtx, sell, *products)
	if err != nil {
		span.RecordError(err)
	}

	err = s.GetDetails(newCtx)
	if err != nil {
		span.RecordError(err)
	}

	for _, good := range s.Details.Cargo {
		if _, ok := (*products)[good.Good]; ok {
			if _, ok := buy[good.Good]; ok {
				log.Printf(
					"Buying necessary to complete lot of %s: %d (total %d)",
					good.Good, buy[good.Good]-good.Quantity, buy[good.Good])
				buy[good.Good] -= good.Quantity
			}
		}
	}

	err = s.BuyAll(newCtx, buy, *products)
	if err != nil {
		span.RecordError(err)
	}

	return nil
}

// SellAll is wrapper to sell all units of products in the provided list.
func (s *Ship) SellAll(ctx context.Context, sell map[string]int, marketplace map[string]Product) error {
	sellCtx, sellSpan := s.tracer.Start(
		ctx,
		"Sell goods",
		trace.WithAttributes(
			attribute.Key("Goods to sell").Int(len(sell))))
	defer sellSpan.End()

	for good, quantity := range sell {
		if _, ok := marketplace[good]; ok {
			log.Printf("Selling lot of %s: %d\n", good, quantity)

			_, err := s.Sell(sellCtx, good, quantity)
			if err != nil {
				sellSpan.RecordError(err)
			}
		} else {
			sellSpan.AddEvent(
				"Cannot sell product",
				trace.WithAttributes(
					attribute.Key("product").String(good)))
		}
	}
	return nil
}

// BuyAll is wrapper to buy the products in the provided list.
func (s *Ship) BuyAll(ctx context.Context, buy map[string]int, marketplace map[string]Product) error {
	buyCtx, buySpan := s.tracer.Start(
		ctx,
		"Buy goods",
		trace.WithAttributes(
			attribute.Key("Goods to buy").Int(len(buy))))
	defer buySpan.End()

	for good, quantity := range buy {
		if _, ok := marketplace[good]; ok {
			if good == "FUEL" {
				log.Printf("Priority purchase of %s: %d\n", good, quantity)
				err := s.ForceBuyFuel(ctx, quantity)
				if err != nil {
					buySpan.RecordError(err)
				}
			} else {
				log.Printf("Buying lot of %s: %d\n", good, quantity)
				actualQuantity := quantity
				if s.Details.SpaceAvailable < quantity*marketplace[good].VolumePerUnit {
					actualQuantity = int(s.Details.SpaceAvailable / marketplace[good].VolumePerUnit)
					log.Printf("Low cargo space! Available: %d / Wanted: %d (Volume Per Unit: %d)\n", s.Details.SpaceAvailable, actualQuantity, marketplace[good].VolumePerUnit)
				} else {
					for _, product := range s.Details.Cargo {
						if product.Good == good {
							actualQuantity = quantity - product.Quantity
							log.Printf("Completing lot of %s: Original(%d) / Additional(%d)\n", good, product.Quantity, actualQuantity)
						}
					}
				}

				// FIXME
				_, err := s.Buy(buyCtx, good, actualQuantity)
				if err != nil {
					buySpan.RecordError(err)
				}
			}
		}
	}
	return nil
}

// NewFlightPlan sets a new destination for the ship to fly to.
func (s *Ship) NewFlightPlan(ctx context.Context, destination string) (*FlightPlan, error) {
	newCtx, span := s.tracer.Start(
		ctx,
		"New Flight Plan",
		trace.WithAttributes(
			attribute.Key("destination").String(destination)))
	defer span.End()

	data, err := s.webProxy.SetNewFlightPlan(destination)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	var fp FlightPlan
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	if err = decoder.Decode(&fp); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	if fp.Error.Code > 0 {
		re := regexp.MustCompile(InsufficientFuelRegex)
		found := re.FindAllStringSubmatch(fp.Error.Message, 1)
		fuel, err := strconv.Atoi(found[0][1])
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}

		// Error from the server, we should still report it.
		err = fmt.Errorf("ERROR FROM SERVER (%d): %s", fp.Error.Code, fp.Error.Message)
		span.RecordError(err)
		// span.SetStatus(codes.Error, err.Error())

		err = s.ForceBuyFuel(newCtx, fuel)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}

		return s.NewFlightPlan(newCtx, destination)
	}

	return &fp, nil
}

// GetFlightPlan retrieves current flight plan, if any.
func (s *Ship) GetFlightPlan(ctx context.Context) (*FlightPlan, error) {
	if err := s.GetDetails(ctx); err != nil {
		return nil, err
	}

	var fp FlightPlan

	if len(s.Details.FlightPlanId) < 1 {
		return nil, nil
	}

	log.Printf("Flight Plan found: %s\n", s.Details.FlightPlanId)
	data, err := s.webProxy.GetFlightPlan(s.Details.FlightPlanId)
	if err != nil {
		return nil, err
	}

	decoder := yaml.NewDecoder(bytes.NewReader(data))
	if err = decoder.Decode(&fp); err != nil {
		return nil, err
	}

	if len(fp.Error.Message) > 0 {
		// Error from the server, we should still report it.
		err = fmt.Errorf("ERROR FROM SERVER (%d): %s", fp.Error.Code, fp.Error.Message)
		return nil, err
	}

	return &fp, nil
}

// Sell sends a sell order to the game.
func (s *Ship) Sell(ctx context.Context, good string, quantity int) (*Trade, error) {
	return s.trade(ctx, "sell", good, quantity)
}

// Buy sends a buy order to the game.
func (s *Ship) Buy(ctx context.Context, good string, quantity int) (*Trade, error) {
	return s.trade(ctx, "buy", good, quantity)
}

// ForceBuyFuel will prioritize the purchase of fuel, selling goods if necessary.
func (s *Ship) ForceBuyFuel(ctx context.Context, fuel int) error {
	newCtx, span := s.tracer.Start(
		ctx,
		"Buy emergency fuel",
		trace.WithAttributes(
			attribute.Key("Extra fuel required").Int(fuel)))
	defer span.End()

	if s.Details.SpaceAvailable > fuel {
		if _, err := s.Buy(newCtx, "FUEL", fuel); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return err
		}

		return nil
	}

	_, products, err := s.GetMarketplaceProducts(newCtx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	log.Printf("Priority purchase of %s issued: %d\n", "FUEL", fuel)

	if s.Details.SpaceAvailable > fuel {
		log.Printf("Enough space in cargo bay. Buying fuel...")
		if _, err = s.Buy(newCtx, "FUEL", fuel); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return err
		}
		return nil
	}

	log.Printf("Not enough free room in cargo bay. Selling other goods to make room: %d", fuel)
	remaining := fuel
	for _, cargo := range s.Details.Cargo {
		if _, ok := (*products)[cargo.Good]; ok {
			log.Printf("Selling %v: max room to free (%d), need(%d)\n", cargo.Good, cargo.TotalVolume, remaining)
			if cargo.TotalVolume > remaining {
				if _, err = s.Sell(newCtx, cargo.Good, remaining/(*products)[cargo.Good].VolumePerUnit); err != nil {
					span.RecordError(err)
					span.SetStatus(codes.Error, err.Error())
					return err
				}

				if _, err = s.Buy(newCtx, "FUEL", fuel); err != nil {
					span.RecordError(err)
					span.SetStatus(codes.Error, err.Error())
					return err
				}

				log.Println("Ship is re-fueled.")
				return nil
			}

			if _, err = s.Sell(newCtx, cargo.Good, cargo.Quantity); err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				return err
			}
			remaining -= cargo.Quantity
		}
	}
	log.Println("ALERT! Not enough room to fuel!")
	return fmt.Errorf("could not purchase fuel - impossible to sell products at location?")
}

// trade is generic call to buy and sell products in game.
func (s *Ship) trade(ctx context.Context, action, good string, quantity int) (*Trade, error) {
	_, span := s.tracer.Start(
		ctx,
		"Trade goods",
		trace.WithAttributes(
			attribute.Key("action").String(action),
			attribute.Key("good").String(good),
			attribute.Key("quantity").Int(quantity)))
	defer span.End()

	var data []byte
	var err error
	switch strings.ToLower(action) {
	case "sell":
		data, err = s.webProxy.SellGood(good, quantity)
	case "buy":
		data, err = s.webProxy.BuyGood(good, quantity)
	}
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	var operation Trade
	decoder := yaml.NewDecoder((bytes.NewReader(data)))
	if err = decoder.Decode(&operation); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	if len(operation.Error.Message) > 0 {
		// Error from the server, we should still report it.
		err = fmt.Errorf("ERROR FROM SERVER (%d): %s", operation.Error.Code, operation.Error.Message)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	switch strings.ToLower(action) {
	case "sell":
		web.MoneyEarned.
			WithLabelValues(s.Details.Id, operation.Order.Good).
			Add(float64(operation.Order.Total))
		web.GoodsSold.
			WithLabelValues(s.Details.Id, operation.Order.Good, operation.Ship.Details.Location).
			Add(float64(operation.Order.Quantity))
	case "buy":
		web.MoneySpent.
			WithLabelValues(s.Details.Id, operation.Order.Good).
			Add(float64(operation.Order.Total))
		web.GoodsBought.
			WithLabelValues(s.Details.Id, operation.Order.Good, operation.Ship.Details.Location).
			Add(float64(operation.Order.Quantity))
	}
	// web.UserCredits.Set(float64(operation.Credits))

	return &operation, nil
}
