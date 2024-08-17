package repository

import (
	"context"
	"github.com/bos-hieu/shortlink/internal/models"
	"github.com/bos-hieu/shortlink/pkg/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

// GetShortLink retrieves a short link from the database
func GetShortLink(ctx context.Context, shortLink string) (*models.ShortLink, error) {
	var item models.ShortLink

	err := mongodb.GetClient().
		Database("short_link").
		Collection("short_links").
		FindOne(ctx, bson.M{"short_link": shortLink}).Decode(&item)

	if err != nil {
		return nil, err
	}

	return &item, nil
}

func GetShortLinkByDefaultUrl(ctx context.Context, defaultURL string) (*models.ShortLink, error) {
	var item models.ShortLink

	err := mongodb.GetClient().
		Database("short_link").
		Collection("short_links").
		FindOne(ctx, bson.M{"default_destination_url": defaultURL}).Decode(&item)

	if err != nil {
		return nil, err
	}

	return &item, nil
}

// CreateShortLink creates a new short link in the database
func CreateShortLink(ctx context.Context, shortLink *models.ShortLink) error {
	_, err := mongodb.GetClient().
		Database("short_link").
		Collection("short_links").
		InsertOne(ctx, shortLink)
	return err
}
