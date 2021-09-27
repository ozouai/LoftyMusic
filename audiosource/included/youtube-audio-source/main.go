package main

import (
	"context"
	"fmt"
	"log"
	"os"

	youtubev2 "github.com/kkdai/youtube/v2"
	"github.com/ozouai/loftymusic/audiosource/audiosourceclient"
	"github.com/ozouai/loftymusic/audiosource/audiosourcepb"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"gopkg.zouai.io/colossus/clog"
)

type YoutubeServer struct {
	youtube *youtube.Service
	search  *youtube.SearchService
}

func main() {
	ctx := context.Background()
	ctx, logger := clog.NewRootLogger(ctx, "YoutubeAudioSource")
	_ = logger
	service, err := youtube.NewService(context.TODO(), option.WithAPIKey(os.Getenv("YOUTUBE_API_KEY")))
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}
	search := youtube.NewSearchService(service)
	server := &YoutubeServer{
		youtube: service,
		search:  search,
	}
	client, err := audiosourceclient.New(context.Background(), "127.0.0.1:8080", "youtube", server)
	if err != nil {
		panic(err)
	}
	client.Wait()
}

func IsYoutubeID(url string) (string, bool) {
	res, err := youtubev2.ExtractVideoID(url)
	if err != nil {
		return res, true
	}
	return url, false
}

func (m *YoutubeServer) Search(ctx context.Context, req *audiosourcepb.Search_Request) (*audiosourcepb.Search_Response, error) {

	if id, ok := IsYoutubeID(req.GetTerm()); ok {
		res, err := m.youtube.Videos.List([]string{"id,snippet"}).Id(id).Do()
		if err != nil {
			return nil, fmt.Errorf("error searching youtube by id: %w", err)
		}
		if len(res.Items) == 0 {
			return &audiosourcepb.Search_Response{Result: &audiosourcepb.Search_Response_NotFound{NotFound: "No results found"}}, nil
		}
		target := res.Items[0]
		return &audiosourcepb.Search_Response{
			Result: &audiosourcepb.Search_Response_Details{
				Details: &audiosourcepb.AudioMetadata{
					Id:   string(target.Id),
					Name: target.Snippet.Title,
				},
			},
		}, nil
	}

	results, err := m.search.List([]string{"id,snippet"}).Q(req.GetTerm()).Type("video").MaxResults(1).Do()
	if err != nil {
		return nil, fmt.Errorf("error searching youtube: %w", err)
	}
	if len(results.Items) == 0 {
		return &audiosourcepb.Search_Response{Result: &audiosourcepb.Search_Response_NotFound{NotFound: "No results found"}}, nil
	}
	target := results.Items[0]
	return &audiosourcepb.Search_Response{
		Result: &audiosourcepb.Search_Response_Details{
			Details: &audiosourcepb.AudioMetadata{
				Id:   string(target.Id.VideoId),
				Name: target.Snippet.Title,
			},
		},
	}, nil
}
