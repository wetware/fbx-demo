package main

import (
	"context"
	"errors"

	capnp "capnproto.org/go/capnp/v3"
	llm_cap "github.com/wetware/fbx-demo/llm/cap/llm"
)

// LLM implements the LLM capability.
type LLM struct {
}

func NewLLM() *LLM {
	return &LLM{}
}

func (_ *LLM) GenerateResponse(ctx context.Context, call llm_cap.LLM_generateResponse) error {
	res, err := call.AllocResults()
	if err != nil {
		return err
	}

	input, err := call.Args().Input()
	if !input.HasMention() {
		return errors.New("GenerateResponse call had no input")
	}

	mention, err := input.Mention()
	if err != nil {
		return err
	}

	var mentionContext []string
	if !input.HasContext() {
		mentionContext = []string{}
	} else {
		tl, err := input.Context()
		if err != nil {
			return err
		}
		mentionContext, err = decodeTextList(tl)
		if err != nil {
			return err
		}
	}

	response, err := askLLMForResponse(ollamaURL, mention, mentionContext)
	if err != nil {
		return err
	}

	return res.SetOutput(response)
}

func decodeTextList(tl capnp.TextList) ([]string, error) {
	textList := make([]string, tl.Len())
	for i := 0; i < tl.Len(); i++ {
		text, err := tl.At(i)
		if err != nil {
			return nil, err
		}
		textList[i] = text
	}
	return textList, nil
}
