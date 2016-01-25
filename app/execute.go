package app

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func Execute(c *cli.Context) {
	config := populateConfig(c)
	InitializeLogging(config)
	LOG.WithFields(logrus.Fields{
		"config": config,
	}).Debug("Printing config")
	command, err := getCommand()
	if err != nil {
		fmt.Println("\nCommand name is missing in the input\n")
		cli.ShowAppHelp(c)
		return
	}

	if config.UpdateCache || !localCacheAvailable() || localCacheExpired() {
		if err := updateCache(config); err != nil {
			log.Fatal(err)
		}
	}
	LOG.WithFields(logrus.Fields{
		"command": command,
	}).Debug("Printing command")
	path, err := getPageLocation(command, config.Platform)

	if err == COMMAND_NOT_FOUND {
		fmt.Printf("TLDR for command %s was not found.\nYou can get your TLDRs added by sending a pull request to https://github.com/tldr-pages/tldr\n", command)
		os.Exit(-1)
	}
	if err != nil {
		LOG.Fatal(err)
	}
	LOG.WithFields(logrus.Fields{
		"path": path,
	}).Debug("Printing path")

	content, err := ioutil.ReadFile(path)
	if err != nil {
		LOG.Fatal(err)
	}

	if config.NoRender {
		fmt.Println(string(content))
	} else {
		Render(config, content)
	}
}

func getCommand() (string, error) {
	if len(os.Args) < 2 {
		return "", COMMAND_NOT_PROVIDED
	}
	command := os.Args[len(os.Args)-1]
	if strings.HasPrefix(command, "-") {
		return "", COMMAND_NOT_PROVIDED
	}
	return command, nil
}
