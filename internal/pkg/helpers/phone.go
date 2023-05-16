package helpers

import (
	"regexp"
)

func CheckPhone(str string) bool {
	isMatch, _ := regexp.MatchString("^\\d{1,11}$", str)
	return isMatch
}
