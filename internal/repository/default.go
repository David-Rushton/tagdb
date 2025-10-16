package repository

import (
	"fmt"
	"log"
	"os"
	"path"

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
		log.Panicf("cannot find repository root: %s", rootDir)
	}

	if !fileInfo.IsDir() {
		log.Panic("cannot find repository root must be a directory")
	}

	if err := createDirIfNotExists(issuesDir); err != nil {
		log.Panicf("cannot create issues directory: %s", issuesDir)
	}

	if err := createDirIfNotExists(tagsDir); err != nil {
		log.Panicf("cannot create issues directory: %s", issuesDir)
	}

	log.Printf("repository root directory set to: %s", rootDir)
	log.Printf("repository issues directory set to: %s", issuesDir)
	log.Printf("repository tags directory set to: %s", tagsDir)
}

func AddIssue(title, markdownBody string) (Issue, error) {
	issue := Issue{
		Id:           uuid.NewString(),
		Title:        title,
		MarkdownBody: markdownBody,
	}

	return issue, writeIssue(Issue{})
}

func UpdateIssue(issue Issue) (Issue, error) {
	if !issueExists(issue.Id) {
		var err = fmt.Errorf("cannot update issue %s, it does not exist", issue.Id)
		return Issue{}, err
	}

	return issue, writeIssue(issue)
}

func GetIssue(id string) (Issue, error) {
	if !issueExists(id) {
		var err = fmt.Errorf("cannot get issue %s, it does not exist", id)
		return Issue{}, err
	}

	return readIssue(id)
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
