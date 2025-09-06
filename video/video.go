package video

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"time"

	"github.com/aaronzipp/deeptube/database"

	_ "modernc.org/sqlite"
)

type VideoType string

// Those types are not officially documented
// see https://stackoverflow.com/a/76602819/8313407
const (
	NormalVideo VideoType = "UULF"
	ShortVideo  VideoType = "UUSH"
	LiveVideo   VideoType = "UULV"
)

const youtubeLinkTemplate = "https://www.youtube.com/watch_popup?v=%s"

const hoursInDay = 24
const hoursInMonth = hoursInDay * 30
const hoursInYear = hoursInDay * 365

type Video struct {
	Title       string
	VideoId     string
	ChannelName string
	Description string
	PublishedAt time.Time
	VideoLength Length
	Thumbnail   string
	WasLive     bool
}
type Videos []Video

func (v Video) YouTubeLink() string {
	return fmt.Sprintf(youtubeLinkTemplate, v.VideoId)
}

func (v Video) TimeSincePublished() string {
	timeDifference := time.Since(v.PublishedAt)

	if timeDifference.Hours() >= hoursInYear {
		years := int(timeDifference.Hours() / hoursInYear)
		if years == 1 {
			return "1 year ago"
		}
		return fmt.Sprintf("%d years ago", years)
	}
	if timeDifference.Hours() >= hoursInMonth {
		months := int(timeDifference.Hours() / hoursInMonth)
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	}
	if timeDifference.Hours() >= hoursInDay {
		days := int(timeDifference.Hours() / hoursInDay)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}
	if timeDifference.Hours() >= 1 {
		hours := int(timeDifference.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	}
	if timeDifference.Minutes() >= 1 {
		minutes := int(timeDifference.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	}
	if timeDifference.Seconds() >= 1 {
		seconds := int(timeDifference.Seconds())
		if seconds == 1 {
			return "1 second ago"
		}
		return fmt.Sprintf("%d seconds ago", seconds)
	}

	return "now"

}

func (v Video) String() string {
	return fmt.Sprintf(
		" --- Id: %s ---\nTitle: %s\nChannel: %s\nPublished At: %s\nLength: %s",
		v.VideoId,
		v.Title,
		v.ChannelName,
		v.PublishedAt,
		v.VideoLength.String(),
	)
}

func VideosFromDB() (Videos, error) {
	ctx := context.Background()
	db, err := sql.Open("sqlite", "videos.db")

	if err != nil {
		return nil, err
	}

	queries := database.New(db)

	dbVideos, err := queries.FetchVideos(ctx)

	if err != nil {
		return nil, err
	}

	vids := make(Videos, len(dbVideos))
	for i, vid := range dbVideos {
		publishedTime, _ := time.Parse("2006-01-02 15:04:05", vid.PublishedAt.String)
		vids[i] = Video{
			Title:       vid.Title.String,
			VideoId:     vid.VideoID,
			ChannelName: vid.ChannelName.String,
			Description: vid.Description.String,
			VideoLength: Length{
				Hours:   int(vid.Hours.Int64),
				Minutes: int(vid.Minutes.Int64),
				Seconds: int(vid.Seconds.Int64),
			},
			PublishedAt: publishedTime,
			WasLive:     vid.WasLive.Int64 == 1,
		}
	}

	return vids, nil
}

func (v Videos) Sort() {
	sort.Slice(v, func(i, j int) bool {
		return v[i].PublishedAt.After(v[j].PublishedAt)
	})
}
