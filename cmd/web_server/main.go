package main

import (
	"net/http"

	"dev.azure.com/trayport/Hackathon/_git/Q/internal/logger"
)

func main() {
	bootstrap()
	logger.Fatal(http.ListenAndServe(":31979", nil))
	logger.Info("web server listening on http://localhost:31979")
}

func bootstrap() {
	logger.Info("bootstrapping web server.")

	addApiEndpoints()
	addStaticSite()
}

func addApiEndpoints() {
	http.HandleFunc("POST /api/issues", PostIssue)
	http.HandleFunc("PATCH /api/issues/{issueId}", PatchIssue)
	http.HandleFunc("GET /api/issues/{issueId}", GetIssue)
	http.HandleFunc("GET /api/issues", GetIssues)
	http.HandleFunc("DELETE /api/issues/{issueId}", DeleteIssue)
	http.HandleFunc("POST /api/issues/{issueId}/tags", PostTags)
	http.HandleFunc("DELETE /api/issues/{issueId}/tags/{tag}", DeleteTag)
}

func addStaticSite() {
	http.HandleFunc("/", serveStaticSite)
}
