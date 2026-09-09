# downward-otel

OpenTelemetry resource detector for Docker Swarm metadata provided by [`swarm-deploy/downward`](https://github.com/swarm-deploy/downward).

```go
res, err := resource.New(
    ctx,
    resource.WithDetectors(downwardotel.Detector{}),
)
```

The detector maps Swarm metadata to standard OpenTelemetry resource attributes and preserves the original values under `docker.swarm.*`.
