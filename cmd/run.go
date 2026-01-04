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

package cmd

import (
	"context"
	"goakt-actors-cluster/actors"
	"goakt-actors-cluster/config"
	"goakt-actors-cluster/discovery"
	"goakt-actors-cluster/service"
	"os"
	"time"

	"github.com/spf13/cobra"
	goakt "github.com/tochemey/goakt/v3/actor"
	"github.com/tochemey/goakt/v3/log"
	"github.com/tochemey/goakt/v3/remote"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the actor system with service discovery",
	Long:  `Run the actor system with etcd-based service discovery. Configuration can be provided via environment variables or command line flags.`,
	Run: func(cmd *cobra.Command, args []string) {
		// create a background context
		ctx := context.Background()

		// get command line parameters
		port, _ := cmd.Flags().GetInt("port")
		serviceName, _ := cmd.Flags().GetString("service-name")
		systemName, _ := cmd.Flags().GetString("system-name")
		gossipPort, _ := cmd.Flags().GetInt("gossip-port")
		peersPort, _ := cmd.Flags().GetInt("peers-port")
		remotingPort, _ := cmd.Flags().GetInt("remoting-port")
		etcdEndpoints, _ := cmd.Flags().GetStringSlice("etcd-endpoints")
		etcdDialTimeout, _ := cmd.Flags().GetInt("etcd-dial-timeout")
		etcdTTL, _ := cmd.Flags().GetInt64("etcd-ttl")
		etcdTimeout, _ := cmd.Flags().GetInt("etcd-timeout")

		// get the configuration with command line overrides
		config, err := config.GetConfigWithOverrides(
			port, serviceName, systemName, gossipPort, peersPort, remotingPort,
			etcdEndpoints, etcdDialTimeout, etcdTTL, etcdTimeout,
		)
		//  handle the error
		if err != nil {
			panic(err)
		}
		// use the address default log. real-life implement the log interface`
		logger := log.New(log.DebugLevel, os.Stdout)

		// grab the host
		host, _ := os.Hostname()

		// define the discovery options
		// Use PeersPort for discovery registration to avoid port conflict with GossipPort
		discoConfig := &discovery.Config{
			Endpoints:       config.EtcdEndpoints,
			ActorSystemName: config.ActorSystemName,
			Host:            host,
			DiscoveryPort:   config.GossipPort,
			TTL:             config.EtcdTTL,
			DialTimeout:     time.Duration(config.EtcdDialTimeout) * time.Second,
			Timeout:         time.Duration(config.EtcdTimeout) * time.Second,
		}
		// instantiate the etcd discovery provider
		disco := discovery.NewDiscovery(discoConfig)

		clusterConfig := goakt.
			NewClusterConfig().
			WithDiscovery(disco).
			WithPartitionCount(19).
			WithDiscoveryPort(config.GossipPort).
			WithPeersPort(config.PeersPort).
			WithClusterBalancerInterval(time.Second).
			WithKinds(new(actors.Account))

		// create the actor system
		actorSystem, err := goakt.NewActorSystem(
			config.ActorSystemName,
			goakt.WithLogger(logger),
			goakt.WithActorInitMaxRetries(3),
			goakt.WithRemote(remote.NewConfig(host, config.RemotingPort)),
			goakt.WithCluster(clusterConfig))

		// handle the error
		if err != nil {
			logger.Panic(err)
		}

		remoting := remote.NewRemoting()
		// create the combined service (AccountService with Ops functionality)
		accountService := service.NewAccountService(actorSystem, remoting, logger, disco, config, config.Port)

		actorSystem.Run(ctx,
			func(ctx context.Context) error {
				// start the combined service
				accountService.Start()
				return nil
			},
			func(ctx context.Context) error {
				remoting.Close()
				newCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
				defer cancel()
				// stop the combined service
				if err := accountService.Stop(newCtx); err != nil {
					logger.Errorf("failed to stop account service %+v", err)
				}
				return nil
			})
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Define command line flags for run command
	runCmd.Flags().IntP("port", "p", 0, "Service port (overrides PORT environment variable)")
	runCmd.Flags().String("service-name", "", "Service name (overrides SERVICE_NAME environment variable)")
	runCmd.Flags().String("system-name", "", "Actor system name (overrides SYSTEM_NAME environment variable)")
	runCmd.Flags().Int("gossip-port", 0, "Gossip port (overrides GOSSIP_PORT environment variable)")
	runCmd.Flags().Int("peers-port", 0, "Peers port (overrides PEERS_PORT environment variable)")
	runCmd.Flags().Int("remoting-port", 0, "Remoting port (overrides REMOTING_PORT environment variable)")
	runCmd.Flags().StringSlice("etcd-endpoints", []string{}, "ETCD endpoints (comma separated, overrides ETCD_ENDPOINTS environment variable)")
	runCmd.Flags().Int("etcd-dial-timeout", 0, "ETCD dial timeout in seconds (overrides ETCD_DIAL_TIMEOUT environment variable)")
	runCmd.Flags().Int64("etcd-ttl", 0, "ETCD TTL in seconds (overrides ETCD_TTL environment variable)")
	runCmd.Flags().Int("etcd-timeout", 0, "ETCD timeout in seconds (overrides ETCD_TIMEOUT environment variable)")
}
