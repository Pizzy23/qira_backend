package util

import (
	"qira/internal/interfaces"
	"regexp"
	"strings"
)

type InputData struct {
	Description string `json:"description"`
	ID          int    `json:"id"`
	InScope     bool   `json:"in_scope"`
	ThreatEvent string `json:"threat_event"`
	ThreatGroup string `json:"threat_group"`
}

func sanitizeString(input string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9 ]+")
	return reg.ReplaceAllString(input, "")
}

func SanitizeInputCatalogue(input *interfaces.InputThreatEventCatalogue) interfaces.InputThreatEventCatalogue {
	input.Description = sanitizeString(input.Description)
	input.ThreatEvent = sanitizeString(input.ThreatEvent)
	input.ThreatGroup = sanitizeString(input.ThreatGroup)
	return *input
}

func CleanString(input string) string {
	cleanedString := strings.ReplaceAll(input, "\t", "")
	cleanedString = strings.ReplaceAll(cleanedString, "\r", "")
	cleanedString = strings.TrimSpace(cleanedString)
	cleanedString = strings.Join(strings.Fields(cleanedString), " ")
	return cleanedString
}
