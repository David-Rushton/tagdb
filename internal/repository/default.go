package repository

import (
	"fmt"
	"os"
	"path"

	"dev.azure.com/trayport/Hackathon/_git/Q/internal/logger"
	"github.com/google/uuid"
)

var (
	rootDir   = os.Getenv("Q_REPOSITORY_ROOT")
	issuesDir = path.Join(rootDir, "issues")
	tagsDir   = path.Join(rootDir, "tags")
)

func init() {
	fileInfo, err := os.Stat(rootDir)
	if err != nil {
		logger.Panicf("cannot find repository root: %s", rootDir)
	}

	if !fileInfo.IsDir() {
		logger.Panic("cannot find repository root must be a directory")
	}

	if err := createDirIfNotExists(issuesDir); err != nil {
		logger.Panicf("cannot create issues directory: %s", issuesDir)
	}

	if err := createDirIfNotExists(tagsDir); err != nil {
		logger.Panicf("cannot create issues directory: %s", issuesDir)
	}

	logger.Infof("repository root directory set to: %s", rootDir)
	logger.Infof("repository issues directory set to: %s", issuesDir)
	logger.Infof("repository tags directory set to: %s", tagsDir)
}

func AddIssue(title, markdownBody string, tags []string) (Issue, error) {
	err := isValidTags(tags)
	if !err {
		var err = fmt.Errorf("cannot add issue, tags %v are invalid, they should contains lowercase letters, numbers and hyphens only - starting with a letter", tags)
		return Issue{}, err
	}

	issue := Issue{
		Id:           uuid.NewString(),
		Title:        title,
		MarkdownBody: markdownBody,
		Tags:         tags,
	}

	return issue, writeIssue(issue.Id, issue.Title, issue.MarkdownBody, issue.Tags)
}

func UpdateIssue(issueId, title, markdownBody string, tags []string) (Issue, error) {
	if !issueExists(issueId) {
		var err = fmt.Errorf("cannot update issue %s, it does not exist", issueId)
		return Issue{}, err
	}

	err := isValidTags(tags)
	if !err {
		var err = fmt.Errorf("cannot add issue, tags %v are invalid, they should contains lowercase letters, numbers and hyphens only - starting with a letter", tags)
		return Issue{}, err
	}

	issue := Issue{
		Id:           issueId,
		Title:        title,
		MarkdownBody: markdownBody,
		Tags:         tags,
	}

	return issue, writeIssue(issue.Id, issue.Title, issue.MarkdownBody, issue.Tags)
}

func GetIssue(issueId string) (Issue, error) {
	if !issueExists(issueId) {
		var err = fmt.Errorf("cannot get issue %s, it does not exist", issueId)
		return Issue{}, err
	}

	return readIssue(issueId)
}

func GetIssues(filterTags []string, includeDeleted bool) ([]Issue, error) {
	return readIssues(filterTags, includeDeleted)
}

func AddTags(issueId string, tags []string) error {
	// TODO: Make atomic.
	for _, tag := range tags {
		if err := AddTag(issueId, tag); err != nil {
			return err
		}
	}

	return nil
}

func AddTag(issueId, tag string) error {
	if !isValidTag(tag) {
		var err = fmt.Errorf("tag %v is invalid, it should contains lowercase letters, numbers and hyphens only - starting with a letter", tag)
		return err
	}

	if !issueExists(issueId) {
		var err = fmt.Errorf("cannot add tag to issue %s, it does not exist", issueId)
		return err
	}

	return writeTag(issueId, tag)

}

func RemoveTag(issueId, tag string) error {
	return deleteTag(issueId, tag)
}
