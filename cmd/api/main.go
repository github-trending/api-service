package main

import (
	"fmt"

	"github.com/github-trending/github-trending"
	"github.com/kataras/iris"

	"github.com/github-trending/online-api/config"
	"github.com/github-trending/online-api"
)

var addr = iris.Addr(":8080")

func main() {
	debug := config.Get("debug")

	app := iris.New()

	if debug == "true" {
		app.Logger().SetLevel("debug")
	}

	app.Use(func(ctx iris.Context) {
		ctx.Application().Logger().Debugf("--> %s %s", ctx.Method(), ctx.Path())
		ctx.Next()
		ctx.Application().Logger().Debugf("<-- %s %s %d", ctx.Method(), ctx.Path(), ctx.GetStatusCode())
	})

	app.OnErrorCode(iris.StatusBadRequest, handleBadRequest)
	app.OnErrorCode(iris.StatusNotFound, handleBadRequest)
	app.OnErrorCode(iris.StatusServiceUnavailable, handleServiceUnavailable)

	app.Get("/", getHATEOAS)
	app.Get(api.RootEndpoint, getHATEOAS)
	app.Get(api.RepositoryEndpoint, getRepos)

	app.Run(addr, iris.WithCharset("UTF-8"))
}

func handleBadRequest(ctx iris.Context) {
   ctx.Application().Logger().Infof("<-- %s %s %d", ctx.Method(), ctx.Path(), ctx.GetStatusCode())

   ctx.JSON(api.ErrorBadRequest)
}

func handleNotFound(ctx iris.Context) {
   ctx.Application().Logger().Infof("<-- %s %s %d", ctx.Method(), ctx.Path(), ctx.GetStatusCode())

   ctx.JSON(api.ErrorNotFound)
}

func handleServiceUnavailable(ctx iris.Context) {
   ctx.Application().Logger().Infof("<-- %s %s %d", ctx.Method(), ctx.Path(), ctx.GetStatusCode())

   ctx.JSON(api.ErrorServiceUnavailable)
}

func getHATEOAS(ctx iris.Context) {
	ctx.JSON(api.HATEOAS)
}

func getRepos(ctx iris.Context) {
	since := ctx.URLParam("since")

	if since == "" {
		since = "daily"
	}

	ctx.Application().Logger().Debugf("request repositories with since param: %s", since)

	t := trending.New()

	data, err := t.Since(since).Repos()

	fmt.Println(t, data)

	if err != nil {
		ctx.StatusCode(iris.StatusServiceUnavailable)
		ctx.Next()
		return
	}

	var result []api.Repository

	for _, repo := range data {
		result = append(result, api.Repository{
			Title: repo.Title,
			Owner: repo.Owner,
			Name: repo.Name,
			Description: repo.Description,
			Language: repo.Language,
			Stars: repo.Stars,
			AdditionalStars: repo.AdditionalStars,
			URL: repo.URL.String(),
		})
	}

	ctx.JSON(result)
}
