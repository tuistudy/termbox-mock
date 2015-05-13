# termbox-mock
a simple way to mock the termbox for testing. 

when i try to test my termbox based application (in golang). i need to find some way to unit test the module. after some tries and discuss with @nsf i finally find a simple way to unit test my application. it's worth to share this idea to others.   

##interface
first, i don't want to expose the termbox API to my application. so i decided to wrap the termbox API with an interface. second, with an interface, it's easy to replace the underlying implementation. here is the interface:

~~~go
i//Terminal is a simple wrapper for the Terminal display purpose
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
~~~

most of the API function in termbox are included, those i never used does not included.

##implementation - the real one
based on the fist opinion, I choose to wrap the termbox.PollEvent function in WaitEvent method. otherwise, the application has to know the termbox API. such as termbox.Event and some constant. that's what i don't want. 

The other choice i made is replacing the termbox.PollEvent with an function value. it's defalut value is termbox.Event. so that in the following tests i have the opportunity to replace this function with new one. this way, I can mock most of the situation. you will see it later.

~~~go
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
~~~

##implementation - the mock one
here is the mock implementation, in TestTermboxImp i replace the proxy with a new one. it's a fuction type ~~~func() termbox.Event ~~~. in this test, the proxy just send the EventResize, EventInterrupt and EventKey in turn. one by one, slep 15ms for every Event. 

~~~go
package view

import (
	"fmt"
	runewidth "github.com/mattn/go-runewidth"
	termbox "github.com/nsf/termbox-go"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

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
~~~

~~~go
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
~~~

~~~go
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
~~~
