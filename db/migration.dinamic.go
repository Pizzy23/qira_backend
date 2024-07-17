package db

type Relevance struct {
	ID           int64  `json:"id" xorm:"pk autoincr 'id' INT"`
	ControlID    int64  `json:"controlId" xorm:"INT notnull"`
	TypeOfAttack string `json:"type_of_attack" xorm:"VARCHAR(255)"`
	Porcent      int64  `json:"porcent" xorm:"INT notnull"`
}

type Control struct {
	ID              int64  `json:"id" xorm:"pk autoincr 'id' INT"`
	ControlID       int64  `json:"controlId" xorm:"INT notnull"`
	TypeOfAttack    string `json:"type_of_attack" xorm:"VARCHAR(255)"`
	Porcent         string `json:"porcent" xorm:"VARCHAR(255) notnull"`
	AggregateTable  string `json:"aggregateTable" xorm:"VARCHAR(255) notnull"`
	Aggregate       string `json:"aggregate" xorm:"VARCHAR(255) notnull"`
	ControlGapTable string `json:"controlGapTable" xorm:"VARCHAR(255) notnull"`
	ControlGap      string `json:"controlGap" xorm:"VARCHAR(255) notnull"`
}

type Propused struct {
	ID              int64  `json:"id" xorm:"pk autoincr 'id' INT"`
	ControlID       int64  `json:"controlId" xorm:"INT notnull"`
	TypeOfAttack    string `json:"type_of_attack" xorm:"VARCHAR(255)"`
	Porcent         string `json:"porcent" xorm:"VARCHAR(255) notnull"`
	AggregateTable  string `json:"aggregateTable" xorm:"VARCHAR(255) notnull"`
	Aggregate       string `json:"aggregate" xorm:"VARCHAR(255) notnull"`
	ControlGapTable string `json:"controlGapTable" xorm:"VARCHAR(255) notnull"`
	ControlGap      string `json:"controlGap" xorm:"VARCHAR(255) notnull"`
}

type RelevanceDinamicInput struct {
	ControlID    int    `json:"controlId" xorm:"INT notnull"`
	TypeOfAttack string `json:"type_of_attack" xorm:"VARCHAR(255)"`
	Porcent      string `json:"porcent" xorm:"VARCHAR(255)notnull"`
}

type ControlDinamicInput struct {
	ControlID    int    `json:"controlId" xorm:"INT notnull"`
	TypeOfAttack string `json:"type_of_attack" xorm:"VARCHAR(255)"`
	Porcent      string `json:"porcent" xorm:"VARCHAR(255)notnull"`
}

type PropusedDinamicInput struct {
	ControlID    int    `json:"controlId" xorm:"INT notnull"`
	TypeOfAttack string `json:"type_of_attack" xorm:"VARCHAR(255)"`
	Porcent      string `json:"porcent" xorm:"VARCHAR(255)notnull"`
}
