package downwardotel

import (
	"context"
	"testing"

	downward "github.com/swarm-deploy/downward/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	res, err := NewDetector().Detect(context.Background())
	require.NoError(t, err)

	attrs := make(map[attribute.Key]attribute.Value, res.Len())
	for _, attr := range res.Attributes() {
		attrs[attr.Key] = attr.Value
	}

	assert.Equal(t, "billing", attrs[semconv.ServiceNamespaceKey].AsString())
	assert.Equal(t, "api", attrs[semconv.ServiceNameKey].AsString())
	assert.Equal(t, "task-id", attrs[semconv.ServiceInstanceIDKey].AsString())
	assert.Equal(t, "node-id", attrs[semconv.HostIDKey].AsString())
	assert.Equal(t, "worker-02", attrs[semconv.HostNameKey].AsString())
	assert.Equal(t, "docker", attrs[semconv.ContainerRuntimeNameKey].AsString())
	assert.Equal(t, "billing_api", attrs[SwarmServiceNameKey].AsString())
	assert.EqualValues(t, 2, attrs[SwarmTaskSlotKey].AsInt64())
}

func TestDetectorDetectInvalidTaskSlot(t *testing.T) {
	t.Setenv(downward.EnvTaskSlot, "invalid")

	_, err := NewDetector().Detect(context.Background())
	require.Error(t, err)
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
		assert.Equal(t, tt.want, serviceName(tt.stack, tt.full))
	}
}
