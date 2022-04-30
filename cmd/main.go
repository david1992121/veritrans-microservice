package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/david1992121/veritrans-microservice/pkg/transport"
	"github.com/go-kit/kit/log"
	"github.com/joho/godotenv"
	"github.com/oklog/oklog/pkg/group"
)

const (
	defaultHTTPPort = "8080"
	defaultGRPCPort = "8081"
)

func main() {
	// set logger
	var (
		logger   log.Logger
		httpAddr = net.JoinHostPort("0.0.0.0", envString("HTTP_PORT", defaultHTTPPort))
		// grpcAddr = net.JoinHostPort("localhost", envString("GRPC_PORT", defaultGRPCPort))
	)

	logger = initLogger()

	// get environment variables if exists
	if err := godotenv.Load(); err != nil {
		logger.Log("read", "env", "err", err)
	}

	var (
		httpHandler = transport.NewHTTPHandler(logger)
		// grpcServer = transport.NewGRPCServer(eps)
	)

	var g group.Group
	{
		httpListener, err := net.Listen("tcp", httpAddr)
		if err != nil {
			logger.Log("transport", "HTTP", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log("transport", "HTTP", "addr", httpAddr)
			return http.Serve(httpListener, httpHandler)
		}, func(error) {
			httpListener.Close()
		})
	}

	// TODO : gRPC Listener and server
	// {
	// 	grpcListener, err := net.Listen("tcp", grpcAddr)
	// 	if err != nil {
	// 		logger.Log("transport", "gRPC", "during", "Listen", "err", err)
	// 		os.Exit(1)
	// 	}
	// 	g.Add(func() error {
	// 		return baseServer.serve(grpcListener)
	// 	}, func(error) {
	// 		grpcListener.Close()
	// 	})
	// }

	{
		// This function just sits and waits for ctrl-C.
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}
	logger.Log("exit", g.Run())
}

func envString(env, defaultValue string) string {
	e := os.Getenv(env)
	if e == "" {
		return defaultValue
	}
	return e
}

func initLogger() log.Logger {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	return logger
}
