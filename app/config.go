package app

import "github.com/codegangsta/cli"

type Config struct {
	SourceURL      string
	Debug          bool
	NoRender       bool
	RenderOnTermUI bool
	UpdateCache    bool
}

func populateConfig(c *cli.Context) (config Config) {
	config.Debug = c.Bool("debug")
	config.SourceURL = c.String("source")
	config.NoRender = c.Bool("no-render")
	config.RenderOnTermUI = c.Bool("termui")
	config.UpdateCache = c.Bool("update-cache")

	return config
}
