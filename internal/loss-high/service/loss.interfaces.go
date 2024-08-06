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
	TypeOfLoss     string  `json:"type_of_loss"`
	MinimumLoss    float64 `json:"minimum_loss"`
	MaximumLoss    float64 `json:"maximum_loss"`
	MostLikelyLoss float64 `json:"most_likely_loss"`
}

type InputLossHigh struct {
	ThreatEvent    string
	Assets         []string
	LossType       string
	Impact         string
	MinimumLoss    float64
	MaximumLoss    float64
	MostLikelyLoss float64
}

type AggregatedLossDetailGranulade struct {
	LossType       string  `json:"LossType"`
	Impact         string  `json:"Impact"`
	LossEditNumber int64   `json:"loss_edit_number"`
	MinimumLoss    float64 `json:"MinimumLoss"`
	MaximumLoss    float64 `json:"MaximumLoss"`
	MostLikelyLoss float64 `json:"MostLikelyLoss"`
}

type AggregatedLossResponseGranulade struct {
	ThreatEventID int64
	ThreatEvent   string
	Assets        string
	Losses        []AggregatedLossDetailGranulade
}
