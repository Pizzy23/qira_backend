package interfaces

type AssetsInventory struct { //post get:all/id
	ID                int     `json:"id"`
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	Location          string  `json:"location"`
	Responsible       string  `json:"responsible"`
	BusinessValue     int     `json:"business_value"`
	ReplacementCost   float64 `json:"replacement_cost"`
	Criticality       string  `json:"criticality"`
	Users             string  `json:"users"`
	TargetEnvironment string  `json:"target_environment"`
}

type ThreatEventCatalogue struct { //Ameaças
	ID          int    `json:"id"`
	ThreatGroup string `json:"threat_group"`
	ThreatEvent string `json:"threat_event"`
	Description string `json:"description"`
	InScope     string `json:"in_scope"`
}

type Frequency struct { //Criar ameaça, ai cria a frequencia junto 0, Editar a frequencia]
	ID                  int    `json:"id"`
	ThreatEventID       int    `json:"threat_event_id"`
	ThreatEvent         string `json:"threat_event"`
	MinFrequency        int    `json:"min_frequency"`
	MaxFrequency        int    `json:"max_frequency"`
	MostCommonFrequency int    `json:"most_common_frequency"`
	SupportInformation  string `json:"support_information"`
}

type ThreatEventAssets struct { // Vai selecionar Assets, Oque foi afetado Pode ser varios afetados
	ID            int      `json:"id"`
	ThreatID      int      `json:"threat_id"`
	ThreatEvent   string   `json:"threat_event"`
	AffectedAsset []string `json:"affected_asset"`
}

type LossHigh struct { //Vinculado ao Events com maximo e min de perca.
	ID             int      `json:"id"`
	ThreatEventID  int      `json:"threat_event_id"`
	ThreatEvent    string   `json:"threat_event"`
	Assets         []string `json:"assets"`
	LossType       string   `json:"loss_type"`
	MinimumLoss    float64  `json:"minimum_loss"`
	MaximumLoss    float64  `json:"maximum_loss"`
	MostLikelyLoss float64  `json:"most_likely_loss"`
}

type RiskCalculator struct { //so da input :)
	ID                int     `json:"id"`
	ThreatEventID     int     `json:"threat_event_id"`
	ThreatEvent       string  `json:"threat_event"`
	MinFrequency      int     `json:"min_frequency"`
	MaxFrequency      int     `json:"max_frequency"`
	ModeFrequency     int     `json:"mode_frequency"`
	EstimateFrequency float64 `json:"estimate_frequency"`
	MinLoss           int     `json:"min_loss"`
	MaxLoss           int     `json:"max_loss"`
	ModeLoss          int     `json:"mode_loss"`
	EstimateLoss      float64 `json:"estimate_loss"`
	MinRisk           int     `json:"min_risk"`
	MaxRisk           int     `json:"max_risk"`
	ModeRisk          int     `json:"mode_risk"`
	EstimateRisk      float64 `json:"estimate_risk"`
}

type ControlLibrary struct {
	ID               int    `json:"id"`
	ControlID        int    `json:"controlId"`
	ControlType      string `json:"control_type"`
	ControlReference string `json:"control_reference"`
	Information      string `json:"information"`
	InScope          string `json:"in_scope"`
}

type ControlImplementation struct {
	ControlID              int     `json:"controlId"`
	CurrentImplementation  int     `json:"current_implementation"`
	CurrentPercentValue    string  `json:"current_percent_value"`
	ProposedImplementation int     `json:"proposed_implementation"`
	ProposedPercentValue   string  `json:"proposed_percent_value"`
	ProjectedCost          float64 `json:"projected_cost"`
}

type AggregatedControlStrength struct { // COnversa com o event
	ID                      int    `json:"id"`
	ThreatEventID           int    `json:"threat_event_id"`
	ThreatEvent             string `json:"threat_event"`
	CurrentControlStrength  string `json:"current_control_strength"`
	ProposedControlStrength string `json:"proposed_control_strength"`
}

type InputAssetsInventory struct {
	Name                    string  `json:"name"`
	Description             string  `json:"description"`
	Location                string  `json:"location"`
	Responsible             string  `json:"responsible"`
	BusinessValue           float64 `json:"business_value"`
	ReplacementCost         float64 `json:"replacement_cost"`
	Criticality             string  `json:"criticality"`
	Users                   string  `json:"users"`
	RoleInTargetEnvironment string  `json:"roleInTargetEnvironment"`
}

type InputThreatEventCatalogue struct {
	ThreatGroup string `json:"threat_group"`
	ThreatEvent string `json:"threat_event"`
	Description string `json:"description"`
	InScope     bool   `json:"in_scope"`
}

type InputFrequency struct {
	ThreatEvent         string  `json:"threat_event"`
	MinFrequency        float64 `json:"min_frequency"`
	MaxFrequency        float64 `json:"max_frequency"`
	MostCommonFrequency float64 `json:"most_common_frequency"`
	SupportInformation  string  `json:"support_information"`
}

type InputThreatEventAssets struct {
	ThreatEvent   string   `json:"threat_event"`
	AffectedAsset []string `json:"affected_asset"`
}

type OutPutThreatEventAssets struct {
	ThreatID      int64    `json:"threat_id" `
	ThreatEvent   string   `json:"threat_event"`
	AffectedAsset []string `json:"affected_asset"`
}

type InputLossHigh struct {
	ThreatEvent    string  `json:"threat_event"`
	LossType       string  `json:"loss_type"`
	MinimumLoss    float64 `json:"minimum_loss"`
	MaximumLoss    float64 `json:"maximum_loss"`
	MostLikelyLoss float64 `json:"most_likely_loss"`
}
type InputLossHighGranulade struct {
	ThreatEvent    string  `json:"threat_event"`
	LossType       string  `json:"loss_type"`
	Impact         string  `json:"impact"`
	LossEditNumber int64   `json:loss_edit_number`
	MinimumLoss    float64 `json:"minimum_loss"`
	MaximumLoss    float64 `json:"maximum_loss"`
	MostLikelyLoss float64 `json:"most_likely_loss"`
}

type InputRiskCalculator struct {
	ThreatEventID int     `json:"threat_event_id"`
	ThreatEvent   string  `json:"threat_event"`
	RiskType      string  `json:"risk_type"`
	Min           int     `json:"min"`
	Max           int     `json:"max"`
	Mode          int     `json:"mode"`
	Estimate      float64 `json:"estimate"`
}

type InputControlLibrary struct {
	ControlType      string `json:"control_type"`
	ControlReference string `json:"control_reference"`
	Information      string `json:"information"`
	InScope          bool   `json:"in_scope"`
}

type InputControlImplementation struct {
	ControlID              int     `json:"controlId"`
	CurrentImplementation  int     `json:"current_implementation"`
	CurrentPercentValue    string  `json:"current_percent_value"`
	ProposedImplementation int     `json:"proposed_implementation"`
	ProposedPercentValue   string  `json:"proposed_percent_value"`
	ProjectedCost          float64 `json:"projected_cost"`
}

type InputAggregatedControlStrength struct {
	ThreatEventID           int    `json:"threat_event_id"`
	ThreatEvent             string `json:"threat_event"`
	CurrentControlStrength  string `json:"current_control_strength"`
	ProposedControlStrength string `json:"proposed_control_strength"`
}

type InputThreatEventAndAsset struct {
	Event InputThreatEventAssets `json:"event"`
	Asset InputAssetsInventory   `json:"Asset"`
}

type ThreatEventAndAsset struct {
	ThreatID    int    `json:"threat_id" xorm:"INT notnull"`
	ThreatEvent string `json:"threat_event" xorm:"VARCHAR(255) notnull"`
	AssetName   string `json:"asset_name" xorm:"JSON notnull"`
}

type RiskCalc struct {
	Metric string `json:"metric"`
}

type CombinedRisk struct {
	ThreatEventID       int64
	ThreatEvent         string
	MinFrequency        float64
	MaxFrequency        float64
	MostLikelyFrequency float64
	MinimumLoss         float64
	MaximumLoss         float64
	MostLikelyLoss      float64
}

type ImplementsInput struct {
	ControlID int64 `json:"controlId"`
	Current   int   `json:"current"`
	Proposed  int   `json:"proposed"`
	Cost      int   `json:"cost"`
}
type ImplementsInputNoID struct {
	Current  int `json:"current"`
	Proposed int `json:"proposed"`
	Cost     int `json:"cost"`
}

type RelevanceDinamicInput struct {
	ControlID  int            `json:"controlId"`
	Attributes map[string]int `json:"attributes"` // dynamic attributes
}

type LossLevel struct {
	Probability    float64 `json:"probability"`
	AcceptableLoss float64 `json:"acceptable_loss"`
}

type LossExceedance struct {
	Risk float64 `json:"risk" `
	Loss int64   `json:"loss" `
}
