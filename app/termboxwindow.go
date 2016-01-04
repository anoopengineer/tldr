package app

import (
	"sync"

	termbox "github.com/nsf/termbox-go"
)

const TOP_LEFT = '┌'
const TOP_RIGHT = '┐'
const VERTICAL_LINE = '│'
const HORIZONTAL_LINE = '─'
const BOTTOM_RIGHT = '┘'
const BOTTOM_LEFT = '└'
const VERTICAL_LEFT = '┤'
const VERTICAL_RIGHT = '├'
const HORIZONTAL_DOWN = '┬'
const HORIZONTAL_UP = '┴'
const QUOTA_LEFT = '«'
const QUOTA_RIGHT = '»'

var screenMutex sync.Mutex

var curWindow *Window

type Widget struct {
	Text   string
	Fg, Bg []termbox.Attribute
}

type Window struct {
	sync.Mutex
	x, y         int
	scrollPos    int
	Fg, Bg       termbox.Attribute
	Widgets      []Widget
	EventHandler func(termbox.Event)
}

func GetWindowInstance() *Window {
	screenMutex.Lock()
	defer screenMutex.Unlock()
	if curWindow != nil {
		return curWindow
	}
	curWindow := &Window{sync.Mutex{}, 0, 0, 0, termbox.ColorWhite, termbox.ColorDefault, nil, nil}
	return curWindow
}

func (w *Window) MainUIHandler() error {
	for {
		e := termbox.PollEvent()

		if e.Type == termbox.EventError {
			return (e.Err)
		}

		// handle resize
		if e.Type == termbox.EventResize {
			err := w.Repaint()
			if err != nil {
				return (err)
			}

			termbox.Flush()
			continue
		}

		if w.EventHandler != nil {
			w.EventHandler(e)
		}
	}
}

func (w *Window) Paint(widgets []Widget) error {
	w.Lock()
	w.Widgets = widgets
	w.Unlock()
	return w.Repaint()
}

func (w *Window) Repaint() error {

	w.Lock()
	err := termbox.Clear(w.Fg, w.Bg)
	width, height := termbox.Size()

	if err != nil {
		return err
	}

	for _, row := range w.Widgets[w.scrollPos:] {
		j := 0
		for _, c := range row.Text {

			termbox.SetCell(w.x, w.y, rune(c), row.Fg[j], row.Bg[j])
			j++

			if w.x += 1; w.x > width {
				break
			}
		}
		w.x = 0

		if w.y += 1; w.y > height {
			break
		}
	}
	w.y = 0
	w.x = 0

	w.Unlock()
	return nil
}

func (w *Window) ScrollUp() {
	if _, height := termbox.Size(); len(w.Widgets) <= height {
		return
	}
	if w.scrollPos > 0 {
		w.scrollPos -= 1
	}
}

func (w *Window) ScrollDown() {
	_, height := termbox.Size()
	if len(w.Widgets) <= height {
		return
	}
	rowLength := len(w.Widgets)
	if w.scrollPos < (rowLength - height) {
		w.scrollPos += 1
	}
}
