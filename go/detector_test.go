package downwardotel

import (
	"context"
	"testing"

	downward "github.com/swarm-deploy/downward/go"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
)

func TestDetectorDetect(t *testing.T) {
	t.Setenv(downward.EnvStackName, "billing")
	t.Setenv(downward.EnvServiceID, "service-id")
	t.Setenv(downward.EnvServiceName, "billing_api")
	t.Setenv(downward.EnvTaskID, "task-id")
	t.Setenv(downward.EnvTaskName, "billing_api.2.task-id")
	t.Setenv(downward.EnvTaskSlot, "2")
	t.Setenv(downward.EnvNodeID, "node-id")
	t.Setenv(downward.EnvNodeName, "worker-02")

	res, err := (Detector{}).Detect(context.Background())
	if err != nil {
		t.Fatalf("Detect() error = %v", err)
	}

	attrs := make(map[attribute.Key]attribute.Value, res.Len())
	for _, attr := range res.Attributes() {
		attrs[attr.Key] = attr.Value
	}

	assertStringAttribute(t, attrs, semconv.ServiceNamespaceKey, "billing")
	assertStringAttribute(t, attrs, semconv.ServiceNameKey, "api")
	assertStringAttribute(t, attrs, semconv.ServiceInstanceIDKey, "task-id")
	assertStringAttribute(t, attrs, semconv.HostIDKey, "node-id")
	assertStringAttribute(t, attrs, semconv.HostNameKey, "worker-02")
	assertStringAttribute(t, attrs, semconv.ContainerRuntimeNameKey, "docker")
	assertStringAttribute(t, attrs, SwarmServiceNameKey, "billing_api")

	if got := attrs[SwarmTaskSlotKey].AsInt64(); got != 2 {
		t.Fatalf("%s = %d, want 2", SwarmTaskSlotKey, got)
	}
}

func TestDetectorDetectInvalidTaskSlot(t *testing.T) {
	t.Setenv(downward.EnvTaskSlot, "invalid")

	_, err := (Detector{}).Detect(context.Background())
	if err == nil {
		t.Fatal("Detect() error = nil, want error")
	}
}

func TestServiceName(t *testing.T) {
	tests := []struct {
		stack string
		full  string
		want  string
	}{
		{stack: "billing", full: "billing_api", want: "api"},
		{stack: "billing", full: "api", want: "api"},
		{full: "billing_api", want: "billing_api"},
	}

	for _, tt := range tests {
		if got := serviceName(tt.stack, tt.full); got != tt.want {
			t.Fatalf("serviceName(%q, %q) = %q, want %q", tt.stack, tt.full, got, tt.want)
		}
	}
}

func assertStringAttribute(t *testing.T, attrs map[attribute.Key]attribute.Value, key attribute.Key, want string) {
	t.Helper()

	value, ok := attrs[key]
	if !ok {
		t.Fatalf("missing attribute %s", key)
	}
	if got := value.AsString(); got != want {
		t.Fatalf("%s = %q, want %q", key, got, want)
	}
}
