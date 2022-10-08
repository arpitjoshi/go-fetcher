package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	webPages = []string{
		"yahoo.com",
		"google.com",
		"bing.com",
		"amazon.com",
		"github.com",
		"gitlab.com",
		"bitnz.com",
	}

	results struct {
		// put here content length of each page
		ContentLength map[string]int

		// accumulate here the content length of all pages
		TotalContentLength int
	}
)

type ItemResult struct {
	Url           string
	ContentLength int
}

func main() {
	// initialize the output structure
	results.ContentLength = make(map[string]int)
	results.TotalContentLength = 0

	// create an http client, a list of urls and a channel to
	// pass url fetch results back
	client := new(http.Client)
	ch := make(chan *ItemResult)

	// fetch the contents of the list of urls in parallel
	for i := 0; i < len(webPages); i++ {
		go worker(client, webPages[i], ch)
	}

	// store the results as soon as they come in
	for i := 0; i < len(webPages); i++ {
		output := <-ch
		// fmt.Printf("%s,\t%d bytes\n", output.Url, output.ContentLength)
		results.ContentLength[output.Url] = output.ContentLength
		results.TotalContentLength = results.TotalContentLength + output.ContentLength
	}

	for key, element := range results.ContentLength {
		fmt.Println(key, "-", element)
	}
	fmt.Println("\nTotalContentLength =", results.TotalContentLength)

}

func worker(client *http.Client, url string, ch chan *ItemResult) {
	result := ItemResult{url, -1}

	// get the url
	response, err := client.Get("https://" + url)
	if err != nil {
		// fmt.Println(err, url)
		ch <- &result
		return
	}

	// check the http server response was OK
	if response.Status != "200 OK" {
		// fmt.Println(response.Status, url)
		ch <- &result
		return
	}

	// read the response body
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		// fmt.Println(err, url)
		ch <- &result
		return
	}

	// pass the result back to caller via a channel
	result.ContentLength = len(string(b))
	ch <- &result
}
