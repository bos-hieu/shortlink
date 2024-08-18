package main

import (
	"github.com/bos-hieu/shortlink/internal/entities"
	"github.com/bos-hieu/shortlink/internal/usecase"
	"github.com/bos-hieu/shortlink/pkg/mongodb"
	"github.com/bos-hieu/shortlink/pkg/redis"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	initClients()
	router := initGinRouter()
	err := router.Run()
	if err != nil {
		panic(err)
	}
}

// initClients initializes clients
func initClients() {
	err := mongodb.InitClient()
	if err != nil {
		panic(err)
	}

	err = redis.InitClient()
	if err != nil {
		panic(err)
	}
}

// initGinRouter initializes gin router
func initGinRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/short-link", func(context *gin.Context) {
		req := &entities.CreateShortLinkRequest{}
		err := context.ShouldBindJSON(req)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		resp, err := usecase.CreateShortLink(context.Request.Context(), req)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.JSON(http.StatusOK, gin.H{
			"short_link": resp.ShortLink,
		})
	})

	router.GET("/404", func(context *gin.Context) {
		context.String(http.StatusOK, "404")
	})

	router.GET("/:short-link", func(context *gin.Context) {
		shortLink := context.Param("short-link")
		if shortLink == "" {
			context.JSON(http.StatusInternalServerError, gin.H{
				"error": "shortlink is empty",
			})
			return
		}

		resp, err := usecase.GetShortLink(context.Request.Context(), &entities.GetShortLinkRequest{
			ShortLink:    shortLink,
			CountryCode:  context.Request.Header.Get("CF-IPCountry"), // example get country code from cloudflare
			LanguageCode: context.Request.Header.Get("Accept-Language"),
		})
		if err != nil {
			context.Redirect(http.StatusTemporaryRedirect, "/404")
		}

		context.Redirect(http.StatusTemporaryRedirect, resp.DestinationURL)
	})

	return router
}
