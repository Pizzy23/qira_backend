package db

type Relevance struct {
	ID           int64  `json:"id" xorm:"pk autoincr 'id'"`
	ControlID    int64  `json:"controlId" xorm:"INT 'control_id'"`
	Information  string `xorm:"TEXT 'information'"`
	TypeOfAttack string `json:"type_of_attack" xorm:"VARCHAR(255) 'type_of_attack'"`
	Porcent      int64  `json:"porcent" xorm:"INT 'porcent'"`
}

type Control struct { //Control Strength
	ID           int64  `json:"id" xorm:"pk autoincr 'id'"`
	ControlID    int64  `json:"controlId" xorm:"INT 'control_id'"`
	TypeOfAttack string `json:"type_of_attack" xorm:"VARCHAR(255) 'type_of_attack'"`
	Porcent      string `json:"porcent" xorm:"VARCHAR(255) 'porcent'"`
	Aggregate    string `json:"aggregate" xorm:"VARCHAR(255) 'aggregate'"`
	ControlGap   string `json:"controlGap" xorm:"VARCHAR(255) 'control_gap'"`
}

type Propused struct { //Propused Strength
	ID           int64  `json:"id" xorm:"pk autoincr 'id'"`
	ControlID    int64  `json:"controlId" xorm:"INT 'control_id'"`
	TypeOfAttack string `json:"type_of_attack" xorm:"VARCHAR(255) 'type_of_attack'"`
	Porcent      string `json:"porcent" xorm:"VARCHAR(255) 'porcent'"`
	Aggregate    string `json:"aggregate" xorm:"VARCHAR(255) 'aggregate'"`
	ControlGap   string `json:"controlGap" xorm:"VARCHAR(255) 'control_gap'"`
}

type RelevanceDinamicInput struct {
	ControlID    int64  `json:"controlId" xorm:"INT 'control_id'"`
	TypeOfAttack string `json:"type_of_attack" xorm:"VARCHAR(255) 'type_of_attack'"`
	Porcent      int64  `json:"porcent" xorm:"INT 'porcent'"`
}

type ControlDinamicInput struct {
	ControlID    int64  `json:"controlId" xorm:"INT 'control_id'"`
	TypeOfAttack string `json:"type_of_attack" xorm:"VARCHAR(255) 'type_of_attack'"`
	Porcent      string `json:"porcent" xorm:"VARCHAR(255) 'porcent'"`
}

type PropusedDinamicInput struct {
	ControlID    int64  `json:"controlId" xorm:"INT 'control_id'"`
	TypeOfAttack string `json:"type_of_attack" xorm:"VARCHAR(255) 'type_of_attack'"`
	Porcent      string `json:"porcent" xorm:"VARCHAR(255) 'porcent'"`
}
