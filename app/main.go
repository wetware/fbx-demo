package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"

	capnp "capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/rpc"

	ww "github.com/wetware/pkg"
	csp_server "github.com/wetware/pkg/cap/csp/server"

	tiktok "github.com/wetware/fbx-demo/app/cap/tiktok"
)

const (
	namespace  = "ww"
	tiktokHost = "TIKTOK_HOST"
	tiktokPort = "TIKTOK_PORT"
)

// go:embed guest/bot.wasm
var bytecode []byte

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ttHost, ttPort := loadEnvVars()

	// Connect to TikTok provider.
	tcpConn, err := dialTcp(ttHost, ttPort)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to dial TikTok provider: %s", err))
		os.Exit(1)
	}

	conn := rpc.NewConn(rpc.NewStreamTransport(tcpConn), nil)
	defer conn.Close()

	tt := tiktok.TikTok(conn.Bootstrap(ctx))
	// Validate TikTok provider.
	if err := tt.Resolve(ctx); err != nil {
		slog.Error(fmt.Sprintf("Failed to resolve TikTok provider: %s", err))
		os.Exit(1)
	}

	// Connect to Wetware cluster.
	wwClient, err := ww.InitClient(ctx, namespace, os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		panic(err)
	}

	// Prepare the TikTok capability for the guest process.
	bootstrap := csp_server.NewProcessBootstrap()
	bootstrap.AddDirect(capnp.Client(tt))

	// Spawn the guest process.
	proc, release := wwClient.Root.Exec().Exec(ctx, bootstrap.ToCap(), bytecode, 0)
	defer release()
	err = proc.Wait(ctx)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while waiting for process to exit: %s", err))
		os.Exit(1)
	}
}

func dialTcp(host, port string) (net.Conn, error) {
	return net.Dial("tcp", net.JoinHostPort(host, port))
}

func loadEnvVars() (host string, port string) {
	// Load env vars.
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
