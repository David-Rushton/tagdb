package repository

import "regexp"

var (
	tagRegexp = regexp.MustCompile(`^[a-z]+[a-z-0-9]*$`)
)

func isValidTags(tags []string) bool {
	for _, tag := range tags {
		if !isValidTag(tag) {
			return false
		}
	}

	return true
}

func isValidTag(tag string) bool {
	return tagRegexp.MatchString(tag)
}
