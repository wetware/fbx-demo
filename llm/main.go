package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"os"

	capnp "capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/rpc"
	llm_cap "github.com/wetware/fbx-demo/llm/cap/llm"
)

const (
	llmHost = "LLM_HOST"
	llmPort = "LLM_PORT"
)

func main() {
	ctx := context.Background()
	if err := runServer(ctx); err != nil {
		slog.Error("Failed to run server", "error", err)
		os.Exit(1)
	}
}

func runServer(ctx context.Context) error {
	host := os.Getenv(llmHost)
	if host == "" {
		return errors.New("Missing environment variable LLM_HOST")
	}
	port := os.Getenv(llmPort)
	if port == "" {
		return errors.New("Missing environment variable LLM_PORT")
	}

	// Setup and serve the capability.
	slog.Info("Building LLM capability server...")
	llm := NewLLM()
	cli := llm_cap.LLM_ServerToClient(llm)
	return rpc.ListenAndServe(ctx, "tcp", net.JoinHostPort(host, port), capnp.Client(cli))

}
