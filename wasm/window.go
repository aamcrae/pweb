package wasm

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"syscall/js"
)

type Window struct {
	window, document, head, body js.Value
	resizeJS                     js.Func
	Width, Height                int
}

func GetWindow() *Window {
	w := &Window{}
	w.window = js.Global()
	w.document = w.window.Get("document")
	w.head = w.document.Get("head")
	w.body = w.document.Get("body")
	w.refreshSize()
	return w
}

func (w *Window) GetById(id string) js.Value {
	return w.document.Call("getElementById", id)
}

func (w *Window) Display(s string) {
	w.body.Set("innerHTML", s)
}

func (w *Window) LoadStyle(s string) {
	link := w.document.Call("createElement", "link")
	link.Set("type", "text/css")
	link.Set("rel", "stylesheet")
	link.Set("href", s)
	w.head.Call("appendChild", link)
}

func (w *Window) AddStyle(s string) {
	style := w.document.Call("createElement", "style")
	textNode := w.document.Call("createTextNode", s)
	style.Set("type", "text/css")
	w.head.Call("appendChild", textNode)
}

func (w *Window) SetTitle(title string) *Window {
	w.document.Set("title", title)
	return w
}

func (w *Window) OnResize(f func()) {
	w.resizeJS = js.FuncOf(func(this js.Value, args []js.Value) any {
		w.refreshSize()
		f()
		return nil
	})
	w.window.Call("addEventListener", "resize", w.resizeJS)
}

func (w *Window) AddJSFunction(name string, f func(js.Value, []js.Value) any) {
	w.window.Set(name, js.FuncOf(f))
}

func (w *Window) refreshSize() {
	w.Width = js.ValueOf(w.window.Get("innerWidth")).Int()
	w.Height = js.ValueOf(w.window.Get("innerHeight")).Int()
}

func (w *Window) Wait() {
	var c chan struct{}
	<-c
}

// Retrieve file from server
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

func (w *Window) GetIntOrDefault(id string, def int) int {
	return def
}
