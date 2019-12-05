package main

import (
	"log"
	"sync"
)

var (
	portMapping = map[string]int{
		"hoge": 18080,
		"foo":  18081,
		"fuga": 18082,
		"bar":  18083,
		"baz":  18084,
	}
)

func main() {
	var wg sync.WaitGroup
	serve(&wg, (&Hoge{}).Serve)
	serve(&wg, (&Fuga{}).Serve)
	serve(&wg, (&Bar{}).Serve)
	serve(&wg, (&Baz{}).Serve)
	serve(&wg, (&Foo{}).Serve)
	wg.Wait()
}

func serve(wg *sync.WaitGroup, f func() error) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := f(); err != nil {
			log.Fatal(err)
		}
	}()
}
