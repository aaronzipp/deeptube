-- name: FetchVideos :many
select *
from videos
where is_hidden = 0
;

-- name: HideVideo :exec
update videos
set is_hidden = 1
where video_id = ?;

-- name: AddThumbnail :exec
INSERT INTO thumbnails(video_id, thumbnail, updated_at)
VALUES (?, ?, CURRENT_TIMESTAMP)
ON CONFLICT(video_id) DO UPDATE SET
    thumbnail = excluded.thumbnail,
    updated_at = CURRENT_TIMESTAMP

-- name: FetchThumbnail :one
SELECT thumbnail FROM thumbnails WHERE video_id = ?;
