package trending

import (
	"encoding/json"
	"fmt"

	"github.com/apex/log"
	"github.com/github-trending/github-trending"

	"github.com/github-trending/api-service"
	"github.com/github-trending/api-service/config"
	"github.com/github-trending/api-service/storage"
)

var attr string = config.Get("redis_addr")
var auth string = config.Get("redis_auth")
var debug string = config.Get("debug")

var Storage = storage.NewStorage(attr, auth, debug)

func init() {
	since := []string{"daily", "weekly", "monthly"}

	// Init storage.
	for _, item := range since {
		if exists, err := Storage.HExists("repositories", item); err != nil {
			log.WithError(err).Fatal("init storage")
		} else if !exists {
			log.WithFields(log.Fields{
				"key":   "repositories",
				"field": item,
				"value": "",
			}).Info("init storage")

			Storage.HSet("repositories", item, "")
		}
	}
}

func Repos(since, language string) ([]api.Repository, error) {
	var key string

	if language != "" {
		key = fmt.Sprintf("%s_%s", language, "repositories")
	} else {
		key = "repositories"
	}

	// try load data from storage.
	value, err := Storage.HGet(key, since)

	if err != nil {
		return nil, err
	}

	var result []api.Repository

	if value != "" {
		if err := json.Unmarshal([]byte(value), &result); err != nil {
			return nil, err
		} else {
			return result, nil
		}
	}

	log.Debug("cache is empty, fetching data from GitHub")

	t := trending.New()

	repositories, err := t.Since(since).Repos()

	if err != nil {
		return nil, err
	}

	for _, repo := range repositories {
		result = append(result, api.Repository{
			Title:           repo.Title,
			Owner:           repo.Owner,
			Name:            repo.Name,
			Description:     repo.Description,
			Language:        repo.Language,
			Stars:           repo.Stars,
			AdditionalStars: repo.AdditionalStars,
			URL:             repo.URL.String(),
		})
	}

	// marshal data and save it in storage.
	if stringify, err := json.Marshal(result); err != nil {
		log.WithError(err).Error("stringify repositories")
	} else {
		status, err := Storage.HSet(key, since, string(stringify))

		if err != nil {
			log.WithError(err).Error("caching repositories")
		} else if !status {
			log.Warn("caching repositories")
		} else {
			log.Info("caching repositories")
		}
	}

	return result, nil
}
