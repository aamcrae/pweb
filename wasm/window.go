package wasm

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"syscall/js"
)

type Window struct {
	document, head, body js.Value
}

func GetWindow() *Window {
	w := &Window{}
	w.document = js.Global().Get("document")
	w.head = w.document.Get("head")
	w.body = w.document.Get("body")
	return w
}

func (w *Window) GetById(id string) js.Value {
	return w.document.Call("getElementById", id)
}

func (w *Window) Display(p *Page) {
	w.body.Set("innerHTML", p.String())
}

func (w *Window) LoadStyle(s string) {
	link := w.document.Call("createElement", "link")
	link.Set("type", "text/css")
	link.Set("rel", "stylesheet")
	link.Set("href", "/pweb/"+s)
	w.head.Call("appendChild", link)
}

func (w *Window) SetTitle(title string) *Window {
	w.document.Set("title", title)
	return w
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
