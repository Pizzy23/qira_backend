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
	"xorm.io/xorm"
)

func MonteCarloSimulationAggregated(c *gin.Context, threatEvent string, reciverEmail string) {
	var losses []db.LossHighTotal

	engine, exists := c.Get("db")
	if !exists {
		c.Set("Response", "Database connection not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	err := engine.(*xorm.Engine).Where("threat_event = ?", threatEvent).Find(&losses)
	if err != nil {
		c.Set("Response", "LossHigh not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	if len(losses) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("No loss data found for threatEvent '%s'", threatEvent),
		})
		return
	}

	// Calculando os valores agregados
	var totalMinimo, totalMaximo, totalMaisProvavel float64
	for _, loss := range losses {
		totalMinimo += loss.MinimumLoss
		totalMaximo += loss.MaximumLoss
		totalMaisProvavel += loss.MostLikelyLoss
	}

	// Validação dos valores
	if totalMinimo == 0 && totalMaximo == 0 && totalMaisProvavel == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Sua perda é igual a 0 no threatEvent '%s', por favor coloque valores válidos.", threatEvent),
		})
		return
	}

	const (
		amostras = 10000 // valores repeticao
		numBins  = 70    // setado
	)

	binWidth := (totalMaximo - totalMinimo) / float64(numBins)

	a := (4*totalMaisProvavel + totalMaximo - 5*totalMinimo) / (totalMaximo - totalMinimo)
	b := (5*totalMaximo - totalMinimo - 4*totalMaisProvavel) / (totalMaximo - totalMinimo)

	source := rand.NewSource(uint64(99))
	distribuicao := distuv.Beta{
		Alpha: a,
		Beta:  b,
		Src:   rand.New(source),
	}
	valores := make([]float64, amostras)
	for i := range valores {
		valores[i] = distribuicao.Rand()*(totalMaximo-totalMinimo) + totalMinimo
	}

	frequencias := make([]int, numBins)
	for _, valor := range valores {
		index := int((valor - totalMinimo) / binWidth)
		if index >= numBins {
			index = numBins - 1
		}
		frequencias[index]++
	}

	fmt.Println("Frequências dos bins:")
	binData := make([]map[string]interface{}, numBins)
	for i, freq := range frequencias {
		lowerBound := totalMinimo + float64(i)*binWidth
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
