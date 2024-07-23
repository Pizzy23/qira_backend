package simulation

import (
	"fmt"
	"net/http"
	"os"
	"qira/db"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gopkg.in/mail.v2"
	"xorm.io/xorm"
)

type MonteCarloRisk struct {
	InherentMin  float64 `json:"inherent_min"`
	InherentMax  float64 `json:"inherent_max"`
	InherentMode float64 `json:"inherent_mode"`
}

func MonteCarloSimulationReport(c *gin.Context, threatEvent string, receiverEmail string) {
	var frequencies []db.Frequency
	var riskCalculations []db.RiskCalculation
	var control []db.Control
	var proposed []db.Propused

	engine, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}

	err := engine.(*xorm.Engine).Where("threat_event = ?", threatEvent).Find(&frequencies)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving frequencies"})
		return
	}

	err = engine.(*xorm.Engine).Where("threat_event = ?", threatEvent).Find(&riskCalculations)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving risk calculations"})
		return
	}

	err = engine.(*xorm.Engine).Where("control_id = ?", -2).Find(&control)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving control data"})
		return
	}

	err = engine.(*xorm.Engine).Find(&proposed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving proposed control data"})
		return
	}

	// Calculate the Inherent Risks
	var totalRiskMin, totalRiskMax, totalRiskMode float64
	for _, risk := range riskCalculations {
		if risk.RiskType == "Risk" {
			totalRiskMin += risk.Min
			totalRiskMax += risk.Max
			totalRiskMode += risk.Mode
		}
	}

	// Aggregate Control Gap
	var aggregateControlGap float64
	for _, ctrl := range control {
		gapStr := strings.TrimSuffix(ctrl.ControlGap, "%")
		gap, err := strconv.ParseFloat(gapStr, 64)
		if err == nil {
			aggregateControlGap += gap / 100.0 // Convertendo de porcentagem para valor decimal
		}
	}

	// Calculate the Inherent Risks adjusted by the Aggregate Control Gap
	inherentRiskMin := totalRiskMin / aggregateControlGap
	inherentRiskMax := totalRiskMax / aggregateControlGap
	inherentRiskMode := totalRiskMode / aggregateControlGap

	// Proposed Risks
	var proposedRiskMin, proposedRiskMax, proposedRiskMode float64
	for _, prop := range proposed {
		gapStr := strings.TrimSuffix(prop.ControlGap, "%")
		gap, err := strconv.ParseFloat(gapStr, 64)
		if err == nil {
			gap /= 100.0 // Convertendo de porcentagem para valor decimal
			proposedRiskMin += inherentRiskMin * gap
			proposedRiskMax += inherentRiskMax * gap
			proposedRiskMode += inherentRiskMode * gap
		}
	}

	// Generate the Monte Carlo Simulation and get bin data
	binData := generateMonteCarloSimulation(inherentRiskMin, inherentRiskMax, inherentRiskMode)

	// Combine all information into a report
	reportPath := "report.txt"
	err = createReport(reportPath, inherentRiskMin, inherentRiskMax, inherentRiskMode, proposedRiskMin, proposedRiskMax, proposedRiskMode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating report"})
		return
	}

	if receiverEmail != "" {
		sendEmailWithAttachmentsReport(receiverEmail, reportPath)
	}

	// Delete the report after sending
	err = os.Remove(reportPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting report"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Simulation report generated and sent via email successfully!",
		"bins":    binData,
	})
}

func generateMonteCarloSimulation(inherentRiskMin, inherentRiskMax, inherentRiskMode float64) []map[string]interface{} {
	const (
		amostras = 10000
		numBins  = 70
	)

	binWidth := (inherentRiskMax - inherentRiskMin) / float64(numBins)

	a := (4*inherentRiskMode + inherentRiskMax - 5*inherentRiskMin) / (inherentRiskMax - inherentRiskMin)
	b := (5*inherentRiskMax - inherentRiskMin - 4*inherentRiskMode) / (inherentRiskMax - inherentRiskMin)

	source := rand.NewSource(uint64(99))
	distribuicao := distuv.Beta{
		Alpha: a,
		Beta:  b,
		Src:   rand.New(source),
	}
	valores := make([]float64, amostras)
	for i := range valores {
		valores[i] = distribuicao.Rand()*(inherentRiskMax-inherentRiskMin) + inherentRiskMin
	}

	frequencias := make([]int, numBins)
	for _, valor := range valores {
		index := int((valor - inherentRiskMin) / binWidth)
		if index >= numBins {
			index = numBins - 1
		}
		frequencias[index]++
	}

	fmt.Println("Frequências dos bins:")
	binData := make([]map[string]interface{}, numBins)
	for i, freq := range frequencias {
		lowerBound := inherentRiskMin + float64(i)*binWidth
		upperBound := lowerBound + binWidth
		fmt.Printf("Bin %d: [%.2f - %.2f], Frequência: %d\n", i, lowerBound, upperBound, freq)
		midPoint := (lowerBound + upperBound) / 2
		binData[i] = map[string]interface{}{
			"midPoint":  midPoint,
			"frequency": freq,
		}
	}

	p := plot.New()
	p.Title.Text = "Distribuição PERT de Perdas Financeiras"
	p.X.Label.Text = "Perdas (R$)"
	p.Y.Label.Text = "Frequência de Aparição"

	hist, err := plotter.NewHist(plotter.Values(valores), numBins)
	if err != nil {
		panic(err)
	}
	hist.Normalize(1)
	p.Add(hist)

	histPath := "hist.png"
	if err := p.Save(12*vg.Inch, 6*vg.Inch, histPath); err != nil {
		panic(err)
	}

	sort.Float64s(valores)
	pLEC := plot.New()
	pLEC.Title.Text = "Curva de Excedência de Perdas (LEC)"
	pLEC.X.Label.Text = "Perdas (R$)"
	pLEC.Y.Label.Text = "Probabilidade de Excedência"

	lec := make(plotter.XYs, amostras)
	for i := range lec {
		lec[i].X = valores[i]
		lec[i].Y = 1 - float64(i)/float64(amostras)
	}

	line, err := plotter.NewLine(lec)
	if err != nil {
		panic(err)
	}
	pLEC.Add(line)
	pLEC.Add(plotter.NewGrid())

	lecPath := "lec.png"
	if err := pLEC.Save(12*vg.Inch, 6*vg.Inch, lecPath); err != nil {
		panic(err)
	}

	fmt.Println("Plots saved as hist.png and lec.png")

	// Excluir os arquivos
	os.Remove(histPath)
	os.Remove(lecPath)

	return binData
}

func createReport(path string, inherentRiskMin, inherentRiskMax, inherentRiskMode, proposedRiskMin, proposedRiskMax, proposedRiskMode float64) error {
	report := fmt.Sprintf("Inherent Risk:\nMin: %f\nMax: %f\nMode: %f\n\nProposed Risk:\nMin: %f\nMax: %f\nMode: %f\n", inherentRiskMin, inherentRiskMax, inherentRiskMode, proposedRiskMin, proposedRiskMax, proposedRiskMode)

	return os.WriteFile(path, []byte(report), 0644)
}

func sendEmailWithAttachmentsReport(recipient, attachmentPath string) error {
	user := os.Getenv("EMAIL")
	pass := os.Getenv("EPASSWORD")
	m := mail.NewMessage()
	m.SetHeader("From", user)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", "Simulation Report")
	m.SetBody("text/plain", "Please find attached the simulation report.")
	m.Attach(attachmentPath)

	d := mail.NewDialer("smtp.gmail.com", 587, user, pass)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	fmt.Println("Email sent successfully")
	return nil
}
