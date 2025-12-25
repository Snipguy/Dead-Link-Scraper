package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/html"
)

// const starting_site string = "https://scrape-me.dreamsofcode.io/"

const starting_site string = "https://crawler-test.com/"

func isDeadLink(url string) bool {
	transport := &http.Transport{
		ForceAttemptHTTP2: false,
	}
	client := &http.Client{
		Transport: transport,
		Timeout: 10 * time.Second,
	}

	resp, err := client.Head(url)
	if err != nil {
		// fallback to GET
		req, _ := http.NewRequest("GET", url, nil)
		resp, err = client.Do(req)
		if err != nil {
			return true // dead
		}
	}
	defer resp.Body.Close()

	return resp.StatusCode >= 400
}

func getLinks(body io.Reader) []string {
	var links []string
	z := html.NewTokenizer(body)

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return links
		case html.StartTagToken:
			t := z.Token()
			if t.Data == "a" {
				for _, attr := range t.Attr {
					if attr.Key == "href" {
						links = append(links, attr.Val)
					}
				}
			}
		}

	}
}

func main() {
	// reader := bufio.NewReader(os.Stdin)
	// starting_site, _ := reader.ReadString('\n')

	fmt.Printf("------start------\n")
	resp, err := http.Get(starting_site)
	if err != nil {
		fmt.Printf("error at fetching the starting site : %v", err)
		fmt.Printf("-------------------------\n")
	}

	defer resp.Body.Close()

	// Links := getLinks(resp.Body)
	// var internal_links []string
	// var external_links []string

	// for _, link := range Links {
	// 	if strings.Contains(link, starting_site) || strings.HasPrefix(link, "/") {
	// 		internal_links = append(internal_links, link)
	// 		for _, link := range internal_links {
	// 			if isDeadLink(link) {
	// 				fmt.Printf("%v ----> %v\n", link, "❌")
	// 				fmt.Printf("-------------gg------------\n")
	// 			} else {
	// 				fmt.Printf("%v ----> %v\n", link, "✔")
	// 				fmt.Printf("-------------------------\n")
	// 				Links = append(Links, link)
	// 			}
	// 		}
	// 	} else {
	// 		external_links = append(external_links, link)
	// 	}
	// }

	// starting_site := "https://scrape-me.dreamsofcode.io/"

	visited := make(map[string]bool)
	queue := []string{starting_site}

	baseURL, _ := url.Parse(starting_site)
	var counter int = 0

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if visited[current] {
			continue
		}
		visited[current] = true

		resp, err := http.Get(current)
		if err != nil {
			counter = counter + 1
			continue
		}

		links := getLinks(resp.Body)
		resp.Body.Close()

		for _, link := range links {
			linkURL, err := url.Parse(link)
			if err != nil {
				counter = counter + 1
				continue
			}

			absURL := baseURL.ResolveReference(linkURL)

			if absURL.Scheme != "http" && absURL.Scheme != "https" {
				counter = counter + 1
				continue
			}

			if isDeadLink(absURL.String()) {
				fmt.Println(absURL.String(), "❌")
				fmt.Printf("-------------------------\n")
				counter = counter + 1
			} else {
				fmt.Println(absURL.String(), "✔")
				fmt.Printf("-------------------------\n")
				counter = counter + 1
			}

			if absURL.Host == baseURL.Host && !visited[absURL.String()] {
				queue = append(queue, absURL.String())
			}
		}
	}

	// for _, link := range external_links {
	// 	status := isDeadLink(link)

	// 	if status {
	// 		fmt.Printf("%v ----> %v\n", link, "✔")
	// 		fmt.Printf("-------------------------\n")
	// 	} else {
	// 		fmt.Printf("%v ----> %v\n", link, "❌")
	// 		fmt.Printf("-------------------------\n")
	// 	}
	// }

	fmt.Printf("all links = %v", counter)
	fmt.Printf("------end------")
}
