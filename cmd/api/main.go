package main

import (
	"os"

	"github.com/apex/log"
	"github.com/sqrthree/debugfmt"
	"github.com/kataras/iris"

	"github.com/github-trending/api-service"
	"github.com/github-trending/api-service/config"
	"github.com/github-trending/api-service/trending"
)

var addr = iris.Addr(":8080")

var trendingDeamon *trending.Deamon

func main() {
	debug := config.Get("debug")

	app := iris.New()

	if debug == "true" {
		log.SetLevel(log.DebugLevel)
		log.SetHandler(debugfmt.New(os.Stdout))
		app.Logger().SetLevel("debug")
	}

	trendingDeamon = trending.StartDeamon()

	// refrech latest data when app is restarted.
	go trendingDeamon.Refrech()

	app.Use(func(ctx iris.Context) {
		ctx.Application().Logger().Debugf("--> %s %s", ctx.Method(), ctx.Path())
		ctx.Next()
		ctx.Application().Logger().Debugf("<-- %s %s %d", ctx.Method(), ctx.Path(), ctx.GetStatusCode())
	})

	// Register custom handler for specific http errors.
	app.OnErrorCode(iris.StatusBadRequest, handleBadRequest)
	app.OnErrorCode(iris.StatusNotFound, handleBadRequest)
	app.OnErrorCode(iris.StatusServiceUnavailable, handleServiceUnavailable)

	// Register routes
	app.Get("/", getHATEOAS)
	app.Get(api.RootEndpoint, getHATEOAS)
	app.Get(api.RepositoryEndpoint, getRepos)

	app.Run(addr, iris.WithCharset("UTF-8"))
}

// handleBadRequest handles 400 request.
func handleBadRequest(ctx iris.Context) {
	ctx.Application().Logger().Infof("<-- %s %s %d", ctx.Method(), ctx.Path(), ctx.GetStatusCode())

	ctx.JSON(api.ErrorBadRequest)
}

// handleBadRequest handles 404 request.
func handleNotFound(ctx iris.Context) {
	ctx.Application().Logger().Infof("<-- %s %s %d", ctx.Method(), ctx.Path(), ctx.GetStatusCode())

	ctx.JSON(api.ErrorNotFound)
}

// handleBadRequest handles 500 request.
func handleServiceUnavailable(ctx iris.Context) {
	ctx.Application().Logger().Infof("<-- %s %s %d", ctx.Method(), ctx.Path(), ctx.GetStatusCode())

	ctx.JSON(api.ErrorServiceUnavailable)
}

// handleBadRequest handles `GET /` and `GET /api` request, it reflects [HATEOAS](https://en.wikipedia.org/wiki/HATEOAS).
func getHATEOAS(ctx iris.Context) {
	ctx.JSON(api.HATEOAS)
}

// getRepos returns repositories from https://github.com/trending
func getRepos(ctx iris.Context) {
	since := ctx.URLParam("since")

	if since == "" {
		since = "daily"
	}

	ctx.Application().Logger().Debugf("request repositories with param <since>: %s", since)

	var data []api.Repository

	err := trendingDeamon.GetJSON("repositories", since, &data)

	if err != nil {
		ctx.Application().Logger().Error(err)
		ctx.StatusCode(iris.StatusServiceUnavailable)
		ctx.Next()
		return
	}

	lastModifiedAt, err := trendingDeamon.Get("last_modified_time_of_repositories", since)

	if err == nil {
		ctx.Header("Last-Modified", lastModifiedAt)
	} else {
		ctx.Application().Logger().Error(err)
	}

	ctx.JSON(data)
}
