package main

import (
	"github.com/aaronzipp/deeptube/video"
	"image"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"os/exec"
	"runtime"
)

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
	image.SetMinSize(fyne.NewSize(200, 112))
	image.FillMode = canvas.ImageFillContain
	return image
}

func main() {
	videos, err := video.VideosFromDB()
	if err != nil {
		panic(err)
	}
	videos.Sort()

	videos = videos[:10]

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
		channel := widget.NewLabel(vid.ChannelName)
		duration := vid.VideoLength.String()
		if vid.WasLive {
			duration += " LIVE"
		}
		durationLabel := widget.NewLabel(duration)
		published := widget.NewLabel(vid.TimeSincePublished())

		card := container.NewVBox(
			thumbnail,
			container.NewVBox(
				title,
				channel,
				durationLabel,
				published,
			),
		)

		videoCard := widget.NewCard("", "", container.NewPadded(card))
		cards = append(cards, videoCard)
	}

	videosContainer := container.NewGridWithColumns(4, cards...)

	scroll := container.NewVScroll(videosContainer)
	w.SetContent(scroll)
	w.Resize(fyne.NewSize(900, 600))
	w.ShowAndRun()
}
