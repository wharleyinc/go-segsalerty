//go:generate mockgen -source=authorizer.go -destination authorizer_mock.go -package auth . Repository, Service, EmailAdapter

package shorthy

import (
	"context"
	"errors"
	"go-segsalerty/common/logger"
	"go-segsalerty/internal/model"
	"go.uber.org/zap"
	"math"
	"net/url"
	"strings"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length   = uint64(len(alphabet))
)

type Service interface {
	ShortenUrl(ctx context.Context, newUrl model.Shortener) (*model.ShortenerDetails, error)
	//AuthenticateAccount(ctx context.Context, existingAccount model.AuthenticateAccount) (*model.AuthToken, error)
}

type Repository interface {
	ShortenSave(ctx context.Context, longUrl, shortUrl string) (*model.ShortenerDetails, error)
}

type short struct {
	repository Repository
}

func NewShortenerService(repository Repository) Service {
	return &short{repository: repository}
}

func (s short) ShortenUrl(ctx context.Context, newUrl model.Shortener) (*model.ShortenerDetails, error) {
	ctx = logger.With(ctx,
		zap.String("newUrl", newUrl.OriginalUrl))

	if !newUrl.ValidUrl() {
		return nil, model.ErrorInvalidUrl
	}
	var newURLLLL *url.URL
	var err error

	if len(newUrl.OriginalUrl) >= 5 {
		newURLLLL, err = url.ParseRequestURI(newUrl.OriginalUrl)
		if err != nil {
			return nil, err
		}
		logger.Info(ctx, "this is gotten here:", zap.String("originalUrl", newURLLLL.String()))
	}

	newShort := model.Shortener{
		OriginalUrl: newURLLLL.Host,
		ShortUrl:    newURLLLL.Host,
	}

	result, err := s.repository.ShortenSave(ctx, newShort.OriginalUrl, newShort.ShortUrl)
	if err != nil {
		logger.Error(ctx, "unable to create account", zap.Error(err))
		if errors.Is(err, model.ErrorAlreadyExist) {
			return nil, model.ErrorAlreadyExist
		}
		return nil, model.ErrorCreatingShortUrl
	}

	return result, nil
}

func Encode(number uint64) string {
	var encodedBuilder strings.Builder
	encodedBuilder.Grow(11)

	for ; number > 0; number = number / length {
		encodedBuilder.WriteByte(alphabet[(number % length)])
	}

	return encodedBuilder.String()
}

func Decode(encoded string) (uint64, error) {
	var number uint64

	for i, symbol := range encoded {
		alphabeticPosition := strings.IndexRune(alphabet, symbol)

		if alphabeticPosition == -1 {
			return uint64(alphabeticPosition), errors.New("invalid character: " + string(symbol))
		}
		number += uint64(alphabeticPosition) * uint64(math.Pow(float64(length), float64(i)))
	}

	return number, nil
}
