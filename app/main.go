package main

import (
	"context"
	_ "embed"
	"io"
	"log/slog"
	"net"
	"os"

	capnp "capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/rpc"

	ww "github.com/wetware/pkg"
	csp_server "github.com/wetware/pkg/cap/csp/server"

	llm "github.com/wetware/fbx-demo/app/cap/llm"
	tiktok "github.com/wetware/fbx-demo/app/cap/tiktok"
)

const (
	namespace = "ww"

	llmHost = "LLM_HOST"
	llmPort = "LLM_PORT"

	tiktokHost = "TIKTOK_HOST"
	tiktokPort = "TIKTOK_PORT"
)

//go:embed guest/bot.wasm
var bytecode []byte

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Connect to LLM provider.
	llmHost, llmPort := loadLLMEnvVars()
	llmConn, llmCloser, err := dialCap(ctx, "tiktok", llmHost, llmPort)
	defer llmCloser.Close()
	if err != nil {
		slog.Error("Failed to dial TikTok service", "error", err)
		os.Exit(1)
	}
	llmCap := llm.LLM(llmConn)
	if err := llmCap.Resolve(ctx); err != nil {
		slog.Error("Failed to resolve TikTok provider", "error", err)
		os.Exit(1)
	}

	// Connect to TikTok provider.
	ttHost, ttPort := loadTikTokEnvVars()
	ttConn, ttCloser, err := dialCap(ctx, "tiktok", ttHost, ttPort)
	defer ttCloser.Close()
	if err != nil {
		slog.Error("Failed to dial TikTok service", "error", err)
		os.Exit(1)
	}
	ttCap := tiktok.TikTok(ttConn)
	if err := ttCap.Resolve(ctx); err != nil {
		slog.Error("Failed to resolve TikTok provider", "error", err)
		os.Exit(1)
	}

	slog.Info("Resolved TikTok capability.")

	// Connect to Wetware cluster.
	wwClient, closer, err := ww.InitClient(ctx, namespace, os.Stdin, os.Stdout, os.Stderr)
	if closer != nil {
		defer closer.Close()
	}
	if err != nil {
		slog.Error("Failed to initialize Wetware client", "error", err)
		os.Exit(1)
	}

	// Prepare the TikTok capability for the guest process.
	bootstrap := csp_server.NewProcessBootstrap()
	bootstrap.AddDirect(capnp.Client(llmCap))
	bootstrap.AddDirect(capnp.Client(ttCap))

	// Spawn the guest process.
	proc, release := wwClient.Root.Exec().Exec(ctx, bootstrap.ToCap(), bytecode, 0)
	defer release()
	err = proc.Wait(ctx)
	if err != nil {
		slog.Error("Error while waiting for process to exit", "error", err)
		os.Exit(1)
	}
}

func dialCap(ctx context.Context, capName, host, port string) (capnp.Client, io.Closer, error) {
	// Connect to a capability provider.
	tcpConn, err := dialTcp(host, port)
	if err != nil {
		slog.Error("Failed to dial provider", "capability", capName, "error", err)
		os.Exit(1)
	}

	slog.Info("Dialed provider", "capability", capName)

	conn := rpc.NewConn(rpc.NewStreamTransport(tcpConn), nil)

	slog.Info("Created RPC transport over dialed connection")

	return conn.Bootstrap(ctx), conn, nil
}

func dialTcp(host, port string) (net.Conn, error) {
	return net.Dial("tcp", net.JoinHostPort(host, port))
}

func loadTikTokEnvVars() (host string, port string) {
	host = os.Getenv(tiktokHost)
	if host == "" {
		host = "127.0.0.1"
	}
	port = os.Getenv(tiktokPort)
	if port == "" {
		port = "6060"
	}
	return host, port
}

func loadLLMEnvVars() (host string, port string) {
	// Load env vars.
	host = os.Getenv(llmHost)
	if host == "" {
		host = "127.0.0.1"
	}
	port = os.Getenv(llmPort)
	if port == "" {
		port = "6061"
	}
	return host, port
}
