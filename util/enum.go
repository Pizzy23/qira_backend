package util

func EnumLoss(loss string) bool {
	switch loss {
	case "Singular":
		return true
	case "LossHigh":
		return true
	case "Granular":
		return true
	default:
		return false
	}
}
