package youtube

import (
	"context"
	"os"
	"time"

	"github.com/aaronzipp/deeptube/video"

	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func YoutubeService() (*youtube.Service, error) {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return youtube.NewService(
		ctx,
		option.WithAPIKey(os.Getenv("YOUTUBE_API_KEY")),
	)
}

func FetchVideoIdsFromPlaylist(playlistId string) ([]string, error) {
	youtubeService, err := YoutubeService()
	if err != nil {
		return nil, err
	}
	call := youtubeService.PlaylistItems.List([]string{"contentDetails"})
	result, err := call.PlaylistId(playlistId).MaxResults(10).Do()

	if err != nil {
		return nil, err
	}

	output := result.Items

	ids := make([]string, len(output))
	for i, item := range output {
		ids[i] = item.ContentDetails.VideoId
	}

	return ids, nil
}

func FetchVideos(ids []string) (video.Videos, error) {
	youtubeService, err := YoutubeService()
	if err != nil {
		return nil, err
	}

	result, err := youtubeService.Videos.List(
		[]string{"contentDetails", "snippet"},
	).Id(ids...).Do()
	if err != nil {
		return nil, err
	}

	output := result.Items

	videos := make(video.Videos, len(output))

	for i, item := range output {
		length, err := video.LengthFromString(item.ContentDetails.Duration)
		if err != nil {
			return nil, err
		}
		publishedAt, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
		if err != nil {
			return nil, err
		}
		videos[i] = video.Video{
			ChannelName: item.Snippet.ChannelTitle,
			Title:       item.Snippet.Title,
			VideoId:     item.Id,
			Thumbnail:   item.Snippet.Thumbnails.Standard.Url,
			Description: item.Snippet.Description,
			PublishedAt: publishedAt,
			VideoLength: length,
		}
	}

	return videos, nil
}
