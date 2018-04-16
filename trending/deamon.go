package trending

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/apex/log"
	"github.com/github-trending/github-trending"

	"github.com/github-trending/api-service"
	"github.com/github-trending/api-service/config"
	"github.com/github-trending/api-service/storage"
)

var attr string = config.Get("redis_addr")
var auth string = config.Get("redis_auth")
var debug string = config.Get("debug")

type Deamon struct {
	storage             storage.Storage
	refreshRequestChan  chan bool
	refreshResponseChan chan bool
}

// StartDeamon return a deamon server.
func StartDeamon() *Deamon {
	storage := storage.NewStorage(attr, auth, debug)

	d := Deamon{
		storage:             storage,
		refreshRequestChan:  make(chan bool),
		refreshResponseChan: make(chan bool),
	}

	d.InitStorage()

	go d.UpdateDaemon(2 * time.Hour)

	return &d
}

// InitStorage initializes storage.
// trying to get a key that doesn't exist will make an error in Redis package.
func (d *Deamon) InitStorage() {
	since := []string{"daily", "weekly", "monthly"}

	// Init storage if the key doesn't exist.
	for _, item := range since {
		if exists, err := d.storage.HExists("repositories", item); err != nil {
			log.WithError(err).Fatal("init storage")
		} else if !exists {
			log.WithFields(log.Fields{
				"key":   "repositories",
				"field": item,
				"value": "",
			}).Info("init storage")

			d.storage.HSet("repositories", item, "")
		}
	}
}

// Get is a thin wrapper around Storage.HGet().
func (d *Deamon) Get(key, field string) (string, error) {
	log.WithFields(log.Fields{
		"key":   key,
		"field": field,
	}).Debugf("get data from storage")

	value, err := d.storage.HGet(key, field)

	return value, err
}

// GetJSON gets the JSON-encoded data, parses it and stores the result in the value pointed to by v.
func (d *Deamon) GetJSON(key, field string, v interface{}) error {
	log.WithFields(log.Fields{
		"key":   key,
		"field": field,
	}).Debugf("get json-data from storage")

	value, err := d.storage.HGet(key, field)

	if err != nil {
		return err
	}

	if value == "" {
		return errors.New("data is empty")
	}

	if err := json.Unmarshal([]byte(value), v); err != nil {
		log.WithError(err).Errorf("stringify json-data from storage")
		return err
	} else {
		return nil
	}
}

// Set is a thin wrapper around Storage.HSet().
func (d *Deamon) Set(key, field, value string) (bool, error) {
	log.WithFields(log.Fields{
		"key":   key,
		"field": field,
	}).Debugf("set data in storage")

	status, err := d.storage.HSet(key, field, value)

	return status, err
}

// Refrech notices UpdateDeamon to refrech the data.
func (d *Deamon) Refrech() {
	log.Debugf("refresh data")

	d.refreshRequestChan <- true

	<-d.refreshResponseChan

	log.Debugf("refresh data done")
}

// UpdateDaemon represents a update deamon server.
func (d *Deamon) UpdateDaemon(tickDuration time.Duration) {
TICK_DEAMON:
	ticker := time.NewTicker(tickDuration)

	for {
		select {
		case <-d.refreshRequestChan:
			d.updateRepos()

			d.refreshResponseChan <- true

			goto TICK_DEAMON
		case <-ticker.C:
			d.updateRepos()

			goto TICK_DEAMON
		}
	}
}

func (d *Deamon) updateRepos() {
	tasks := []string{"daily", "weekly", "monthly"}

	for _, task := range tasks {
		repos, err := fetchRepos(task)

		if err != nil {
			log.WithError(err).Errorf("fetch %s repos failed", task)
			continue
		}

		err = d.cacheRepos(task, repos)

		if err != nil {
			log.WithError(err).Errorf("cache %s repos failed", task)
			continue
		}
	}
}

// fetchRepos fetches repositories from GitHub.
func fetchRepos(since string) ([]api.Repository, error) {
	log.Infof("load data from GitHub with since = %s", since)

	t := trending.New()

	repos, err := t.Since(since).Repos()

	if err != nil {
		return nil, err
	}

	var result []api.Repository

	for _, repo := range repos {
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

	return result, nil
}

// cacheRepos stores the data to storage.
func (d *Deamon) cacheRepos(field string, value []api.Repository) error {
	result, err := json.Marshal(value)

	if err != nil {
		log.WithError(err).Errorf("stringify %s repositories", field)
		return err
	}

	_, err = d.Set("repositories", field, string(result))

	if err != nil {
		return err
	}

	_, err = d.Set("last_modified_time_of_repositories", field, time.Now().Format(time.RFC1123))

	if err != nil {
		return err
	}

	return nil
}
