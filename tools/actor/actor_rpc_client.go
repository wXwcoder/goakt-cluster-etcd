package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"connectrpc.com/connect"
	"github.com/tochemey/goakt/v3/actor"
	"github.com/tochemey/goakt/v3/log"
	"github.com/tochemey/goakt/v3/remote"
	"github.com/tochemey/goakt/v3/supervisor"

	"goakt-actors-cluster/internal/opspb"
	"goakt-actors-cluster/internal/opspb/opspbconnect"
	"goakt-actors-cluster/internal/samplepb"
	"goakt-actors-cluster/internal/samplepb/samplepbconnect"
)

// ActorRPCClient 使用 Actor RPC 接口与集群交互的客户端
type ActorRPCClient struct {
	actorSystem      actor.ActorSystem
	remoting         remote.Remoting
	interactionActor *actor.PID
	logger           log.Logger
	accountClient    samplepbconnect.AccountServiceClient
	opsClient        opspbconnect.OpsServiceClient
}

// NewActorRPCClient 创建新的 Actor RPC 客户端
func NewActorRPCClient() (*ActorRPCClient, error) {
	ctx := context.Background()

	// 配置 Actor 系统
	logger := log.New(log.InfoLevel, os.Stdout)

	actorSystem, err := actor.NewActorSystem(
		"TestClientSystem",
		actor.WithRemote(remote.NewConfig("127.0.0.1", 0)), // 使用随机端口
		actor.WithLogger(logger),
	)
	if err != nil {
		return nil, fmt.Errorf("创建 Actor 系统失败: %v", err)
	}

	// 启动 Actor 系统
	if err := actorSystem.Start(ctx); err != nil {
		return nil, fmt.Errorf("启动 Actor 系统失败: %v", err)
	}

	// 等待 Actor 系统启动
	time.Sleep(time.Second)

	// 创建 remoting 实例
	remoting := remote.NewRemoting()

	// 创建集群交互 Actor
	interactionActor, err := actorSystem.Spawn(
		ctx,
		"ClusterInteractionActor",
		NewClusterInteractionActor(),
		actor.WithSupervisor(
			supervisor.NewSupervisor(
				supervisor.WithStrategy(supervisor.OneForOneStrategy),
				supervisor.WithAnyErrorDirective(supervisor.ResumeDirective),
			)),
	)
	if err != nil {
		return nil, fmt.Errorf("创建集群交互 Actor 失败: %v", err)
	}

	// 等待 Actor 启动
	time.Sleep(time.Second)

	// 创建 HTTP 客户端
	httpClient := &http.Client{Timeout: 10 * time.Second}

	// 创建 Connect RPC 客户端
	accountClient := samplepbconnect.NewAccountServiceClient(httpClient, "http://localhost:14001")
	opsClient := opspbconnect.NewOpsServiceClient(httpClient, "http://localhost:14001")

	return &ActorRPCClient{
		actorSystem:      actorSystem,
		remoting:         remoting,
		interactionActor: interactionActor,
		logger:           logger,
		accountClient:    accountClient,
		opsClient:        opsClient,
	}, nil
}

// Close 关闭客户端
func (c *ActorRPCClient) Close() error {
	ctx := context.Background()
	c.remoting.Close()
	return c.actorSystem.Stop(ctx)
}

// TestAccountActorInteractions 测试 Account Actor 的 RPC 交互
func (c *ActorRPCClient) TestAccountActorInteractions() {
	fmt.Println("=== 测试 Account Actor RPC 交互 ===")

	// 测试创建账户
	c.testCreateAccount("account-001", 1000.0)

	// 测试查询账户
	c.testGetAccount("account-001")

	// 测试存款操作
	c.testCreditAccount("account-001", 500.0)

	// 再次查询账户
	c.testGetAccount("account-001")

	// 测试另一个账户
	c.testCreateAccount("account-002", 2000.0)
	c.testGetAccount("account-002")
}

// testCreateAccount 测试创建账户
func (c *ActorRPCClient) testCreateAccount(accountID string, balance float64) {
	fmt.Printf("创建账户: %s, 初始余额: %.2f\n", accountID, balance)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 创建请求
	request := &samplepb.CreateAccountRequest{
		CreateAccount: &samplepb.CreateAccount{
			AccountId:      accountID,
			AccountBalance: balance,
		},
	}

	// 使用 Connect RPC 客户端发送请求
	resp, err := c.accountClient.CreateAccount(ctx, connect.NewRequest(request))
	if err != nil {
		fmt.Printf("创建账户失败: %v\n", err)
		return
	}

	// 获取响应
	account := resp.Msg.GetAccount()
	fmt.Printf("账户创建成功: ID=%s, 余额=%.2f\n", account.AccountId, account.AccountBalance)
	fmt.Println()
}

// testGetAccount 测试查询账户
func (c *ActorRPCClient) testGetAccount(accountID string) {
	fmt.Printf("查询账户: %s\n", accountID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 创建请求
	request := &samplepb.GetAccountRequest{
		AccountId: accountID,
	}

	// 使用 Connect RPC 客户端发送请求
	resp, err := c.accountClient.GetAccount(ctx, connect.NewRequest(request))
	if err != nil {
		fmt.Printf("查询账户失败: %v\n", err)
		return
	}

	// 获取响应
	account := resp.Msg.GetAccount()
	fmt.Printf("账户查询成功: ID=%s, 余额=%.2f\n", account.AccountId, account.AccountBalance)
	fmt.Println()
}

// testCreditAccount 测试存款操作
func (c *ActorRPCClient) testCreditAccount(accountID string, amount float64) {
	fmt.Printf("存款到账户: %s, 金额: %.2f\n", accountID, amount)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 创建请求
	request := &samplepb.CreditAccountRequest{
		CreditAccount: &samplepb.CreditAccount{
			AccountId: accountID,
			Balance:   amount,
		},
	}

	// 使用 Connect RPC 客户端发送请求
	resp, err := c.accountClient.CreditAccount(ctx, connect.NewRequest(request))
	if err != nil {
		fmt.Printf("存款操作失败: %v\n", err)
		return
	}

	// 获取响应
	account := resp.Msg.GetAccount()
	fmt.Printf("存款操作成功: ID=%s, 新余额=%.2f\n", account.AccountId, account.AccountBalance)
	fmt.Println()
}

// TestOpsActorInteractions 测试 OpsActor 的 RPC 交互
func (c *ActorRPCClient) TestOpsActorInteractions() {
	fmt.Println("=== 测试 OpsActor RPC 交互 ===")

	// 测试健康检查
	c.testHealthCheck()

	// 测试获取集群节点
	c.testGetClusterNodes()

	// 测试获取集群统计信息
	c.testGetClusterStats()

	// 测试获取 Actor 详情
	c.testGetActorDetails("account-001")
}

// testHealthCheck 测试健康检查
func (c *ActorRPCClient) testHealthCheck() {
	fmt.Println("执行健康检查")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 创建请求
	request := &opspb.HealthCheckRequest{}

	// 使用 Connect RPC 客户端发送请求
	resp, err := c.opsClient.HealthCheck(ctx, connect.NewRequest(request))
	if err != nil {
		fmt.Printf("健康检查失败: %v\n", err)
		return
	}

	// 获取响应
	response := resp.Msg
	fmt.Printf("健康检查成功: 状态=%s\n", response.Status)
	fmt.Println()
}

// testGetClusterNodes 测试获取集群节点
func (c *ActorRPCClient) testGetClusterNodes() {
	fmt.Println("获取集群节点")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 创建请求
	request := &opspb.GetClusterNodesRequest{}

	// 使用 Connect RPC 客户端发送请求
	resp, err := c.opsClient.GetClusterNodes(ctx, connect.NewRequest(request))
	if err != nil {
		fmt.Printf("获取集群节点失败: %v\n", err)
		return
	}

	// 获取响应
	response := resp.Msg
	// 检查集群信息是否存在
	if response.ClusterInfo == nil {
		fmt.Printf("获取集群节点成功: 集群信息为空\n")
	} else {
		fmt.Printf("获取集群节点成功: 集群名称=%s\n", response.ClusterInfo.ClusterName)
	}
	fmt.Println()
}

// testGetClusterStats 测试获取集群统计信息
func (c *ActorRPCClient) testGetClusterStats() {
	fmt.Println("获取集群统计信息")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 创建请求
	request := &opspb.ClusterStatsRequest{}

	// 使用 Connect RPC 客户端发送请求
	resp, err := c.opsClient.GetClusterStats(ctx, connect.NewRequest(request))
	if err != nil {
		fmt.Printf("获取集群统计信息失败: %v\n", err)
		return
	}

	// 获取响应
	response := resp.Msg
	fmt.Printf("获取集群统计信息成功: 总Actor数=%d\n", response.TotalActors)
	fmt.Println()
}

// testGetActorDetails 测试获取 Actor 详情
func (c *ActorRPCClient) testGetActorDetails(actorID string) {
	fmt.Printf("获取 Actor 详情: %s\n", actorID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 创建请求
	request := &opspb.GetActorDetailsRequest{
		ActorId: actorID,
	}

	// 使用 Connect RPC 客户端发送请求
	resp, err := c.opsClient.GetActorDetails(ctx, connect.NewRequest(request))
	if err != nil {
		fmt.Printf("获取 Actor 详情失败: %v\n", err)
		return
	}

	// 获取响应
	response := resp.Msg
	fmt.Printf("获取 Actor 详情成功: 总数=%d\n", response.TotalCount)
	fmt.Println()
}

// main 函数 - 测试 Actor RPC 客户端
func main() {
	fmt.Println("=== 启动 Actor RPC 客户端 ===")

	// 创建客户端
	client, err := NewActorRPCClient()
	if err != nil {
		fmt.Printf("创建客户端失败: %v\n", err)
		return
	}
	defer client.Close()

	fmt.Println("客户端已启动，开始测试...")

	// 测试 OpsActor 交互
	client.TestOpsActorInteractions()

	// 测试 Account Actor 交互
	client.TestAccountActorInteractions()

	fmt.Println("=== 测试完成 ===")

	// 等待用户中断
	fmt.Println("按 Ctrl+C 退出...")
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	fmt.Println("客户端已关闭")
}
