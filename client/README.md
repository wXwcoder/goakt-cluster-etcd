# GoAkt Cluster Client

This is a client tool for interacting with the GoAkt cluster using the official goakt client implementation patterns. The client provides a command-line interface and a Go API for interacting with the cluster.

## Features

Based on the official GoAkt client documentation, this client provides:

- **Kinds**: List all the actor kinds in the cluster
- **Spawn**: Spawn an actor in the cluster
- **SpawnWithBalancer**: Spawn an actor in the cluster with a given balancer strategy
- **Stop**: Kill/stop an actor in the cluster
- **Ask**: Send a message to a given actor in the cluster and expect a response
- **Tell**: Send a fire-forget message to a given actor in the cluster
- **Whereis**: Locate and get the address of a given actor
- **ReSpawn**: Restart a given actor
- **AskGrain**: Send a request/response to a Grain
- **TellGrain**: Send a fire-forget message to a given Grain

## Installation

```bash
cd client
go mod tidy
```

## Usage

### Command Line Interface

The client provides a command-line interface for common operations:

#### Create an Account
```bash
go run main.go create-account --id "account-001" --balance 1000.0 --endpoints "localhost:14001"
```

#### Get Account Details
```bash
go run main.go get-account --id "account-001" --endpoints "localhost:14001"
```

#### Credit an Account
```bash
go run main.go credit-account --id "account-001" --amount 500.0 --endpoints "localhost:14001"
```

#### Debit an Account
```bash
go run main.go debit-account --id "account-001" --amount 200.0 --endpoints "localhost:14001"
```

#### Get Cluster Nodes Information
```bash
go run main.go get-nodes --endpoints "localhost:14001"
```

#### Get Cluster Statistics
```bash
go run main.go get-stats --endpoints "localhost:14001"
```

#### Perform Health Check
```bash
go run main.go health-check --endpoints "localhost:14001"
```

### Go API Usage

You can also use the client as a Go library in your own applications:

```go
package main

import (
    "context"
    "fmt"
    "goakt-actors-cluster/client"
)

func main() {
    // Create a new client
    clusterClient, err := client.NewClusterClient([]string{"localhost:14001"})
    if err != nil {
        panic(err)
    }

    // Connect to the cluster
    ctx := context.Background()
    if err := clusterClient.Connect(ctx); err != nil {
        panic(err)
    }
    defer clusterClient.Disconnect(ctx)

    // Create an account
    if err := clusterClient.CreateAccount(ctx, "test-account", 1000.0); err != nil {
        panic(err)
    }

    // Get account details
    account, err := clusterClient.GetAccount(ctx, "test-account")
    if err != nil {
        panic(err)
    }

    fmt.Printf("Account: %s, Balance: %.2f\n", account.AccountId, account.AccountBalance)
}
```

## Configuration

The client can be configured with:

- **Endpoints**: A list of cluster endpoints to connect to
- **Timeout**: Request timeout duration
- **Logger**: Custom logger implementation

## Architecture

The client follows the official GoAkt client patterns:

- **Mini Load Balancer**: The client includes a mini load balancer that helps route requests to the appropriate node using strategies like:
  - Round Robin: A given node is chosen using the round-robin strategy
  - Random: A given node is chosen randomly
  - Least Load: The node with the least number of actors is chosen

## Limitations

- Currently, actors created using the cluster client are created using the default mailbox and supervisor strategy (custom spawning not supported yet)
- Some operations return mock data in this example implementation
- Production implementations would need to properly handle remote communication with the cluster