package rss

import (
	"context"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	// create client and request
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}
	// set request header
	req.Header.Set("User-Agent", "gator")
	// make the request check response status code
	resp, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}
	if resp.StatusCode > 299 {
		return &RSSFeed{}, errors.New("Response failed")
	}
	// read and store the body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RSSFeed{}, err
	}
	resp.Body.Close()
	// Unmarshal the response body into RSSFeed{}
	contentRSS := RSSFeed{}
	err = xml.Unmarshal(body, &contentRSS)
	if err != nil {
		return &RSSFeed{}, err
	}

	return &contentRSS, nil
}
