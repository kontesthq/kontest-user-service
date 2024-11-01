package utils

import (
	"fmt"
	"github.com/kontesthq/go-load-balancer/loadbalancer"
	"log/slog"
)

var (
	clients = make(map[string]loadbalancer.Client)
)

func GetOrCreateClient(serviceName string) (loadbalancer.Client, error) {
	// Check if we already have a load balancer client for the service
	if _, ok := clients[serviceName]; !ok {
		newClient, err := loadbalancer.NewConsulClientWithCustomRule(ConsulHost, ConsulPort, serviceName, loadbalancer.NewRetryRuleWithDefaults())

		if err != nil {
			slog.Error(fmt.Sprintf("Error creating load balancer for service: %s, Error: %s\n", serviceName, err))
			return nil, err
		}

		clients[serviceName] = newClient
	}

	return clients[serviceName], nil
}
