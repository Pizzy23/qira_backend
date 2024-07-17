package losshigh

type AggregatedLossControl struct {
	ThreatEventId  int64   `json:"threat_event_id"`
	ThreatEvent    string  `json:"threat_event"`
	Assets         string  `json:"assets"`
	LossType       string  `json:"loss_type"`
	MinimumLoss    float64 `json:"minimum_loss"`
	MaximumLoss    float64 `json:"maximum_loss"`
	MostLikelyLoss float64 `json:"most_likely_loss"`
}

type AggregatedLoss struct {
	ThreatEventID  int64   `json:"threat_event_id"`
	ThreatEvent    string  `json:"threat_event"`
	Assets         string  `json:"assets"`
	LossType       string  `json:"loss_type"`
	MinimumLoss    float64 `json:"minimum_loss"`
	MaximumLoss    float64 `json:"maximum_loss"`
	MostLikelyLoss float64 `json:"most_likely_loss"`
}

type AggregatedLossResponse struct {
	ThreatEventID int64                  `json:"threat_event_id"`
	ThreatEvent   string                 `json:"threat_event"`
	Assets        string                 `json:"assets"`
	Losses        []AggregatedLossDetail `json:"losses"`
}

type AggregatedLossDetail struct {
	LossType       string  `json:"loss_type"`
	MinimumLoss    float64 `json:"minimum_loss"`
	MaximumLoss    float64 `json:"maximum_loss"`
	MostLikelyLoss float64 `json:"most_likely_loss"`
}
