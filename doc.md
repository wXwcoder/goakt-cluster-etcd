# How to run it?

## Prerequisites

1. Install [Earthly](https://earthly.dev/get-earthly)
2. Clone the repository
3. Start etcd cluster (see below for setup instructions)

## Using Etcd for Dynamic Service Discovery

This example uses etcd for dynamic service discovery instead of static configuration. This allows for dynamic node scaling without restarting the entire cluster.

### Environment Variables

Copy `.env.example` to `.env` and configure the following variables:

```bash
# Service Configuration
PORT=50051
SERVICE_NAME=account-service
SYSTEM_NAME=goakt-cluster
GOSSIP_PORT=8558
PEERS_PORT=2552
REMOTING_PORT=8080

# Etcd Configuration
ETCD_ENDPOINTS=localhost:2379
ETCD_DIAL_TIMEOUT=5
ETCD_TTL=30
ETCD_TIMEOUT=5
```

### Starting Etcd Cluster

You can use the provided docker-compose file to start a local etcd cluster:

```bash
docker-compose up -d etcd
```

### Running the Application

1. Build the application:
   ```bash
   earthly +etcd-image
   ```

2. Start the cluster nodes:
   ```bash
   docker-compose up node1 node2 node3
   ```

3. To stop the cluster:
   ```bash
   docker-compose down -v --remove-orphans
   ```

4. Check running instances:
   ```bash
   docker-compose ps
   ```

### Dynamic Scaling

With etcd-based discovery, you can dynamically add or remove nodes:

- Add a new node: Start a new container with unique environment variables
- Remove a node: Stop the container, and etcd will automatically remove it from discovery
- The cluster will automatically rebalance actors across available nodes

### Service Access

With any gRPC client you can access the service. The service definitions are [here](../../../protos/sample/pb/v1)
