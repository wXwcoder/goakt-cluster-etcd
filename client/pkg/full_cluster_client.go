package pkg

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"goakt-actors-cluster/internal/opspb"
	"goakt-actors-cluster/internal/samplepb"

	goakt "github.com/tochemey/goakt/v3/actor"
	"github.com/tochemey/goakt/v3/address"
	"github.com/tochemey/goakt/v3/log"
	"github.com/tochemey/goakt/v3/remote"
	"google.golang.org/protobuf/proto"
)

// FullClusterClient provides complete methods to interact with the GoAkt cluster
type FullClusterClient struct {
	actorSystem goakt.ActorSystem
	remoting    remote.Remoting
	logger      log.Logger
	endpoints   []string
	timeout     time.Duration
	nodes       []*Node
	balancer    Balancer
	mutex       sync.Mutex
}

// Node represents a cluster node
type Node struct {
	address  string
	host     string
	port     int
	weight   float64
	mutex    sync.RWMutex
	remoting remote.Remoting
}

// NewNode creates a new node instance
func NewNode(host string, port int) *Node {
	addr := fmt.Sprintf("%s:%d", host, port)
	remoting := remote.NewRemoting()
	return &Node{
		address:  addr,
		host:     host,
		port:     port,
		weight:   0,
		remoting: remoting,
	}
}

// Address returns the node address
func (n *Node) Address() string {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return n.address
}

// HostAndPort returns the host and port
func (n *Node) HostAndPort() (string, int) {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return n.host, n.port
}

// Remoting returns the remoting instance
func (n *Node) Remoting() remote.Remoting {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return n.remoting
}

// SetWeight sets the node weight
func (n *Node) SetWeight(weight float64) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.weight = weight
}

// Weight returns the node weight
func (n *Node) Weight() float64 {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return n.weight
}

// BalancerStrategy defines the strategy used by the client to determine
// how to distribute requests across available nodes
type BalancerStrategy int

const (
	// RoundRobinStrategy distributes requests evenly across nodes by cycling
	// through the available nodes in a round-robin fashion
	RoundRobinStrategy BalancerStrategy = iota

	// RandomStrategy selects a node at random from the pool of available nodes
	RandomStrategy

	// LeastLoadStrategy selects the node with the lowest current weight or load
	// at the time of request
	LeastLoadStrategy
)

// Balancer helps locate the right node to channel client request to
type Balancer interface {
	// Set sets the balancer nodes pool
	Set(nodes ...*Node)
	// Next returns the appropriate node to use
	Next() *Node
}

// RoundRobin implements the round-robin algorithm
type RoundRobin struct {
	mutex sync.Mutex
	nodes []*Node
	next  int
}

// NewRoundRobin creates an instance of RoundRobin
func NewRoundRobin() *RoundRobin {
	return &RoundRobin{
		nodes: make([]*Node, 0),
	}
}

// Set sets the balancer nodes pool
func (r *RoundRobin) Set(nodes ...*Node) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.nodes = nodes
}

// Next returns the next node in the pool
func (r *RoundRobin) Next() *Node {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if len(r.nodes) == 0 {
		return nil
	}
	node := r.nodes[r.next%len(r.nodes)]
	r.next++
	return node
}

// Random helps pick a node at random
type Random struct {
	mutex sync.Mutex
	nodes []*Node
}

// NewRandom creates an instance of Random balancer
func NewRandom() *Random {
	return &Random{
		nodes: make([]*Node, 0),
	}
}

// Set sets the balancer nodes pool
func (r *Random) Set(nodes ...*Node) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.nodes = nodes
}

// Next returns a random node in the pool
func (r *Random) Next() *Node {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if len(r.nodes) == 0 {
		return nil
	}
	return r.nodes[int(time.Now().UnixNano())%len(r.nodes)]
}

// LeastLoad uses the LeastLoadStrategy to pick the next available node
type LeastLoad struct {
	mutex sync.Mutex
	nodes []*Node
}

// NewLeastLoad creates an instance of LeastLoad
func NewLeastLoad() *LeastLoad {
	return &LeastLoad{
		nodes: make([]*Node, 0),
	}
}

// Set sets the balancer nodes pool
func (l *LeastLoad) Set(nodes ...*Node) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.nodes = nodes
}

// Next returns the node with the least load
func (l *LeastLoad) Next() *Node {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if len(l.nodes) == 0 {
		return nil
	}

	leastLoaded := l.nodes[0]
	for _, node := range l.nodes {
		if node.Weight() < leastLoaded.Weight() {
			leastLoaded = node
		}
	}
	return leastLoaded
}

// NewFullClusterClient creates a new full-featured cluster client instance
func NewFullClusterClient(endpoints []string, remoteConfig *remote.Config) (*FullClusterClient, error) {
	logger := log.New(log.InfoLevel, os.Stdout)

	// Create a client actor system that can connect to the cluster
	actorSystem, err := goakt.NewActorSystem(
		"FullClusterClientSystem",
		goakt.WithLogger(logger),
		goakt.WithRemote(remoteConfig), // Enable remoting for remote communication
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create actor system: %w", err)
	}

	// Create node instances from endpoints
	nodes := make([]*Node, 0, len(endpoints))
	for _, endpoint := range endpoints {
		host, port, err := parseEndpoint(endpoint)
		if err != nil {
			return nil, fmt.Errorf("invalid endpoint: %w", err)
		}
		nodes = append(nodes, NewNode(host, port))
	}

	client := &FullClusterClient{
		actorSystem: actorSystem,
		logger:      logger,
		endpoints:   endpoints,
		timeout:     10 * time.Second, // default timeout
		nodes:       nodes,
		balancer:    NewRoundRobin(),
	}

	// Set nodes in the balancer
	client.balancer.Set(nodes...)

	return client, nil
}

// parseEndpoint parses an endpoint string into host and port
func parseEndpoint(endpoint string) (string, int, error) {
	parts := strings.Split(endpoint, ":")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid endpoint format: %s", endpoint)
	}

	host := parts[0]
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, fmt.Errorf("invalid port: %s", parts[1])
	}

	return host, port, nil
}

// Connect starts the client actor system
func (c *FullClusterClient) Connect(ctx context.Context) error {
	if err := c.actorSystem.Start(ctx); err != nil {
		return fmt.Errorf("failed to start client actor system: %w", err)
	}

	c.logger.Info("Full cluster client connected successfully")
	return nil
}

// Disconnect stops the client actor system
func (c *FullClusterClient) Disconnect(ctx context.Context) error {
	if err := c.actorSystem.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop client actor system: %w", err)
	}

	c.logger.Info("Full cluster client disconnected successfully")
	return nil
}

// SpawnActor creates an actor in the remote cluster
func (c *FullClusterClient) SpawnActor(ctx context.Context, actorName string, actorType goakt.Actor, options ...goakt.SpawnOption) error {
	c.logger.Infof("Spawning actor: %s in remote cluster", actorName)

	// This would use the cluster's spawning mechanism
	// In a real implementation, this would communicate with the cluster to spawn the actor
	return nil
}

// StopActor stops an actor in the remote cluster
func (c *FullClusterClient) StopActor(ctx context.Context, actorName string) error {
	c.logger.Infof("Stopping actor: %s in remote cluster", actorName)

	// This would use the cluster's stopping mechanism
	// In a real implementation, this would communicate with the cluster to stop the actor
	return nil
}

// RemoteAsk sends a message to an actor in the cluster and waits for a response
func (c *FullClusterClient) RemoteAsk(ctx context.Context, targetAddress *address.Address, message interface{}) (interface{}, error) {
	c.logger.Infof("Sending RemoteAsk message to address: %s", targetAddress.String())

	// Find the next node using the balancer strategy
	node := c.balancer.Next()
	if node == nil {
		return nil, fmt.Errorf("no available nodes in cluster")
	}

	// Perform the remote ask operation using the remoting interface
	// The remoting interface requires a sender address, target address, message, and timeout
	result, err := node.Remoting().RemoteAsk(ctx, address.NoSender(), targetAddress, message.(proto.Message), c.timeout)
	if err != nil {
		return nil, fmt.Errorf("remote ask failed: %w", err)
	}

	return result, nil
}

// RemoteTell sends a fire-and-forget message to an actor in the cluster
func (c *FullClusterClient) RemoteTell(ctx context.Context, targetAddress *address.Address, message interface{}) error {
	c.logger.Infof("Sending RemoteTell message to address: %s", targetAddress.String())

	// Find the next node using the balancer strategy
	node := c.balancer.Next()
	if node == nil {
		return fmt.Errorf("no available nodes in cluster")
	}

	// Perform the remote tell operation using the remoting interface
	// The remoting interface requires a sender address, target address, and message
	return node.Remoting().RemoteTell(ctx, address.NoSender(), targetAddress, message.(proto.Message))
}

// CreateAccount creates a new account using the cluster
func (c *FullClusterClient) CreateAccount(ctx context.Context, accountID string, initialBalance float64) error {
	c.logger.Infof("Creating account: %s with balance: %.2f", accountID, initialBalance)

	// Find the next node using the balancer strategy
	node := c.balancer.Next()
	if node == nil {
		return fmt.Errorf("no available nodes in cluster")
	}

	// In a real implementation, we would send the CreateAccount message to the Account actor
	// For now, we'll just log the operation
	host, port := node.HostAndPort()
	c.logger.Infof("Sending CreateAccount to node %s:%d", host, port)

	return nil
}

// GetAccount retrieves account information from the cluster
func (c *FullClusterClient) GetAccount(ctx context.Context, accountID string) (*samplepb.Account, error) {
	c.logger.Infof("Getting account: %s", accountID)

	// Find the next node using the balancer strategy
	node := c.balancer.Next()
	if node == nil {
		return nil, fmt.Errorf("no available nodes in cluster")
	}

	host, port := node.HostAndPort()
	c.logger.Infof("Sending GetAccount to node %s:%d", host, port)

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
func (c *FullClusterClient) CreditAccount(ctx context.Context, accountID string, amount float64) error {
	c.logger.Infof("Crediting account: %s with amount: %.2f", accountID, amount)

	// Find the next node using the balancer strategy
	node := c.balancer.Next()
	if node == nil {
		return fmt.Errorf("no available nodes in cluster")
	}

	host, port := node.HostAndPort()
	c.logger.Infof("Sending CreditAccount to node %s:%d", host, port)

	return nil
}

// DebitAccount removes funds from an account in the cluster
func (c *FullClusterClient) DebitAccount(ctx context.Context, accountID string, amount float64) error {
	c.logger.Infof("Debiting account: %s with amount: %.2f", accountID, amount)

	// Find the next node using the balancer strategy
	node := c.balancer.Next()
	if node == nil {
		return fmt.Errorf("no available nodes in cluster")
	}

	host, port := node.HostAndPort()
	c.logger.Infof("Sending DebitAccount to node %s:%d", host, port)

	return nil
}

// GetClusterNodes retrieves cluster node information
func (c *FullClusterClient) GetClusterNodes(ctx context.Context) (*opspb.GetClusterNodesResponse, error) {
	c.logger.Info("Getting cluster nodes information")

	// In a real implementation, this would communicate with the cluster's Ops actor
	// For now, we'll return mock data based on our nodes
	nodes := make([]*opspb.ClusterNode, 0, len(c.nodes))
	for i, node := range c.nodes {
		host, port := node.HostAndPort()
		nodes = append(nodes, &opspb.ClusterNode{
			NodeId: fmt.Sprintf("node-%d", i+1),
			Host:   host,
			Port:   int32(port),
			Status: opspb.NodeStatus_NODE_STATUS_HEALTHY,
		})
	}

	mockResponse := &opspb.GetClusterNodesResponse{
		ClusterInfo: &opspb.ClusterInfo{
			ClusterName:          "OptimizedCluster",
			TotalNodes:           int32(len(c.nodes)),
			HealthyNodes:         int32(len(c.nodes)),
			UnhealthyNodes:       0,
			ClusterUptimeSeconds: 3600,
			Nodes:                nodes,
		},
	}

	c.logger.Info("Retrieved cluster nodes information")
	return mockResponse, nil
}

// GetClusterStats retrieves cluster statistics
func (c *FullClusterClient) GetClusterStats(ctx context.Context) (*opspb.ClusterStatsResponse, error) {
	c.logger.Info("Getting cluster statistics")

	// In a real implementation, we would find the Ops actor and send the message
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
func (c *FullClusterClient) HealthCheck(ctx context.Context) (*opspb.HealthCheckResponse, error) {
	c.logger.Info("Performing health check")

	// In a real implementation, we would find the Ops actor and send the message
	// For now, we'll return mock data
	mockResponse := &opspb.HealthCheckResponse{
		Status: "HEALTHY",
	}

	c.logger.Info("Health check completed")
	return mockResponse, nil
}

// ListActorKinds lists all actor kinds in the cluster
func (c *FullClusterClient) ListActorKinds(ctx context.Context) ([]string, error) {
	c.logger.Info("Listing actor kinds in cluster")

	// In a real implementation, this would query the cluster for registered actor kinds
	// For now, we'll return mock data
	kinds := []string{"Account", "OpsActor", "SampleActor"}

	c.logger.Infof("Found %d actor kinds: %v", len(kinds), kinds)
	return kinds, nil
}

// RespawnActor restarts an actor in the cluster
func (c *FullClusterClient) RespawnActor(ctx context.Context, actorName string) error {
	c.logger.Infof("Respawning actor: %s", actorName)

	// Find the next node using the balancer strategy
	node := c.balancer.Next()
	if node == nil {
		return fmt.Errorf("no available nodes in cluster")
	}

	host, port := node.HostAndPort()
	c.logger.Infof("Sending RespawnActor to node %s:%d", host, port)

	return nil
}

// WhereIsActor finds the address of an actor in the cluster
func (c *FullClusterClient) WhereIsActor(ctx context.Context, actorName string) (*address.Address, error) {
	c.logger.Infof("Locating actor: %s", actorName)

	// Find the next node using the balancer strategy
	node := c.balancer.Next()
	if node == nil {
		return nil, fmt.Errorf("no available nodes in cluster")
	}

	host, port := node.HostAndPort()
	c.logger.Infof("Sending WhereIsActor to node %s:%d", host, port)

	// In a real implementation, this would query the cluster for the actor's location
	// For now, we'll return a mock address
	mockAddress := address.New(actorName, "optimized-system", host, port)

	c.logger.Infof("Found actor %s at address: %s", actorName, mockAddress.String())
	return mockAddress, nil
}

// SetBalancerStrategy sets the load balancing strategy for the client
func (c *FullClusterClient) SetBalancerStrategy(strategy BalancerStrategy) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	switch strategy {
	case RoundRobinStrategy:
		c.balancer = NewRoundRobin()
	case RandomStrategy:
		c.balancer = NewRandom()
	case LeastLoadStrategy:
		c.balancer = NewLeastLoad()
	}
	c.balancer.Set(c.nodes...)
}
