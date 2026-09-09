# downward-otel

OpenTelemetry resource detector for Docker Swarm metadata provided by [`swarm-deploy/downward`](https://github.com/swarm-deploy/downward).

```go
res, err := resource.New(
    ctx,
    resource.WithDetectors(downwardotel.NewDetector()),
)
```

The detector maps Swarm metadata to standard OpenTelemetry resource attributes and preserves the original values under `docker.swarm.*`.

## OpenTelemetry setup

```go
package main

import (
    "context"
    "log"

    downwardotel "github.com/swarm-deploy/downward-otel/go"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
    ctx := context.Background()

    res, err := resource.New(
        ctx,
        resource.WithDetectors(downwardotel.NewDetector()),
    )
    if err != nil {
        log.Fatal(err)
    }

    tracerProvider := sdktrace.NewTracerProvider(
        sdktrace.WithResource(res),
    )

    otel.SetTracerProvider(tracerProvider)
}
```
