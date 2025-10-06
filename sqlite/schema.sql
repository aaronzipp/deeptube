CREATE TABLE videos (
	video_id TEXT PRIMARY KEY,
	title TEXT,
	thumbnail_url TEXT,
	channel_name TEXT,
	description TEXT,
	published_at TEXT,
	hours INTEGER,
	minutes INTEGER,
	seconds INTEGER,
	was_live INTEGER,
	is_hidden INTEGER
);

CREATE TABLE thumbnails (
	video_id TEXT PRIMARY KEY,
	thumbnail BLOB,
	updated_at TEXT,
	FOREIGN KEY(video_id) REFERENCES videos(video_id) ON DELETE CASCADE
);
