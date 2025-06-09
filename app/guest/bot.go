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

func (b bot) runLoop(ctx context.Context) error {
	for {
		// End loop of the context was canceled
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Get the next mention
		fmt.Println("Getting mention...")
		m, err := b.getMention(ctx)
		if err != nil {
			if errors.Is(err, errNoComment) {
				fmt.Println("No comment found, waiting 5s before retrying")
				time.Sleep(5 * time.Second)
				continue
			} else {
				return fmt.Errorf("Error getting mention: %w", err)
			}
		}
		mention, err := unmarshalComment(m)
		if err != nil {
			return fmt.Errorf("Error unmarshaling comment: %w", err)
		}
		fmt.Printf("Got mention %v\n", mention)

		// Get additional context for the mention.
		// In this case, the rest of the comments in the post.
		fmt.Println("Getting comments...")
		cs, err := b.getComments(ctx, mention.mediaId)
		if err != nil {
			if errors.Is(err, errNoComments) {
				fmt.Printf("No additional comments found in post %s\n", mention.mediaId)
				cs = []tiktok.Comment{}
			} else {
				return fmt.Errorf("Error getting comments: %w", err)
			}
		}

		comments := make([]comment, len(cs))
		for i, c := range cs {
			comment, err := unmarshalComment(c)
			if err != nil {
				return fmt.Errorf("Error unmarshaling comment: %w", err)
			}
			comments[i] = comment
		}

		fmt.Printf("Got comments %v", comments)

		// Call the LLM model
		// TODO mikel: add LLM capability. Modify the text for now.
		reply := fmt.Sprintf("%s in bed", mention.text)

		// Send the reply to the user
		fmt.Println("Replying to comment...")
		err = b.replyToComment(ctx, mention, reply)
		if err != nil {
			return fmt.Errorf("Error replying to comment: %w", err)
		}
		fmt.Printf("Replied `%s` to comment", reply)

		// Sanity sleep, don't burn through our API tokens accidentally...
		time.Sleep(30 * time.Second)
	}
}

func (b bot) getMention(ctx context.Context) (tiktok.Comment, error) {
	f, release := b.tt.Mention(ctx, nil)
	defer release()
	select {
	case <-ctx.Done():
	case <-f.Done():
	}
	if ctx.Err() != nil {
		return tiktok.Comment{}, ctx.Err()
	}

	res, err := f.Struct()
	if err != nil {
		return tiktok.Comment{}, err
	}

	if !res.HasComment() {
		return tiktok.Comment{}, errNoComment
	}
	return res.Comment()
}

func (b bot) getComments(ctx context.Context, mediaId string) ([]tiktok.Comment, error) {
	f, release := b.tt.Comments(ctx, func(tt tiktok.TikTok_comments_Params) error {
		return tt.SetMediaId(mediaId)
	})
	defer release()

	select {
	case <-ctx.Done():
	case <-f.Done():
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	res, err := f.Struct()
	if err != nil {
		return nil, err
	}

	if !res.HasComments() {
		return nil, errNoComments
	}

	cl, err := res.Comments()
	if err != nil {
		return nil, err
	}

	comments := make([]tiktok.Comment, cl.Len())
	for i := 0; i < cl.Len(); i++ {
		comments[i] = cl.At(i)
	}

	return comments, nil
}

func (b bot) replyToComment(ctx context.Context, source comment, reply string) error {
	f, release := b.tt.Reply(ctx, func(tt tiktok.TikTok_reply_Params) error {
		err := tt.SetMediaId(source.mediaId)
		if err != nil {
			return err
		}
		err = tt.SetCommendId(source.id)
		if err != nil {
			return err
		}
		return tt.SetResponse(reply)
	})
	defer release()

	select {
	case <-ctx.Done():
	case <-f.Done():
	}

	return ctx.Err()
}
