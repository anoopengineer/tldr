package app

import (
	"archive/zip"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/Sirupsen/logrus"
)

type Cache struct {
	LastUpdateDate time.Time
	CacheLocation  string
}

func updateCache(config Config) error {
	LOG.Info("Downloading fresh tldr's from " + config.SourceURL)
	currentCacheLocation, _ := getCurrentCacheLocation() // May be nil
	LOG.Debug("Current cache location is " + currentCacheLocation)
	newCacheLocation := getNewCacheLocation()
	LOG.Debug("New cache location is " + newCacheLocation)

	zipLocation := filepath.Join(newCacheLocation, "master.zip")
	out, err := os.Create(zipLocation)
	if err != nil {
		LOG.Fatal(err)
	}
	defer out.Close()
	resp, err := http.Get(config.SourceURL)
	if err != nil {
		LOG.Fatal(err)
	}
	defer resp.Body.Close()
	io.Copy(out, resp.Body)
	LOG.Debug("Download complete")
	if err := unzip(zipLocation, newCacheLocation); err != nil {
		LOG.Fatal(err)
	}
	os.Remove(zipLocation)
	LOG.Debug("Unzip completed")

	//NOW update the cache.json with the new path and delete the old path
	// if everything was successfull so far
	newCache := Cache{
		LastUpdateDate: time.Now(),
		CacheLocation:  newCacheLocation,
	}
	err = updateCacheMetaData(newCache)
	if err != nil {
		LOG.Fatal(err)
	}
	if currentCacheLocation != "" {
		os.RemoveAll(currentCacheLocation)
	}
	return nil
}

func localCacheAvailable() bool {
	LOG.Debug("In localCacheAvailable()")
	path, err := getCurrentCacheLocation()
	if err != nil {
		return false
	}
	path = filepath.Join(path, "tldr-master", "pages", "index.json")

	LOG.WithFields(logrus.Fields{
		"index.json.path": path,
	}).Debug("Printing index.json.path")

	if _, err := os.Stat(path); err == nil {
		LOG.Debug("Local cache exists. Returning true")
		return true
	}
	LOG.Info("Local cached tldr's doesn't exist. Will have to download now.")
	return false
}

func localCacheExpired() bool {
	cache, err := getCacheMetaData()
	if err != nil {
		LOG.Info("Got error from getCacheMetaData. Returning true")
		return true
	}
	cacheCreatedTime := cache.LastUpdateDate
	if time.Now().Sub(cacheCreatedTime) > 30*24*time.Hour {
		LOG.Info("Cache was updated before 30 days. Returning true")
		return true
	}
	return false
}

func getCurrentCacheLocation() (string, error) {
	cache, err := getCacheMetaData()
	if err != nil {
		return "", err
	}
	return cache.CacheLocation, nil
}

func getNewCacheLocation() string {
	usr, err := user.Current()
	if err != nil {
		LOG.Fatal(err)
	}
	path := filepath.Join(usr.HomeDir, ".tldr")
	if err := os.MkdirAll(path, 0755); err != nil {
		LOG.Fatal(err)
	}

	tempDir, err := ioutil.TempDir(path, "cache")
	if err != nil {
		LOG.Fatal(err)
	}
	LOG.WithFields(logrus.Fields{
		"tempDir": tempDir,
	}).Debug("Printing new Cache Location")
	return tempDir
}

func getCacheMetaData() (*Cache, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}
	LOG.WithFields(logrus.Fields{
		"usr.HomeDir": usr.HomeDir,
	}).Debug("Printing usr.HomeDir")

	path := filepath.Join(usr.HomeDir, ".tldr", "cache.json")
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cache Cache
	err = json.Unmarshal(content, &cache)
	return &cache, err
}

func updateCacheMetaData(cache Cache) error {
	usr, err := user.Current()
	if err != nil {
		return err
	}
	LOG.WithFields(logrus.Fields{
		"usr.HomeDir": usr.HomeDir,
	}).Debug("Printing usr.HomeDir")

	path := filepath.Join(usr.HomeDir, ".tldr", "cache.json")
	b, err := json.Marshal(cache)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, b, 0644)
}

func getPageLocation(command, platform string) (string, error) {
	path, err := getCurrentCacheLocation()
	if err != nil {
		return "", err
	}
	path = filepath.Join(path, "tldr-master", "pages", platform, command+".md")
	if _, err := os.Stat(path); err == nil {
		LOG.WithFields(logrus.Fields{
			"path": path,
		}).Debug("Page available")
		return path, nil
	}
	return "", COMMAND_NOT_FOUND
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}
