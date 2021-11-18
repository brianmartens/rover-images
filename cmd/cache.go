package cmd

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"strings"
)

const (
	CacheFile  = "cache_file"
	RoverCache = "rover-cache"
)

/*
	Type cache will be used to store and retrieve mars rover images using the following hierarchy:
		Top level: rover name
		Mid level: camera name
		Bottom level: image date
		values: image urls
*/
type cache map[string]map[string]map[string][]string

// initialize the cache and ensure that maps are being properly initialized prior to writing to them
func initCache() {
	c := getCache()
	if _, ok := c[rover]; !ok {
		c[rover] = make(map[string]map[string][]string)
	}
	if _, ok := c[rover][camera]; !ok {
		c[rover][camera] = make(map[string][]string)
	}
	viper.Set(RoverCache, c)
}

// getCache will attempt to retrieve the cache from the configured filepath and will return the cache if it exists, or an empty cache if it does not
func getCache() cache {
	var c cache
	filename := viper.GetString(CacheFile)

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		switch {
		// cache has yet to be saved to the filesystem
		case strings.Contains(err.Error(), "no such file or directory"):
			c = make(cache)
			return c
		case err != nil:
			logrus.Fatalln("error getting cache from file:", err)
		}
	}

	// return existing cache retrieved from the filesystem
	if err := json.Unmarshal(file, &c); err != nil {
		logrus.Fatalln("error getting cache from file:", err)
	}
	return c
}

// putCache will attempt to save the cache back to the filesystem.
func putCache() {
	c := viper.Get(RoverCache).(cache)
	filename := viper.GetString(CacheFile)

	b, err := json.Marshal(c)
	if err != nil {
		logrus.Fatalln("error saving cache to file:", err)
	}

	if err := ioutil.WriteFile(filename, b, 0666); err != nil {
		logrus.Fatalln("error saving cache to file:", err)
	}
}

// TODO: Query is intended to be a more advanced query functionality for the cache
func (c cache) Query() []string {
	return nil
}
