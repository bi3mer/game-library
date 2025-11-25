package main

import (
	"fmt"
	"image/color"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func github_release_url(repo, filename string) string {

	switch runtime.GOOS {
	case "darwin":
		return fmt.Sprintf("https://github.com/bi3mer/%s/releases/latest/download/mac-%s",
			repo, filename)
	case "windows":
		return fmt.Sprintf("https://github.com/bi3mer/%s/releases/latest/download/win-%s.exe",
			repo, filename)
	}

	fmt.Printf("Unsupported OS type: %s\n", runtime.GOOS)
	os.Exit(1)
	return "FAIL"
}

type Game struct {
	Name     string
	URL      string
	FileName string
}

func NewGame(name, repo, exeName string) Game {
	filename := exeName
	if runtime.GOOS == "windows" {
		filename += ".exe"
	}

	return Game{
		Name:     name,
		URL:      github_release_url(repo, exeName),
		FileName: filepath.Join("builds", filename),
	}
}

func GameFetch(g Game) bool {
	resp, err := http.Get(g.URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fetch error fetching %s: %v\n", g.URL, err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Bad status for %s: %d\n", g.URL, resp.StatusCode)
		return false
	}

	out, err := os.Create(g.FileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create file %s: %v\n", g.FileName, err)
		return false
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		fmt.Fprintf(os.Stderr, "Failed writing %s: %v\n", g.FileName, err)
		return false
	}

	if runtime.GOOS == "darwin" {
		if err := os.Chmod(g.FileName, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to chmod %s: %v\n", g.FileName, err)
			return false
		}
	}

	return true
}

func isDownloaded(g Game) bool {
	_, err := os.Stat(g.FileName)
	return err == nil
}

func main() {
	err := os.Mkdir("builds", 0750)
	if err != nil && !os.IsExist(err) {
		fmt.Println("ERROR: Unable to make the 'builds' directory.")
		os.Exit(1)
	}

	games := []Game{
		NewGame("Wordle", "c-wordle", "wordle"),
		NewGame("Tic-Tac-Toe", "c-tic-tac-toe", "tic-tac-toe"),
		NewGame("Snake", "c-snake", "snake"),
		NewGame("Pong", "raylib-pong", "pong"),
	}

	buttons := make([]widget.Clickable, len(games))
	go func() {
		w := new(app.Window)
		w.Option(app.Title("Colan's Game Library"))
		w.Option(app.Size(400, 500))

		th := material.NewTheme()

		var ops op.Ops
		var list widget.List
		list.Axis = layout.Vertical

		for {
			e := w.Event()
			switch e := e.(type) {
			case app.DestroyEvent:
				os.Exit(0)
			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)

				for i := range games {
					if buttons[i].Clicked(gtx) {
						if isDownloaded(games[i]) {
							path := games[i].FileName
							if runtime.GOOS != "windows" {
								path = "./" + path
							}
							exec.Command(path).Start()
						} else {
							GameFetch(games[i])
						}
					}
				}

				layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{
							Top:    unit.Dp(16),
							Bottom: unit.Dp(16),
						}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							title := material.H4(th, "Colan's Game Library")
							title.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
							title.Alignment = text.Middle
							return title.Layout(gtx)
						})
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return material.List(th, &list).Layout(gtx, len(games), func(gtx layout.Context, i int) layout.Dimensions {
							return layout.Inset{
								Top:    unit.Dp(8),
								Bottom: unit.Dp(8),
							}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return layout.Flex{
										Axis:      layout.Horizontal,
										Alignment: layout.Middle,
									}.Layout(gtx,
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											gtx.Constraints.Min.X = gtx.Dp(unit.Dp(150))
											gtx.Constraints.Max.X = gtx.Dp(unit.Dp(150))
											label := material.Body1(th, games[i].Name)
											label.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
											return label.Layout(gtx)
										}),
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											btnText := "Download"
											if isDownloaded(games[i]) {
												btnText = "Play"
											}
											btn := material.Button(th, &buttons[i], btnText)
											return btn.Layout(gtx)
										}),
									)
								})
							})
						})
					}),
				)

				e.Frame(gtx.Ops)
			}
		}
	}()
	app.Main()
}
