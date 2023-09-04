package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/quintui/rssagg/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ApiKey    string    `json:"api_key"`
}

func databaseUserToUser(databaseUser database.User) User {
	return User{
		ID:        databaseUser.ID,
		Username:  databaseUser.Username,
		CreatedAt: databaseUser.CreatedAt,
		UpdatedAt: databaseUser.UpdatedAt,
		ApiKey:    databaseUser.ApiKey,
	}
}

type Feed struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	AuthorID  uuid.UUID `json:"author_id"`
}

func databaseFeedToFeed(databaseFeed database.Feed) Feed {
	return Feed{
		ID:        databaseFeed.ID,
		Name:      databaseFeed.Name,
		Url:       databaseFeed.Url,
		CreatedAt: databaseFeed.CreatedAt,
		UpdatedAt: databaseFeed.UpdatedAt,
		AuthorID:  databaseFeed.AuthorID,
	}
}

func databaseFeedsToFeeds(databaseFeeds []database.Feed) []Feed {
	feeds := make([]Feed, len(databaseFeeds))
	for _, feed := range databaseFeeds {
		feeds = append(feeds, databaseFeedToFeed(feed))
	}
	return feeds
}

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func databaseFeedFollowToFeedFollow(databaseFeedFollow database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID:        databaseFeedFollow.ID,
		UserID:    databaseFeedFollow.UserID,
		FeedID:    databaseFeedFollow.FeedID,
		CreatedAt: databaseFeedFollow.CreatedAt,
		UpdatedAt: databaseFeedFollow.UpdatedAt,
	}
}

func databaseFeedFollowsToFeedFollows(databaseFeedFollows []database.FeedFollow) []FeedFollow {
	feedFollows := make([]FeedFollow, len(databaseFeedFollows))

	for _, feedFollow := range databaseFeedFollows {
		feedFollows = append(feedFollows, databaseFeedFollowToFeedFollow(feedFollow))
	}

	return feedFollows
}
