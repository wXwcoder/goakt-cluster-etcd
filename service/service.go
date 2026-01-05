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

package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"
	"github.com/pkg/errors"
	goakt "github.com/tochemey/goakt/v3/actor"
	"github.com/tochemey/goakt/v3/address"
	gerrors "github.com/tochemey/goakt/v3/errors"
	"github.com/tochemey/goakt/v3/log"
	"github.com/tochemey/goakt/v3/remote"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/protobuf/proto"

	"goakt-actors-cluster/actors"
	"goakt-actors-cluster/config"
	clusterdiscovery "goakt-actors-cluster/discovery"
	"goakt-actors-cluster/internal/opspb"
	"goakt-actors-cluster/internal/opspb/opspbconnect"
	"goakt-actors-cluster/internal/samplepb"
	"goakt-actors-cluster/internal/samplepb/samplepbconnect"
)

const askTimeout = 5 * time.Second

type AccountService struct {
	actorSystem goakt.ActorSystem
	logger      log.Logger
	port        int
	server      *http.Server
	remoting    remote.Remoting
	discovery   *clusterdiscovery.Discovery
	config      *config.Config
	opsActor    *goakt.PID
}

var _ samplepbconnect.AccountServiceHandler = &AccountService{}
var _ opspbconnect.OpsServiceHandler = &AccountService{}

// NewAccountService creates an instance of AccountService
func NewAccountService(
	system goakt.ActorSystem,
	remoting remote.Remoting,
	logger log.Logger,
	discovery *clusterdiscovery.Discovery,
	config *config.Config,
	port int,
) *AccountService {
	return &AccountService{
		actorSystem: system,
		logger:      logger,
		port:        port,
		remoting:    remoting,
		discovery:   discovery,
		config:      config,
	}
}

// CreateAccount helps create an account
func (s *AccountService) CreateAccount(ctx context.Context, c *connect.Request[samplepb.CreateAccountRequest]) (*connect.Response[samplepb.CreateAccountResponse], error) {
	// grab the actual request
	req := c.Msg
	// grab the account id
	accountID := req.GetCreateAccount().GetAccountId()
	// create the pid and send the command create account
	accountEntity := &actors.Account{}
	// create the given pid
	pid, err := s.actorSystem.Spawn(ctx, accountID, accountEntity, goakt.WithLongLived())
	if err != nil {
		return nil, err
	}
	// send the create command to the pid
	reply, err := goakt.Ask(ctx, pid, &samplepb.CreateAccount{
		AccountId:      accountID,
		AccountBalance: req.GetCreateAccount().GetAccountBalance(),
	}, time.Second)

	// handle the error
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// pattern match on the reply
	switch x := reply.(type) {
	case *samplepb.Account:
		// return the appropriate response
		return connect.NewResponse(&samplepb.CreateAccountResponse{Account: x}), nil
	default:
		// create the error message to send
		err := fmt.Errorf("invalid reply=%s", reply.ProtoReflect().Descriptor().FullName())
		return nil, connect.NewError(connect.CodeInternal, err)
	}
}

// CreditAccount helps credit a given account
func (s *AccountService) CreditAccount(ctx context.Context, c *connect.Request[samplepb.CreditAccountRequest]) (*connect.Response[samplepb.CreditAccountResponse], error) {
	req := c.Msg
	accountID := req.GetCreditAccount().GetAccountId()

	addr, pid, err := s.actorSystem.ActorOf(ctx, accountID)
	if err != nil {
		// check whether it is not found error
		if !errors.Is(err, gerrors.ErrActorNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}

		// return not found
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	var message proto.Message
	command := &samplepb.CreditAccount{
		AccountId: accountID,
		Balance:   req.GetCreditAccount().GetBalance(),
	}

	if pid != nil {
		s.logger.Info("actor is found locally...")
		message, err = goakt.Ask(ctx, pid, command, time.Second)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	if pid == nil {
		s.logger.Info("actor is not found locally...")
		reply, err := s.remoting.RemoteAsk(ctx, address.NoSender(), addr, command, askTimeout)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		message, _ = reply.UnmarshalNew()
	}

	switch x := message.(type) {
	case *samplepb.Account:
		return connect.NewResponse(&samplepb.CreditAccountResponse{Account: x}), nil
	default:
		err := fmt.Errorf("invalid reply=%s", message.ProtoReflect().Descriptor().FullName())
		return nil, connect.NewError(connect.CodeInternal, err)
	}
}

// GetAccount helps get an account
func (s *AccountService) GetAccount(ctx context.Context, c *connect.Request[samplepb.GetAccountRequest]) (*connect.Response[samplepb.GetAccountResponse], error) {
	// grab the actual request
	req := c.Msg
	// grab the account id
	accountID := req.GetAccountId()

	// locate the given actor
	addr, pid, err := s.actorSystem.ActorOf(ctx, accountID)
	// handle the error
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var message proto.Message
	command := &samplepb.GetAccount{
		AccountId: accountID,
	}

	if pid != nil {
		s.logger.Info("actor is found locally...")
		message, err = goakt.Ask(ctx, pid, command, time.Second)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	if pid == nil {
		s.logger.Info("actor is not found locally...")
		reply, err := s.remoting.RemoteAsk(ctx, address.NoSender(), addr, command, askTimeout)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		message, _ = reply.UnmarshalNew()
	}

	// pattern match on the reply
	switch x := message.(type) {
	case *samplepb.Account:
		return connect.NewResponse(&samplepb.GetAccountResponse{Account: x}), nil
	default:
		err := fmt.Errorf("invalid reply=%s", message.ProtoReflect().Descriptor().FullName())
		return nil, connect.NewError(connect.CodeInternal, err)
	}
}

// OpsService implementation

// GetNodeDetails returns detailed information about a specific node

// GetClusterNodes returns the list of all nodes in the cluster
func (s *AccountService) GetClusterNodes(ctx context.Context, req *connect.Request[opspb.GetClusterNodesRequest]) (*connect.Response[opspb.GetClusterNodesResponse], error) {
	// Forward the request to the OpsActor
	actorMessage := &opspb.OpsActorMessage{
		Message: &opspb.OpsActorMessage_GetClusterNodes{
			GetClusterNodes: req.Msg,
		},
	}

	response, err := goakt.Ask(ctx, s.opsActor, actorMessage, time.Second)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	actorResponse, ok := response.(*opspb.OpsActorResponse)
	if !ok {
		return nil, connect.NewError(connect.CodeInternal,
			fmt.Errorf("invalid response type from OpsActor"))
	}

	clusterNodesResponse := actorResponse.GetClusterNodes()
	if clusterNodesResponse == nil {
		return nil, connect.NewError(connect.CodeInternal,
			fmt.Errorf("OpsActor returned nil response for GetClusterNodes"))
	}

	return connect.NewResponse(clusterNodesResponse), nil
}

// GetNodeDetails returns detailed information about a specific node
func (s *AccountService) GetNodeDetails(ctx context.Context, req *connect.Request[opspb.GetNodeDetailsRequest]) (*connect.Response[opspb.GetNodeDetailsResponse], error) {
	// Forward the request to the OpsActor
	actorMessage := &opspb.OpsActorMessage{
		Message: &opspb.OpsActorMessage_GetNodeDetails{
			GetNodeDetails: req.Msg,
		},
	}

	response, err := goakt.Ask(ctx, s.opsActor, actorMessage, time.Second)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	actorResponse, ok := response.(*opspb.OpsActorResponse)
	if !ok {
		return nil, connect.NewError(connect.CodeInternal,
			fmt.Errorf("invalid response type from OpsActor"))
	}

	nodeDetailsResponse := actorResponse.GetNodeDetails()
	if nodeDetailsResponse == nil {
		return nil, connect.NewError(connect.CodeNotFound,
			fmt.Errorf("node %s not found", req.Msg.GetNodeId()))
	}

	return connect.NewResponse(nodeDetailsResponse), nil
}

// HealthCheck performs a health check on the current node
func (s *AccountService) HealthCheck(ctx context.Context, req *connect.Request[opspb.HealthCheckRequest]) (*connect.Response[opspb.HealthCheckResponse], error) {
	// Forward the request to the OpsActor
	actorMessage := &opspb.OpsActorMessage{
		Message: &opspb.OpsActorMessage_HealthCheck{
			HealthCheck: req.Msg,
		},
	}

	response, err := goakt.Ask(ctx, s.opsActor, actorMessage, time.Second)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	actorResponse, ok := response.(*opspb.OpsActorResponse)
	if !ok {
		return nil, connect.NewError(connect.CodeInternal,
			fmt.Errorf("invalid response type from OpsActor"))
	}

	healthCheckResponse := actorResponse.GetHealthCheck()
	if healthCheckResponse == nil {
		return nil, connect.NewError(connect.CodeInternal,
			fmt.Errorf("OpsActor returned nil response for HealthCheck"))
	}

	return connect.NewResponse(healthCheckResponse), nil
}

// GetClusterStats returns cluster-wide statistics
func (s *AccountService) GetClusterStats(ctx context.Context, req *connect.Request[opspb.ClusterStatsRequest]) (*connect.Response[opspb.ClusterStatsResponse], error) {
	// Forward the request to the OpsActor
	actorMessage := &opspb.OpsActorMessage{
		Message: &opspb.OpsActorMessage_GetClusterStats{
			GetClusterStats: req.Msg,
		},
	}

	response, err := goakt.Ask(ctx, s.opsActor, actorMessage, time.Second)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	actorResponse, ok := response.(*opspb.OpsActorResponse)
	if !ok {
		return nil, connect.NewError(connect.CodeInternal,
			fmt.Errorf("invalid response type from OpsActor"))
	}

	clusterStatsResponse := actorResponse.GetClusterStats()
	if clusterStatsResponse == nil {
		return nil, connect.NewError(connect.CodeInternal,
			fmt.Errorf("OpsActor returned nil response for GetClusterStats"))
	}

	return connect.NewResponse(clusterStatsResponse), nil
}

// GetActorDetails returns detailed information about actors in the cluster
func (s *AccountService) GetActorDetails(ctx context.Context, req *connect.Request[opspb.GetActorDetailsRequest]) (*connect.Response[opspb.GetActorDetailsResponse], error) {
	// Forward the request to the OpsActor to get actors from the entire cluster
	actorMessage := &opspb.OpsActorMessage{
		Message: &opspb.OpsActorMessage_GetActorDetails{
			GetActorDetails: req.Msg,
		},
	}

	response, err := goakt.Ask(ctx, s.opsActor, actorMessage, time.Second)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	actorResponse, ok := response.(*opspb.OpsActorResponse)
	if !ok {
		return nil, connect.NewError(connect.CodeInternal,
			fmt.Errorf("invalid response type from OpsActor"))
	}

	actorDetailsResponse := actorResponse.GetActorDetails()
	if actorDetailsResponse == nil {
		return nil, connect.NewError(connect.CodeInternal,
			fmt.Errorf("OpsActor returned nil response for GetActorDetails"))
	}

	return connect.NewResponse(actorDetailsResponse), nil
}

// Start starts the service
func (s *AccountService) Start() {
	go func() {
		// Wait for actor system to be fully running with retry mechanism
		maxRetries := 10
		retryInterval := 1 * time.Second

		for i := 0; i < maxRetries; i++ {
			time.Sleep(retryInterval)

			// Check if actor system is running by attempting to spawn a test actor
			testActor := &testActor{}
			_, err := s.actorSystem.Spawn(
				context.Background(),
				"test-actor-"+fmt.Sprintf("%d", i),
				testActor,
				goakt.WithLongLived(), // Keep actor alive for testing
			)

			if err == nil {
				// Actor system is running, proceed with OpsActor
				break
			}

			if i == maxRetries-1 {
				s.logger.Panicf("actor system failed to start after %d retries: %v", maxRetries, err)
			}

			s.logger.Warnf("actor system not ready yet, retrying in %v (attempt %d/%d)", retryInterval, i+1, maxRetries)
		}

		// Start the OpsActor for this node
		opsActor := actors.NewOpsActor(s.remoting, s.logger, s.discovery, s.config)

		// Spawn the OpsActor with a unique name based on the service name
		opsActorName := "ops-actor"
		//opsActorName := fmt.Sprintf("ops-actor-%s", s.config.ServiceName)
		s.logger.Infof("Attempting to create OpsActor with name: %s (ServiceName: %s)", opsActorName, s.config.ServiceName)

		// Check if OpsActor already exists before spawning
		ctx := context.Background()
		exists, err := s.actorSystem.ActorExists(ctx, opsActorName)
		if err == nil && exists {
			// OpsActor already exists, use the existing one
			s.logger.Infof("OpsActor already exists, using existing actor: %s", opsActorName)
		} else {
			// Spawn a new OpsActor
			pid, err := s.actorSystem.Spawn(
				ctx,
				opsActorName,
				opsActor,
				goakt.WithLongLived(), // Keep actor alive indefinitely
			)
			if err != nil {
				s.logger.Panicf("failed to spawn OpsActor: %v", err)
			}

			s.opsActor = pid
			s.logger.Infof("OpsActor started successfully with name: %s", opsActorName)
		}
	}()

	go func() {
		s.listenAndServe()
	}()
}

// Stop stops the service
func (s *AccountService) Stop(ctx context.Context) error {
	// Stop the OpsActor
	if s.opsActor != nil {
		s.actorSystem.Stop(context.Background())
	}

	return s.server.Shutdown(ctx)
}

// listenAndServe starts the http server
func (s *AccountService) listenAndServe() {
	// create a http service mux
	mux := http.NewServeMux()
	// create an interceptor
	interceptor, err := otelconnect.NewInterceptor()
	if err != nil {
		s.logger.Panic(err)
	}

	// Register AccountService handler
	accountPath, accountHandler := samplepbconnect.NewAccountServiceHandler(s,
		connect.WithInterceptors(interceptor))
	mux.Handle(accountPath, accountHandler)

	// Register OpsService handler
	opsPath, opsHandler := opspbconnect.NewOpsServiceHandler(s)
	mux.Handle(opsPath, opsHandler)

	// create the address
	serverAddr := fmt.Sprintf(":%d", s.port)
	// create a http server instance
	server := &http.Server{
		Addr:              serverAddr,
		ReadTimeout:       3 * time.Second,
		ReadHeaderTimeout: time.Second,
		WriteTimeout:      time.Second,
		IdleTimeout:       1200 * time.Second,
		Handler: h2c.NewHandler(mux, &http2.Server{
			IdleTimeout: 1200 * time.Second,
		}),
	}

	// set the server
	s.server = server
	// listen and service requests
	if err := s.server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return
		}
		s.logger.Panic(errors.Wrap(err, "failed to start actor-remoting service"))
	}
}

// testActor is a simple actor used to check if the actor system is running
type testActor struct{}

func (t *testActor) PreStart(*goakt.Context) error {
	return nil
}

func (t *testActor) Receive(ctx *goakt.ReceiveContext) {
	// Do nothing - this actor is only used for testing system readiness
}

func (t *testActor) PostStop(*goakt.Context) error {
	return nil
}
