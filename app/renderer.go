package app

import (
	"bytes"

	"github.com/Sirupsen/logrus"
)

const WIDTH = 80

func Render(config Config, input []byte) error {
	tldr := NewTLDR(input)
	LOG.WithFields(logrus.Fields{
		"tldr": tldr,
	}).Debug("Printing tldr")

	if config.RenderOnTermUI {
		return RenderOnTermUI(tldr)
	} else {
		return RenderOnTermbox(tldr)
	}
}

func getNumberOfLines(input []byte) int {
	lineSep := []byte{'\n'}
	count := bytes.Count(input, lineSep)
	if count == 0 {
		return 1
	}
	return count
}
