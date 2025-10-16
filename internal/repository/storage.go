package repository

import (
	"encoding/json"
	"os"
	"path"
)

func issueExists(id string) bool {
	_, err := os.Stat(getIssuePath(id))
	return err == nil
}

func tagExists(issueId, tag string) bool {
	// TODO: Check both paths for existence.
	// TODO: Fix is only one exists.
	_, err := os.Stat(getTagPaths(issueId, tag)[0])
	return err == nil
}

func readIssue(id string) (Issue, error) {
	data, err := os.ReadFile(getIssuePath(id))
	if err != nil {
		return Issue{}, err
	}

	var issue Issue
	err = json.Unmarshal(data, &issue)
	if err != nil {
		return Issue{}, err
	}

	return issue, nil
}

func writeIssue(issue Issue) error {
	data, err := issue.toJson()
	if err != nil {
		return err
	}

	path := getIssuePath(issue.Id)
	return os.WriteFile(path, []byte(data), 0644)
}

func writeTag(issueId, tag string) error {
	if !tagExists(issueId, tag) {
		if err := createDirIfNotExists(getTagDir(tag)); err != nil {
			return err
		}

		// TODO: Should be atomic.
		for _, path := range getTagPaths(issueId, tag) {
			if err := os.WriteFile(path, []byte{}, 0644); err != nil {
				return err
			}
		}
	}

	return nil
}

func deleteTag(issueId, tag string) error {
	if tagExists(issueId, tag) {
		// TODO: Should be atomic.
		for _, path := range getTagPaths(issueId, tag) {
			if err := os.Remove(path); err != nil {
				return err
			}
		}
	}

	return nil
}

func getIssuePath(issueId string) string {
	return path.Join(getIssueDir(issueId), "issue.json")
}

func getIssueDir(issueId string) string {
	return path.Join(issuesDir, issueId)
}

func getTagPaths(issueId, tag string) []string {
	return []string{
		path.Join(getTagDir(tag), issueId+".tag"),
		path.Join(getIssueDir(issueId), issueId+".tag"),
	}
}

func getTagDir(tag string) string {
	return path.Join(tagsDir, tag)
}

func createDirIfNotExists(dir string) error {
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(dir, 0755)
		}

		return err
	}

	return nil
}
