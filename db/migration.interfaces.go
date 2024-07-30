package db

type AssetInventory struct {
	ID                      int64   `json:"id" xorm:"pk autoincr"`
	Name                    string  `json:"name" xorm:"VARCHAR(255)"`
	Description             string  `json:"description" xorm:"TEXT"`
	Location                string  `json:"location" xorm:"VARCHAR(255)"`
	Responsible             string  `json:"responsible" xorm:"VARCHAR(255)"`
	BusinessValue           float64 `json:"business_value" xorm:"FLOAT"`
	ReplacementCost         float64 `json:"replacement_cost" xorm:"FLOAT"`
	Criticality             string  `json:"criticality" xorm:"VARCHAR(255)"`
	Users                   string  `json:"users" xorm:"VARCHAR(255)"`
	RoleInTargetEnvironment string  `json:"roleInTargetEnvironment" xorm:"VARCHAR(255)"`
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
type LossHighGranular struct {
	ID             int64   `xorm:"pk autoincr"`
	ThreatEventID  int64   `xorm:"INT"`
	ThreatEvent    string  `xorm:"VARCHAR(255)"`
	Assets         string  `xorm:"TEXT"`
	LossType       string  `xorm:"VARCHAR(255)"`
	Impact         string  `xorm:"VARCHAR(255)"`
	MinimumLoss    float64 `xorm:"FLOAT"`
	MaximumLoss    float64 `xorm:"FLOAT"`
	MostLikelyLoss float64 `xorm:"FLOAT"`
}
type LossHighTotal struct { //Isso seria o total da Alto-Custo que é a soma de minimo, maximo e mode.
	ID             int64   `xorm:"pk autoincr"`
	ThreatEventID  int64   `xorm:"INT"`
	ThreatEvent    string  `xorm:"VARCHAR(255)"`
	Name           string  `xorm:"TEXT"`
	TypeOfLoss     string  `xorm:"VARCHAR(255)"`
	MinimumLoss    float64 `xorm:"FLOAT"`
	MaximumLoss    float64 `xorm:"FLOAT"`
	MostLikelyLoss float64 `xorm:"FLOAT"`
}

type RiskCalculation struct {
	ID            int64   `json:"id" xorm:"pk autoincr"`
	ThreatEventID int64   `json:"threat_event_id" xorm:"INT"`
	ThreatEvent   string  `json:"threat_event" xorm:"VARCHAR(255)"`
	RiskType      string  `json:"risk_type" xorm:"VARCHAR(255)"` //RiskType pode ser "risk", "loss" ou "Frequencia"
	Min           float64 `json:"min" xorm:"INT"`
	Max           float64 `json:"max" xorm:"INT"`
	Mode          float64 `json:"mode" xorm:"INT"`
	Estimate      float64 `json:"estimate" xorm:"FLOAT"`
}

type ControlLibrary struct {
	ID               int64  `xorm:"pk autoincr"`
	ControlType      string `xorm:"VARCHAR(255)"`
	ControlReference string `xorm:"VARCHAR(255)"`
	Information      string `xorm:"TEXT"`
	InScope          bool   `xorm:"BOOL"`
}

type RiskController struct {
	ID   int    `json:"id" xorm:"pk autoincr 'id' INT"`
	Name string `json:"name" xorm:"VARCHAR(255) notnull"`
}

type ThreatEventAssets struct { // Vai selecionar Assets, Oque foi afetado Pode ser varios afetados
	ID            int    `json:"id" xorm:"pk autoincr 'id' INT"`
	ThreatID      int64  `json:"threat_id" xorm:"INT notnull"`
	ThreatEvent   string `json:"threat_event" xorm:"VARCHAR(255) notnull"`
	AffectedAsset string `json:"affected_asset" xorm:"JSON notnull"`
}

type Implements struct {
	ID              int    `json:"id" xorm:"pk autoincr 'id' INT"`
	ControlID       int64  `json:"controlId" xorm:"INT notnull"`
	Current         int    `json:"current" xorm:"INT notnull"`
	Proposed        int    `json:"proposed" xorm:"v notnull"`
	PercentCurrent  string `json:"percentCurrent" xorm:"VARCHAR(255) notnull"`
	PercentProposed string `json:"percentProposed" xorm:"VARCHAR(255) notnull"`
	Cost            int    `json:"cost" xorm:"INT notnull"`
}

type AggregatedStrength struct {
	ID          int    `json:"id" xorm:"pk autoincr 'id' INT"`
	ThreatID    int64  `json:"threat_id" xorm:"INT notnull"`
	ThreatEvent string `json:"threat_event" xorm:"VARCHAR(255) notnull"`
	Current     string `json:"current" xorm:"VARCHAR(255) notnull"`
	Proposed    string `json:"proposed" xorm:"VARCHAR(255) notnull"`
}
type LossExceedance struct {
	Risk string `json:"risk" xorm:"VARCHAR(255) notnull"`
	Loss int64  `json:"loss" xorm:"INT notnull"`
}
