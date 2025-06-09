//go:generate env GOOS=wasip1 GOARCH=wasm go build -o bot.wasm main.go bot.go errs.go models.go
package main

import (
	"context"
	"fmt"
	"os"

	ww "github.com/wetware/pkg/guest/system"

	"github.com/wetware/fbx-demo/app/cap/tiktok"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	caps, releases, err := ww.Bootstrap(ctx)
	if err != nil {
		panic(err)
	}
	defer func() {
		for _, release := range releases {
			release()
		}
	}()

	if len(caps) == 0 {
		fmt.Println("No capabilities found in bootstrap.")
		os.Exit(1)
	}

	tt := tiktok.TikTok(caps[0])
	if err := tt.Resolve(ctx); err != nil {
		panic(err)
	}

	fmt.Println("Bootstrapped TikTok capability.")

	bot := bot{
		tt: tt,
	}
	err = bot.runLoop(ctx)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error running bot loop: %s", err)
		os.Exit(1)
	}
}
