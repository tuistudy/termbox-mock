package view

import (
	"fmt"
	runewidth "github.com/mattn/go-runewidth"
	termbox "github.com/nsf/termbox-go"
	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

/*
func FakeInput(s []byte) {
	if len(inbuf) > 0 && len(s) <= len(inbuf) {
		copy(inbuf, s)
	} else {
		inbuf = append(inbuf, s...)
	}
}
*/
func TestTermboxImp(t *testing.T) {

	x := new(termboxImpl)
	var term Terminal = x

	term.Init()
	term.SetInputMode()
	term.Clear()

	//term.Interrupt()
	count := 0
	x.proxy = func() termbox.Event {
		time.Sleep(15 * time.Millisecond)
		count++
		switch count {
		case 1:
			return termbox.Event{Type: termbox.EventResize}
		case 2:
			return termbox.Event{Type: termbox.EventInterrupt}
		default:
			return termbox.Event{Type: termbox.EventKey, Key: termbox.KeyEsc}
		}
	}

	for {
		term.Clear()
		printMsgf(term, 1, 1, uint16(termbox.ColorWhite), uint16(termbox.ColorDefault), "Press ESC to stop the test.%d", count)
		term.Flush()
		if term.WaitEvent() {
			break
		}
	}
	term.Clear()
	//term.Flush()
	term.Close()
}

func TestInterrupt(t *testing.T) {
	x := new(termboxImpl)
	var term Terminal = x

	term.Init()
	term.SetInputMode()
	term.Clear()
	//term.Flush()

	go wakeMeUp(term)

	for i := 0; i < 1; i++ {
		term.WaitEvent()
	}

	term.Clear()
	//term.Flush()
	term.Close()
}

func TestEventErrCase(t *testing.T) {
	x := new(termboxImpl)
	var term Terminal = x

	term.Init()
	term.SetInputMode()
	term.Clear()
	//term.Flush()

	//term.Interrupt()
	x.proxy = func() termbox.Event {
		time.Sleep(10 * time.Millisecond)
		return termbox.Event{Type: termbox.EventError, Err: fmt.Errorf("after 10 ms, fake an error", 10)}
	}

	defer func() {
		e := recover()
		assert.NotNil(t, e, "should recover from a panic %v", e)
		term.Clear()
		term.Close()
	}()

	term.WaitEvent()
	t.Errorf("EventErr test error")
}

func wakeMeUp(term Terminal) {
	ticker := time.NewTicker(150 * time.Millisecond)
	for _ = range ticker.C {
		term.Interrupt()
	}
}

func printMsgf(term Terminal, x, y int, fg, bg uint16, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	for _, c := range msg {
		term.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}
