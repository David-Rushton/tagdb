package repository

import (
	"encoding/json"
	"os"
	"path"
	"slices"
)

func issueExists(id string) bool {
	_, err := os.Stat(getIssuePath(id))
	return err == nil
}

func tagExists(issueId, tag string) bool {
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

	tags, err := readTags(id)
	if err != nil {
		return Issue{}, err
	}

	issue.Tags = tags

	return issue, nil
}

func readIssues(filterTags []string, includeDeleted bool) ([]Issue, error) {
	// Read from filesystem.
	entries, err := os.ReadDir(issuesDir)
	if err != nil {
		return []Issue{}, err
	}

	// Rehydrate.
	issues := []Issue{}
	for _, entry := range entries {
		if entry.IsDir() {
			if issueExists(entry.Name()) {
				issue, err := readIssue(entry.Name())
				if err != nil {
					return []Issue{}, err
				}

				// Apply filters.
				shouldInclude := true

				if !includeDeleted && slices.Contains(issue.Tags, "deleted") {
					shouldInclude = false
				}

				if shouldInclude && len(filterTags) > 0 {
					matched := 0
					for _, filterTag := range filterTags {
						if slices.Contains(issue.Tags, filterTag) {
							matched++
						}
					}

					shouldInclude = matched == len(filterTags)
				}

				if shouldInclude {
					issues = append(issues, issue)
				}
			}
		}
	}

	return issues, nil
}

func writeIssue(issueId, title, markdownBody string, tags []string) error {
	issue := Issue{
		Id:           issueId,
		Title:        title,
		MarkdownBody: markdownBody,
	}

	data, err := issue.toJson()
	if err != nil {
		return err
	}

	dir := getIssueDir(issue.Id)
	if err := createDirIfNotExists(dir); err != nil {
		return err
	}

	oldTags, err := readTags(issue.Id)
	if err != nil {
		return err
	}

	for _, oldTag := range oldTags {
		if !slices.Contains(tags, oldTag) {
			if err := deleteTag(issue.Id, oldTag); err != nil {
				return err
			}
		}
	}

	for _, tag := range tags {
		if !slices.Contains(oldTags, tag) {
			if err := writeTag(issue.Id, tag); err != nil {
				return err
			}
		}
	}

	path := getIssuePath(issue.Id)
	return os.WriteFile(path, []byte(data), 0644)
}

func writeTag(issueId, tag string) error {
	if !tagExists(issueId, tag) {
		if err := createDirIfNotExists(getTagDir(tag)); err != nil {
			return err
		}

		// TODO: Make atomic.
		for _, path := range getTagPaths(issueId, tag) {
			if err := os.WriteFile(path, []byte{}, 0644); err != nil {
				return err
			}
		}
	}

	return nil
}

func readTags(issueId string) ([]string, error) {
	issueDir := getIssueDir(issueId)
	entries, err := os.ReadDir(issueDir)
	if err != nil {
		return []string{}, err
	}

	var tags []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if path.Ext(name) != ".tag" {
			continue
		}

		tag := name[:len(name)-len(".tag")]
		tags = append(tags, tag)
	}

	return tags, nil
}

func deleteTag(issueId, tag string) error {
	if tagExists(issueId, tag) {
		// TODO: Make atomic.
		for _, path := range getTagPaths(issueId, tag) {
			if err := os.Remove(path); err != nil {
				return err
			}
		}
	}

	return nil
}

func getIssuePath(issueId string) string {
	return path.Join(getIssueDir(issueId), issueId+".json")
}

func getIssueDir(issueId string) string {
	return path.Join(issuesDir, issueId)
}

func getTagPaths(issueId, tag string) []string {
	return []string{
		path.Join(getTagDir(tag), issueId+".tag"),
		path.Join(getIssueDir(issueId), tag+".tag"),
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
