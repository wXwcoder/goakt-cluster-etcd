package main

import (
	"os"
	"time"

	goakt "github.com/tochemey/goakt/v3/actor"
	"github.com/tochemey/goakt/v3/goaktpb"
	"github.com/tochemey/goakt/v3/log"

	"goakt-actors-cluster/internal/opspb"
	"goakt-actors-cluster/internal/samplepb"
)

// ClusterInteractionActor 是一个简单的集群交互 Actor，用于演示目的
// 现在客户端直接使用 RemoteAsk 与目标 Actor 通信，这个 Actor 主要用于演示
type ClusterInteractionActor struct {
	logger log.Logger
}

// 强制编译时检查，确保实现了 Actor 接口
var _ goakt.Actor = (*ClusterInteractionActor)(nil)

// NewClusterInteractionActor 创建新的集群交互 Actor 实例
func NewClusterInteractionActor() *ClusterInteractionActor {
	return &ClusterInteractionActor{}
}

// PreStart 在 Actor 启动前执行
func (c *ClusterInteractionActor) PreStart(ctx *goakt.Context) error {
	logger := log.New(log.InfoLevel, os.Stdout)
	c.logger = logger
	c.logger.Info("Cluster Interaction Actor 已启动")
	return nil
}

// Receive 处理接收到的消息
func (c *ClusterInteractionActor) Receive(ctx *goakt.ReceiveContext) {
	switch msg := ctx.Message().(type) {
	case *goaktpb.PostStart:
		c.logger.Info("Cluster Interaction Actor 成功启动")

	// 处理直接接收到的消息（用于演示）
	case *samplepb.CreateAccount:
		c.logger.Infof("收到创建账户请求: %s", msg.AccountId)
		c.handleCreateAccount(ctx, msg)

	case *samplepb.GetAccount:
		c.logger.Infof("收到获取账户请求: %s", msg.AccountId)
		c.handleGetAccount(ctx, msg)

	case *samplepb.CreditAccount:
		c.logger.Infof("收到存款操作请求: %s", msg.AccountId)
		c.handleCreditAccount(ctx, msg)

	case *opspb.HealthCheckRequest:
		c.logger.Info("收到健康检查请求")
		c.handleHealthCheck(ctx, msg)

	case *opspb.GetClusterNodesRequest:
		c.logger.Info("收到获取集群节点请求")
		c.handleGetClusterNodes(ctx, msg)

	case *opspb.ClusterStatsRequest:
		c.logger.Info("收到获取集群统计信息请求")
		c.handleGetClusterStats(ctx, msg)

	case *opspb.GetActorDetailsRequest:
		c.logger.Infof("收到获取 Actor 详情请求: %s", msg.ActorId)
		c.handleGetActorDetails(ctx, msg)

	case *opspb.GetNodeDetailsRequest:
		c.logger.Info("收到获取节点详情请求")
		c.handleGetNodeDetails(ctx, msg)

	default:
		ctx.Unhandled()
	}
}

// PostStop 在 Actor 停止后执行
func (c *ClusterInteractionActor) PostStop(ctx *goakt.Context) error {
	c.logger.Info("Cluster Interaction Actor 已停止")
	return nil
}

// 处理直接接收到的消息（用于演示）
func (c *ClusterInteractionActor) handleCreateAccount(ctx *goakt.ReceiveContext, msg *samplepb.CreateAccount) {
	c.logger.Infof("处理创建账户请求: %s, 余额: %.2f", msg.AccountId, msg.AccountBalance)

	// 这里可以添加处理逻辑，例如转发到目标 Actor 或直接处理
	// 目前只是记录日志并返回成功响应
	response := &samplepb.AccountCreated{
		AccountId:      msg.AccountId,
		AccountBalance: msg.AccountBalance,
	}

	ctx.Sender().Tell(ctx.Context(), ctx.Self(), response)
}

func (c *ClusterInteractionActor) handleGetAccount(ctx *goakt.ReceiveContext, msg *samplepb.GetAccount) {
	c.logger.Infof("处理获取账户请求: %s", msg.AccountId)

	// 这里可以添加处理逻辑，例如从数据库查询账户信息
	// 目前只是返回模拟数据
	response := &samplepb.Account{
		AccountId:      msg.AccountId,
		AccountBalance: 1000.0, // 模拟余额
	}

	ctx.Sender().Tell(ctx.Context(), ctx.Self(), response)
}

func (c *ClusterInteractionActor) handleCreditAccount(ctx *goakt.ReceiveContext, msg *samplepb.CreditAccount) {
	c.logger.Infof("处理存款操作请求: %s, 金额: %.2f", msg.AccountId, msg.Balance)

	// 这里可以添加处理逻辑，例如更新账户余额
	// 目前只是返回模拟数据
	response := &samplepb.AccountCredited{
		AccountId:      msg.AccountId,
		AccountBalance: 1500.0, // 模拟新余额
	}

	ctx.Sender().Tell(ctx.Context(), ctx.Self(), response)
}

func (c *ClusterInteractionActor) handleHealthCheck(ctx *goakt.ReceiveContext, msg *opspb.HealthCheckRequest) {
	c.logger.Info("处理健康检查请求")

	response := &opspb.HealthCheckResponse{
		Status: "HEALTHY",
	}

	ctx.Sender().Tell(ctx.Context(), ctx.Self(), response)
}

func (c *ClusterInteractionActor) handleGetClusterNodes(ctx *goakt.ReceiveContext, msg *opspb.GetClusterNodesRequest) {
	c.logger.Info("处理获取集群节点请求")

	// 返回模拟的集群节点信息
	nodes := []*opspb.ClusterNode{
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
	}

	// 创建集群信息
	clusterInfo := &opspb.ClusterInfo{
		ClusterName:          "TestCluster",
		TotalNodes:           2,
		HealthyNodes:         2,
		UnhealthyNodes:       0,
		ClusterUptimeSeconds: 3600,
		Nodes:                nodes,
	}

	response := &opspb.GetClusterNodesResponse{
		ClusterInfo: clusterInfo,
	}

	ctx.Sender().Tell(ctx.Context(), ctx.Self(), response)
}

func (c *ClusterInteractionActor) handleGetClusterStats(ctx *goakt.ReceiveContext, msg *opspb.ClusterStatsRequest) {
	c.logger.Info("处理获取集群统计信息请求")

	response := &opspb.ClusterStatsResponse{
		TotalActors:             10,
		TotalMessagesProcessed:  100,
		AverageMessageLatencyMs: 5.0,
		ActivePartitions:        5,
		ActorTypeCounts:         map[string]int32{"Account": 8, "OpsActor": 2},
	}

	ctx.Sender().Tell(ctx.Context(), ctx.Self(), response)
}

func (c *ClusterInteractionActor) handleGetActorDetails(ctx *goakt.ReceiveContext, msg *opspb.GetActorDetailsRequest) {
	c.logger.Infof("处理获取 Actor 详情请求: %s", msg.ActorId)

	// 创建模拟的 Actor 详情
	actorDetails := &opspb.ActorDetails{
		ActorId:      msg.ActorId,
		ActorType:    "AccountActor",
		NodeId:       "node-1",
		Status:       opspb.ActorStatus_ACTOR_STATUS_RUNNING,
		CreatedAt:    time.Now().Unix(),
		LastActivity: time.Now().Unix(),
		MessageCount: 10,
		ErrorCount:   0,
	}

	response := &opspb.GetActorDetailsResponse{
		Actors:     []*opspb.ActorDetails{actorDetails},
		TotalCount: 1,
	}

	ctx.Sender().Tell(ctx.Context(), ctx.Self(), response)
}

func (c *ClusterInteractionActor) handleGetNodeDetails(ctx *goakt.ReceiveContext, msg *opspb.GetNodeDetailsRequest) {
	c.logger.Info("处理获取节点详情请求")

	// 创建模拟的节点信息
	node := &opspb.ClusterNode{
		NodeId:          msg.NodeId,
		Host:            "127.0.0.1",
		Port:            14001,
		Status:          opspb.NodeStatus_NODE_STATUS_HEALTHY,
		LastHeartbeat:   time.Now().Unix(),
		UptimeSeconds:   3600,
		ActorSystemName: "ClusterSystem",
		ServiceName:     "account-service",
		GossipPort:      8558,
		PeersPort:       2552,
		RemotingPort:    14001,
		ActorCount:      5,
		Version:         "1.0.0",
	}

	response := &opspb.GetNodeDetailsResponse{
		Node: node,
	}

	ctx.Sender().Tell(ctx.Context(), ctx.Self(), response)
}
