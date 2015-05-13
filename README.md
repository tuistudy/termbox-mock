# termbox-mock
a simple way to mock the termbox for testing. 

when i try to test my termbox based application (in golang). i need to find some way to unit test the module. after some tries and discuss with @nsf i finally find a simple way to unit test my application. it's worth to share this idea to others.   

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

based on the fist opion, I choose to wrap the termbox.PollEvent function in WaitEvent method. otherwise, the interface has to know the termbox API. such as termbox.Event and some constant. The other choise i made is replace the termbox.PollEvent with an function value. it's defalut value is termbox.Event. so that in the following test i have the opportunity to replace this function with mine. with this help, I can mock most of the situation. you will see it later.
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
