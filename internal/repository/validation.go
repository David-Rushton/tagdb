package repository

import "regexp"

var (
	tagRegexp = regexp.MustCompile(`^[a_z]+[a-z-0-9]*$`)
)

func isValidTag(tag string) bool {
	return tagRegexp.MatchString(tag)
}
