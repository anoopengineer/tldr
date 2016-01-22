package main

import (
	"os"

	app "github.com/anoopengineer/tldr/app"
)

var (
	VERSION = "" //AUTO Populated by goxc
)

func main() {
	tldr := app.MainApp(VERSION)
	tldr.Run(os.Args)
}
