package main

import (
	"regexp"
)

var ffmpegTimeRegex = regexp.MustCompile(`^\d{2}:\d{2}:\d{2}(\.\d{1,2})?$`)

func IsValidFFmpegTime(s string) bool {
	return ffmpegTimeRegex.MatchString(s)
}

func contains(s []string, target string) bool {
	for _, v := range s {
		if v == target {
			return true
		}
	}
	return false
}
