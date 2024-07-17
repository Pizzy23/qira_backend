package db

type RelevanceDinamic struct {
	ID        int64 `json:"id" xorm:"pk autoincr 'id' INT"`
	ControlID int   `json:"controlId" xorm:"INT notnull"`
    Attacktxt string `json:"Attacktxt" xorm:"VARCHAR(255) notnull"`
    yourbaseIsattack string `json:"yourbaseIsattack" xorm:"VARCHAR(255) notnull"`
}

type ControlDinamic struct {
	ID              int64  `json:"id" xorm:"pk autoincr 'id' INT"`
	ControlID       int    `json:"controlId" xorm:"INT notnull"`
	AggregateTable  string `json:"aggregateTable" xorm:"VARCHAR(255) notnull"`
	Aggregate       string `json:"aggregate" xorm:"VARCHAR(255) notnull"`
	ControlGapTable string `json:"controlGapTable" xorm:"VARCHAR(255) notnull"`
	ControlGap      string `json:"controlGap" xorm:"VARCHAR(255) notnull"`
    Attacktxt string `json:"Attacktxt" xorm:"VARCHAR(255) notnull"`
    yourbaseIsattack string `json:"yourbaseIsattack" xorm:"VARCHAR(255) notnull"`
}

type PropusedDinamic struct {
	ID              int64  `json:"id" xorm:"pk autoincr 'id' INT"`
	ControlID       int    `json:"controlId" xorm:"INT notnull"`
	AggregateTable  string `json:"aggregateTable" xorm:"VARCHAR(255) notnull"`
	Aggregate       string `json:"aggregate" xorm:"VARCHAR(255) notnull"`
	ControlGapTable string `json:"controlGapTable" xorm:"VARCHAR(255) notnull"`
	ControlGap      string `json:"controlGap" xorm:"VARCHAR(255) notnull"`
    Attacktxt string `json:"Attacktxt" xorm:"VARCHAR(255) notnull"`
    yourbaseIsattack string `json:"yourbaseIsattack" xorm:"VARCHAR(255) notnull"`
}

type RelevanceDinamicInput struct {
	ControlID int `json:"controlId" xorm:"INT notnull"`
    Attacktxt string `json:"Attacktxt" xorm:"VARCHAR(255) notnull"`
    yourbaseIsattack string `json:"yourbaseIsattack" xorm:"VARCHAR(255) notnull"`
}

type ControlDinamicInput struct {
	ControlID int `json:"controlId" xorm:"INT notnull"`
    Attacktxt string `json:"Attacktxt" xorm:"VARCHAR(255) notnull"`
    yourbaseIsattack string `json:"yourbaseIsattack" xorm:"VARCHAR(255) notnull"`
}

type PropusedDinamicInput struct {
	ControlID int `json:"controlId" xorm:"INT notnull"`
    Attacktxt string `json:"Attacktxt" xorm:"VARCHAR(255) notnull"`
    yourbaseIsattack string `json:"yourbaseIsattack" xorm:"VARCHAR(255) notnull"`
}
