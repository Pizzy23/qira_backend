package simulation

import (
	"fmt"
	"net/http"
	"os"
	"qira/db"
	"sort"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gopkg.in/mail.v2"
	"xorm.io/xorm"
)

func MonteCarloSimulationReport(c *gin.Context, threatEvent string, receiverEmail string) {
	var loss db.LossHighTotal
	var frequencies []db.Frequency
	var riskCalculations []db.RiskCalculation

	engine, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}

	found, err := engine.(*xorm.Engine).Where("threat_event = ?", threatEvent).Get(&loss)
	if err != nil || !found {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "LossHigh not found"})
		return
	}

	err = engine.(*xorm.Engine).Where("threat_event = ?", threatEvent).Find(&frequencies)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving frequencies"})
		return
	}

	err = engine.(*xorm.Engine).Where("threat_event = ?", threatEvent).Find(&riskCalculations)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving risk calculations"})
		return
	}

	// Generate the Monte Carlo Simulation
	generateMonteCarloSimulation(c, loss)

	// Combine all information into a report
	reportPath := "report.txt"
	err = createReport(reportPath, loss, frequencies, riskCalculations)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating report"})
		return
	}

	// Send the report by email
	sendEmailWithAttachmentsReport(receiverEmail, reportPath)

	// Delete the report after sending
	err = os.Remove(reportPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting report"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Simulation report generated and sent via email successfully!"})
}

func generateMonteCarloSimulation(c *gin.Context, loss db.LossHighTotal) {
	// Usando os valores do banco de dados
	minimo := loss.MinimumLoss
	maximo := loss.MaximumLoss
	maisProvavel := loss.MostLikelyLoss

	// Validação dos valores
	if minimo == 0 && maximo == 0 && maisProvavel == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Sua perda é igual a 0 no threatEvent '%s', por favor coloque valores válidos.", loss.ThreatEvent),
		})
		return
	}

	const (
		amostras = 10000 // valores repeticao
		numBins  = 70    // setado
	)

	binWidth := (maximo - minimo) / float64(numBins)

	a := (4*maisProvavel + maximo - 5*minimo) / (maximo - minimo)
	b := (5*maximo - minimo - 4*maisProvavel) / (maximo - minimo)

	source := rand.NewSource(uint64(99))
	distribuicao := distuv.Beta{
		Alpha: a,
		Beta:  b,
		Src:   rand.New(source),
	}
	valores := make([]float64, amostras)
	for i := range valores {
		valores[i] = distribuicao.Rand()*(maximo-minimo) + minimo
	}

	frequencias := make([]int, numBins)
	for _, valor := range valores {
		index := int((valor - minimo) / binWidth)
		if index >= numBins {
			index = numBins - 1
		}
		frequencias[index]++
	}

	fmt.Println("Frequências dos bins:")
	binData := make([]map[string]interface{}, numBins)
	for i, freq := range frequencias {
		lowerBound := minimo + float64(i)*binWidth
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

	// Enviar os arquivos por email
	sendEmailWithAttachments("recipient@example.com", histPath, lecPath)

	// Excluir os arquivos
	os.Remove(histPath)
	os.Remove(lecPath)
}

func createReport(path string, loss db.LossHighTotal, frequencies []db.Frequency, riskCalculations []db.RiskCalculation) error {
	report := fmt.Sprintf("Loss High Total:\nMinimum Loss: %f\nMaximum Loss: %f\nMost Likely Loss: %f\n\nFrequencies:\n", loss.MinimumLoss, loss.MaximumLoss, loss.MostLikelyLoss)

	for _, freq := range frequencies {
		report += fmt.Sprintf("Threat Event: %s, Min Frequency: %f, Max Frequency: %f, Most Likely Frequency: %f\n", freq.ThreatEvent, freq.MinFrequency, freq.MaxFrequency, freq.MostLikelyFrequency)
	}

	report += "\nRisk Calculations:\n"
	for _, risk := range riskCalculations {
		report += fmt.Sprintf("Threat Event: %s, Min: %f, Max: %f, Mode: %f, Estimate: %f\n", risk.ThreatEvent, risk.Min, risk.Max, risk.Mode, risk.Estimate)
	}

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
