package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	webPages = []string{
		"yahoo.com",
		"google.com",
		"bing.com",
		"amazon.com",
		"github.com",
		"gitlab.com",
	}

	results struct {
		// put here content length of each page
		ContentLength map[string]int

		// accumulate here the content length of all pages
		TotalContentLength int
	}
)

func main() {
	fmt.Println("Arpit")
	fmt.Println(webPages)
	start := time.Now()
	ch := make(chan string)

	for i := 0; i < len(webPages); i++ {
		go fetch(webPages[i], ch)
	}
	for i := 0; i < len(webPages); i++ {
		fmt.Println(<-ch)
	}

	fmt.Println(results)
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())

}

func fetch(url string, ch chan<- string) {
	fmt.Println(url)

	start := time.Now()
	resp, err := http.Get("https://" + url)
	if err != nil {
		ch <- fmt.Sprint(err) // send to channel ch
		return
	}
	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close() // don't leak resources
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %7d %s", secs, nbytes, url)

}
