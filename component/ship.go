package component

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/otaviokr/spacetraders-ship/kafka"
	"github.com/otaviokr/spacetraders-ship/web"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gopkg.in/yaml.v3"
)

// Ship contains the essential information to authenticate in the game, but also to map the response from ship details.
type Ship struct {
	tracer trace.Tracer
	// webProxy web.Proxy
	webProxy kafka.Proxy
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
// func NewShip(ctx context.Context, tracer trace.Tracer, id, token string) (*Ship, error) {
func NewShip(
	ctx context.Context, tracer trace.Tracer,
	id, connectionType, connectionString,
	topicRead string, partitionRead int,
	topicWrite string, partitionWrite int) (*Ship, error) {
	// return NewShipCustomProxy(ctx, tracer, web.NewWebProxy(id, token), id, token)
	return NewShipCustomProxy(ctx, tracer,
		kafka.NewKafkaProxy(
			ctx,
			id,
			connectionType,
			connectionString,
			topicRead,
			partitionRead,
			topicWrite,
			partitionWrite),
		id)
}

// NewShipCustomProxy creates a new instance of component.Ship, using a provided custom web.WebProxy.
// func NewShipCustomProxy(ctx context.Context, tracer trace.Tracer, proxy web.Proxy, id, token string) (*Ship, error) {
func NewShipCustomProxy(ctx context.Context, tracer trace.Tracer, proxy kafka.Proxy, id string) (*Ship, error) {
	shipCtx, span := tracer.Start(ctx, "Activate Ship")
	defer span.End()
	ship := Ship{
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

	log.Println("Getting ship details...")
	data, err := s.webProxy.GetShipInfo()
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Key("data").String(string(data)))
		span.SetStatus(codes.Error, err.Error())
		log.Println("could not get response:", err)
		return err
	} else {
		s.Error.Code = -1
		s.Error.Message = ""
	}

	decoder := yaml.NewDecoder(bytes.NewReader(data))
	if err = decoder.Decode(&s); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		log.Println("could not decode response:", err)
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
		if found == nil {
			err = fmt.Errorf("UNEXPECTED ERROR (%d): %s", fp.Error.Code, fp.Error.Message)
			log.Println(err.Error())
			span.RecordError(err)
		} else {
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
