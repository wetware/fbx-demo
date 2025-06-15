package main

import (
	capnp "capnproto.org/go/capnp/v3"
	"github.com/wetware/fbx-demo/app/cap/llm"
)

type input struct {
	mention string
	context []string
}

func marshalInput(input input) (llm.Input, error) {
	_, seg := capnp.NewSingleSegmentMessage(nil)
	in, err := llm.NewInput(seg)
	if err != nil {
		return llm.Input{}, err
	}

	if err := in.SetMention(input.mention); err != nil {
		return llm.Input{}, err
	}

	inCtx, err := encodeTextList(input.context)
	if err != nil {
		return llm.Input{}, err
	}

	err = in.SetContext(inCtx)
	return in, err
}

func encodeTextList(tl []string) (capnp.TextList, error) {
	_, seg := capnp.NewSingleSegmentMessage(nil)
	etl, err := capnp.NewTextList(seg, int32(len(tl)))
	if err != nil {
		return capnp.TextList{}, err
	}

	for i, text := range tl {
		if err := etl.Set(i, text); err != nil {
			return capnp.TextList{}, err
		}
	}

	return etl, nil
}
