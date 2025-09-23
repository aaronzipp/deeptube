-- name: AddVideo :exec
INSERT INTO videos (video_id, title, thumbnail, channel_name, description, published_at, hours, minutes, seconds, was_live) 
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    ON CONFLICT(video_id) DO UPDATE SET
	title = excluded.title,
	thumbnail = excluded.thumbnail,
	channel_name = excluded.channel_name,
	description = excluded.description,
	published_at = excluded.published_at,
	hours = excluded.hours,
	minutes = excluded.minutes,
	seconds = excluded.seconds,
	was_live = excluded.was_live,
	is_hidden = 0;

-- name: FetchVideos :many
select *
from videos
where is_hidden = 0
;

-- name: HideVideo :exec
update videos
set is_hidden = 1
where video_id = ?;
