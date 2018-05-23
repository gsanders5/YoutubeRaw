package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"github.com/BrianAllred/goydl"
	"path/filepath"
	"os"
	"time"
)

var tmpDir = os.TempDir()

func main() {

	t, err := ioutil.TempDir(os.TempDir(), "ytraw-")
	if err != nil {
		log.Fatal(err)
	}
	tmpDir = t
	defer os.Remove(t)


	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/v/", youtubeHandler)
	fs := http.FileServer(http.Dir(tmpDir))
	http.Handle("/s/", http.StripPrefix("/s/", fs))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func youtubeHandler(w http.ResponseWriter, r *http.Request) {

	// Validate the path has "arguments"
	if r.URL.Path == "/v/" {
		youtubeErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	// Get the video ID
	videoId := strings.TrimPrefix(r.URL.Path, "/v/")
	videoFileName := ("ytraw-" + videoId + ".mp4")
	videoFilePath := filepath.Join(tmpDir, videoFileName)
	pidFilePath := filepath.Join(tmpDir, (videoFileName +".pid"))

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
		time.Sleep(6 * time.Hour)
		os.Remove(videoFilePath)
	}(videoFilePath)
}
