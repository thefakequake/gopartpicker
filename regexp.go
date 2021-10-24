package gopartpicker

import (
	"strings"

	"github.com/dlclark/regexp2"
)

var (
	listURLConversionRequired = regexp2.MustCompile(`(https|http):\/\/([a-z].{2})?pcpartpicker.com\/user\/[a-zA-Z0-9]*\/saved\/#view=([a-zA-Z0-9]){4-8}`, 0)
	pcppURLCheck              = regexp2.MustCompile(`(https|http):\/\/([a-z].{2})?pcpartpicker.com\/?`, 0)
	productURLCheck           = regexp2.MustCompile(`(https|http):\/\/([a-z].{2})?pcpartpicker.com\/product\/[a-zA-Z0-9]{4,8}\/[\S]*`, 0)
	partListURLCheck          = regexp2.MustCompile(`(http|https):\/\/([a-z]{2}\.)?pcpartpicker.com\/((list\/[a-zA-Z0-9]{4,8})|((user\/\w*\/saved\/(#view=)?[a-zA-Z0-9]{4,8})))`, 0)
	vendorNameCheck           = regexp2.MustCompile(`(?<=pcpartpicker.com\/mr\/).*(?=\/)`, 0)
)

// Extracts the name of a vendor from a PCPartPicker affiliate link.
func ExtractVendorName(URL string) string {
	if URL == "" {
		return ""
	}
	m, err := vendorNameCheck.FindStringMatch(URL)
	if err != nil {
		return ""
	}
	return m.String()
}

// Converts certain PCPartPicker list URLs to a specific format in order to prevent client side JS loading.
func ConvertListURL(URL string) string {
	match, _ := listURLConversionRequired.MatchString(URL)

	if !match {
		return URL
	}

	return strings.Replace(URL, "#view=", "", 1)
}

// Checks if a URL is a PCPartPicker URL, making sure to check all regional subdomains.
func MatchPCPPURL(URL string) bool {
	match, _ := pcppURLCheck.MatchString(URL)

	return match
}

// Checks if a URL is PCPartPicker product URL.
func MatchProductURL(URL string) bool {
	match, _ := productURLCheck.MatchString(URL)

	return match
}

// Check if a URL is a PCPartPicker part list URL.
func MatchPartListURL(URL string) bool {
	match, _ := partListURLCheck.MatchString(URL)

	return match
}

// Returns a list of all valid PCPartPicker part list URLs in the provided text.
func ExtractPartListURLs(text string) []string {
	return regexp2FindAllString(partListURLCheck, text)
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
