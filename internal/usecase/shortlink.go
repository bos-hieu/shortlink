package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/bos-hieu/shortlink/internal/entities"
	"github.com/bos-hieu/shortlink/internal/models"
	"github.com/bos-hieu/shortlink/internal/repository"
	"github.com/bos-hieu/shortlink/pkg/redis"
	"github.com/bos-hieu/shortlink/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func CreateShortLink(ctx context.Context, req *entities.CreateShortLinkRequest) (*entities.CreateShortLinkResponse, error) {
	currentShortLink, err := repository.GetShortLinkByDefaultUrl(ctx, req.DefaultURL)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}
	} else {
		return &entities.CreateShortLinkResponse{
			ShortLink: fmt.Sprintf("localhost:8080/%s", currentShortLink.ShortLink),
		}, nil
	}

	shortLink := &models.ShortLink{}
	shortLink.SetFields(
		shortLink.WithID(),
		shortLink.WithShortLink(utils.GenUniqueValue()),
		shortLink.WithDestinationURLsByCountry(req.CountriesURLs),
		shortLink.WithDestinationURLsByLanguage(req.LanguagesURLs),
		shortLink.WithDefaultDestinationURL(req.DefaultURL),
	)

	err = repository.CreateShortLink(ctx, shortLink)
	if err != nil {
		return nil, err
	}

	result := &entities.CreateShortLinkResponse{
		ShortLink: fmt.Sprintf("localhost:8080/%s", shortLink.ShortLink),
	}
	return result, nil
}

func GetShortLink(ctx context.Context, req *entities.GetShortLinkRequest) (*entities.GetShortLinkResponse, error) {
	if destinationURL, err := redis.GetClient().Get(ctx, req.GetRedisKey()).Result(); err == nil {
		if destinationURL == "" {
			return nil, errors.New("the short link is not set")
		}

		return &entities.GetShortLinkResponse{
			DestinationURL: destinationURL,
		}, nil
	}

	shortLink, err := repository.GetShortLink(ctx, req.ShortLink)
	if err != nil {
		return nil, err
	}

	result := &entities.GetShortLinkResponse{
		DestinationURL: shortLink.GetDestinationURL(req.CountryCode, req.LanguageCode),
	}

	err = redis.GetClient().Set(ctx, req.GetRedisKey(), result.DestinationURL, 0).Err()
	if err != nil {
		log.Printf("Failed to set short link to redis: %s", err.Error())
	}
	return result, nil
}
