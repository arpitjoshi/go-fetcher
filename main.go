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
	// create an http client, a list of urls and a channel to
	// pass url fetch results back
	client := new(http.Client)
	ch := make(chan *ItemResult)

	results.ContentLength = make(map[string]int)

	// fetch the contents of the list of urls in parallel
	for i := 0; i < len(webPages); i++ {
		go fetch(client, webPages[i], ch)
	}

	// print out the results as soon as they come in
	//  (here we could write the data to disk if required)
	for i := 0; i < len(webPages); i++ {
		output := <-ch
		fmt.Printf("%s,\t%d bytes\n", output.Url, output.ContentLength)
		results.ContentLength[output.Url] = output.ContentLength
		results.TotalContentLength = results.TotalContentLength + output.ContentLength
	}

	fmt.Print(results)

}

func fetch(client *http.Client, url string, ch chan *ItemResult) {
	result := ItemResult{url, -1}

	// get the url
	response, err := client.Get("https://" + url)
	if err != nil {
		fmt.Println(err, url)
		ch <- &result
		return
	}

	// check the http server response was OK
	if response.Status != "200 OK" {
		fmt.Println(response.Status, url)
		ch <- &result
		return
	}

	// read the response body
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err, url)
		ch <- &result
		return
	}

	// pass the result back to caller via a channel
	result.ContentLength = len(string(b))
	ch <- &result
}
