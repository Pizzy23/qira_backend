package simulation

import (
	"fmt"
	"image/color"
	"net/http"
	"os"
	"qira/db"
	"qira/internal/interfaces"
	"sort"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"xorm.io/xorm"
)

func MonteCarloSimulationAppetite(c *gin.Context, threatEvent string, reciverEmail string) {
	var losses []db.LossHighTotal
	var lossData []db.LossExceedance

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

	err = db.GetAll(engine.(*xorm.Engine), &lossData)
	if err != nil {
		c.Set("Response", "LossExceedance not found")
		c.Status(http.StatusInternalServerError)
		return
	}

	if len(losses) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("No loss data found for threatEvent '%s'", threatEvent),
		})
		return
	}

	var totalMinimo, totalMaximo, totalMaisProvavel float64
	for _, loss := range losses {
		totalMinimo += loss.MinimumLoss
		totalMaximo += loss.MaximumLoss
		totalMaisProvavel += loss.MostLikelyLoss
	}

	if totalMinimo == 0 && totalMaximo == 0 && totalMaisProvavel == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Sua perda é igual a 0 no threatEvent '%s', por favor coloque valores válidos.", threatEvent),
		})
		return
	}

	const (
		amostras = 10000
		numBins  = 70
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
	line.LineStyle.Color = color.RGBA{R: 255, A: 255} // Vermelho
	pLEC.Add(line)

	userLEC := make(plotter.XYs, len(lossData))
	for i, ld := range lossData {
		userLEC[i].X = float64(ld.Loss)
		userLEC[i].Y = parseRiskToFloat(ld.Risk)
	}

	scatter, err := plotter.NewScatter(userLEC)
	if err != nil {
		panic(err)
	}
	scatter.GlyphStyle.Shape = draw.CircleGlyph{}
	scatter.GlyphStyle.Radius = vg.Points(2)
	scatter.GlyphStyle.Color = color.RGBA{B: 255, A: 255} // Azul

	pLEC.Add(scatter)
	pLEC.Add(plotter.NewGrid())

	lecPath := "lec.png"
	if err := pLEC.Save(12*vg.Inch, 6*vg.Inch, lecPath); err != nil {
		panic(err)
	}

	fmt.Println("Plots saved as hist.png and lec.png")

	if reciverEmail != "" {
		//sendEmailWithAttachments(reciverEmail, histPath, lecPath)
	}

	// Excluir os arquivos
	os.Remove(histPath)
	os.Remove(lecPath)

	c.JSON(200, gin.H{
		"bins": binData,
	})
}

func parseRiskToFloat(risk string) float64 {
	switch risk {
	case "100%":
		return 1.0
	case "75%":
		return 0.75
	case "50%":
		return 0.5
	case "25%":
		return 0.25
	case "0%":
		return 0.0
	default:
		return 0.0
	}
}

func UploadLossData(c *gin.Context, lossData []interfaces.LossExceedance) {

	engine, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection not found",
		})
		return
	}

	for _, ld := range lossData {
		existing := &db.LossExceedance{}
		has, err := engine.(*xorm.Engine).Where("risk = ? AND loss = ?", ld.Risk, ld.Loss).Get(existing)
		if err == nil && !has {
			newLoss := db.LossExceedance{
				Risk: ld.Risk,
				Loss: ld.Loss,
			}
			_, err := engine.(*xorm.Engine).Insert(newLoss)
			if err != nil {
				fmt.Printf("Error inserting loss data: %v\n", err)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Loss data uploaded successfully",
	})
}
