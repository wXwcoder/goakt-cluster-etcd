package pkg

/*
import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tochemey/goakt/v3/remote"
)

func TestClusterClient(t *testing.T) {
	// Test creating a new client
	endpoints := []string{"localhost:14001"}
	remoteConfig := remote.NewConfig("localhost", 14001)

	client, err := NewFullClusterClient(endpoints, remoteConfig)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	// Test connecting and disconnecting
	ctx := context.Background()
	err = client.Connect(ctx)
	assert.NoError(t, err)

	err = client.Disconnect(ctx)
	assert.NoError(t, err)
}

func TestAccountOperations(t *testing.T) {
	endpoints := []string{"localhost:14001"}
	remoteConfig := remote.NewConfig("localhost", 14001)

	client, err := NewFullClusterClient(endpoints, remoteConfig)
	assert.NoError(t, err)

	ctx := context.Background()
	err = client.Connect(ctx)
	assert.NoError(t, err)
	defer client.Disconnect(ctx)

	// Test creating an account
	accountID := "test-account-123"
	initialBalance := 1000.0

	err = client.CreateAccount(ctx, accountID, initialBalance)
	assert.NoError(t, err)

	// Test getting the account
	account, err := client.GetAccount(ctx, accountID)
	assert.NoError(t, err)
	assert.Equal(t, accountID, account.AccountId)
	assert.Equal(t, initialBalance, account.AccountBalance)

	// Test crediting the account
	creditAmount := 500.0
	err = client.CreditAccount(ctx, accountID, creditAmount)
	assert.NoError(t, err)

	// Test debiting the account
	debitAmount := 200.0
	err = client.DebitAccount(ctx, accountID, debitAmount)
	assert.NoError(t, err)
}

func TestClusterOperations(t *testing.T) {
	endpoints := []string{"localhost:14001"}
	remoteConfig := remote.NewConfig("localhost", 14001)

	client, err := NewFullClusterClient(endpoints, remoteConfig)
	assert.NoError(t, err)

	ctx := context.Background()
	err = client.Connect(ctx)
	assert.NoError(t, err)
	defer client.Disconnect(ctx)

	// Test getting cluster nodes
	nodes, err := client.GetClusterNodes(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, nodes)
	assert.NotNil(t, nodes.ClusterInfo)

	// Test getting cluster stats
	stats, err := client.GetClusterStats(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, stats)

	// Test health check
	health, err := client.HealthCheck(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, health)
	assert.NotEmpty(t, health.Status)

	// Test listing actor kinds
	kinds, err := client.ListActorKinds(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, kinds)
}
*/
