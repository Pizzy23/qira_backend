package simulation

import "qira/db"

type ThreatEventRequest struct {
	MinFreq  float64 `json:"minfreq,omitempty"`
	PertFreq float64 `json:"pertfreq,omitempty"`
	MaxFreq  float64 `json:"maxfreq,omitempty"`
	MinLoss  float64 `json:"minloss,omitempty"`
	PertLoss float64 `json:"pertloss,omitempty"`
	MaxLoss  float64 `json:"maxloss,omitempty"`
}

type FrontEndResponse struct {
	FrequencyMax      float64 `json:"FrequencyMax"`
	FrequencyMin      float64 `json:"FrequencyMin"`
	FrequencyEstimate float64 `json:"FrequencyEstimate"`
	LossMax           float64 `json:"LossMax"`
	LossMin           float64 `json:"LossMin"`
	LossEstimate      float64 `json:"LossEstimate"`
}

type AcceptableLoss struct {
	Risk string  `json:"risk"`
	Loss float64 `json:"loss"`
}

type FrontEndResponseAppLoss struct {
	FrequencyMax      float64             `json:"FrequencyMax"`
	FrequencyMin      float64             `json:"FrequencyMin"`
	FrequencyEstimate float64             `json:"FrequencyEstimate"`
	LossMax           float64             `json:"LossMax"`
	LossMin           float64             `json:"LossMin"`
	LossEstimate      float64             `json:"LossEstimate"`
	LossExceedance    []db.LossExceedance `json:"LossExceedance"`
}

type FrontEndResponseAgg struct {
	FrequencyMax      float64 `json:"FrequencyMax"`
	FrequencyMin      float64 `json:"FrequencyMin"`
	FrequencyEstimate float64 `json:"FrequencyEstimate"`
	LossMax           float64 `json:"LossMax"`
	LossMin           float64 `json:"LossMin"`
	LossEstimate      float64 `json:"LossEstimate"`
}

type FrontEndResponseAppReport struct {
	ProposedMin    float64             `json:"ProposedMin"`
	ProposedMax    float64             `json:"ProposedMax"`
	ProposedPert   float64             `json:"ProposedPert"`
	LossExceedance []db.LossExceedance `json:"LossExceedance"`
}

type OutputProcess struct {
	FrequencyMax      float64 `json:"FrequencyMax"`
	FrequencyMin      float64 `json:"FrequencyMin"`
	FrequencyEstimate float64 `json:"FrequencyEstimate"`
	LossMax           float64 `json:"LossMax"`
	LossMin           float64 `json:"LossMin"`
	LossEstimate      float64 `json:"LossEstimate"`
}
