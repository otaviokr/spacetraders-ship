package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/otaviokr/spacetraders-ship/component"
	"github.com/otaviokr/spacetraders-ship/web"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"

	traceSdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	TracerName = "spacetrader-ship"
)

// main is just the starting point, but we keep just the bare minimum here.
//
// https://pace.dev/blog/2020/02/12/why-you-shouldnt-use-func-main-in-golang-by-mat-ryer.html
func main() {
	token := os.Getenv("USER_TOKEN")
	shipId := os.Getenv("SHIP_ID")
	filePath := os.Getenv("CONFIG_FILE_PATH")
	jaegerUrl := os.Getenv("JAEGER_URL")

	metricsPort := os.Getenv("METRICS_PORT")
	if len(metricsPort) < 1 {
		metricsPort = "9090"
	}

	// This is function to expose the metrics to Prometheus.
	go exposeMetrics(metricsPort)

	// The main loop is actually inside the run function.
	if err := run(token, shipId, filePath, jaegerUrl); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

// run contains the main loop of the program. It will collect data from the Space Traders game and
// expose them to Prometheus.
func run(token, shipId, configFilePath, jaegerUrl string) error {
	bgCtx := context.Background()
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerUrl)))
	if err != nil {
		log.Fatal(err)
	}

	tp := traceSdk.NewTracerProvider(
		traceSdk.WithBatcher(exp),
		traceSdk.WithResource(newResource()))
	defer func() {
		if err := tp.Shutdown(bgCtx); err != nil {
			log.Fatal(err)
		}
	}()

	otel.SetTracerProvider(tp)
	tracer := otel.Tracer(TracerName)

	// Defining the ship we will use.
	ship, err := component.NewShip(bgCtx, tracer, shipId, token)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Registered new ship wth ID %s\n", shipId)

	if len(ship.Details.FlightPlanId) > 0 {
		flightPlan, err := ship.GetFlightPlan(bgCtx)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Ship is en route: %ds to reach %s (%+v)\n",
			flightPlan.Details.TimeRemainingInSeconds,
			flightPlan.Details.Destination,
			flightPlan.Details.ArrivesAt)
		time.Sleep(time.Duration(flightPlan.Details.TimeRemainingInSeconds) * time.Second)
	}

	// Trading routes are supposed to be cyclical, so we are locked in an eternal loop.
	// If the ship is not in the right location when we start the application, the first step
	// is to take the ship to the right location and start from there.
	// FIXME we need to catch Ctrl+C and other termination commands to do a clean stop!
	for {
		// Read the trading route from file.
		routes, err := component.ReadRouteFile(configFilePath)
		if err != nil {
			log.Fatal(err)
		}

		totalStops := len(routes.Route)
		rootCtx, span := tracer.Start(
			context.Background(),
			"Route",
			trace.WithAttributes(
				attribute.Key("ship.id").String(shipId),
				attribute.Key("ship.route.total_stops").Int(totalStops)))
		log.Printf("Starting new trading route cycle with %d stops\n", totalStops)

		// Each step of the route requires the same procedure:
		//	- If we are not at location, we travel to it;
		//	- Sell the goods;
		// 	- Buy the goods (including FUEL).
		for i, route := range routes.Route {
			routeCtx, routeSpan := tracer.Start(
				rootCtx,
				"Sprint",
				trace.WithAttributes(
					attribute.Key("route.location").String(route.Station),
					attribute.Key("route.sell").Int(len(route.Sell)),
					attribute.Key("route.buy").Int(len(route.Buy))))
			log.Printf("Route Step %d/%d: %s\n", i+1, totalStops, route.Station)

			if err = ship.GetDetails(routeCtx); err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			log.Printf("Ship current status: Flight Plan(%s) Location(%s) Cargo(%+v)\n",
				ship.Details.FlightPlanId,
				ship.Details.Location,
				ship.Details.Cargo)

			if ship.Details.Location != route.Station {
				log.Printf("Setting new coordinates: %s to %s\n", ship.Details.Location, route.Station)
				err = ship.Fly(routeCtx, route.Station)
				if err != nil {
					span.RecordError(err)
					span.SetStatus(codes.Error, err.Error())
					panic(err)
				}
			}

			if ship.Details.Location == route.Station {
				log.Printf("Ship reached %s\n", ship.Details.Location)
				dockCtx, dockSpan := tracer.Start(
					routeCtx,
					"Docked",
					trace.WithAttributes(
						attribute.Key("Location").String(ship.Details.Location)))
				err = ship.DoCommerce(dockCtx, route.Sell, route.Buy)
				if err != nil {
					dockSpan.RecordError(err)
					dockSpan.SetStatus(codes.Error, err.Error())
				}
				dockSpan.End()
			}
			routeSpan.End()
		}
		log.Println("Trading route finished. Receiving new orders...")
		span.End()
		web.TradeCycles.
			WithLabelValues(shipId).
			Inc()
	}
}

// exposeMetrics is a very simple web server that Prometheus can access to collect the metrics.
//
// port is the port where the web server is listening.
func exposeMetrics(port string) {
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

// newResource returns a resource describing this application to be used in the traces for Jaeger.
func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("spacetraders-ship"),
			semconv.ServiceVersionKey.String("v0.1.0"),
			attribute.String("environment", "demo"),
		),
	)
	return r
}
