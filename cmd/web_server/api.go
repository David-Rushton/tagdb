package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"dev.azure.com/trayport/Hackathon/_git/Q/internal/logger"
	"dev.azure.com/trayport/Hackathon/_git/Q/internal/repository"
)

type IssueParams struct {
	Title        string `json:"title"`
	MarkdownBody string `json:"markdownBody"`
	Tags         []string
}

type TagParams struct {
	Tags []string `json:"tags"`
}

// Adds a new issue.
// Example request body:
//
//	{
//	  "title": "title",
//	  "markdownBody": "body, supports markdown",
//	  "tags": ["tag1", "tag2"]
//	}
func PostIssue(w http.ResponseWriter, r *http.Request) {
	logger.Info("POST api/issues requested")

	var issueParams IssueParams
	if err := json.NewDecoder(r.Body).Decode(&issueParams); err != nil {
		logger.Warnf("cannot read body because %v", err)
		http.Error(w, "400", http.StatusBadRequest)
		return
	}

	issue, err := repository.AddIssue(issueParams.Title, issueParams.MarkdownBody, issueParams.Tags)
	if err != nil {
		logger.Errorf("cannot create issue because %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(&issue)
	if err != nil {
		logger.Errorf("cannot convert issue to json because %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// Updates an existing issue.
func PatchIssue(w http.ResponseWriter, r *http.Request) {
	issueId := r.PathValue("issueId")
	logger.Infof("PATCH api/issues/%v requested", issueId)

	_, err := repository.GetIssue(issueId)
	if err != nil {
		http.Error(w, "404", http.StatusNotFound)
		return
	}

	var issueParams IssueParams
	if err := json.NewDecoder(r.Body).Decode(&issueParams); err != nil {
		logger.Warnf("cannot read body because %v", err)
		http.Error(w, "400", http.StatusBadRequest)
		return
	}

	issue, err := repository.UpdateIssue(
		issueId,
		issueParams.Title,
		issueParams.MarkdownBody,
		issueParams.Tags)
	if err != nil {
		logger.Errorf("cannot update issue because %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(&issue)
	if err != nil {
		logger.Errorf("cannot convert issue to json because %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// Returns an existing issue.
func GetIssue(w http.ResponseWriter, r *http.Request) {
	issueId := r.PathValue("issueId")
	logger.Infof("GET api/issues/%v requested", issueId)

	issue, err := repository.GetIssue(issueId)
	if err != nil {
		logger.Warnf("cannot find issue because %v", err)
		http.Error(w, "404", http.StatusNotFound)
		return
	}

	data, err := json.Marshal(&issue)
	if err != nil {
		logger.Errorf("cannot convert issue to json because %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// Returns all issues.
//
// Params:
//   - tags: optional, comma-separated list of tags to filter by
//   - includeDeleted: optional, when "true" deleted issues are included in the response
func GetIssues(w http.ResponseWriter, r *http.Request) {
	logger.Info("GET api/issues requested")

	// Read query string.
	queryString := r.URL.Query()
	includeDeleted := queryString.Get("includeDeleted") == "true"
	filterTags := strings.Split(queryString.Get("tags"), ",")
	if len(filterTags) == 1 && filterTags[0] == "" {
		filterTags = []string{}
	}

	// Get issues.
	issue, err := repository.GetIssues(filterTags, includeDeleted)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Serialize issues.
	data, err := json.Marshal(&issue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write response.
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// Deletes an issue.
//
// Deletes are soft.  The issue is tagged as deleted but not removed from the database.  Deleted
// issues are excluded from GET api/issues, etc.
func DeleteIssue(w http.ResponseWriter, r *http.Request) {
	issueId := r.PathValue("issueId")
	logger.Infof("DELETE api/issues/%v requested", issueId)

	// Get issue.
	_, err := repository.GetIssue(issueId)
	if err != nil {
		http.Error(w, "404", http.StatusNotFound)
		return
	}

	// Add deleted tag.
	if err := repository.AddTag(issueId, "deleted"); err != nil {
		logger.Errorf("cannot add deleted tag to issue because %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Tags an issue.
func PostTags(w http.ResponseWriter, r *http.Request) {
	issueId := r.PathValue("issueId")
	logger.Infof("POST api/issues/%v/tags requested", issueId)

	// Get issue.
	_, err := repository.GetIssue(issueId)
	if err != nil {
		http.Error(w, "404", http.StatusNotFound)
		return
	}

	// Deserialize request body.
	var tagParams TagParams
	if err := json.NewDecoder(r.Body).Decode(&tagParams); err != nil {
		logger.Warnf("cannot read body because %v", err)
		http.Error(w, "400", http.StatusBadRequest)
		return
	}

	// Add tags.
	if err := repository.AddTags(issueId, tagParams.Tags); err != nil {
		logger.Errorf("cannot add tags %v to issue because %v", tagParams.Tags, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the updated issue.
	issue, err := repository.GetIssue(issueId)
	if err != nil {
		logger.Errorf("cannot get issue %v after adding tags because %v", issueId, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(&issue)
	if err != nil {
		logger.Errorf("cannot convert issue to json because %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// Deletes a tag from an issue.
func DeleteTag(w http.ResponseWriter, r *http.Request) {
	issueId := r.PathValue("issueId")
	tag := r.PathValue("tag")
	logger.Infof("DELETE api/issues/%v/tags/%v requested", issueId, tag)

	// Get issue.
	_, err := repository.GetIssue(issueId)
	if err != nil {
		http.Error(w, "404", http.StatusNotFound)
		return
	}

	// Remove tag.
	if err := repository.RemoveTag(issueId, tag); err != nil {
		logger.Errorf("cannot remove tag %v because %v", tag, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the updated issue.
	issue, err := repository.GetIssue(issueId)
	if err != nil {
		logger.Errorf("cannot get issue %v after adding tags because %v", issueId, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(&issue)
	if err != nil {
		logger.Errorf("cannot convert issue to json because %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
