package main

import (
	"context"
	"fmt"
	"os"
	"time"

	goakt "github.com/tochemey/goakt/v3/actor"
	"github.com/tochemey/goakt/v3/address"
	"github.com/tochemey/goakt/v3/log"
	"github.com/tochemey/goakt/v3/remote"

	"goakt-actors-cluster/actors"
	"goakt-actors-cluster/config"
	clusterdiscovery "goakt-actors-cluster/discovery"
	"goakt-actors-cluster/internal/opspb"
)

func main() {
	// Create a new actor system for testing
	logger := log.New(log.InfoLevel, os.Stdout)

	remoting := remote.NewRemoting()

	// Create a mock config for testing
	cfg := &config.Config{
		ServiceName:     "test-service",
		ActorSystemName: "Accounts",
		GossipPort:      8558,
		PeersPort:       2552,
		RemotingPort:    14001,
		EtcdEndpoints:   []string{"localhost:2379"},
		EtcdTTL:         30,
		EtcdTimeout:     5,
		EtcdDialTimeout: 5,
		Port:            50051,
	}

	// Create a mock discovery config
	discoveryConfig := &clusterdiscovery.Config{
		Context:         context.Background(),
		Endpoints:       []string{"localhost:2379"},
		ActorSystemName: cfg.ActorSystemName,
		Host:            "localhost",
		DiscoveryPort:   cfg.GossipPort,
		TTL:             cfg.EtcdTTL,
		DialTimeout:     time.Duration(cfg.EtcdDialTimeout) * time.Second,
		Timeout:         time.Duration(cfg.EtcdTimeout) * time.Second,
	}

	// Create a mock discovery
	discovery := clusterdiscovery.NewDiscovery(discoveryConfig)
	clusterConfig := goakt.
		NewClusterConfig().
		WithDiscovery(discovery).
		WithPartitionCount(19).
		WithDiscoveryPort(cfg.GossipPort).
		WithPeersPort(cfg.PeersPort).
		WithClusterBalancerInterval(time.Second).
		WithKinds(new(actors.OpsActor))

	actorSystem, err := goakt.NewActorSystem(
		cfg.ActorSystemName,
		goakt.WithLogger(logger),
		goakt.WithActorInitMaxRetries(3),
		goakt.WithRemote(remote.NewConfig("localhost", cfg.RemotingPort)),
		goakt.WithCluster(clusterConfig))

	if err != nil {
		fmt.Printf("Failed to create actor system: %v\n", err)
		return
	}

	// Start the actor system
	ctx := context.Background()
	if err := actorSystem.Start(ctx); err != nil {
		fmt.Printf("Failed to start actor system: %v\n", err)
		return
	}

	// Wait for actor system to start
	time.Sleep(time.Second * 2)

	// Spawn the ops actor
	opsActorName := "ops-actor"
	addr, _, err := actorSystem.ActorOf(ctx, opsActorName)
	if err != nil {
		fmt.Printf("Failed to spawn ops actor: %v\n", err)
		return
	}

	fmt.Printf("Ops actor find successfully: %s\n", addr)

	// Wait a bit for the actor to be ready
	time.Sleep(time.Second)

	// Test all ops functions using actor system
	testHealthCheck(ctx, actorSystem, remoting, addr)
	testGetClusterNodes(ctx, actorSystem, remoting, addr)
	testGetClusterStats(ctx, actorSystem, remoting, addr)
	testGetActorDetails(ctx, actorSystem, remoting, addr)

	// Stop the actor system
	if err := actorSystem.Stop(ctx); err != nil {
		fmt.Printf("Failed to stop actor system: %v\n", err)
	}
}

// testHealthCheck tests the health check functionality using actor system
func testHealthCheck(ctx context.Context, actorSystem goakt.ActorSystem, remoting remote.Remoting, addr *address.Address) {
	fmt.Println("=== Testing Health Check via Actor System ===")

	// Create health check request
	req := &opspb.HealthCheckRequest{}

	// Send request to ops actor using Ask
	resp, err := remoting.RemoteAsk(ctx, address.NoSender(), addr, req, 10*time.Second)
	if err != nil {
		fmt.Printf("Health check failed: %v\n", err)
		return
	}
	message, _ := resp.UnmarshalNew()

	// Cast response to expected type
	healthResp, ok := message.(*opspb.OpsActorResponse)
	if !ok {
		fmt.Printf("Invalid response type: %T\n", message)
		return
	}

	// Extract health check response
	hcResp := healthResp.GetHealthCheck()
	if hcResp != nil {
		fmt.Printf("Health check successful: Status=%s\n", hcResp.Status)
	} else {
		fmt.Printf("Health check response is nil\n")
	}
	fmt.Println()
}

// testGetClusterNodes tests getting cluster nodes using actor system
func testGetClusterNodes(ctx context.Context, actorSystem goakt.ActorSystem, remoting remote.Remoting, addr *address.Address) {
	fmt.Println("=== Testing Get Cluster Nodes via Actor System ===")

	// Create get cluster nodes request
	req := &opspb.GetClusterNodesRequest{
		IncludeDetails: true,
	}

	// Send request to ops actor using Ask
	resp, err := remoting.RemoteAsk(ctx, address.NoSender(), addr, req, 10*time.Second)
	if err != nil {
		fmt.Printf("Get cluster nodes failed: %v\n", err)
		return
	}
	message, _ := resp.UnmarshalNew()

	// Cast response to expected type
	nodesResp, ok := message.(*opspb.OpsActorResponse)
	if !ok {
		fmt.Printf("Invalid response type: %T\n", message)
		return
	}

	// Extract cluster nodes response
	clusterNodesResp := nodesResp.GetClusterNodes()
	if clusterNodesResp != nil {
		fmt.Printf("Get cluster nodes successful: Total nodes=%d, Healthy nodes=%d\n",
			clusterNodesResp.ClusterInfo.GetTotalNodes(),
			clusterNodesResp.ClusterInfo.GetHealthyNodes())

		if clusterNodesResp.ClusterInfo.GetNodes() != nil {
			fmt.Printf("Number of nodes returned: %d\n", len(clusterNodesResp.ClusterInfo.GetNodes()))
		}
	} else {
		fmt.Printf("Get cluster nodes response is nil\n")
	}
	fmt.Println()
}

// testGetClusterStats tests getting cluster stats using actor system
func testGetClusterStats(ctx context.Context, actorSystem goakt.ActorSystem, remoting remote.Remoting, addr *address.Address) {
	fmt.Println("=== Testing Get Cluster Stats via Actor System ===")

	// Create get cluster stats request
	req := &opspb.ClusterStatsRequest{}

	// Send request to ops actor using Ask
	resp, err := remoting.RemoteAsk(ctx, address.NoSender(), addr, req, 10*time.Second)
	if err != nil {
		fmt.Printf("Get cluster stats failed: %v\n", err)
		return
	}
	message, _ := resp.UnmarshalNew()

	// Cast response to expected type
	statsResp, ok := message.(*opspb.OpsActorResponse)
	if !ok {
		fmt.Printf("Invalid response type: %T\n", message)
		return
	}

	// Extract cluster stats response
	clusterStatsResp := statsResp.GetClusterStats()
	if clusterStatsResp != nil {
		fmt.Printf("Get cluster stats successful: Total actors=%d, Active partitions=%d\n",
			clusterStatsResp.GetTotalActors(),
			clusterStatsResp.GetActivePartitions())

		fmt.Printf("Actor type counts: %v\n", clusterStatsResp.GetActorTypeCounts())
	} else {
		fmt.Printf("Get cluster stats response is nil\n")
	}
	fmt.Println()
}

// testGetActorDetails tests getting actor details using actor system
func testGetActorDetails(ctx context.Context, actorSystem goakt.ActorSystem, remoting remote.Remoting, addr *address.Address) {
	fmt.Println("=== Testing Get Actor Details via Actor System ===")

	// Create get actor details request
	req := &opspb.GetActorDetailsRequest{
		ActorId: "account-001",
	}

	// Send request to ops actor using Ask
	resp, err := remoting.RemoteAsk(ctx, address.NoSender(), addr, req, 10*time.Second)
	if err != nil {
		fmt.Printf("Get actor details failed: %v\n", err)
		return
	}
	message, _ := resp.UnmarshalNew()

	// Cast response to expected type
	actorDetailsResp, ok := message.(*opspb.OpsActorResponse)
	if !ok {
		fmt.Printf("Invalid response type: %T\n", message)
		return
	}

	// Extract actor details response
	getActorDetailsResp := actorDetailsResp.GetActorDetails()
	if getActorDetailsResp != nil {
		fmt.Printf("Get actor details successful: Total count=%d\n",
			getActorDetailsResp.GetTotalCount())

		if getActorDetailsResp.GetActors() != nil {
			fmt.Printf("Number of actors returned: %d\n", len(getActorDetailsResp.GetActors()))
			for _, actor := range getActorDetailsResp.GetActors() {
				fmt.Printf("  Actor ID: %s, Type: %s, Status: %s\n",
					actor.GetActorId(), actor.GetActorType(), actor.GetStatus().String())
			}
		}
	} else {
		fmt.Printf("Get actor details response is nil\n")
	}
	fmt.Println()
}
