package app

import (
	"strings"

	"github.com/Sirupsen/logrus"
)

type TLDR struct {
	Name        string
	Description []string
	Items       []Example
}

type Example struct {
	Description string
	Command     []string
}

func NewTLDR(input []byte) TLDR {
	var tldr TLDR
	tldr.Description = make([]string, 0)
	tldr.Items = make([]Example, 0)

	lines := strings.Split(string(input), "\n")
	// LOG.WithFields(logrus.Fields{
	// 	"lines":     strings.Join(lines, "\n"),
	// 	"lines_len": len(tldr.Items),
	// }).Debug("Splitted the data into individual lines")

	for _, v := range lines {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		if strings.HasPrefix(v, "#") {
			tldr.Name = strings.TrimSpace(strings.TrimPrefix(v, "#"))
			LOG.WithFields(logrus.Fields{
				"tldr": tldr,
			}).Debug("Got Name")
		} else if strings.HasPrefix(v, ">") {
			desc := strings.TrimSpace(strings.TrimPrefix(v, ">"))
			tldr.Description = append(tldr.Description, desc)
			LOG.WithFields(logrus.Fields{
				"tldr": tldr,
			}).Debug("Got Description")
		} else if strings.HasPrefix(v, "-") {
			var ex Example
			ex.Command = make([]string, 0)
			ex.Description = strings.TrimSpace(strings.TrimPrefix(v, "-"))
			LOG.WithFields(logrus.Fields{
				"ex":        ex,
				"items_len": len(tldr.Items),
			}).Debug("Adding example")
			tldr.Items = append(tldr.Items, ex)
			LOG.WithFields(logrus.Fields{
				"tldr": tldr,
			}).Debug("Got Command Description")
		} else if strings.HasPrefix(v, "`") {
			command := strings.TrimPrefix(v, "`")
			command = strings.TrimSuffix(command, "`")
			command = strings.TrimSpace(command)
			lastItemIndex := len(tldr.Items) - 1
			lastItem := tldr.Items[lastItemIndex]
			lastItem.Command = append(lastItem.Command, command)
			tldr.Items[lastItemIndex] = lastItem
			LOG.WithFields(logrus.Fields{
				"tldr": tldr,
			}).Debug("Got Command")
		} else {
			LOG.WithFields(logrus.Fields{
				"line": v,
			}).Error("Got unknown line")
		}

	}
	return tldr
}

func (t *TLDR) toArray() []string {
	retVal := make([]string, 0)
	retVal = append(retVal, t.Name)
	retVal = append(retVal, "")
	for _, v := range t.Description {
		retVal = append(retVal, v)
	}
	retVal = append(retVal, "")

	for _, v := range t.Items {
		retVal = append(retVal, "")
		retVal = append(retVal, strings.Title(v.Description))
		for _, c := range v.Command {
			retVal = append(retVal, "    "+c)
		}
	}
	return retVal
}
