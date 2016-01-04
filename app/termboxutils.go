package app

import (
	"os"
	"strings"

	"github.com/nsf/termbox-go"
)

func RenderOnTermbox(tldr TLDR) error {
	if err := termbox.Init(); err != nil {
		return err
	}
	defer termbox.Close()
	mainWindow := GetWindowInstance()
	rows := getHeaderWidget()
	rows = append(rows, tldr.toWidgets()...)
	if err := mainWindow.Paint(rows); err != nil {
		return err
	}

	termbox.Flush()

	mainWindow.EventHandler = func(event termbox.Event) {
		if event.Ch == 'j' || event.Key == termbox.KeyArrowDown {
			mainWindow.ScrollDown()
		} else if event.Ch == 'k' || event.Key == termbox.KeyArrowUp {
			mainWindow.ScrollUp()
		} else if event.Ch == 'q' || event.Key == termbox.KeyEsc {
			termbox.Close()
			termbox.Flush()
			os.Exit(0)
		}
		mainWindow.Repaint()
		termbox.Flush()
	}
	return mainWindow.MainUIHandler()
}

func getHeaderWidget() []Widget {
	var header = make([]Widget, 0)
	header = append(header, getHeaderHeader("┌TLDR──────────────────────────────────────────────────────────────────────────┐"))
	header = append(header, getHeaderBody("| PRESS q or ESC to QUIT                                                       |"))
	header = append(header, getHeaderBody("|                                                                              |"))
	header = append(header, getHeaderSpecialBody("| TLDR - Simplified and community-driven man pages                             |"))
	header = append(header, getHeaderFooter("└──────────────────────────────────────────────────────────────────────────────┘"))
	header = append(header, getHeaderHeader(""))
	return header
}

func getHeaderHeader(str string) Widget {
	fgColors := make([]termbox.Attribute, WIDTH)
	bgColors := make([]termbox.Attribute, WIDTH)
	for i, _ := range fgColors {
		fgColors[i] = termbox.ColorCyan
		bgColors[i] = termbox.ColorDefault
	}
	fgColors[1] = termbox.ColorGreen
	fgColors[2] = termbox.ColorGreen
	fgColors[3] = termbox.ColorGreen
	fgColors[4] = termbox.ColorGreen
	return Widget{str, fgColors, bgColors}
}

func getHeaderFooter(str string) Widget {
	fgColors := make([]termbox.Attribute, WIDTH)
	bgColors := make([]termbox.Attribute, WIDTH)
	for i, _ := range fgColors {
		fgColors[i] = termbox.ColorCyan
		bgColors[i] = termbox.ColorDefault
	}
	return Widget{str, fgColors, bgColors}
}

func getHeaderBody(str string) Widget {
	fgColors := make([]termbox.Attribute, WIDTH)
	bgColors := make([]termbox.Attribute, WIDTH)
	for i, _ := range fgColors {
		fgColors[i] = termbox.ColorDefault
		bgColors[i] = termbox.ColorDefault
	}
	fgColors[0] = termbox.ColorCyan
	fgColors[WIDTH-1] = termbox.ColorCyan
	return Widget{str, fgColors, bgColors}
}
func getHeaderSpecialBody(str string) Widget {
	fgColors := make([]termbox.Attribute, WIDTH)
	bgColors := make([]termbox.Attribute, WIDTH)
	for i, _ := range fgColors {
		fgColors[i] = termbox.ColorDefault
		bgColors[i] = termbox.ColorDefault
	}
	fgColors[0] = termbox.ColorCyan
	fgColors[WIDTH-1] = termbox.ColorCyan

	fgColors[2] = termbox.ColorGreen
	fgColors[3] = termbox.ColorGreen
	fgColors[4] = termbox.ColorGreen
	fgColors[5] = termbox.ColorGreen
	return Widget{str, fgColors, bgColors}
}

func getTitleWidget(str string) Widget {
	fgColors := make([]termbox.Attribute, len(str))
	bgColors := make([]termbox.Attribute, len(str))
	for i, _ := range fgColors {
		fgColors[i] = termbox.ColorCyan
		bgColors[i] = termbox.ColorDefault
	}

	return Widget{str, fgColors, bgColors}
}

func getTextWidget(str string) Widget {
	fgColors := make([]termbox.Attribute, len(str))
	bgColors := make([]termbox.Attribute, len(str))
	for i, _ := range fgColors {
		fgColors[i] = termbox.ColorDefault
		bgColors[i] = termbox.ColorDefault
	}

	return Widget{str, fgColors, bgColors}
}

func getCommandDescriptionWidget(str string) Widget {
	fgColors := make([]termbox.Attribute, len(str))
	bgColors := make([]termbox.Attribute, len(str))
	for i, _ := range fgColors {
		fgColors[i] = termbox.ColorDefault
		bgColors[i] = termbox.ColorDefault
	}
	fgColors[1] = termbox.ColorCyan
	return Widget{str, fgColors, bgColors}
}

func (t *TLDR) toWidgets() []Widget {
	retVal := make([]Widget, 0)
	retVal = append(retVal, getTitleWidget(strings.ToUpper(t.Name)))
	retVal = append(retVal, getTextWidget(""))
	for _, v := range t.Description {
		retVal = append(retVal, getTextWidget(v))
	}
	retVal = append(retVal, getTextWidget(""))

	for _, v := range t.Items {
		retVal = append(retVal, getTextWidget(""))
		retVal = append(retVal, getCommandDescriptionWidget(" "+string(QUOTA_RIGHT)+" "+strings.Title(v.Description)))
		for _, c := range v.Command {
			retVal = append(retVal, getTextWidget("    "+c))
		}
	}
	retVal = append(retVal, getTextWidget(""))
	retVal = append(retVal, getTextWidget(""))

	return retVal

}
