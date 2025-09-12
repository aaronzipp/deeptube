package youtube

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aaronzipp/deeptube/video"

	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"gopkg.in/yaml.v3"
)

type Subscription struct {
	Channel         string   `yaml:"channel"`
	ID              string   `yaml:"id"`
	Categories      []string `yaml:"categories"`
	Live            bool     `yaml:"live,omitempty"`
	ExcludeKeywords []string `yaml:"exclude_keywords,omitempty"`
	Shorts          bool     `yaml:"shorts,omitempty"`
}

type Playlist struct {
	Playlist   string   `yaml:"playlist"`
	ID         string   `yaml:"id"`
	Categories []string `yaml:"categories"`
}

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
		thumbnail := ""
		if item.Snippet.Thumbnails.Standard != nil {
			thumbnail = item.Snippet.Thumbnails.Standard.Url
		} else if item.Snippet.Thumbnails.High != nil {
			thumbnail = item.Snippet.Thumbnails.High.Url
		} else if item.Snippet.Thumbnails.Medium != nil {
			thumbnail = item.Snippet.Thumbnails.Medium.Url
		} else if item.Snippet.Thumbnails.Default != nil {
			thumbnail = item.Snippet.Thumbnails.Default.Url
		}

		videos[i] = video.Video{
			ChannelName: item.Snippet.ChannelTitle,
			Title:       item.Snippet.Title,
			VideoId:     item.Id,
			Thumbnail:   thumbnail,
			Description: item.Snippet.Description,
			PublishedAt: publishedAt,
			VideoLength: length,
		}
	}

	return videos, nil
}

func shouldExcludeVideo(title string, excludeKeywords []string) bool {
	if len(excludeKeywords) == 0 {
		return false
	}

	titleLower := strings.ToLower(title)
	for _, keyword := range excludeKeywords {
		if strings.Contains(titleLower, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

func FetchAllVideos(subscriptions []Subscription, playlists []Playlist) (video.Videos, error) {
	playlistExcludes := make(map[string][]string)

	playlistIds := []string{}
	for _, subscription := range subscriptions {
		playlistId := strings.Replace(subscription.ID, "UC", string(video.NormalVideo), 1)
		playlistIds = append(playlistIds, playlistId)
		playlistExcludes[playlistId] = subscription.ExcludeKeywords

		if subscription.Live {
			playlistId = strings.Replace(subscription.ID, "UC", string(video.LiveVideo), 1)
			playlistIds = append(playlistIds, playlistId)
			playlistExcludes[playlistId] = subscription.ExcludeKeywords
		}
		if subscription.Shorts {
			playlistId = strings.Replace(subscription.ID, "UC", string(video.ShortVideo), 1)
			playlistIds = append(playlistIds, playlistId)
			playlistExcludes[playlistId] = subscription.ExcludeKeywords
		}
	}
	for _, playlist := range playlists {
		playlistIds = append(playlistIds, playlist.ID)
		// Playlists don't have exclude_keywords, so they get an empty array
		playlistExcludes[playlist.ID] = []string{}
	}

	vids := video.Videos{}
	for _, playlistId := range playlistIds {
		videoIds, err := FetchVideoIdsFromPlaylist(playlistId)
		if err != nil {
			return nil, fmt.Errorf(
				"failed fetching video ids from playlist %q: %!s",
				playlistId,
				err,
			)
		}
		playlistVids, err := FetchVideos(videoIds)
		if len(playlistVids) == 0 {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed fetching videos with ids %+v: %!s", videoIds, err)
		}

		excludes := playlistExcludes[playlistId]
		filteredVids := make(video.Videos, 0, len(playlistVids))
		for _, vid := range playlistVids {
			if !shouldExcludeVideo(vid.Title, excludes) {
				filteredVids = append(filteredVids, vid)
			}
		}

		vids = append(vids, filteredVids...)
	}
	return vids, nil
}

func ParseSubscriptions(filename string) ([]Subscription, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var subs []Subscription
	err = yaml.Unmarshal(data, &subs)
	if err != nil {
		return nil, err
	}

	return subs, nil
}

func ParsePlaylists(filename string) ([]Playlist, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var playlists []Playlist
	err = yaml.Unmarshal(data, &playlists)
	if err != nil {
		return nil, err
	}

	return playlists, nil
}

func RefreshVideos() error {

	subscriptions, err := ParseSubscriptions("subscriptions.yaml")
	if err != nil {
		return err
	}
	playlists, err := ParsePlaylists("playlists.yaml")
	if err != nil {
		return err
	}
	videos, err := FetchAllVideos(subscriptions, playlists)

	if err != nil {
		return err
	}

	videos.Sort()
	videos.WriteToDB()

	return nil
}
