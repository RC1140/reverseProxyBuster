package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
)

var streamURL = ""

func main() {
	u := getStreamURL(streamURL)
	log.Printf(u)
	u2, err := url.Parse(u)
	if err != nil {
		panic(err)
	}

	log.Printf(u2.Host)
	log.Printf(fmt.Sprintf("http://localhost:8080%s", u2.Path))
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = "https"
			req.URL.Host = u2.Host
			req.Header.Set("Referer", streamURL)
		},
	}
	log.Fatal(http.ListenAndServe(":8080", proxy))
}

func getStreamURL(scrapeLocation string) string {
	resp, _ := http.Get(scrapeLocation)

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	b := string(body)
	lines := strings.Split(b, "\n")
	for _, l := range lines {
		if strings.Contains(l, "source:") {
			var validID = regexp.MustCompile(`".*"`)
			src := validID.FindString(l)
			end := len(src) - 1
			//if we have a valid start use it as is otherwise add it
			if strings.HasPrefix(src[1:], "https://") {
				return src[1 : len(src)-1]
			}
			return fmt.Sprintf("https://%s", src[3:end])
		}
	}
	return ""
}
