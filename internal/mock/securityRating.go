package mock

import "fmt"

type SecurityRating struct {
	Score   int    `json:"score"`
	Range   string `json:"range"`
	Min     string `json:"min"`
	Max     string `json:"max"`
	Average string `json:"average"`
}

func getSecurityRatings() []SecurityRating {
	return []SecurityRating{
		{Score: 0, Range: "N/A", Min: "N/A", Max: "N/A", Average: "N/A"},
		{Score: 1, Range: "1-35%", Min: "1%", Max: "35%", Average: "18%"},
		{Score: 2, Range: "36-65%", Min: "36%", Max: "65%", Average: "51%"},
		{Score: 3, Range: "66-95%", Min: "66%", Max: "95%", Average: "81%"},
		{Score: 4, Range: "96-100%", Min: "96%", Max: "100%", Average: "98%"},
	}
}

func FindAverageByScore(score int) (string, error) {
	ratings := getSecurityRatings()
	for _, rating := range ratings {
		if rating.Score == score {
			return rating.Average, nil
		}
	}
	return "", fmt.Errorf("score %d not found", score)
}
