package downwardotel

import "go.opentelemetry.io/otel/attribute"

const (
	SwarmStackNameKey   attribute.Key = "docker.swarm.stack.name"
	SwarmServiceIDKey   attribute.Key = "docker.swarm.service.id"
	SwarmServiceNameKey attribute.Key = "docker.swarm.service.name"
	SwarmTaskIDKey      attribute.Key = "docker.swarm.task.id"
	SwarmTaskNameKey    attribute.Key = "docker.swarm.task.name"
	SwarmTaskSlotKey    attribute.Key = "docker.swarm.task.slot"
	SwarmNodeIDKey      attribute.Key = "docker.swarm.node.id"
	SwarmNodeNameKey    attribute.Key = "docker.swarm.node.name"
)
