package db

type AssetInventory struct {
	ID                      int64   `json:"id" xorm:"pk autoincr 'id'"`
	Name                    string  `json:"name" xorm:"VARCHAR(255) 'name'"`
	Description             string  `json:"description" xorm:"TEXT 'description'"`
	Location                string  `json:"location" xorm:"VARCHAR(255) 'location'"`
	Responsible             string  `json:"responsible" xorm:"VARCHAR(255) 'responsible'"`
	BusinessValue           float64 `json:"business_value" xorm:"FLOAT 'business_value'"`
	ReplacementCost         float64 `json:"replacement_cost" xorm:"FLOAT 'replacement_cost'"`
	Criticality             string  `json:"criticality" xorm:"VARCHAR(255) 'criticality'"`
	Users                   string  `json:"users" xorm:"VARCHAR(255) 'users'"`
	RoleInTargetEnvironment string  `json:"roleInTargetEnvironment" xorm:"VARCHAR(255) 'role_in_target_environment'"`
}

type ThreatEventCatalog struct {
	ID          int64  `xorm:"pk autoincr 'id'"`
	ThreatGroup string `xorm:"VARCHAR(255) 'threat_group'"`
	ThreatEvent string `xorm:"VARCHAR(255) 'threat_event'"`
	Description string `xorm:"TEXT 'description'"`
	InScope     bool   `xorm:"BOOL 'in_scope'"`
}

type Frequency struct {
	ID                    int64   `xorm:"pk autoincr 'id'"`
	ThreatEventID         int64   `xorm:"INT 'threat_event_id'"`
	ThreatEvent           string  `xorm:"VARCHAR(255) 'threat_event'"`
	MinFrequency          float64 `xorm:"FLOAT 'min_frequency'"`
	MaxFrequency          float64 `xorm:"FLOAT 'max_frequency'"`
	MostLikelyFrequency   float64 `xorm:"FLOAT 'most_likely_frequency'"`
	SupportingInformation string  `xorm:"TEXT 'supporting_information'"`
}

type LinkThreat struct {
	ID            int64  `xorm:"pk autoincr 'id'"`
	ThreatEventID int64  `xorm:"INT 'threat_event_id'"`
	ThreatEvent   string `xorm:"VARCHAR(255) 'threat_event'"`
	Assets        string `xorm:"TEXT 'assets'"`
}

type LossHigh struct {
	ID             int64   `xorm:"pk autoincr 'id'"`
	ThreatEventID  int64   `xorm:"INT 'threat_event_id'"`
	ThreatEvent    string  `xorm:"VARCHAR(255) 'threat_event'"`
	Assets         string  `xorm:"TEXT 'assets'"`
	LossType       string  `xorm:"VARCHAR(255) 'loss_type'"`
	MinimumLoss    float64 `xorm:"FLOAT 'minimum_loss'"`
	MaximumLoss    float64 `xorm:"FLOAT 'maximum_loss'"`
	MostLikelyLoss float64 `xorm:"FLOAT 'most_likely_loss'"`
}

type LossHighGranular struct {
	ID             int64   `xorm:"pk autoincr 'id'"`
	ThreatEventID  int64   `xorm:"INT 'threat_event_id'"`
	ThreatEvent    string  `xorm:"VARCHAR(255) 'threat_event'"`
	Assets         string  `xorm:"TEXT 'assets'"`
	LossType       string  `xorm:"VARCHAR(255) 'loss_type'"`
	Impact         string  `xorm:"VARCHAR(255) 'impact'"`
	MinimumLoss    float64 `xorm:"FLOAT 'minimum_loss'"`
	MaximumLoss    float64 `xorm:"FLOAT 'maximum_loss'"`
	MostLikelyLoss float64 `xorm:"FLOAT 'most_likely_loss'"`
}

type LossHighTotal struct {
	ID             int64   `xorm:"pk autoincr 'id'"`
	ThreatEventID  int64   `xorm:"INT 'threat_event_id'"`
	ThreatEvent    string  `xorm:"VARCHAR(255) 'threat_event'"`
	Name           string  `xorm:"TEXT 'name'"`
	TypeOfLoss     string  `xorm:"VARCHAR(255) 'type_of_loss'"`
	MinimumLoss    float64 `xorm:"FLOAT 'minimum_loss'"`
	MaximumLoss    float64 `xorm:"FLOAT 'maximum_loss'"`
	MostLikelyLoss float64 `xorm:"FLOAT 'most_likely_loss'"`
}

type RiskCalculation struct {
	ID            int64   `json:"id" xorm:"pk autoincr 'id'"`
	ThreatEventID int64   `json:"threat_event_id" xorm:"INT 'threat_event_id'"`
	ThreatEvent   string  `json:"threat_event" xorm:"VARCHAR(255) 'threat_event'"`
	Categorie     string  `json:"categorie" xorm:"VARCHAR(255) 'categorie'"`
	RiskType      string  `json:"risk_type" xorm:"VARCHAR(255) 'risk_type'"`
	Min           float64 `json:"min" xorm:"INT 'min'"`
	Max           float64 `json:"max" xorm:"INT 'max'"`
	Mode          float64 `json:"mode" xorm:"INT 'mode'"`
	Estimate      float64 `json:"estimate" xorm:"FLOAT 'estimate'"`
}

type ControlLibrary struct {
	ID               int64  `xorm:"pk autoincr 'id'"`
	ControlType      string `xorm:"VARCHAR(255) 'control_type'"`
	ControlReference string `xorm:"VARCHAR(255) 'control_reference'"`
	Information      string `xorm:"TEXT 'information'"`
	InScope          bool   `xorm:"BOOL 'in_scope'"`
}

type RiskController struct {
	ID   int    `json:"id" xorm:"pk autoincr 'id'"`
	Name string `json:"name" xorm:"VARCHAR(255) 'name'"`
}

type ThreatEventAssets struct {
	ID            int    `json:"id" xorm:"pk autoincr 'id'"`
	ThreatID      int64  `json:"threat_id" xorm:"INT 'threat_id'"`
	ThreatEvent   string `json:"threat_event" xorm:"VARCHAR(255) 'threat_event'"`
	AffectedAsset string `json:"affected_asset" xorm:"JSON 'affected_asset'"`
}

type Implements struct {
	ID              int    `json:"id" xorm:"pk autoincr 'id'"`
	ControlID       int64  `json:"controlId" xorm:"INT 'control_id'"`
	Current         int    `json:"current" xorm:"INT 'current'"`
	Proposed        int    `json:"proposed" xorm:"INT 'proposed'"`
	PercentCurrent  string `json:"percentCurrent" xorm:"VARCHAR(255) 'percent_current'"`
	PercentProposed string `json:"percentProposed" xorm:"VARCHAR(255) 'percent_proposed'"`
	Cost            int    `json:"cost" xorm:"INT 'cost'"`
}

type AggregatedStrength struct {
	ID          int    `json:"id" xorm:"pk autoincr 'id'"`
	ThreatID    int64  `json:"threat_id" xorm:"INT 'threat_id'"`
	ThreatEvent string `json:"threat_event" xorm:"VARCHAR(255) 'threat_event'"`
	Current     string `json:"current" xorm:"VARCHAR(255) 'current'"`
	Proposed    string `json:"proposed" xorm:"VARCHAR(255) 'proposed'"`
}

type LossExceedance struct {
	Risk float64 `json:"risk" xorm:"VARCHAR(255) 'risk'"`
	Loss int64   `json:"loss" xorm:"INT 'loss'"`
}
