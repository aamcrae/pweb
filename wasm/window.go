package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"syscall/js"
)

// Minimum movement to consider a touch event to be a swipe.
const minSwipe = 30
// Maximum time before a touch event stops being a swipe.
const maxSwipeTime = time.Millisecond * 300

// Direction is an enum indicating the swipe direction.
type Direction int

const (
	Right Direction = iota
	Left
	Up
	Down
)

// Window is the main structure for interfacing to the browser
type Window struct {
	window, document, head, body js.Value
	Width, Height                int
	// touch values
	startTime time.Time
	startX, startY int
	endX, endY     int
	multiTouch     bool
}

// GetWindow creates a new Window ready to interface to the browser.
func GetWindow() *Window {
	w := &Window{}
	w.window = js.Global()
	w.document = w.window.Get("document")
	w.head = w.document.Get("head")
	w.body = w.document.Get("body")
	w.refreshSize()
	return w
}

// GetById returns the named JS element.
func (w *Window) GetById(id string) js.Value {
	return w.document.Call("getElementById", id)
}

// Display sets the HTML onto the window.
func (w *Window) Display(s string) {
	w.body.Set("innerHTML", s)
}

// LoadStyle adds a link to read the CSS file indicated
func (w *Window) LoadStyle(s string) {
	link := w.document.Call("createElement", "link")
	link.Set("type", "text/css")
	link.Set("rel", "stylesheet")
	link.Set("href", s)
	w.head.Call("appendChild", link)
}

// AddStyle adds the CSS string to the window directly.
func (w *Window) AddStyle(s string) {
	style := w.document.Call("createElement", "style")
	style.Set("type", "text/css")
	ss := style.Get("styleSheet")
	if ss.Truthy() {
		ss.Set("cssText", s)
	} else {
		style.Call("appendChild", w.document.Call("createTextNode", s))
	}
	w.head.Call("appendChild", style)
}

// SetTitle sets the window title
func (w *Window) SetTitle(title string) *Window {
	w.document.Set("title", title)
	return w
}

// Goto navigates the browser to the URL.
func (w *Window) Goto(url string) {
	w.window.Get("location").Set("href", js.ValueOf(url))
}

// OnSwipe registers a callback to be called for swipe events.
// If the callback handles the event, it returns true.
func (w *Window) OnSwipe(f func(Direction) bool) {
	touchStartJS := js.FuncOf(func(this js.Value, args []js.Value) any {
		t := args[0].Get("touches")
		if t.IsUndefined() {
			fmt.Printf("cannot find touches\n")
			return nil
		}
		w.startTime = time.Now()
		w.startX = t.Index(0).Get("clientX").Int()
		w.startY = t.Index(0).Get("clientY").Int()
		w.endX = w.startX
		w.endY = w.startY
		// Ignore multi-touch gestures
		w.multiTouch = t.Length() > 1
		return nil
	})
	touchMoveJS := js.FuncOf(func(this js.Value, args []js.Value) any {
		e := args[0].Get("targetTouches")
		if e.IsUndefined() {
			fmt.Printf("targetTouches is undefined\n")
			return nil
		}
		if e.Length() == 1 {
			w.endX = e.Index(0).Get("clientX").Int()
			w.endY = e.Index(0).Get("clientY").Int()
		} else {
			w.multiTouch = true
		}
		return nil
	})
	touchEndJS := js.FuncOf(func(this js.Value, args []js.Value) any {
		if w.multiTouch {
			return nil
		}
		// If the swipe event lasted more than the preset max, do not
		// consider this a swipe event.
		if time.Now().Sub(w.startTime) > maxSwipeTime {
			return nil
		}
		e := args[0]
		x := w.startX - w.endX
		y := w.startY - w.endY
		var d Direction
		ax := abs(x)
		ay := abs(y)
		// Figure out up/down or left/right
		if ax > ay {
			if ax < minSwipe {
				return nil
			}
			if x > 0 {
				d = Left
			} else {
				d = Right
			}
		} else {
			if ay < minSwipe {
				return nil
			}
			if y > 0 {
				d = Up
			} else {
				d = Down
			}
		}
		if f(d) {
			// Don't process the default action if the
			// swipe was handled by the callback
			e.Call("preventDefault")
		}
		return nil
	})
	touchCancelJS := js.FuncOf(func(this js.Value, args []js.Value) any {
		return nil
	})
	w.window.Call("addEventListener", "touchstart", touchStartJS)
	w.window.Call("addEventListener", "touchend", touchEndJS)
	w.window.Call("addEventListener", "touchmove", touchMoveJS)
	w.window.Call("addEventListener", "touchcancel", touchCancelJS)
}

// OnResize registers a callback to be invoked when the window changes size.
func (w *Window) OnResize(f func()) {
	resizeJS := js.FuncOf(func(this js.Value, args []js.Value) any {
		w.refreshSize()
		f()
		return nil
	})
	w.window.Call("addEventListener", "resize", resizeJS)
}

// OnKey registers a callback to be invoked when a key is pressed.
func (w *Window) OnKey(f func(key string)) {
	onKey := js.FuncOf(func(this js.Value, args []js.Value) any {
		f(js.ValueOf(args[0].Get("key")).String())
		return nil
	})
	w.document.Call("addEventListener", "keydown", onKey)
}

// AddJSFunction registers a javascript function that invokes the Go function passed.
func (w *Window) AddJSFunction(name string, f func(js.Value, []js.Value) any) {
	w.window.Set(name, js.FuncOf(f))
}

// refreshSize updates the width/height of the window.
func (w *Window) refreshSize() {
	w.Width = js.ValueOf(w.body.Get("clientWidth")).Int()
	w.Height = js.ValueOf(w.window.Get("innerHeight")).Int()
}

// Wait forces this thread to wait.
func (w *Window) Wait() {
	var c chan struct{}
	<-c
}

// GetContent retrieves a file from the server
func GetContent(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Read body: %v", err)
	}

	return data, nil
}

// abs returns the absolute value of a value
func abs(a int) int {
	if a < 0 {
		return -a
	} else {
		return a
	}
}
