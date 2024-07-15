package db

type RelevanceDinamic struct {
	ID        int `json:"id" xorm:"pk autoincr 'id' INT"`
	ControlID int `json:"controlId" xorm:"INT notnull"`
}

type StrengthDinamic struct {
	ID        int `json:"id" xorm:"pk autoincr 'id' INT"`
	ControlID int `json:"controlId" xorm:"INT notnull"`
}

type PropusedDinamic struct {
	ID        int `json:"id" xorm:"pk autoincr 'id' INT"`
	ControlID int `json:"controlId" xorm:"INT notnull"`
}

type RelevanceDinamicInput struct {
	ControlID int `json:"controlId" xorm:"INT notnull"`
}

type StrengthDinamicInput struct {
	ControlID int `json:"controlId" xorm:"INT notnull"`
}

type PropusedDinamicInput struct {
	ControlID int `json:"controlId" xorm:"INT notnull"`
}
