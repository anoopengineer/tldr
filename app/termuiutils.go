package app

import (
	"strings"

	"github.com/gizak/termui"
)

func RenderOnTermUI(tldr TLDR) error {
	err := termui.Init()
	if err != nil {
		return err
	}
	defer termui.Close()

	buffers := make([]termui.Bufferer, 0)

	yOffset := 0
	header := termui.NewPar(" PRESS q or ESC TO QUIT\n\n[TLDR](fg-green) - Simplified and community-driven man pages")
	header.Height = 5
	header.Width = WIDTH
	header.TextFgColor = termui.ColorWhite
	header.BorderLabel = "TLDR"
	header.BorderFg = termui.ColorCyan
	header.Y = yOffset
	yOffset += (header.Height + 1)
	buffers = append(buffers, header)

	tldrName := termui.NewPar("[" + strings.ToUpper(tldr.Name) + "](fg-cyan)")
	tldrName.Height = 1
	tldrName.Width = WIDTH
	tldrName.TextFgColor = termui.ColorWhite
	tldrName.Border = false
	tldrName.Y = yOffset
	yOffset += (tldrName.Height + 1)
	buffers = append(buffers, tldrName)

	tldrDesc := termui.NewPar(strings.Join(tldr.Description, "\n"))
	tldrDesc.Height = len(tldr.Description)
	tldrDesc.Width = WIDTH
	tldrDesc.TextFgColor = termui.ColorWhite
	tldrDesc.Border = false
	tldrDesc.Y = yOffset
	yOffset += (tldrDesc.Height + 1)
	buffers = append(buffers, tldrDesc)

	for _, item := range tldr.Items {
		commandDesc := termui.NewPar(" [*](fg-cyan) " + strings.Title(item.Description))
		commandDesc.Height = getNumberOfLines([]byte(item.Description))
		commandDesc.Width = WIDTH
		commandDesc.TextFgColor = termui.ColorWhite
		commandDesc.Border = false
		commandDesc.Y = yOffset
		yOffset += (commandDesc.Height)
		buffers = append(buffers, commandDesc)

		for _, c := range item.Command {
			command := termui.NewPar(c)
			command.Height = getNumberOfLines([]byte(c))
			command.Width = WIDTH - 4
			command.TextFgColor = termui.ColorWhite
			command.Border = false
			command.Y = yOffset
			command.X = 4
			yOffset += (command.Height)
			buffers = append(buffers, command)
		}
		yOffset++
	}

	termui.Render(buffers...) // feel free to call Render, it's async and non-block
	// for _, v := range buffers {
	// 	termui.Body.AddRows(termui.NewRow(termui.NewCol(12, 0, v)))
	// }
	// termui.Body.Align()
	// termui.Render(termui.Body)

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})
	termui.Handle("/sys/kbd/C-c", func(termui.Event) {
		termui.StopLoop()
	})
	termui.Handle("/sys/kbd/<escape>", func(termui.Event) {
		termui.StopLoop()
	})
	// termui.Handle("/sys/kbd", func(e termui.Event) {
	// 	LOG.WithFields(logrus.Fields{
	// 		"tldr": e,
	// 	}).Debug("Printing key")
	// })
	// termui.Handle("/sys/kbd/<down>", func(e termui.Event) {
	// 	maxx, maxy := termbox.Size()
	// 	fmt.Println(maxx, maxy, yOffset)
	//
	// })
	termui.Handle("/sys/wnd/resize", func(e termui.Event) {
		termui.Body.Width = termui.TermWidth()
		termui.Body.Align()
		termui.Render(termui.Body)
	})
	termui.Loop()
	return nil
}
