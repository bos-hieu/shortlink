package entities

import "fmt"

type CreateShortLinkRequest struct {
	DefaultURL    string            `json:"url"`
	CountriesURLs map[string]string `json:"countries_urls"`
	LanguagesURLs map[string]string `json:"languages_urls"`
}

type CreateShortLinkResponse struct {
	ShortLink string `json:"short_link"`
}

type GetShortLinkRequest struct {
	ShortLink    string `json:"short_link"`
	CountryCode  string `json:"country_code"`
	LanguageCode string `json:"language_code"`
}

func (a GetShortLinkRequest) GetRedisKey() string {
	return fmt.Sprintf("short_link_%s_%s_%s", a.ShortLink, a.CountryCode, a.LanguageCode)
}

type GetShortLinkResponse struct {
	DestinationURL string `json:"destination_url"`
}