# spacetraders-ship
Space Traders monitoring using Grafana, Prometheus and Golang

## Overview

This repository is part of a bigger solution to play [Space Traders](https://spacetraders.io) with a different approach: instead of taking the actions directly, you are the CEO of the company and don't have time for the micromanagement. As expected you are only interested in numbers, charts and cool graphics showing you how much money your company is profiting.

Using "real-world" applications to monitor "real-world" business, we will be able to see how the company is progressing.

Hopefully, in the future, it will be possible to send "management orders" to tweak your business (buy new ships, change routes etc.), but at the moment, just let your employees do their thing.

## What this application does

This application is written in Golang, and its sole purpose is to make requests to Space Traders API for information; but, differently from [spacetraders-inventory](https://github.com/otaviokr/spacetraders-inventory), this application represents a ship, and as a loyal employee, it will comply with the daily routine defined by you. So, it will report things like:

- How much fuel are you using?
- How much _material_ have you bought/sold?
- How long until you reach the next planet?

You should have an instance of this application running for each ship you own, meaning this can scale to a very big solution. For starters, you can manage it with simple solutions like docker-compose, Portainer etc., but as your game grows, you should consider to move to Kubernetes or something similar. (I hope to have instructions for this in the future, once my game becomes that big too!).

While the ship is flying to planet to planet, trading all sorts of good, it will report to Prometheus the data you expect in your reports in Grafana, all automatically and no action needed - remember you are the boss, you are not supposed to work!

Because you are the ultimate responsible for what you company does, you should provide the route the ship should work, what it should sell and what it should buy. This trading route can be as small or large as you want, and you will need to define it in a YAML file (there is a template/example in the etc/routes folders, if you need help to get started). Having a good tradding route is the difference between struggling with money and unstoppable success, so dedicate some time thinking what the ship should do.

## What this application does not do

It is not very interactive: the idea is to provide information, not an interface to interact directly with the Space Traders API. Other than defining the trading route for the ship, there is not much you can do at the moment.

I want to implement some more interaction, but I'm not yet sure if this will be implemented on this project or if this will be a project on its own.

## How to run this application

As mentioned, this is part of a bigger application, so this should not be started individually; it will be invoked, compiled, started and monitored together with the other components. For more information, refer to the [spacetraders-hq repo](https://github.com/otaviokr/spacetraders-hq).

That being said, you could run this application isolated using the docker-compose template file. Notice that it will start Prometheus, Grafana and Jaeger too. If this is needed or not, and the configuration details, will depend on your needs and environment.

If you have nothing like that already in your environment (no Prometheus, Grafana or Jaeger running), the docker-compose template file contains the defaults values and should work out of the box. So, for a quick-and-dirty test, you can simply rename the file to `docker-compose.yml` and run the following command:

```shell
docker-compose up -d --build
```

You should find the GUIs in the following URLs:

- Prometheus: http://localhost:9090/graph
- Application exported metrics: http://localhost:9091/metrics
- Grafana: http://localhost:3000
- Jaeger: http://localhost:16686

When you are done playing, just run `docker-compose down`.

## I have no idea what you are talking about

I will try to cover all the important topics and terms here, but if something is still not clear, let me know and I'll try to elaborate on it.

- **API**: simply put, an API is a set of URLs that you can access to reach an application running somewhere else. In our case, the Space Traders game is running on a server we don't have direct access, so if we want to interact with the game, we can use its API to ask information about ships, trade goods, fly to other planets etc.

- **Space Traders**: using their own words, it is "a unique multiplayer game built on a free Web API". I just want to clarify that I am not involved directly in the Space Traders game development or with the company behind it. I just find the idea amazing and wanted to build something on top of it. More info at https://spacetraders.io

- **Golang (aka Go)**: a programming language that generates a small, fast application that can be easily run in Docker. I'm using this language to collect all the required information from Space Traders and feed Prometheus.

- **Prometheus**: a monitoring application, to which you can provide metrics, generate alerts or be the datasource for other applications (e.g., Grafana). In our case, the application built with Golang feeds the data collected from Space Traders game to Prometheus, which will forward that data to Grafana. More info at https://prometheus.com

- **Grafana**: a solution to create dashboard with data coming from a wide variety of sources. In our case, we are getting data from Prometheus. I have created some dashboards with the data, but feel free to use those as a starting point for your own custom dashboards. More info at https://grafana.com

- **Jaeger**: an application to monitor and troubleshoot distributed tracing. In other words, the Golang application I created sends to Jaeger tracing information, indicating when specific processes start, when they end, and some details during its execution, in case we want to troubleshoot why a particular step in the execution took longer than expected. More info at https://www.jaegertracing.io

## Unit Testing

There is an ongoing implementation of unit testing on this repository. This is far from being complete, but you can already have a taste of what it is like, if you are interested.

Run the command below to generate the expected mocks:

```shell
mockgen -source=web/request.go -destination=mocks/request_mock.go -package=mocks
```

To run the tests:

```shell
go test -timeout 30s github.com/otaviokr/spacetraders-ship/component
```

Since this is a work in progress, tests may be temporarily broken... sorry!
