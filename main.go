package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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
		// "bitnz.com",
		// "random-website-name.com",
	}

	results struct {
		// put here content length of each page
		ContentLength map[string]int

		// accumulate here the content length of all pages
		TotalContentLength int

		// extra information for error
	}
)

type ItemResult struct {
	Url           string
	ContentLength int
	Err           string
}

const time_out = 30

func main() {
	// initialize the output structure
	results.ContentLength = make(map[string]int)
	results.TotalContentLength = 0

	// http client
	client := new(http.Client)
	ch := make(chan *ItemResult)

	// fetch the contents in parallel
	for i := 0; i < len(webPages); i++ {
		go worker(client, webPages[i], ch)
	}

	timer := time.AfterFunc(time.Second*time_out, func() {
		fmt.Println("Out of time!Could fetch partial pages only...\n")
		close(ch)
		print_results()
		os.Exit(0)
	})
	defer timer.Stop()

	for i := 0; i < len(webPages); i++ {
		output, ok := <-ch
		if !ok {
			fmt.Printf("channel is closed\n")
			return
		} else {
			// fmt.Printf("%s,\t%d bytes\n", output.Url, output.ContentLength)
			results.ContentLength[output.Url] = output.ContentLength
			results.TotalContentLength = results.TotalContentLength + output.ContentLength
		}
	}

	print_results()

}

func worker(client *http.Client, url string, ch chan *ItemResult) {
	result := ItemResult{url, -1, ""}

	// get the url
	response, err := client.Get("https://" + url)
	if err != nil {
		// fmt.Println(err, url)
		result.Err = err.Error()
		ch <- &result
		return
	}

	// check the http server response was OK
	if response.Status != "200 OK" {
		// fmt.Println(response.Status, url)
		result.Err = err.Error()
		ch <- &result
		return
	}

	// read the response body
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		// fmt.Println(err, url)
		result.Err = err.Error()
		ch <- &result
		return
	}

	// pass the result back to caller via a channel
	result.ContentLength = len(string(b))
	ch <- &result
	return
}

func print_results() {
	for key, element := range results.ContentLength {
		fmt.Println(key, "-", element)
	}
	fmt.Println("\nTotalContentLength =", results.TotalContentLength, "\n")
	fmt.Println(results)
	fmt.Println("\n")

}
