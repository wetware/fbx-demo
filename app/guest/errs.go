package main

import "errors"

var (
	errNoComment  = errors.New("TikTok capability had no new comments")
	errNoComments = errors.New("Post had no other comments")
)
