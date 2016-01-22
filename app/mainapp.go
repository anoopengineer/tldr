package app

import (
	"github.com/codegangsta/cli"
)

func MainApp(version string) *cli.App {
	tldr := cli.NewApp()
	tldr.Name = "tldr"
	tldr.Author = "Anoop Kunjuraman <anoopengineer@gmail.com"
	tldr.Version = version
	tldr.Usage = "Simplified and community-driven man pages"

	tldr.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "Enable debug mode.",
		},
		cli.BoolFlag{
			Name:  "no-render, r",
			Usage: "Show the content without rendering it",
		},
		cli.BoolFlag{
			Name:  "termui",
			Usage: "Show the content rendered using termui rather than termbox",
		},
		cli.BoolFlag{
			Name:  "update-cache, c",
			Usage: "Update the local cache",
		},
		cli.BoolFlag{
			Name:  "help, h",
			Usage: "Display help",
		},
		cli.StringFlag{
			Name:  "source, s",
			Value: "https://github.com/tldr-pages/tldr/archive/master.zip",
			Usage: "Source URL pointing to the zip file",
			// EnvVar: "LEGACY_COMPAT_LANG,APP_LANG,LANG",
		},
	}
	tldr.Action = Execute
	tldr.HideHelp = true
	return tldr
}
