# DeepTube

DeepTube is a desktop application for managing YouTube subscriptions and playlists
without having the drawbacks of [algorithmically](https://calnewport.com/on-disruption-and-distraction)
[curated](https://calnewport.com/back-to-the-internet-future)
[content](https://calnewport.com/tiktoks-poison-pill/).<br>
It fetches recent videos from specified subscriptions and playlists,
stores them in a local SQLite database, and displays them in a YT-like subscription box.
Users can watch videos directly or hide them to declutter the view.
The app runs in the system tray and refreshes videos automatically every 30 minutes.

## Configuration

1. Create `subscriptions.yaml` and `playlists.yaml` files in the project root with your YouTube subscriptions and playlists.

   Example `subscriptions.yaml`:
   ```yaml
   - channel: "Channel Name"
     id: "UCxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
     categories: ["Tech", "News"]
     live: true
     exclude_keywords: ["sponsored", "ad"]
     shorts: false
   ```

   Example `playlists.yaml`:
   ```yaml
   - playlist: "Playlist Name"
     id: "PLxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
     categories: ["Music"]
   ```

2. To obtain channel/playlist IDs:
   - For channels:
      - Visit the channel's YouTube page and click onto "...more" ![click onto "...more" on the channel page](https://github.com/aaronzipp/deeptube/blob/main/assets/channel_main_page.png?raw=true)
      - Scroll down and click on "Share channel" ![the bottom of the channel info](https://github.com/aaronzipp/deeptube/blob/main/assets/channel_description.png?raw=true)
      - Click on "Copy channel ID"<br> ![the channel ID](https://github.com/aaronzipp/deeptube/blob/main/assets/channel_id.png?raw=true)
   - For playlists: Open the playlist, copy the ID from the URL (e.g., `youtube.com/playlist?list=PL...`).

3. Set up Google YouTube Data API credentials:
   - Visit [YouTube Data API](https://developers.google.com/youtube/v3/getting-started) to enable the API and obtain an API key.
   - Create a `.env` file with `YOUTUBE_API_KEY=your_api_key_here`.
   - This makes it possible to get the video information.

4. Create a `videos.db` file and create the [sqlite](https://sqlite.org/index.html) tables defined in `sqlite/schema.sql`

## Building the Executable

Ensure you have [GO](https://go.dev/dl/) installed with at least version 1.24.5

- **Windows**: Run `go build -ldflags -H=windowsgui` to build a GUI executable.
- **Mac**: Additional configuration is required for system tray functionality. Refer to the [fyne-io/systray README](https://github.com/fyne-io/systray?tab=readme-ov-file#macos) for setup instructions. Then run `go build`.
- **Linux**: On Linux running `go build` should be enough. If you are using an older desktop environment and run into problems refer to [this link](https://github.com/fyne-io/systray?tab=readme-ov-file#linuxbsd).
