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

func MonteCarloSimulation(c *gin.Context, threatEvent string, reciverEmail string) {
	var loss db.LossHighTotal

	engine, exists := c.Get("db")
	if !exists {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	found, err := engine.(*xorm.Engine).Where("threat_event = ?", threatEvent).Get(&loss)
	if err != nil || !found {
		c.Set("Response", "LossHigh not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	// Usando os valores do banco de dados
	minimo := loss.MinimumLoss
	maximo := loss.MaximumLoss
	maisProvavel := loss.MostLikelyLoss

	// Validação dos valores
	if minimo == 0 && maximo == 0 && maisProvavel == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Sua perda é igual a 0 no threatEvent '%s', por favor coloque valores válidos.", threatEvent),
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
	if reciverEmail != "" {
		sendEmailWithAttachments(reciverEmail, histPath, lecPath)
	}

	// Excluir os arquivos
	os.Remove(histPath)
	os.Remove(lecPath)

	c.JSON(200, gin.H{
		"bins": binData,
	})
}

func sendEmailWithAttachments(recipient, histPath, lecPath string) {
	user := os.Getenv("EMAIL")
	pass := os.Getenv("EPASSWORD")
	m := mail.NewMessage()
	m.SetHeader("From", user)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", "Monte Carlo Simulation Results")
	m.SetBody("text/plain", "Please find attached the results of the Monte Carlo Simulation.")
	m.Attach(histPath)
	m.Attach(lecPath)

	d := mail.NewDialer("smtp.gmail.com", 587, user, pass)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	fmt.Println("Email sent successfully")
}
