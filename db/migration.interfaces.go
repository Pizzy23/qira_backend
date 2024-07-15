package db

type AssetInventory struct {
	ID                      int64   `xorm:"pk autoincr"`
	Name                    string  `xorm:"VARCHAR(255)"`
	Description             string  `xorm:"TEXT"`
	Location                string  `xorm:"VARCHAR(255)"`
	Responsible             string  `xorm:"VARCHAR(255)"`
	BusinessValue           int     `xorm:"INT"`
	ReplacementCost         float64 `xorm:"FLOAT"`
	Criticality             string  `xorm:"VARCHAR(255)"`
	Users                   string  `xorm:"VARCHAR(255)"`
	RoleInTargetEnvironment string  `xorm:"VARCHAR(255)"`
}

type ThreatEventCatalog struct {
	ID          int64  `xorm:"pk autoincr"`
	ThreatGroup string `xorm:"VARCHAR(255)"`
	ThreatEvent string `xorm:"VARCHAR(255)"`
	Description string `xorm:"TEXT"`
	InScope     bool   `xorm:"BOOL"`
}

type Frequency struct {
	ID                    int64   `xorm:"pk autoincr"`
	ThreatEventID         int64   `xorm:"INT"`
	ThreatEvent           string  `xorm:"VARCHAR(255)"`
	MinFrequency          float64 `xorm:"FLOAT"`
	MaxFrequency          float64 `xorm:"FLOAT"`
	MostLikelyFrequency   float64 `xorm:"FLOAT"`
	SupportingInformation string  `xorm:"TEXT"`
}

type LinkThreat struct {
	ID            int64  `xorm:"pk autoincr"`
	ThreatEventID int64  `xorm:"INT"`
	ThreatEvent   string `xorm:"VARCHAR(255)"`
	Assets        string `xorm:"TEXT"`
}

type LossHigh struct {
	ID             int64   `xorm:"pk autoincr"`
	ThreatEventID  int64   `xorm:"INT"`
	ThreatEvent    string  `xorm:"VARCHAR(255)"`
	Assets         string  `xorm:"TEXT"`
	LossType       string  `xorm:"VARCHAR(255)"`
	MinimumLoss    float64 `xorm:"FLOAT"`
	MaximumLoss    float64 `xorm:"FLOAT"`
	MostLikelyLoss float64 `xorm:"FLOAT"`
}

type RiskCalculation struct {
	ID            int64   `xorm:"pk autoincr"`
	ThreatEventID int64   `xorm:"INT"`
	ThreatEvent   string  `xorm:"VARCHAR(255)"`
	RiskType      string  `xorm:"VARCHAR(255)"`
	Min           float64 `xorm:"INT"`
	Max           float64 `xorm:"INT"`
	Mode          float64 `xorm:"INT"`
	Estimate      float64 `xorm:"FLOAT"`
}

type ControlLibrary struct {
	ID               int64  `xorm:"pk autoincr"`
	ControlType      string `xorm:"VARCHAR(255)"`
	ControlReference string `xorm:"VARCHAR(255)"`
	Information      string `xorm:"TEXT"`
	InScope          bool   `xorm:"BOOL"`
}

type Relevance struct { //Remover essa tabela
	ID        int `json:"id" xorm:"pk autoincr 'id' INT"`
	ControlID int `json:"controlId" xorm:"INT notnull"`
}

type Strength struct { //Remover essa tabela
	ID        int `json:"id" xorm:"pk autoincr 'id' INT"`
	ControlID int `json:"controlId" xorm:"INT notnull"`
}

type Propused struct { //Remover essa tabela
	ID        int `json:"id" xorm:"pk autoincr 'id' INT"`
	ControlID int `json:"controlId" xorm:"INT notnull"`
}

type RiskController struct {
	ID   int    `json:"id" xorm:"pk autoincr 'id' INT"`
	Name string `json:"name" xorm:"VARCHAR(255) notnull"`
}

type ThreatEventAssets struct { // Vai selecionar Assets, Oque foi afetado Pode ser varios afetados
	ID            int    `json:"id" xorm:"pk autoincr 'id' INT"`
	ThreatID      int    `json:"threat_id" xorm:"INT notnull"`
	ThreatEvent   string `json:"threat_event" xorm:"VARCHAR(255) notnull"`
	AffectedAsset string `json:"affected_asset" xorm:"JSON notnull"`
}

type Implements struct {
	ID              int    `json:"id" xorm:"pk autoincr 'id' INT"`
	ControlID       int    `json:"controlId" xorm:"INT notnull"`
	Current         int    `json:"current" xorm:"INT notnull"`
	Proposed        int    `json:"proposed" xorm:"v notnull"`
	PercentCurrent  string `json:"percentCurrent" xorm:"VARCHAR(255) notnull"`
	PercentProposed string `json:"percentProposed" xorm:"VARCHAR(255) notnull"`
	Cost            int    `json:"cost" xorm:"INT notnull"`
}

type AggregatedStrength struct {
	ID          int    `json:"id" xorm:"pk autoincr 'id' INT"`
	ThreatID    int    `json:"threat_id" xorm:"INT notnull"`
	ThreatEvent string `json:"threat_event" xorm:"VARCHAR(255) notnull"`
	Current     string `json:"current" xorm:"VARCHAR(255) notnull"`
	Proposed    string `json:"proposed" xorm:"VARCHAR(255) notnull"`
}
