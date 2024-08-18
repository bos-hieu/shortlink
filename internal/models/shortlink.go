package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type (
	ShortLink struct {
		// ID is the mongo object ID
		ID primitive.ObjectID `bson:"_id"`

		// ShortLink is the short link identifier from the default URL
		// This field will be indexed by unique to increase read performance.
		ShortLink string `bson:"short_link"`

		// DefaultDestinationURL is the default URL to redirect to if no other rules are matched
		// This field will be indexed by unique to increase read performance.
		DefaultDestinationURL string `bson:"default_destination_url"`

		// DestinationURLsByCountry is a map of country codes to destination URLs
		// The country codes are in ISO 3166-1 alpha-2 format, for example: US, GB, FR, etc.
		DestinationURLsByCountry map[string]string `bson:"destination_urls_by_country"`

		// DestinationURLsByLanguageCountry is a map of language codes to destination URLs
		// The language codes are in ISO 639-1 format, for example: en-US, en-GB, fr-FR, etc.
		DestinationURLsByLanguageCountry map[string]string `bson:"destination_urls_by_language_country"`
	}

	// ShortLinkOption is a function that modifies a ShortLink
	ShortLinkOption func(*ShortLink)
)

func (a *ShortLink) SetFields(options ...ShortLinkOption) {
	for _, option := range options {
		option(a)
	}
}

// WithID sets the ID
func (a ShortLink) WithID() ShortLinkOption {
	return func(a *ShortLink) {
		a.ID = primitive.NewObjectID()
	}
}

// WithShortLink sets the short link
func (a ShortLink) WithShortLink(shortLink string) ShortLinkOption {
	return func(a *ShortLink) {
		a.ShortLink = shortLink
	}
}

// WithDefaultDestinationURL sets the default destination URL for the short link
func (a ShortLink) WithDefaultDestinationURL(defaultDestinationURL string) ShortLinkOption {
	return func(a *ShortLink) {
		a.DefaultDestinationURL = defaultDestinationURL
	}
}

// WithDestinationURLsByCountry sets the destination URL for specific countries
func (a ShortLink) WithDestinationURLsByCountry(destinationURLsByCountry map[string]string) ShortLinkOption {
	return func(a *ShortLink) {
		a.DestinationURLsByCountry = destinationURLsByCountry
	}
}

// WithDestinationURLsByLanguage sets the destination URL for specific languages
func (a ShortLink) WithDestinationURLsByLanguage(destinationURLsByLanguageCountry map[string]string) ShortLinkOption {
	return func(a *ShortLink) {
		a.DestinationURLsByLanguageCountry = destinationURLsByLanguageCountry
	}
}

func (a ShortLink) GetDestinationURL(countryCode, languageCode string) string {
	if a.DestinationURLsByLanguageCountry != nil && languageCode != "" {
		destinationURL, hasURL := a.DestinationURLsByLanguageCountry[languageCode]
		if hasURL {
			return destinationURL
		}
	}

	if a.DestinationURLsByCountry != nil && countryCode != "" {
		destinationURL, hasURL := a.DestinationURLsByCountry[countryCode]
		if hasURL {
			return destinationURL
		}
	}

	return a.DefaultDestinationURL
}
