package main

import (
	"github.com/aaronzipp/deeptube/video"
	"image"
	"net/http"
	"os/exec"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const numVideos = 60

func openBrowser(url string) {
	switch runtime.GOOS {
	case "windows":
		exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		exec.Command("open", url).Start()
	default:
		exec.Command("xdg-open", url).Start()
	}
}

func loadImage(url string) *canvas.Image {
	resp, err := http.Get(url)
	if err != nil {
		return canvas.NewImageFromResource(theme.FyneLogo())
	}
	defer resp.Body.Close()
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return canvas.NewImageFromResource(theme.FyneLogo())
	}
	image := canvas.NewImageFromImage(img)
	image.SetMinSize(fyne.NewSize(210, 118)) // 16:9 ratio, like YouTube
	image.FillMode = canvas.ImageFillContain
	return image
}

func main() {
	videos, err := video.VideosFromDB()
	if err != nil {
		panic(err)
	}
	videos.Sort()
	if len(videos) > numVideos {
		videos = videos[:numVideos]
	}

	a := app.New()
	w := a.NewWindow("DeepTube")

	var cards []fyne.CanvasObject

	for _, vid := range videos {
		thumbnail := loadImage(vid.Thumbnail)

		title := widget.NewLabelWithStyle(
			vid.Title,
			fyne.TextAlignLeading,
			fyne.TextStyle{Bold: true},
		)
		title.Wrapping = fyne.TextWrapWord

		channel := widget.NewLabel(vid.ChannelName)
		duration := vid.VideoLength.String()
		if vid.WasLive {
			duration += " LIVE"
		}
		durationLabel := widget.NewLabel(duration)
		published := widget.NewLabel(vid.TimeSincePublished())

		infoBox := container.NewVBox(
			title,
			channel,
			container.NewHBox(durationLabel, widget.NewLabel(" • "), published),
		)

		card := container.NewVBox(
			thumbnail,
			infoBox,
		)

		btn := widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {
			openBrowser(vid.YouTubeLink())
		})

		videoCard := widget.NewCard("", "", container.NewBorder(nil, btn, nil, nil, card))
		cards = append(cards, videoCard)
	}

	grid := container.NewGridWithColumns(4, cards...)
	scroll := container.NewVScroll(grid)

	w.SetContent(scroll)
	w.SetFullScreen(true)
	w.ShowAndRun()
}
