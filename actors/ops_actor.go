/*
 * MIT License
 *
 * Copyright (c) 2022-2025 Arsene Tochemey Gandote
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package actors

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	goakt "github.com/tochemey/goakt/v3/actor"
	"github.com/tochemey/goakt/v3/log"
	"github.com/tochemey/goakt/v3/remote"

	"goakt-actors-cluster/config"
	clusterdiscovery "goakt-actors-cluster/discovery"
	"goakt-actors-cluster/internal/opspb"
)

const (
	version = "1.0.0"
)

// OpsActor represents the global cluster operations actor
// This actor runs as a singleton in the cluster and provides cluster management APIs
type OpsActor struct {
	logger    log.Logger
	remoting  remote.Remoting
	discovery *clusterdiscovery.Discovery
	config    *config.Config
	startTime time.Time
}

// NewOpsActor creates a new instance of OpsActor
func NewOpsActor(
	remoting remote.Remoting,
	logger log.Logger,
	discovery *clusterdiscovery.Discovery,
	config *config.Config,
) *OpsActor {
	return &OpsActor{
		logger:    logger,
		remoting:  remoting,
		discovery: discovery,
		config:    config,
		startTime: time.Now(),
	}
}

// enforce compilation error
var _ goakt.Actor = (*OpsActor)(nil)

// PreStart is used to pre-set initial values for the actor
func (p *OpsActor) PreStart(ctx *goakt.Context) error {
	p.logger.Infof("OpsActor started on node: %s", ctx.ActorName())
	return nil
}

// Receive handles the messages sent to the actor
func (p *OpsActor) Receive(ctx *goakt.ReceiveContext) {
	switch msg := ctx.Message().(type) {
	case *opspb.OpsActorMessage:
		p.handleOpsMessage(ctx, msg)
	default:
		ctx.Unhandled()
	}
}

// handleOpsMessage processes OpsActorMessage requests
func (p *OpsActor) handleOpsMessage(ctx *goakt.ReceiveContext, msg *opspb.OpsActorMessage) {
	switch {
	case msg.GetGetClusterNodes() != nil:
		p.handleGetClusterNodes(ctx, msg.GetGetClusterNodes())
	case msg.GetGetNodeDetails() != nil:
		p.handleGetNodeDetails(ctx, msg.GetGetNodeDetails())
	case msg.GetHealthCheck() != nil:
		p.handleHealthCheck(ctx, msg.GetHealthCheck())
	case msg.GetGetClusterStats() != nil:
		p.handleGetClusterStats(ctx, msg.GetGetClusterStats())
	case msg.GetGetActorDetails() != nil:
		p.handleGetActorDetails(ctx, msg.GetGetActorDetails())
	default:
		ctx.Unhandled()
	}
}

// handleGetClusterNodes handles GetClusterNodes request
func (p *OpsActor) handleGetClusterNodes(ctx *goakt.ReceiveContext, req *opspb.GetClusterNodesRequest) {
	nodes, err := p.getClusterNodes()
	if err != nil {
		ctx.Self().Logger().Errorf("failed to get cluster nodes: %v", err)
		return
	}

	clusterInfo := &opspb.ClusterInfo{
		ClusterName:          p.config.ActorSystemName,
		TotalNodes:           int32(len(nodes)),
		HealthyNodes:         p.countHealthyNodes(nodes),
		UnhealthyNodes:       int32(len(nodes)) - p.countHealthyNodes(nodes),
		ClusterUptimeSeconds: int64(time.Since(p.startTime).Seconds()),
		Nodes:                nodes,
	}

	response := &opspb.OpsActorResponse{
		Response: &opspb.OpsActorResponse_ClusterNodes{
			ClusterNodes: &opspb.GetClusterNodesResponse{
				ClusterInfo: clusterInfo,
			},
		},
	}

	ctx.Response(response)
}

// handleGetNodeDetails handles GetNodeDetails request
func (p *OpsActor) handleGetNodeDetails(ctx *goakt.ReceiveContext, req *opspb.GetNodeDetailsRequest) {
	nodeID := req.GetNodeId()
	nodes, err := p.getClusterNodes()
	if err != nil {
		ctx.Self().Logger().Errorf("failed to get cluster nodes: %v", err)
		return
	}

	for _, node := range nodes {
		if node.NodeId == nodeID {
			response := &opspb.OpsActorResponse{
				Response: &opspb.OpsActorResponse_NodeDetails{
					NodeDetails: &opspb.GetNodeDetailsResponse{
						Node: node,
					},
				},
			}
			ctx.Response(response)
			return
		}
	}

	ctx.Self().Logger().Warnf("node %s not found", nodeID)
}

// handleHealthCheck handles HealthCheck request
func (p *OpsActor) handleHealthCheck(ctx *goakt.ReceiveContext, req *opspb.HealthCheckRequest) {
	hostname, _ := os.Hostname()
	nodeID := fmt.Sprintf("%s:%d", hostname, p.config.GossipPort)

	response := &opspb.OpsActorResponse{
		Response: &opspb.OpsActorResponse_HealthCheck{
			HealthCheck: &opspb.HealthCheckResponse{
				Status:    "healthy",
				Timestamp: time.Now().Format(time.RFC3339),
				Version:   version,
				NodeId:    nodeID,
			},
		},
	}

	ctx.Response(response)
}

// handleGetClusterStats handles GetClusterStats request
func (p *OpsActor) handleGetClusterStats(ctx *goakt.ReceiveContext, req *opspb.ClusterStatsRequest) {
	stats := &opspb.ClusterStatsResponse{
		TotalActors:             0,   // TODO: Implement actor count tracking
		TotalMessagesProcessed:  0,   // TODO: Implement message tracking
		AverageMessageLatencyMs: 0.0, // TODO: Implement latency tracking
		ActivePartitions:        19,  // Fixed partition count from cluster config
		ActorTypeCounts:         p.getActorTypeCounts(),
	}

	response := &opspb.OpsActorResponse{
		Response: &opspb.OpsActorResponse_ClusterStats{
			ClusterStats: stats,
		},
	}

	ctx.Response(response)
}

// handleGetActorDetails handles GetActorDetails request
func (p *OpsActor) handleGetActorDetails(ctx *goakt.ReceiveContext, req *opspb.GetActorDetailsRequest) {
	actorID := req.GetActorId()
	nodeID := req.GetNodeId()

	// Get all actors from the cluster
	actors := p.getClusterActors()

	// Filter actors based on request parameters
	var filteredActors []*opspb.ActorDetails
	for _, actor := range actors {
		// Filter by node ID if specified
		if nodeID != "" && actor.NodeId != nodeID {
			continue
		}
		// Filter by actor ID if specified
		if actorID != "" && actor.ActorId != actorID {
			continue
		}
		filteredActors = append(filteredActors, actor)
	}

	actorDetails := &opspb.GetActorDetailsResponse{
		Actors:     filteredActors,
		TotalCount: int32(len(filteredActors)),
	}

	response := &opspb.OpsActorResponse{
		Response: &opspb.OpsActorResponse_ActorDetails{
			ActorDetails: actorDetails,
		},
	}

	ctx.Response(response)
}

// getClusterNodes retrieves all nodes from service discovery
func (p *OpsActor) getClusterNodes() ([]*opspb.ClusterNode, error) {
	if p.discovery == nil {
		return nil, errors.New("discovery not initialized")
	}

	peers, err := p.discovery.DiscoverPeers()
	if err != nil {
		return nil, errors.Wrap(err, "failed to discover peers")
	}

	var nodes []*opspb.ClusterNode
	hostname, _ := os.Hostname()
	currentNodeID := fmt.Sprintf("%s:%d", hostname, p.config.GossipPort)

	// Add current node
	nodes = append(nodes, &opspb.ClusterNode{
		NodeId:          currentNodeID,
		Host:            hostname,
		Port:            int32(p.config.Port),
		Status:          opspb.NodeStatus_NODE_STATUS_HEALTHY,
		LastHeartbeat:   time.Now().Unix(),
		UptimeSeconds:   int64(time.Since(p.startTime).Seconds()),
		ActorSystemName: p.config.ActorSystemName,
		ServiceName:     p.config.ServiceName,
		GossipPort:      int32(p.config.GossipPort),
		PeersPort:       int32(p.config.PeersPort),
		RemotingPort:    int32(p.config.RemotingPort),
		ActorCount:      0, // TODO: Implement actor count tracking
		Version:         version,
	})

	// Add discovered peers
	for _, peer := range peers {
		if peer != currentNodeID {
			nodes = append(nodes, &opspb.ClusterNode{
				NodeId:          peer,
				Host:            peer,
				Port:            int32(p.config.Port),
				Status:          opspb.NodeStatus_NODE_STATUS_HEALTHY,
				LastHeartbeat:   time.Now().Unix(),
				UptimeSeconds:   int64(time.Since(p.startTime).Seconds()),
				ActorSystemName: p.config.ActorSystemName,
				ServiceName:     p.config.ServiceName,
				GossipPort:      int32(p.config.GossipPort),
				PeersPort:       int32(p.config.PeersPort),
				RemotingPort:    int32(p.config.RemotingPort),
				ActorCount:      0, // TODO: Get actual count from remote nodes
				Version:         version,
			})
		}
	}

	return nodes, nil
}

// getClusterActors retrieves all actors from the cluster
func (p *OpsActor) getClusterActors() []*opspb.ActorDetails {
	var actors []*opspb.ActorDetails

	// Get current node's actors
	currentNodeActors := p.getCurrentNodeActors()

	actors = append(actors, currentNodeActors...)

	// TODO: Get actors from remote nodes in the cluster
	// This would require inter-node communication to query remote actor systems

	return actors
}

// getCurrentNodeActors retrieves actors from the current node
func (p *OpsActor) getCurrentNodeActors() []*opspb.ActorDetails {
	var actors []*opspb.ActorDetails

	// Get hostname for node ID
	hostname, _ := os.Hostname()
	nodeID := fmt.Sprintf("%s:%d", hostname, p.config.GossipPort)

	// TODO: Implement actual actor retrieval from the actor system
	// For now, we'll create sample data to demonstrate the functionality

	// Sample account actors (simulating actual account actors in the system)
	sampleActors := []string{
		"account-001",
		"account-002",
		"account-003",
		"account-004",
		"account-005",
	}

	for i, actorID := range sampleActors {
		actor := &opspb.ActorDetails{
			ActorId:       actorID,
			ActorType:     "Account",
			NodeId:        nodeID,
			Status:        opspb.ActorStatus_ACTOR_STATUS_RUNNING,
			CreatedAt:     time.Now().Add(-time.Duration(i) * time.Hour).Unix(),
			LastActivity:  time.Now().Add(-time.Duration(i) * time.Minute).Unix(),
			MessageCount:  int32(i * 10),
			ErrorCount:    int32(i),
			ParentActorId: "", // Account actors typically don't have parents
			ChildActorIds: []string{},
			Metadata: map[string]string{
				"created_by": "system",
				"version":    "1.0.0",
			},
		}
		actors = append(actors, actor)
	}

	// Add the OpsActor itself
	opsActor := &opspb.ActorDetails{
		ActorId:       "ops-actor",
		ActorType:     "OpsActor",
		NodeId:        nodeID,
		Status:        opspb.ActorStatus_ACTOR_STATUS_RUNNING,
		CreatedAt:     p.startTime.Unix(),
		LastActivity:  time.Now().Unix(),
		MessageCount:  0, // TODO: Track actual message count
		ErrorCount:    0,
		ParentActorId: "",
		ChildActorIds: []string{},
		Metadata: map[string]string{
			"role":        "cluster-management",
			"description": "Cluster operations and monitoring actor",
		},
	}
	actors = append(actors, opsActor)

	return actors
}

// countHealthyNodes counts the number of healthy nodes
func (p *OpsActor) countHealthyNodes(nodes []*opspb.ClusterNode) int32 {
	var count int32
	for _, node := range nodes {
		if node.Status == opspb.NodeStatus_NODE_STATUS_HEALTHY {
			count++
		}
	}
	return count
}

// getActorTypeCounts returns counts of different actor types
func (p *OpsActor) getActorTypeCounts() map[string]int32 {
	// For now, we only have Account actors
	return map[string]int32{
		"Account": 0, // TODO: Implement actor type count tracking
	}
}

// PostStop is used to free-up resources when the actor stops
func (p *OpsActor) PostStop(ctx *goakt.Context) error {
	p.logger.Info("OpsActor stopped")
	return nil
}
