package simulation

import (
	"bytes"
	"encoding/json"
	"net/http"
	"qira/db"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

type AcceptableLoss struct {
	Risk string  `json:"risk"`
	Loss float64 `json:"loss"`
}

type FrontEndResponseApp struct {
	FrequencyMax     float64          `json:"FrequencyMax"`
	FrequencyMin     float64          `json:"FrequencyMin"`
	FrequencyMode    float64          `json:"FrequencyMode"`
	LossMax          float64          `json:"LossMax"`
	LossMin          float64          `json:"LossMin"`
	LossMode         float64          `json:"LossMode"`
	Bins             []Bin            `json:"bins"`
	Lecs             []float64        `json:"lecs"`
	CumLecs          []float64        `json:"cum_lecs"`
	AcceptableLosses []AcceptableLoss `json:"acceptable_losses"`
}

func MonteCarloSimulationAppetite(c *gin.Context, threatEvent string, reciverEmail string) {
	var riskCalculations []db.RiskCalculation

	engine, exists := c.Get("db")
	if !exists {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	err := engine.(*xorm.Engine).Where("threat_event = ?", threatEvent).Find(&riskCalculations)
	if err != nil {
		c.Set("Response", "Error retrieving risk calculations")
		c.Status(http.StatusInternalServerError)
		return
	}

	var totalMinFreq, totalPertFreq, totalMaxFreq float64
	var totalMinLoss, totalPertLoss, totalMaxLoss float64

	frequencyRequests := make([]ThreatEventRequest, len(riskCalculations))
	lossRequests := make([]ThreatEventRequest, len(riskCalculations))

	for i, risk := range riskCalculations {
		if risk.RiskType == "Frequency" {
			totalMinFreq += risk.Min
			totalPertFreq += risk.Mode
			totalMaxFreq += risk.Max
			frequencyRequests[i] = ThreatEventRequest{
				MinFreq:  risk.Min,
				PertFreq: risk.Mode,
				MaxFreq:  risk.Max,
			}
		} else if risk.RiskType == "Loss" {
			totalMinLoss += risk.Min
			totalPertLoss += risk.Mode
			totalMaxLoss += risk.Max
			lossRequests[i] = ThreatEventRequest{
				MinLoss:  risk.Min,
				PertLoss: risk.Mode,
				MaxLoss:  risk.Max,
			}
		}
	}

	threatEventRequests := make([]ThreatEventRequest, len(frequencyRequests))
	for i := range frequencyRequests {
		threatEventRequests[i] = ThreatEventRequest{
			MinFreq:  frequencyRequests[i].MinFreq,
			PertFreq: frequencyRequests[i].PertFreq,
			MaxFreq:  frequencyRequests[i].MaxFreq,
			MinLoss:  lossRequests[i].MinLoss,
			PertLoss: lossRequests[i].PertLoss,
			MaxLoss:  lossRequests[i].MaxLoss,
		}
	}

	analyzeRequest := AnalyzeRequest{
		ThreatEvents: threatEventRequests,
	}

	requestBody, err := json.Marshal(analyzeRequest)
	if err != nil {
		c.Set("Response", "Error marshaling request body")
		c.Status(http.StatusInternalServerError)
		return
	}

	response, err := http.Post("https://qira-bellujrb-test.replit.app/analyze", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		c.Set("Response", "Error sending request")
		c.Status(http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	var analyzeResponse AnalyzeResponse
	if err := json.NewDecoder(response.Body).Decode(&analyzeResponse); err != nil {
		c.Set("Response", "Error decoding response")
		c.Status(http.StatusInternalServerError)
		return
	}

	bins := make([]Bin, len(analyzeResponse.Bins)-1)
	for i := 0; i < len(analyzeResponse.Bins)-1; i++ {
		midPoint := (analyzeResponse.Bins[i] + analyzeResponse.Bins[i+1]) / 2
		bins[i] = Bin{
			Frequency: analyzeResponse.Freqs[i],
			MidPoint:  midPoint,
		}
	}

	acceptableLosses := []AcceptableLoss{
		{"100%", totalMaxLoss},
		{"75%", totalMaxLoss * 0.75},
		{"50%", totalMaxLoss * 0.5},
		{"25%", totalMaxLoss * 0.25},
		{"0%", totalMaxLoss * 0},
	}

	finalResponse := FrontEndResponseApp{
		FrequencyMax:     totalMaxFreq,
		FrequencyMin:     totalMinFreq,
		FrequencyMode:    totalPertFreq,
		LossMax:          totalMaxLoss,
		LossMin:          totalMinLoss,
		LossMode:         totalPertLoss,
		Bins:             bins,
		Lecs:             analyzeResponse.Lecs,
		CumLecs:          analyzeResponse.CumFreqs,
		AcceptableLosses: acceptableLosses,
	}

	c.JSON(http.StatusOK, finalResponse)
}
