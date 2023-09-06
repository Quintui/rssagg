package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/quintui/rssagg/internal/database"
)

func startScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Starting scraping with concurrency %d and time between request %s", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))

		if err != nil {
			log.Println("Could not get next feed to fetch")
			continue
		}

		log.Printf("Found %d feeds to fetch", len(feeds))

		wg := &sync.WaitGroup{}

		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, wg, feed)
		}

		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)

	if err != nil {
		log.Printf("Could not mark feed %s as fetched: %s", feed.Name, err)
		return
	}

	// Implement data parsing

	feedData, err := fetchFeed(feed.Url)

	for _, item := range feedData.Channel.Item {
		publishedAt := sql.NullTime{}
		if pubAt, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  pubAt,
				Valid: true,
			}
		}
		_, err := db.CreatePost(context.Background(), database.CreatePostParams{
			ID:    uuid.New(),
			Title: item.Title,
			Url:   item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})

		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Could not create post: %s", err)
			continue
		}
	}

	if err != nil {
		log.Printf("Could not fetch feed %s: %s", feed.Name, err)
	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))

}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(feedUrl string) (RSSFeed, error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := httpClient.Get(feedUrl)

	if err != nil {
		return RSSFeed{}, err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)

	if err != nil {
		return RSSFeed{}, err
	}

	var rssFeed RSSFeed

	err = xml.Unmarshal(dat, &rssFeed)

	if err != nil {
		return RSSFeed{}, err

	}
	return rssFeed, nil

}
