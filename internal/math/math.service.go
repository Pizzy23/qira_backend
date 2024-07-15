package calculations

func CalcRisks(A float64, B float64, C float64) float64 {
	return (A + 4*B + C) / 6
}

func CalcLoss(A int, B int) int {
	return A + B
}
