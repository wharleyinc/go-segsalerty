package model

import (
	wrapping "go-segsalerty/common/error"
)

const (
	ErrorInvalidUrl       = wrapping.Error("error: invalid url address")
	ErrorAlreadyExist     = wrapping.Error("email already exist")
	ErrorCreatingShortUrl = wrapping.Error("error: creating short url")
)
