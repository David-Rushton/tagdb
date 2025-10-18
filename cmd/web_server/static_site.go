package main

import (
	"net/http"
	"path"

	"dev.azure.com/trayport/Hackathon/_git/Q/internal/logger"
)

func serveStaticSite(w http.ResponseWriter, r *http.Request) {
	filePath := path.Join("wwwRoot", r.URL.Path+".html")
	if _, err := http.Dir("wwwRoot").Open(r.URL.Path + ".html"); err != nil {
		// TODO: 500.
		logger.Errorf("cannot read web asset because %v", err)
	}

	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, filePath)
}
