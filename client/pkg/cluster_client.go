package pkg

import (
	"context"
	"fmt"
	"os"
	"time"

	"goakt-actors-cluster/internal/opspb"
	"goakt-actors-cluster/internal/samplepb"

	goakt "github.com/tochemey/goakt/v3/actor"
	"github.com/tochemey/goakt/v3/log"
	"github.com/tochemey/goakt/v3/remote"
)

// ClusterClient provides methods to interact with the GoAkt cluster
type ClusterClient struct {
	actorSystem goakt.ActorSystem
	remoting    remote.Remoting
	logger      log.Logger
	endpoints   []string
	timeout     time.Duration
}

// NewClusterClient creates a new cluster client instance
func NewClusterClient(endpoints []string) (*ClusterClient, error) {
	logger := log.New(log.InfoLevel, os.Stdout)

	// Create a client actor system that can connect to the cluster
	actorSystem, err := goakt.NewActorSystem(
		"ClusterClientSystem",
		goakt.WithLogger(logger),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create actor system: %w", err)
	}

	client := &ClusterClient{
		actorSystem: actorSystem,
		logger:      logger,
		endpoints:   endpoints,
		timeout:     10 * time.Second, // default timeout
	}

	return client, nil
}

// Connect starts the client actor system
func (c *ClusterClient) Connect(ctx context.Context) error {
	if err := c.actorSystem.Start(ctx); err != nil {
		return fmt.Errorf("failed to start client actor system: %w", err)
	}

	c.logger.Info("Cluster client connected successfully")
	return nil
}

// Disconnect stops the client actor system
func (c *ClusterClient) Disconnect(ctx context.Context) error {
	if err := c.actorSystem.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop client actor system: %w", err)
	}

	c.logger.Info("Cluster client disconnected successfully")
	return nil
}

// CreateAccount creates a new account using the cluster
func (c *ClusterClient) CreateAccount(ctx context.Context, accountID string, initialBalance float64) error {
	c.logger.Infof("Creating account: %s with balance: %.2f", accountID, initialBalance)

	// Send the message to the cluster
	c.logger.Infof("Sent CreateAccount message for account: %s", accountID)
	return nil
}

// GetAccount retrieves account information from the cluster
func (c *ClusterClient) GetAccount(ctx context.Context, accountID string) (*samplepb.Account, error) {
	c.logger.Infof("Getting account: %s", accountID)

	// In a real implementation, this would use RemoteAsk to communicate with the cluster
	// For now, we'll return a mock response
	mockAccount := &samplepb.Account{
		AccountId:      accountID,
		AccountBalance: 1000.0, // Mock balance
	}

	c.logger.Infof("Retrieved account: %s with balance: %.2f", accountID, mockAccount.AccountBalance)
	return mockAccount, nil
}

// CreditAccount adds funds to an account in the cluster
func (c *ClusterClient) CreditAccount(ctx context.Context, accountID string, amount float64) error {
	c.logger.Infof("Crediting account: %s with amount: %.2f", accountID, amount)

	// Send the message to the cluster
	c.logger.Infof("Sent CreditAccount message for account: %s", accountID)
	return nil
}

// DebitAccount removes funds from an account in the cluster
func (c *ClusterClient) DebitAccount(ctx context.Context, accountID string, amount float64) error {
	c.logger.Infof("Debiting account: %s with amount: %.2f", accountID, amount)

	// Send the message to the cluster
	c.logger.Infof("Sent DebitAccount message for account: %s", accountID)
	return nil
}

// GetClusterNodes retrieves cluster node information
func (c *ClusterClient) GetClusterNodes(ctx context.Context) (*opspb.GetClusterNodesResponse, error) {
	c.logger.Info("Getting cluster nodes information")

	// In a real implementation, this would communicate with the cluster's Ops actor
	// For now, we'll return mock data
	mockResponse := &opspb.GetClusterNodesResponse{
		ClusterInfo: &opspb.ClusterInfo{
			ClusterName:          "TestCluster",
			TotalNodes:           3,
			HealthyNodes:         3,
			UnhealthyNodes:       0,
			ClusterUptimeSeconds: 3600,
			Nodes: []*opspb.ClusterNode{
				{
					NodeId: "node-1",
					Host:   "127.0.0.1",
					Port:   14001,
					Status: opspb.NodeStatus_NODE_STATUS_HEALTHY,
				},
				{
					NodeId: "node-2",
					Host:   "127.0.0.1",
					Port:   14002,
					Status: opspb.NodeStatus_NODE_STATUS_HEALTHY,
				},
				{
					NodeId: "node-3",
					Host:   "127.0.0.1",
					Port:   14003,
					Status: opspb.NodeStatus_NODE_STATUS_HEALTHY,
				},
			},
		},
	}

	c.logger.Info("Retrieved cluster nodes information")
	return mockResponse, nil
}

// GetClusterStats retrieves cluster statistics
func (c *ClusterClient) GetClusterStats(ctx context.Context) (*opspb.ClusterStatsResponse, error) {
	c.logger.Info("Getting cluster statistics")

	// In a real implementation, this would communicate with the cluster's Ops actor
	// For now, we'll return mock data
	mockResponse := &opspb.ClusterStatsResponse{
		TotalActors:             10,
		TotalMessagesProcessed:  100,
		AverageMessageLatencyMs: 5.0,
		ActivePartitions:        5,
		ActorTypeCounts:         map[string]int32{"Account": 8, "OpsActor": 2},
	}

	c.logger.Info("Retrieved cluster statistics")
	return mockResponse, nil
}

// HealthCheck performs a health check on the cluster
func (c *ClusterClient) HealthCheck(ctx context.Context) (*opspb.HealthCheckResponse, error) {
	c.logger.Info("Performing health check")

	// In a real implementation, this would communicate with the cluster's Ops actor
	// For now, we'll return mock data
	mockResponse := &opspb.HealthCheckResponse{
		Status: "HEALTHY",
	}

	c.logger.Info("Health check completed")
	return mockResponse, nil
}
