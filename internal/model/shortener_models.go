package model

import (
	"net/url"
	"time"
)

type Shortener struct {
	OriginalUrl string
	ShortUrl    string
}

func (ep Shortener) ValidUrl() bool {
	_, err := url.Parse(ep.OriginalUrl)
	return err == nil
}

type ShortenerDetails struct {
	ID string
	Shortener
	CreatedAt  time.Time
	Visits     int
	Expiration time.Time
}
