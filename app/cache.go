package app

import (
	"archive/zip"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"

	"github.com/Sirupsen/logrus"
)

func UpdateCache(config Config) error {
	LOG.Info("Downloading fresh tldr's from " + config.SourceURL)
	cacheLocation := getCacheLocation()
	if err := os.RemoveAll(cacheLocation); err != nil {
		LOG.Fatal(err)
	}
	LOG.Debug("Deleted cache location, if it had existed ")
	if err := os.MkdirAll(cacheLocation, 0755); err != nil {
		LOG.Fatal(err)
	}
	LOG.Debug("Created cache location")

	zipLocation := filepath.Join(cacheLocation, "master.zip")
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
	if err := Unzip(zipLocation, cacheLocation); err != nil {
		LOG.Fatal(err)
	}
	os.Remove(zipLocation)
	LOG.Debug("Unzip completed")
	return nil
}

func LocalCacheAvailable() bool {
	LOG.Debug("In localCacheAvailable()")
	path := filepath.Join(getCacheLocation(), "tldr-master", "pages", "index.json")

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

func getCacheLocation() string {
	usr, err := user.Current()
	if err != nil {
		LOG.Fatal(err)
	}
	LOG.WithFields(logrus.Fields{
		"usr.HomeDir": usr.HomeDir,
	}).Debug("Printing usr.HomeDir")

	path := filepath.Join(usr.HomeDir, ".tldr", "cache")
	return path
}

func getPageLocation(command, platform string) (string, error) {
	path := filepath.Join(getCacheLocation(), "tldr-master", "pages", platform, command+".md")
	if _, err := os.Stat(path); err == nil {
		LOG.WithFields(logrus.Fields{
			"path": path,
		}).Debug("Page available")
		return path, nil
	}
	return "", COMMAND_NOT_FOUND
}

func Unzip(src, dest string) error {
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
