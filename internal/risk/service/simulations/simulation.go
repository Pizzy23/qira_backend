package simulation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"qira/db"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

type ThreatEventRequest struct {
	MinFreq  float64 `json:"minfreq,omitempty"`
	PertFreq float64 `json:"pertfreq,omitempty"`
	MaxFreq  float64 `json:"maxfreq,omitempty"`
	MinLoss  float64 `json:"minloss,omitempty"`
	PertLoss float64 `json:"pertloss,omitempty"`
	MaxLoss  float64 `json:"maxloss,omitempty"`
}

type AnalyzeRequest struct {
	ThreatEvents []ThreatEventRequest `json:"threat_events"`
}

type AnalyzeResponse struct {
	Bins     []float64 `json:"bins"`
	CumFreqs []float64 `json:"cum_freqs"`
	Freqs    []int     `json:"freqs"`
	Lecs     []float64 `json:"lecs"`
}

type Bin struct {
	Frequency int     `json:"frequency"`
	MidPoint  float64 `json:"midPoint"`
}

type FrontEndResponse struct {
	FrequencyMax  float64   `json:"FrequencyMax"`
	FrequencyMin  float64   `json:"FrequencyMin"`
	FrequencyMode float64   `json:"FrequencyMode"`
	LossMax       float64   `json:"LossMax"`
	LossMin       float64   `json:"LossMin"`
	LossMode      float64   `json:"LossMode"`
	Bins          []Bin     `json:"bins"`
	Lecs          []float64 `json:"lecs"`
	CumLecs       []float64 `json:"cum_lecs"`
}

func MonteCarloSimulation(c *gin.Context, threatEvent string, reciverEmail string) {
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

	frequencyRequests := []ThreatEventRequest{}
	lossRequests := []ThreatEventRequest{}

	for _, risk := range riskCalculations {
		if risk.RiskType == "Frequency" {
			frequencyRequests = append(frequencyRequests, ThreatEventRequest{
				MinFreq:  risk.Min,
				PertFreq: risk.Mode,
				MaxFreq:  risk.Max,
			})
		} else if risk.RiskType == "Loss" {
			lossRequests = append(lossRequests, ThreatEventRequest{
				MinLoss:  risk.Min,
				PertLoss: risk.Mode,
				MaxLoss:  risk.Max,
			})
		}
	}

	threatEventRequests := []ThreatEventRequest{}
	for i := range frequencyRequests {
		te := ThreatEventRequest{
			MinFreq:  frequencyRequests[i].MinFreq,
			PertFreq: frequencyRequests[i].PertFreq,
			MaxFreq:  frequencyRequests[i].MaxFreq,
		}
		if i < len(lossRequests) {
			te.MinLoss = lossRequests[i].MinLoss
			te.PertLoss = lossRequests[i].PertLoss
			te.MaxLoss = lossRequests[i].MaxLoss
		}
		threatEventRequests = append(threatEventRequests, te)
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

	if response.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(response.Body)
		c.Set("Response", fmt.Sprintf("Received non-OK response: %s", string(bodyBytes)))
		c.Status(http.StatusInternalServerError)
		return
	}

	var analyzeResponse AnalyzeResponse
	if err := json.NewDecoder(response.Body).Decode(&analyzeResponse); err != nil {
		bodyBytes, _ := ioutil.ReadAll(response.Body)
		c.Set("Response", fmt.Sprintf("Error decoding response: %s", string(bodyBytes)))
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

	finalResponse := FrontEndResponse{
		FrequencyMax:  threatEventRequests[0].MaxFreq,
		FrequencyMin:  threatEventRequests[0].MinFreq,
		FrequencyMode: threatEventRequests[0].PertFreq,
		LossMax:       threatEventRequests[0].MaxLoss,
		LossMin:       threatEventRequests[0].MinLoss,
		LossMode:      threatEventRequests[0].PertLoss,
		Bins:          bins,
		Lecs:          analyzeResponse.Lecs,
		CumLecs:       analyzeResponse.CumFreqs,
	}

	c.JSON(http.StatusOK, finalResponse)
}
