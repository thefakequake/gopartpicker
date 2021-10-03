package gopartpicker

import (
	"strings"

	"github.com/dlclark/regexp2"
)

// Converts certain PCPartPicker list URLs to a specific format in order to prevent client side JS loading.
func ConvertListURL(URL string) string {
	re := regexp2.MustCompile(`(https|http):\/\/([a-z].{2})?(pcpartpicker|partpicker).com\/user\/[a-zA-Z0-9]*\/saved\/#view=([a-zA-Z0-9]){4-8}`, 0)
	match, _ := re.MatchString(URL)

	if !match {
		return URL
	}

	return strings.Replace(URL, "#view=", "", 1)
}

// Checks if a URL is a PCPartPicker URL, making sure to check all regional subdomains
func MatchPCPPURL(URL string) bool {
	re := regexp2.MustCompile(`(https|http):\/\/([a-z].{2})?(pcpartpicker|partpicker).com\/?`, 0)

	match, _ := re.MatchString(URL)

	return match
}

// Checks if a URL is PCPartPicker product URL
func MatchProductURL(URL string) bool {
	re := regexp2.MustCompile(`(https|http):\/\/([a-z].{2})?(pcpartpicker|partpicker).com\/product\/[a-zA-Z0-9]{4,8}\/[\S]*`, 0)

	match, _ := re.MatchString(URL)

	return match
}

func regexp2FindAllString(re *regexp2.Regexp, s string) []string {
	var matches []string
	m, _ := re.FindStringMatch(s)
	for m != nil {
		matches = append(matches, m.String())
		m, _ = re.FindNextMatch(m)
	}
	return matches
}
