package main

import (
	"context"
	"fmt"
	"os"

	"goakt-actors-cluster/client/pkg"

	"github.com/spf13/cobra"
	"github.com/tochemey/goakt/v3/remote"
)

var (
	endpoints    []string = []string{"localhost:14001"}
	accountID    string   = "account01"
	balance      float64  = 1000.00
	amount       float64  = 500.00
	host         string   = "localhost"
	remotingPort int      = 14002
	balancerType string   = "random"
)

var rootCmd = &cobra.Command{
	Use:   "cluster-client",
	Short: "A client for interacting with the GoAkt cluster",
	Long:  `A command-line client tool for interacting with the GoAkt cluster using the official goakt client implementation.`,
}

var createAccountCmd = &cobra.Command{
	Use:   "create-account",
	Short: "Create a new account in the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		if accountID == "" {
			fmt.Println("Error: account ID is required")
			os.Exit(1)
		}

		// Create remote configuration
		remoteConfig := remote.NewConfig(host, remotingPort)

		client, err := pkg.NewFullClusterClient(endpoints, remoteConfig)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		// Set the balancer strategy based on the flag
		strategy := getBalancerStrategy(balancerType)
		client.SetBalancerStrategy(strategy)

		ctx := context.Background()
		if err := client.Connect(ctx); err != nil {
			fmt.Printf("Error connecting to cluster: %v\n", err)
			os.Exit(1)
		}
		defer client.Disconnect(ctx)

		if err := client.CreateAccount(ctx, accountID, balance); err != nil {
			fmt.Printf("Error creating account: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Account %s created successfully with balance: %.2f\n", accountID, balance)
	},
}

var getAccountCmd = &cobra.Command{
	Use:   "get-account",
	Short: "Get account details from the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		if accountID == "" {
			fmt.Println("Error: account ID is required")
			os.Exit(1)
		}

		// Create remote configuration
		remoteConfig := remote.NewConfig(host, remotingPort)

		client, err := pkg.NewFullClusterClient(endpoints, remoteConfig)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		// Set the balancer strategy based on the flag
		strategy := getBalancerStrategy(balancerType)
		client.SetBalancerStrategy(strategy)

		ctx := context.Background()
		if err := client.Connect(ctx); err != nil {
			fmt.Printf("Error connecting to cluster: %v\n", err)
			os.Exit(1)
		}
		defer client.Disconnect(ctx)

		account, err := client.GetAccount(ctx, accountID)
		if err != nil {
			fmt.Printf("Error getting account: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Account Details:\n")
		fmt.Printf("  ID: %s\n", account.AccountId)
		fmt.Printf("  Balance: %.2f\n", account.AccountBalance)
	},
}

var creditAccountCmd = &cobra.Command{
	Use:   "credit-account",
	Short: "Credit funds to an account in the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		if accountID == "" {
			fmt.Println("Error: account ID is required")
			os.Exit(1)
		}

		// Create remote configuration
		remoteConfig := remote.NewConfig(host, remotingPort)

		client, err := pkg.NewFullClusterClient(endpoints, remoteConfig)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		// Set the balancer strategy based on the flag
		strategy := getBalancerStrategy(balancerType)
		client.SetBalancerStrategy(strategy)

		ctx := context.Background()
		if err := client.Connect(ctx); err != nil {
			fmt.Printf("Error connecting to cluster: %v\n", err)
			os.Exit(1)
		}
		defer client.Disconnect(ctx)

		if err := client.CreditAccount(ctx, accountID, amount); err != nil {
			fmt.Printf("Error crediting account: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Account %s credited with: %.2f\n", accountID, amount)
	},
}

var debitAccountCmd = &cobra.Command{
	Use:   "debit-account",
	Short: "Debit funds from an account in the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		if accountID == "" {
			fmt.Println("Error: account ID is required")
			os.Exit(1)
		}

		// Create remote configuration
		remoteConfig := remote.NewConfig(host, remotingPort)

		client, err := pkg.NewFullClusterClient(endpoints, remoteConfig)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		// Set the balancer strategy based on the flag
		strategy := getBalancerStrategy(balancerType)
		client.SetBalancerStrategy(strategy)

		ctx := context.Background()
		if err := client.Connect(ctx); err != nil {
			fmt.Printf("Error connecting to cluster: %v\n", err)
			os.Exit(1)
		}
		defer client.Disconnect(ctx)

		if err := client.DebitAccount(ctx, accountID, amount); err != nil {
			fmt.Printf("Error debiting account: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Account %s debited with: %.2f\n", accountID, amount)
	},
}

var getClusterNodesCmd = &cobra.Command{
	Use:   "get-nodes",
	Short: "Get cluster node information",
	Run: func(cmd *cobra.Command, args []string) {
		// Create remote configuration
		remoteConfig := remote.NewConfig(host, remotingPort)

		client, err := pkg.NewFullClusterClient(endpoints, remoteConfig)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		// Set the balancer strategy based on the flag
		strategy := getBalancerStrategy(balancerType)
		client.SetBalancerStrategy(strategy)

		ctx := context.Background()
		if err := client.Connect(ctx); err != nil {
			fmt.Printf("Error connecting to cluster: %v\n", err)
			os.Exit(1)
		}
		defer client.Disconnect(ctx)

		nodes, err := client.GetClusterNodes(ctx)
		if err != nil {
			fmt.Printf("Error getting cluster nodes: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Cluster Information:\n")
		fmt.Printf("  Name: %s\n", nodes.ClusterInfo.ClusterName)
		fmt.Printf("  Total Nodes: %d\n", nodes.ClusterInfo.TotalNodes)
		fmt.Printf("  Healthy Nodes: %d\n", nodes.ClusterInfo.HealthyNodes)
		fmt.Printf("  Unhealthy Nodes: %d\n", nodes.ClusterInfo.UnhealthyNodes)

		fmt.Printf("\nNodes:\n")
		for _, node := range nodes.ClusterInfo.Nodes {
			fmt.Printf("  - ID: %s, Host: %s:%d, Status: %s\n",
				node.NodeId, node.Host, node.Port, node.Status.String())
		}
	},
}

var getClusterStatsCmd = &cobra.Command{
	Use:   "get-stats",
	Short: "Get cluster statistics",
	Run: func(cmd *cobra.Command, args []string) {
		// Create remote configuration
		remoteConfig := remote.NewConfig(host, remotingPort)

		client, err := pkg.NewFullClusterClient(endpoints, remoteConfig)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		// Set the balancer strategy based on the flag
		strategy := getBalancerStrategy(balancerType)
		client.SetBalancerStrategy(strategy)

		ctx := context.Background()
		if err := client.Connect(ctx); err != nil {
			fmt.Printf("Error connecting to cluster: %v\n", err)
			os.Exit(1)
		}
		defer client.Disconnect(ctx)

		stats, err := client.GetClusterStats(ctx)
		if err != nil {
			fmt.Printf("Error getting cluster stats: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Cluster Statistics:\n")
		fmt.Printf("  Total Actors: %d\n", stats.TotalActors)
		fmt.Printf("  Total Messages Processed: %d\n", stats.TotalMessagesProcessed)
		fmt.Printf("  Average Message Latency: %.2f ms\n", stats.AverageMessageLatencyMs)
		fmt.Printf("  Active Partitions: %d\n", stats.ActivePartitions)
		fmt.Printf("  Actor Type Counts: %v\n", stats.ActorTypeCounts)
	},
}

var healthCheckCmd = &cobra.Command{
	Use:   "health-check",
	Short: "Perform a health check on the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		// Create remote configuration
		remoteConfig := remote.NewConfig(host, remotingPort)

		client, err := pkg.NewFullClusterClient(endpoints, remoteConfig)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		// Set the balancer strategy based on the flag
		strategy := getBalancerStrategy(balancerType)
		client.SetBalancerStrategy(strategy)

		ctx := context.Background()
		if err := client.Connect(ctx); err != nil {
			fmt.Printf("Error connecting to cluster: %v\n", err)
			os.Exit(1)
		}
		defer client.Disconnect(ctx)

		result, err := client.HealthCheck(ctx)
		if err != nil {
			fmt.Printf("Error performing health check: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Cluster Health: %s\n", result.Status)
	},
}

var listKindsCmd = &cobra.Command{
	Use:   "list-kinds",
	Short: "List all actor kinds in the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		// Create remote configuration
		remoteConfig := remote.NewConfig(host, remotingPort)

		client, err := pkg.NewFullClusterClient(endpoints, remoteConfig)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		// Set the balancer strategy based on the flag
		strategy := getBalancerStrategy(balancerType)
		client.SetBalancerStrategy(strategy)

		ctx := context.Background()
		if err := client.Connect(ctx); err != nil {
			fmt.Printf("Error connecting to cluster: %v\n", err)
			os.Exit(1)
		}
		defer client.Disconnect(ctx)

		kinds, err := client.ListActorKinds(ctx)
		if err != nil {
			fmt.Printf("Error listing actor kinds: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Actor Kinds in Cluster:\n")
		for _, kind := range kinds {
			fmt.Printf("  - %s\n", kind)
		}
	},
}

var whereIsCmd = &cobra.Command{
	Use:   "whereis",
	Short: "Find the location of an actor in the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		if accountID == "" {
			fmt.Println("Error: actor name is required")
			os.Exit(1)
		}

		// Create remote configuration
		remoteConfig := remote.NewConfig(host, remotingPort)

		client, err := pkg.NewFullClusterClient(endpoints, remoteConfig)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		// Set the balancer strategy based on the flag
		strategy := getBalancerStrategy(balancerType)
		client.SetBalancerStrategy(strategy)

		ctx := context.Background()
		if err := client.Connect(ctx); err != nil {
			fmt.Printf("Error connecting to cluster: %v\n", err)
			os.Exit(1)
		}
		defer client.Disconnect(ctx)

		address, err := client.WhereIsActor(ctx, accountID)
		if err != nil {
			fmt.Printf("Error finding actor: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Actor %s is located at: %s\n", accountID, address.String())
	},
}

func getBalancerStrategy(strategy string) pkg.BalancerStrategy {
	switch strategy {
	case "random":
		return pkg.RandomStrategy
	case "least-load":
		return pkg.LeastLoadStrategy
	default:
		return pkg.RoundRobinStrategy
	}
}

func init() {
	// Add flags for all commands
	flags := []struct {
		cmd *cobra.Command
	}{
		{createAccountCmd},
		{getAccountCmd},
		{creditAccountCmd},
		{debitAccountCmd},
		{getClusterNodesCmd},
		{getClusterStatsCmd},
		{healthCheckCmd},
		{listKindsCmd},
		{whereIsCmd},
	}

	for _, f := range flags {
		f.cmd.PersistentFlags().StringSliceVar(&endpoints, "endpoints", []string{"localhost:14001"}, "Cluster endpoints")
		f.cmd.PersistentFlags().StringVar(&host, "host", "localhost", "Host for remoting")
		f.cmd.PersistentFlags().IntVar(&remotingPort, "port", 14001, "Port for remoting")
	}

	// Add specific flags for account commands
	createAccountCmd.Flags().StringVar(&accountID, "id", "", "Account ID (required)")
	createAccountCmd.Flags().Float64Var(&balance, "balance", 0.0, "Initial balance")

	getAccountCmd.Flags().StringVar(&accountID, "id", "", "Account ID (required)")

	creditAccountCmd.Flags().StringVar(&accountID, "id", "", "Account ID (required)")
	creditAccountCmd.Flags().Float64Var(&amount, "amount", 0.0, "Amount to credit")

	debitAccountCmd.Flags().StringVar(&accountID, "id", "", "Account ID (required)")
	debitAccountCmd.Flags().Float64Var(&amount, "amount", 0.0, "Amount to debit")

	whereIsCmd.Flags().StringVar(&accountID, "name", "", "Actor name (required)")

	// Add balancer type flag to all commands
	balancerFlags := []struct {
		cmd *cobra.Command
	}{
		{createAccountCmd},
		{getAccountCmd},
		{creditAccountCmd},
		{debitAccountCmd},
		{getClusterNodesCmd},
		{getClusterStatsCmd},
		{healthCheckCmd},
		{listKindsCmd},
		{whereIsCmd},
	}

	for _, f := range balancerFlags {
		f.cmd.Flags().StringVar(&balancerType, "balancer", "round-robin", "Load balancer strategy (round-robin, random, least-load)")
	}

	// Add commands to root
	rootCmd.AddCommand(createAccountCmd)
	rootCmd.AddCommand(getAccountCmd)
	rootCmd.AddCommand(creditAccountCmd)
	rootCmd.AddCommand(debitAccountCmd)
	rootCmd.AddCommand(getClusterNodesCmd)
	rootCmd.AddCommand(getClusterStatsCmd)
	rootCmd.AddCommand(healthCheckCmd)
	rootCmd.AddCommand(listKindsCmd)
	rootCmd.AddCommand(whereIsCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
