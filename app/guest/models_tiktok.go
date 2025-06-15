package main

import (
	"github.com/wetware/fbx-demo/app/cap/tiktok"
)

type user struct {
	id       string
	uid      string
	nickname string
}

type comment struct {
	id      string
	mediaId string
	text    string
	author  user
	replies []comment
}

func unmarshalUser(u tiktok.User) (user, error) {
	id, err := u.Id()
	if err != nil {
		return user{}, err
	}

	uid, err := u.Uid()
	if err != nil {
		return user{}, err
	}

	nickname, err := u.Nickname()
	if err != nil {
		return user{}, err
	}

	return user{
		id:       id,
		uid:      uid,
		nickname: nickname,
	}, nil
}

func unmarshalComment(c tiktok.Comment) (comment, error) {
	id, err := c.Id()
	if err != nil {
		return comment{}, err
	}

	mediaId, err := c.MediaId()
	if err != nil {
		return comment{}, err
	}

	text, err := c.Text()
	if err != nil {
		return comment{}, err
	}

	a, err := c.Author()
	if err != nil {
		return comment{}, err
	}

	author, err := unmarshalUser(a)
	if err != nil {
		return comment{}, err
	}

	r, err := c.Replies()
	if err != nil {
		return comment{}, err
	}

	replies := make([]comment, r.Len())
	for i := 0; i < r.Len(); i++ {
		replies[i], err = unmarshalComment(r.At(i))
		if err != nil {
			return comment{}, err
		}
	}

	return comment{
		id:      id,
		mediaId: mediaId,
		text:    text,
		author:  author,
		replies: replies,
	}, nil
}
