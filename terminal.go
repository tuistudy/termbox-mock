package view

import (
	//"fmt"
	termbox "github.com/nsf/termbox-go"
)

//Terminal is a simple wrapper for the Terminal display purpose
type Terminal interface {
	Init()
	Interrupt()
	Close()
	Flush()
	SetCell(x, y int, c rune, fg, bg uint16)
	WaitEvent() bool
	SetInputMode()
	Clear()
}

//Terminal implementation
type termboxImpl struct {
	ver   string
	proxy func() termbox.Event
}

func (tb *termboxImpl) Init() {
	termbox.Init()
	//this should be the default , but you can change the proxy when you need it
	tb.proxy = termbox.PollEvent
}

func (tb termboxImpl) Clear() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

func (tb termboxImpl) Close() {
	termbox.Close()
}

func (tb termboxImpl) Flush() {
	termbox.Flush()
}

func (tb termboxImpl) Interrupt() {
	termbox.Interrupt()
}

func (tb termboxImpl) SetInputMode() {
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
}

func (tb termboxImpl) SetCell(x, y int, c rune, fg, bg uint16) {
	termbox.SetCell(x, y, c, termbox.Attribute(fg), termbox.Attribute(bg))
}

func (tb termboxImpl) WaitEvent() bool {
	switch ev := tb.proxy(); ev.Type { // PollEvent will be blocked
	case termbox.EventKey:
		if ev.Key == termbox.KeyEsc {
			return true
		}
	case termbox.EventResize:
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	case termbox.EventError:
		panic(ev.Err)
	}
	return false
}
