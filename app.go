package main

import (
	"image"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/aaronzipp/deeptube/video"
	"github.com/aaronzipp/deeptube/youtube"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const numVideos = 60
const logoPath = "assets/logo.png"

const applicationName = "DeepTube"

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
		return canvas.NewImageFromFile(logoPath)
	}
	defer resp.Body.Close()
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return canvas.NewImageFromFile(logoPath)
	}
	image := canvas.NewImageFromImage(img)
	image.SetMinSize(fyne.NewSize(210, 118)) // 16:9 ratio, like YouTube
	image.FillMode = canvas.ImageFillContain
	return image
}

func launchGUI(a fyne.App) {
	videos, err := video.VideosFromDB()
	if err != nil {
		panic(err)
	}
	videos.Sort()
	if len(videos) > numVideos {
		videos = videos[:numVideos]
	}

	w := a.NewWindow(applicationName)
	w.SetIcon(a.Icon())

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

		watchBtn := widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {
			openBrowser(vid.YouTubeLink())
		})

		hideBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
			vid.Hide()
		})

		videoCard := widget.NewCard("", "", container.NewBorder(nil, container.NewHBox(watchBtn, hideBtn), nil, nil, card))
		cards = append(cards, videoCard)
	}

	grid := container.NewGridWithColumns(4, cards...)
	scroll := container.NewVScroll(grid)

	w.SetContent(scroll)
	w.Resize(fyne.NewSize(1200, 800))
	w.Show()
}

func main() {
	a := app.New()

	logo, err := fyne.LoadResourceFromPath(logoPath)
	if err == nil {
		a.SetIcon(logo)
	}

	launchItem := fyne.NewMenuItem("Launch", func() {
		launchGUI(a)
	})

	refreshItem := fyne.NewMenuItem("Refresh", func() {
		err := youtube.RefreshVideos()
		if err != nil {
			// TODO: handle this error by showing the user
		}
	})

	menu := fyne.NewMenu(applicationName, launchItem, refreshItem)

	if desk, ok := a.(desktop.App); ok {
		desk.SetSystemTrayMenu(menu)
	}

	go func() {
		ticker := time.NewTicker(3 * time.Hour)
		for range ticker.C {
			// TODO: log any potential errors
			youtube.RefreshVideos()
		}
	}()

	a.Run()
}
