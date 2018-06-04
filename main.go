package main

import (
	"fmt"
	"github.com/BrianAllred/goydl"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var version = "Unknown Version"
var tmpDir = os.TempDir()

func main() {
	app := cli.NewApp()
	app.Name = "ytraw"
	app.Usage = "Serve YouTube videos as raw mp4 files on demand"
	app.Version = version

	var bindAddress string

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "bindAddress, b",
			Value:       "127.0.0.1:8080",
			Usage:       "Set the IP address and port to bind the server to",
			Destination: &bindAddress,
		},
		cli.StringFlag{
			Name:        "temporaryDirectory, d",
			Value:       os.TempDir(),
			Usage:       "Set a directory to store files in",
			Destination: &tmpDir,
		},
	}

	app.Action = func(c *cli.Context) error {
		startServer(bindAddress)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func startServer(bindAddress string) {
	t, err := ioutil.TempDir(tmpDir, "ytraw-")
	if err != nil {
		log.Fatal(err)
	}
	tmpDir = t
	defer os.Remove(t)

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/v/", youtubeVideoHandler)
	http.HandleFunc("/t/", youtubeThumbnailHandler)
	http.HandleFunc("/e/", youtubeEmbedHandler)
	fs := http.FileServer(http.Dir(tmpDir))
	http.Handle("/s/", http.StripPrefix("/s/", fs))

	log.Fatal(http.ListenAndServe((bindAddress), nil))
}

func youtubeVideoHandler(w http.ResponseWriter, r *http.Request) {

	// Validate the path has "arguments"
	if r.URL.Path == "/v/" {
		youtubeErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	// Get the video ID
	videoId := strings.TrimPrefix(r.URL.Path, "/v/")
	videoFileName := ("ytraw-" + videoId + ".mp4")
	videoFilePath := filepath.Join(tmpDir, videoFileName)
	pidFilePath := filepath.Join(tmpDir, (videoFileName + ".pid"))

	// If the file exists, serve it without running youtube-dl again.
	if _, err := os.Stat(videoFilePath); !os.IsNotExist(err) {
		fmt.Println("Redirecting, file existed.")
		http.Redirect(w, r, "/s/"+videoFileName, http.StatusSeeOther)
		return
	}

	// ToDo: Validate existence of video.

	fmt.Println("Running youtube-dl.")
	youtubeDl := goydl.NewYoutubeDl()
	//youtubeDl.YoutubeDlPath = "/usr/bin/youtube-dl"
	youtubeDl.Options.Output.Value = videoFilePath
	youtubeDl.Options.Format.Value = "bestvideo[ext=mp4]+bestaudio[ext=m4a]/bestvideo+bestaudio"
	youtubeDl.Options.MergeOutputFormat.Value = "mp4"
	youtubeDl.Download("https://www.youtube.com/watch?v=" + videoId)
	cmd, dlErr := youtubeDl.Download("https://www.youtube.com/watch?v=" + videoId)
	ioutil.WriteFile(pidFilePath, []byte(fmt.Sprintf("%d", cmd.Process.Pid)), 0664)
	cmd.Wait()
	os.Remove(pidFilePath)
	if dlErr != nil {
		log.Fatal(dlErr)
	}

	fmt.Println("Redirecting, file did not exist but was created.")
	http.Redirect(w, r, "/s/"+videoFileName, http.StatusSeeOther)

	// Remove any files created by youtube-dl after an interval of time has passed.
	go func(videoFilePath string) {
		time.Sleep(3 * time.Hour)
		os.Remove(videoFilePath)
	}(videoFilePath)
}

func youtubeThumbnailHandler(w http.ResponseWriter, r *http.Request) {

	// Validate the path has "arguments"
	if r.URL.Path == "/t/" {
		youtubeErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	// Get the video ID
	videoId := strings.TrimPrefix(r.URL.Path, "/t/")
	thumbnailFileName := ("ytraw-" + videoId + ".jpg")
	thumbnailFilePath := filepath.Join(tmpDir, thumbnailFileName)
	thumbnailFileURL := "https://img.youtube.com/vi/" + videoId + "/hqdefault.jpg"

	// If the file exists, serve it without running youtube-dl again.
	if _, err := os.Stat(thumbnailFilePath); !os.IsNotExist(err) {
		fmt.Println("Redirecting, file existed.")
		http.Redirect(w, r, "/s/"+thumbnailFileName, http.StatusSeeOther)
		return
	}

	// ToDo: Validate existence of video.

	fmt.Println("Grabbing Thumbnail.")

	err := DownloadFile(thumbnailFilePath, thumbnailFileURL)
	if err != nil {
		panic(err)
	}

	fmt.Println("Redirecting, file did not exist but was created.")
	http.Redirect(w, r, "/s/"+thumbnailFileName, http.StatusSeeOther)

	// Remove any files created by youtube-dl after an interval of time has passed.
	go func(videoFilePath string) {
		time.Sleep(3 * time.Hour)
		os.Remove(videoFilePath)
	}(thumbnailFilePath)
}
