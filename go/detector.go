package downwardotel

import (
	"context"
	"fmt"
	"strings"

	downward "github.com/swarm-deploy/downward/go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
)

var _ resource.Detector = (*Detector)(nil)

type Detector struct{}

func NewDetector() *Detector {
	return &Detector{}
}

func (*Detector) Detect(_ context.Context) (*resource.Resource, error) {
	info, err := downward.Load()
	if err != nil {
		return nil, fmt.Errorf("load downward info: %w", err)
	}

	attrs := make([]attribute.KeyValue, 0, 14)

	if info.Stack.Name != "" {
		attrs = append(attrs,
			semconv.ServiceNamespaceKey.String(info.Stack.Name),
			SwarmStackNameKey.String(info.Stack.Name),
		)
	}

	if info.Service.ID != "" {
		attrs = append(attrs, SwarmServiceIDKey.String(info.Service.ID))
	}
	if info.Service.Name != "" {
		attrs = append(attrs,
			semconv.ServiceNameKey.String(serviceName(info.Stack.Name, info.Service.Name)),
			SwarmServiceNameKey.String(info.Service.Name),
		)
	}

	if info.Task.ID != "" {
		attrs = append(attrs,
			semconv.ServiceInstanceIDKey.String(info.Task.ID),
			SwarmTaskIDKey.String(info.Task.ID),
		)
	}
	if info.Task.Name != "" {
		attrs = append(attrs, SwarmTaskNameKey.String(info.Task.Name))
	}
	if info.Task.Slot != nil {
		attrs = append(attrs, SwarmTaskSlotKey.Int(*info.Task.Slot))
	}

	if info.Node.ID != "" {
		attrs = append(attrs,
			semconv.HostIDKey.String(info.Node.ID),
			SwarmNodeIDKey.String(info.Node.ID),
		)
	}
	if info.Node.Name != "" {
		attrs = append(attrs,
			semconv.HostNameKey.String(info.Node.Name),
			SwarmNodeNameKey.String(info.Node.Name),
		)
	}

	if len(attrs) == 0 {
		return resource.Empty(), nil
	}

	attrs = append(attrs, semconv.ContainerRuntimeNameKey.String("docker"))

	return resource.NewSchemaless(attrs...), nil
}

func serviceName(stackName, fullName string) string {
	if stackName == "" {
		return fullName
	}

	return strings.TrimPrefix(fullName, stackName+"_")
}
