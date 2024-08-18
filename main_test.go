package main

import (
	"bytes"
	"encoding/json"
	"github.com/bos-hieu/shortlink/internal/entities"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestAPIs(t *testing.T) {
	initClients()
	router := initGinRouter()

	if reflect.TypeOf(router) != reflect.TypeOf(&gin.Engine{}) {
		t.Errorf("Expected type of router is *gin.Engine, but got %T", router)
	}

	// Test GET /ping and response message is "pong"
	t.Run("GET /ping", func(t *testing.T) {
		// Create a response recorder
		w := performRequest(router, "GET", "/ping")
		// Check the status code
		if w.Code != 200 {
			t.Errorf("Expected status code is 200, but got %d", w.Code)
		}
		// Check the response body
		if w.Body.String() != "{\"message\":\"pong\"}" {
			t.Errorf("Expected response body is {\"message\":\"pong\"}, but got %s", w.Body.String())
		}
	})

	// Test POST /short-link with request body is from entities.CreateShortLinkRequest
	// Then call the GET /:short-link to check the response
	t.Run("POST /short-link", func(t *testing.T) {
		requestBody := &entities.CreateShortLinkRequest{
			DefaultURL: "https://google.com",
			CountriesURLs: map[string]string{
				"VN": "https://google.com/vn",
				"SG": "https://google.com/sg",
			},
			LanguagesURLs: map[string]string{
				"vi-VN": "https://google.com/vi-VN",
				"en-SG": "https://google.com/en-SG",
				"zh-SG": "https://google.com/zh-SG",
			},
		}

		// create a post request with request body
		w := performRequestWithBody(router, "POST", "/short-link", requestBody)

		// Check the status code
		postStatusCode := 200
		if w.Code != postStatusCode {
			t.Errorf("Expected status code is %d, but got %d", postStatusCode, w.Code)
			t.Log("Response body: ", w.Body.String())
			t.FailNow()
		}

		// parse the response body to entities.CreateShortLinkResponse
		var responseBody entities.CreateShortLinkResponse
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		if err != nil {
			t.Errorf("Failed to parse response body: %s", err.Error())
		}

		if responseBody.ShortLink == "" {
			t.Fatal("Expected short link is not empty, but got empty")
		}

		// Get short link hash from response body
		shortLink := responseBody.ShortLink[len("localhost:8080/"):]
		log.Println("Short link hash: ", shortLink)

		// Test cases for getting short link
		testCases := []*struct {
			name                      string
			shortLink                 string
			countryHeader             string
			languageHeader            string
			expectedRedirectURL       string
			expectedStatusCode        int
			firstTimePerformance      time.Duration
			secondTimePerformance     time.Duration
			notNeedToCheckPerformance bool
		}{
			{
				name:                "Get short link with country header is VN",
				shortLink:           shortLink,
				countryHeader:       "VN",
				expectedRedirectURL: "https://google.com/vn",
				expectedStatusCode:  307,
			},
			{
				name:                "Get short link with country header is SG",
				shortLink:           shortLink,
				countryHeader:       "SG",
				expectedRedirectURL: "https://google.com/sg",
				expectedStatusCode:  307,
			},
			{
				name:                "Get short link with language header is vi-VN",
				shortLink:           shortLink,
				countryHeader:       "VN",
				languageHeader:      "vi-VN",
				expectedRedirectURL: "https://google.com/vi-VN",
				expectedStatusCode:  307,
			},
			{
				name:                "Get short link with language header is en-SG",
				shortLink:           shortLink,
				countryHeader:       "SG",
				languageHeader:      "en-SG",
				expectedRedirectURL: "https://google.com/en-SG",
				expectedStatusCode:  307,
			},
			{
				name:                "Get short link with language header is zh-SG",
				shortLink:           shortLink,
				countryHeader:       "SG",
				languageHeader:      "zh-SG",
				expectedRedirectURL: "https://google.com/zh-SG",
				expectedStatusCode:  307,
			},
			{
				name:                "Get short link without country header and language header",
				shortLink:           shortLink,
				expectedRedirectURL: "https://google.com",
				expectedStatusCode:  307,
			},
			{
				name:                      "Get short link with invalid short link",
				shortLink:                 "invalid-short-link",
				expectedRedirectURL:       "/404",
				expectedStatusCode:        307,
				notNeedToCheckPerformance: true,
			},
		}

		// Test getting short link with different headers - first time - data will be fetched from database
		for _, tc := range testCases {
			t.Run(tc.name+" first time", func(t *testing.T) {
				now := time.Now()
				// Create a request with short link with country header and language header
				w := performRequestWithHeader(router, "GET", "/"+tc.shortLink, map[string]string{
					"CF-IPCountry":    tc.countryHeader,
					"Accept-Language": tc.languageHeader,
				})
				tc.firstTimePerformance = time.Since(now)

				// Check the status code
				if w.Code != tc.expectedStatusCode {
					t.Errorf("Expected status code is %d, but got %d", tc.expectedStatusCode, w.Code)
				}

				// Check the redirect URL
				if w.Header().Get("Location") != tc.expectedRedirectURL {
					t.Errorf("Expected redirect URL is %s, but got %s", tc.expectedRedirectURL, w.Header().Get("Location"))
				}
			})
		}

		// Test getting short link with different headers - second time - data will be fetched from cache
		for _, tc := range testCases {
			t.Run(tc.name+" second time", func(t *testing.T) {
				now := time.Now()
				// Create a request with short link with country header and language header
				w := performRequestWithHeader(router, "GET", "/"+tc.shortLink, map[string]string{
					"CF-IPCountry":    tc.countryHeader,
					"Accept-Language": tc.languageHeader,
				})
				tc.secondTimePerformance = time.Since(now)

				// Check the status code
				if w.Code != tc.expectedStatusCode {
					t.Errorf("Expected status code is %d, but got %d", tc.expectedStatusCode, w.Code)
				}

				// Check the redirect URL
				if w.Header().Get("Location") != tc.expectedRedirectURL {
					t.Errorf("Expected redirect URL is %s, but got %s", tc.expectedRedirectURL, w.Header().Get("Location"))
				}
			})
		}

		// Check the performance of getting short link between first time and second time
		for _, tc := range testCases {
			if tc.notNeedToCheckPerformance {
				continue
			}
			t.Run(tc.name+" performance", func(t *testing.T) {
				if tc.firstTimePerformance < tc.secondTimePerformance {
					t.Errorf("Expected the first time performance is greater than the second time performance, but got %v < %v", tc.firstTimePerformance, tc.secondTimePerformance)
				}
			})
		}
	})

}

// performRequest is a helper function to perform a request to the router
func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, nil)
	r.ServeHTTP(w, req)
	return w
}

// performRequestWithBody is a helper function to perform a request to the router with request body
func performRequestWithBody(r http.Handler, method, path string, body interface{}) *httptest.ResponseRecorder {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, bytes.NewReader(bodyBytes))
	r.ServeHTTP(w, req)
	return w
}

// performRequestWithHeader is a helper function to perform a request to the router with request header
func performRequestWithHeader(r http.Handler, method, path string, headers map[string]string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, nil)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	r.ServeHTTP(w, req)
	return w
}
