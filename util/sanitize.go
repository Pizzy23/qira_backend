package util

import (
	"qira/internal/interfaces"
	"regexp"
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
