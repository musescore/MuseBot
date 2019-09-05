package utils

import (
	"regexp"
	"strconv"
	"strings"
)

func SanitizeText(text string) string {
	text = strings.Replace(text, "<", "&lt;", -1)
	text = strings.Replace(text, ">", "&gt;", -1)
	text = strings.Replace(text, "\n", " ", -1)
	return text
}

func StripToLength(text string, size int, end string) string {
	result := ""
	for i, r := range text {
		if i == size {
			return result + end
		}
		result += string(r)
	}
	return result
}

func InSlice(slice []string, val ...string) bool {
	for _, sv := range slice {
		for _, vv := range val {
			if sv == vv {
				return true
			}
		}
	}
	return false
}

func HashtagToInt(s string) int {
	re, _ := regexp.Compile(`([#])([\d]+)`)
	match := re.FindStringSubmatch(s)
	if len(match) < 2 {
		return 0
	} else {
		i, _ := strconv.Atoi(match[2])
		return i
	}
}