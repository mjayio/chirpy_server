package util

import (
	"regexp"
	"strings"
)

// ReplaceSubstring replaces occurrences of old with new in the input string.
// If count is -1, it replaces all occurrences.
func ReplaceSubstring(input, old, new string, count int) string {
	return strings.Replace(input, old, new, count)
}

// ReplaceAllSubstrings replaces all occurrences of old with new in the input string.
func ReplaceAllSubstrings(input, old, new string) string {
	return strings.ReplaceAll(input, old, new)
}

// ReplaceWholeWord replaces occurrences of a word with a new word in case-sensitive manner.
// It only replaces complete words (not parts of other words).
func ReplaceWholeWord(input, oldWord, newWord string) string {
	pattern := `\b` + regexp.QuoteMeta(oldWord) + `\b`
	regex := regexp.MustCompile(pattern)
	return regex.ReplaceAllString(input, newWord)
}

// ReplaceWholeWordCaseInsensitive replaces occurrences of a word with a new word in case-insensitive manner.
// It only replaces complete words (not parts of other words).
func ReplaceWholeWordCaseInsensitive(input, oldWord, newWord string) string {
	pattern := `(?i)\b` + regexp.QuoteMeta(oldWord) + `\b`
	regex := regexp.MustCompile(pattern)
	return regex.ReplaceAllString(input, newWord)
}
