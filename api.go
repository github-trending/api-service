package api

import (
	"net/http"

	"github.com/github-trending/online-api/config"
)

const (
	DocumentationURL = "/docs"
	RootEndpoint = "/api"
	RepositoryEndpoint = "/api/repositories"
)

var Host = config.Get("host")

type hateoas struct {
	DocumentationURL string `json:"documentation_url"`
	RootEndpoint string `json:"root_endpoint"`
	RepositoryURL string `json:"repository_url"`
}

var HATEOAS = hateoas{
	DocumentationURL: Host + DocumentationURL,
	RootEndpoint: Host + RootEndpoint,
	RepositoryURL: Host + RepositoryEndpoint,
}

type ErrorResponse struct {
	Message string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
}

var ErrorBadRequest = ErrorResponse{
	Message: http.StatusText(http.StatusBadRequest),
	DocumentationURL: Host + DocumentationURL,
}

var ErrorNotFound = ErrorResponse{
	Message: http.StatusText(http.StatusNotFound),
	DocumentationURL: Host + DocumentationURL,
}

var ErrorServiceUnavailable = ErrorResponse{
	Message: http.StatusText(http.StatusServiceUnavailable),
	DocumentationURL: Host + DocumentationURL,
}

type Repository struct {
	Title           string `json:"title"`
	Owner           string `json:"owner"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Language        string `json:""`
	Stars           int	`json:"stars"`
	AdditionalStars int `json:"additional_stars"`
	URL             string `json:"url"`
}
