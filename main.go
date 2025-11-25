package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

func main() {
	fmt.Printf("Colan's Game Library\n")

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

	GameFetch(games[3])
}
