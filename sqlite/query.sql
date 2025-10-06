-- name: AddVideo :exec
INSERT INTO videos (video_id, title, thumbnail_url, channel_name, description, published_at, hours, minutes, seconds, was_live, is_hidden)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0)
    ON CONFLICT(video_id) DO UPDATE SET
	title = excluded.title,
	thumbnail_url = excluded.thumbnail_url,
	channel_name = excluded.channel_name,
	description = excluded.description,
	published_at = excluded.published_at,
	hours = excluded.hours,
	minutes = excluded.minutes,
	seconds = excluded.seconds,
	was_live = excluded.was_live;

-- name: FetchVideos :many
select *
from videos
where is_hidden = 0
;

-- name: HideVideo :exec
update videos
set is_hidden = 1
where video_id = ?;

-- AddThumbnail :exec
INSERT INTO thumbnails(video_id, image, updated_at)
VALUES (?, ?, CURRENT_TIMESTAMP)
ON CONFLICT(video_id) DO UPDATE SET
    image = excluded.image,
    updated_at = CURRENT_TIMESTAMP;
