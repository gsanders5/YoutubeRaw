package main

import (
	"fmt"
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {

	// 404 on invalid URL
	if r.URL.Path != "/" {
		genericErrorHandler(w, r, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	data := "<!DOCTYPE html><html lang=en><meta content=\"text/html; charset=utf-8\"http-equiv=Content-Type><title>Rvnx Feed Internals</title><meta content=\"width=device-width,initial-scale=1\"name=viewport><style>body{margin:40px auto;max-width:650px;line-height:1.6;font-size:18px;color:#444;padding:0 10px}h1,h2,h3{line-height:1.2}</style><main role=main><h1>Rvnx Feed Internals</h1><p>This page hosts automatically generated content for a local instance of the <a href=https://ttrss.rvnx.org/ >Rvnx Feed Reader</a>. If you stumbled upon this page, there is nothing to see here except temporarily cached media from <i>*public*</i> RSS feeds.<footer role=contentinfo><small><a href=https://www.rvnx.org/ >Copyright \u00A9 <time datetime=2018>2018</time> Rvnx</small></footer></main>"
	w.Header().Set("Content-Length", fmt.Sprint(len(data)))
	fmt.Fprint(w, string(data))
}

func genericErrorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		data := "<!DOCTYPE html><html lang=en><meta content=\"text/html; charset=utf-8\"http-equiv=Content-Type><title>Rvnx Feed Internals</title><meta content=\"width=device-width,initial-scale=1\"name=viewport><style>body{margin:40px auto;max-width:650px;line-height:1.6;font-size:18px;color:#444;padding:0 10px}h1,h2,h3{line-height:1.2}</style><main role=main><h1>Error: Page Not Found</h1><p>This server is meant for internal uses... Are you <i>really</i> meant to be here?<footer role=contentinfo><small><a href=https://www.rvnx.org/ >Copyright \u00A9 <time datetime=2018>2018</time> Rvnx</small></footer></main>"
		w.Header().Set("Content-Length", fmt.Sprint(len(data)))
		fmt.Fprint(w, string(data))
	}
}

func youtubeErrorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusBadRequest {
		w.Header().Set("Content-Type", "text/html")
		data := "<!DOCTYPE html><html lang=en><meta content=\"text/html; charset=utf-8\"http-equiv=Content-Type><title>Rvnx Feed Internals</title><meta content=\"width=device-width,initial-scale=1\"name=viewport><style>body{margin:40px auto;max-width:650px;line-height:1.6;font-size:18px;color:#444;padding:0 10px}h1,h2,h3{line-height:1.2}</style><main role=main><h1>Error: Invalid Request</h1><p>Error, you should have specified additional arguments in your URL. Are you <i>really</i> meant to be here?<footer role=contentinfo><small><a href=https://www.rvnx.org/ >Copyright \u00A9 <time datetime=2018>2018</time> Rvnx</small></footer></main>"
		w.Header().Set("Content-Length", fmt.Sprint(len(data)))
		fmt.Fprint(w, string(data))
	}
}