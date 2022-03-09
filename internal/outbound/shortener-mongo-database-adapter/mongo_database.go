package shortenermongodatabaseadapter

import (
	"context"
	"errors"
	"go-segsalerty/common"
	appmongo "go-segsalerty/common/database/mongo"
	"go-segsalerty/common/logger"
	"go-segsalerty/internal/config"
	"go-segsalerty/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"time"
)

const (
	timeout         = 30
	databaseName    = "url-shortener"
	shortCollection = "shorted"
)

type userLink struct {
	ID        primitive.ObjectID `bson:"_id"`
	LongUrl   string             `bson:"long_url"`
	ShortUrl  string             `bson:"short_url"`
	ExpireAt  primitive.DateTime `bson:"expire_at"`
	CreatedAt primitive.DateTime `bson:"created_at"`
}

type adapter struct {
	shortColl *mongo.Collection
}

func NewDatabaseAdapter(config config.Config) (*adapter, error) {
	db, err := appmongo.NewDriver(appmongo.Config{
		URI:     config.MongoURI,
		Timeout: timeout,
	})
	if err != nil {
		return nil, err
	}
	mongodb := db.Database(databaseName)
	return &adapter{
		shortColl: mongodb.Collection(shortCollection),
	}, nil
}

func (a adapter) ShortenSave(ctx context.Context, longUrl, shortUrl string) (*model.ShortenerDetails, error) {

	objId, err := primitive.ObjectIDFromHex(longUrl)
	if err != nil {
		logger.Error(ctx, "error converting longUrl address", zap.Error(errors.New(longUrl)))
		logger.Error(ctx, "error converting longUrl address", zap.Error(err))
		return nil, model.ErrorCreatingShortUrl.Wrap(err)
	}

	var shortLink = userLink{
		ID:        objId,
		LongUrl:   longUrl,
		ShortUrl:  shortUrl,
		ExpireAt:  primitive.NewDateTimeFromTime(common.TimeNow().AddDate(0, 0, 5)),
		CreatedAt: primitive.NewDateTimeFromTime(common.TimeNow()),
	}

	_, err = a.shortColl.InsertOne(ctx, shortLink)
	if err != nil {
		logger.Error(ctx, "error inserting record to db", zap.Error(err))
		if isDup(err) {
			return nil, model.ErrorAlreadyExist.Wrap(err)
		}
		return nil, model.ErrorCreatingShortUrl.Wrap(err)
	}

	return &model.ShortenerDetails{
		ID: shortLink.ID.Hex(),
		Shortener: model.Shortener{
			OriginalUrl: shortLink.LongUrl,
			ShortUrl:    shortLink.ShortUrl,
		},
		CreatedAt:  shortLink.CreatedAt.Time().UTC(),
		Visits:     0,
		Expiration: time.Time{},
	}, err
}

func isDup(err error) bool {
	var e mongo.ServerError
	if errors.As(err, &e) {
		return e.HasErrorCode(11000)
	}
	return false
}
