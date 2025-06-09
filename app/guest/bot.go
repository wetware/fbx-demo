package main

import (
	"context"
	"fmt"
	"os"

	ww "github.com/wetware/pkg/guest/system"

	"github.com/wetware/fbx-demo/app/cap/tiktok"
)

func main() {
	ctx := context.Background()

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
		fmt.Println("No capabilities found in bootstrap")
		os.Exit(1)
	}

	_ = tiktok.TikTok(caps[0])
	// tt := tiktok.TikTok(caps[0])
	// tt.Mention(ctx, params func(tiktok.TikTok_mention_Params) error)
}
