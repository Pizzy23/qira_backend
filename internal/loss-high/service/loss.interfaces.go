package losshigh

type AggregatedLoss struct {
	ThreatEvent    string  `json:"threat_event"`
	Assets         string  `json:"assets"`
	LossType       string  `json:"loss_type"`
	MinimumLoss    float64 `json:"minimum_loss"`
	MaximumLoss    float64 `json:"maximum_loss"`
	MostLikelyLoss float64 `json:"most_likely_loss"`
}
