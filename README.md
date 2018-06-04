# YoutubeRaw

A Golang based server to serve YouTube videos as raw mp4 files on demand. This was specifically created to have an easy way to convert YouTube embeds into HTML5 video embeds in TT-RSS.

## Building

To build a release version, ensure the app version is changed to the current short git commit identifier.

```go build -ldflags "-X main.version=$(git rev-parse --short HEAD)"  -o build/ytraw *.go```
