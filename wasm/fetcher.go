package main

import (
	"fmt"
	"sync"

	"syscall/js"
)

// Fetcher is an interface to the JS fetch API
type Fetcher struct {
	w *Window
	wg sync.WaitGroup
	data []byte
	err error
	errF js.Func
	respF js.Func
	readyF js.Func
}

func NewFetcher(w *Window, url string) * Fetcher{
	f := &Fetcher{w: w}
	fPromise := w.window.Call("fetch", js.ValueOf(url))
	f.wg.Add(1)
	f.respF = js.FuncOf(f.response)
	f.errF = js.FuncOf(f.reject)
	f.readyF = js.FuncOf(f.ready)
	fPromise.Call("then", f.respF, f.errF)
	return f
}

func (f *Fetcher) Get() ([]byte, error) {
	f.wg.Wait()
	f.errF.Release()
	f.respF.Release()
	f.readyF.Release()
	return f.data, f.err
}

func (f *Fetcher) reject(this js.Value, args []js.Value) any {
	f.err = fmt.Errorf("Rejected")
	f.wg.Done()
	return nil
}

func (f *Fetcher) response(this js.Value, args []js.Value) any {
	ok := args[0].Get("ok")
	if js.ValueOf(ok).Bool() == true {
		p2 := args[0].Call("text")
		p2.Call("then", f.readyF, f.errF)
		return p2
	}
	f.err = fmt.Errorf("Resp: %s", js.ValueOf(args[0].Get("status")).String())
	f.wg.Done()
	return nil
}

func (f *Fetcher) ready(this js.Value, args []js.Value) any {
	f.data = []byte(js.ValueOf(args[0]).String())
	f.wg.Done()
	return nil
}
