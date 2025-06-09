package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/wetware/fbx-demo/app/cap/tiktok"
)

type bot struct {
	tt tiktok.TikTok
}

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
