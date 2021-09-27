package main

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/ozouai/loftymusic/audiosource/audiosourcepb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func TestSearch(t *testing.T) {
	service, err := youtube.NewService(context.TODO(), option.WithAPIKey(os.Getenv("YOUTUBE_API_KEY")))
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}
	search := youtube.NewSearchService(service)
	srv := &YoutubeServer{
		youtube: service,
		search:  search,
	}
	res, err := srv.Search(context.TODO(), &audiosourcepb.Search_Request{Term: "Go Proverbs - Rob Pike - Gopherfest - November 18, 2015"})
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "PAAkCSZUG1c", res.GetDetails().GetId())

	res, err = srv.Search(context.TODO(), &audiosourcepb.Search_Request{Term: "https://www.youtube.com/watch?v=PAAkCSZUG1c"})
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "PAAkCSZUG1c", res.GetDetails().GetId())
}
