# flow example

In this example a demo application is run the computes random Fibonacci
numbers. The application is traced using OpenTelemetry with a `flow`
`SpanProcessor`. The metrics produced with by `flow` are collected with an
OpenTelemetry collector and ultimately shipped to a Prometheus instance.

## Prerequisites

This demo expects [docker-compose](https://docs.docker.com/compose/) is
installed and working.

## Start the demo

Start all services.

```sh
docker-compose up --detach
```

This should do a few things.

1. Build the demo Fibonacci application and package it as a Docker image.
2. Start the Fibonacci application.
3. Start an OpenTelemetry Collector instance. This collector will forward
   traces from the application to a Jaeger instance and metrics to a Prometheus
   instance.
4. Start a Jaeger instance to capture traces from the collector.
5. Start a Prometheus instance to scrape metrics from the collector.

## Verify

### Fibonacci Application

First, make sure the Fibonacci application is reporting metrics.

```sh
$ curl -s http://localhost:41820/metrics | grep 'spans_total'

# HELP spans_total The total number of processed spans
# TYPE spans_total counter
spans_total{state="ended"} 762
spans_total{state="started"} 762
```

The total number of spans will be different, but you should see something like above.

### Jaeger Traces

Next, make sure the traces are making it to the Jaeger instance. Navigate
through the [Jaeger UI] you are running locally to ensure traces from the
application are there.

### Prometheus Metrics

Finally, verify the `flow` `SpanProcessor` is working. Go to the [Prometheus
targets page] and make sure all targets are up.

From here you can go to the [Prometheus graph page] and visualize that data
being scraped.

For example, run the following query to see the exported span rate.

```prometheus
rate(spans_total[1m])
```

## Stop and Clean Up

Stop all services.

```sh
docker-compose down
```

If you are all done with the demo you can remove everything as well.

```sh
docker-compose rm
```

[Jaeger UI]: http://localhost:16686/
[Prometheus targets page]: http://localhost:9090/targets
[Prometheus graph page]: http://localhost:9090/graph
