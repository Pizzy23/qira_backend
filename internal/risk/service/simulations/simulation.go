package simulation

import (
	"bytes"
	"encoding/json"
	"net/http"
	"qira/db"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

type ThreatEventRequest struct {
	MinFreq  float64 `json:"minfreq"`
	PertFreq float64 `json:"pertfreq"`
	MaxFreq  float64 `json:"maxfreq"`
	MinLoss  float64 `json:"minloss"`
	PertLoss float64 `json:"pertloss"`
	MaxLoss  float64 `json:"maxloss"`
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
	Max  float64 `json:"Max"`
	Min  float64 `json:"Min"`
	Mode float64 `json:"Mode"`
	Bins []Bin   `json:"bins"`
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

	// Montar o corpo da requisição
	threatEventRequests := []ThreatEventRequest{}
	for _, risk := range riskCalculations {
		if risk.RiskType == "Frequency" || risk.RiskType == "Loss" {
			threatEventRequests = append(threatEventRequests, ThreatEventRequest{
				MinFreq:  risk.Min,
				PertFreq: risk.Mode,
				MaxFreq:  risk.Max,
				MinLoss:  risk.Min,
				PertLoss: risk.Mode,
				MaxLoss:  risk.Max,
			})
		}
	}

	analyzeRequest := AnalyzeRequest{
		ThreatEvents: threatEventRequests,
	}

	// Enviar a requisição POST
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

	// Processar a resposta
	var analyzeResponse AnalyzeResponse
	if err := json.NewDecoder(response.Body).Decode(&analyzeResponse); err != nil {
		c.Set("Response", "Error decoding response")
		c.Status(http.StatusInternalServerError)
		return
	}

	// Transformar os dados da resposta
	bins := make([]Bin, len(analyzeResponse.Bins)-1)
	for i := 0; i < len(analyzeResponse.Bins)-1; i++ {
		midPoint := (analyzeResponse.Bins[i] + analyzeResponse.Bins[i+1]) / 2
		bins[i] = Bin{
			Frequency: analyzeResponse.Freqs[i],
			MidPoint:  midPoint,
		}
	}

	// Preparar a resposta final
	finalResponse := FrontEndResponse{
		Max:  708499968,
		Min:  354249984,
		Mode: 508500000,
		Bins: bins,
	}

	// Enviar a resposta para o cliente
	c.JSON(http.StatusOK, finalResponse)
}
