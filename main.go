package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
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

func fetch_game(g Game) bool {
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

func main() {
	fmt.Printf("Colan's Game Library\n")

	games := []Game{
		{
			Name:     "Snake",
			URL:      github_release_url("c-snake", "snake"),
			FileName: "snake",
		},
		{
			Name:     "Pong",
			URL:      github_release_url("raylib-pong", "pong"),
			FileName: "pong",
		},
	}

	fetch_game(games[0])

}
