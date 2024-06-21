package db

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
	ID            int    `json:"id"`
	ThreatID      int    `json:"threat_id"`
	ThreatEvent   string `json:"threat_event"`
	AffectedAsset string `json:"affected_asset"`
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

type Relevance struct { //Cria tudo RELEVANCE
	ID                         int `json:"id"`
	AuthenticationAttack       int `json:"authentication_attack"`
	AuthorisationAttack        int `json:"authorisation_attack"`
	CommunicationAttack        int `json:"communication_attack"`
	DenialOfServiceAttack      int `json:"denial_of_service_attack"`
	InformationLeakageAttack   int `json:"information_leakage_attack"`
	MalwareAttack              int `json:"malware_attack"`
	MisconfigurationAttack     int `json:"misconfiguration_attack"`
	MisuseAttack               int `json:"misuse_attack"`
	PhysicalAttack             int `json:"physical_attack"`
	ReconnaissanceActivities   int `json:"reconnaissance_activities"`
	SocialEngineeringAttack    int `json:"social_engineering_attack"`
	SoftwareExploitationAttack int `json:"software_exploitation_attack"`
	SupplyChainAttack          int `json:"supply_chain_attack"`
	PeopleFailure              int `json:"people_failure"`
	ProcessFailure             int `json:"process_failure"`
	TechnologyFailure          int `json:"technology_failure"`
	BiologicalEvent            int `json:"biological_event"`
	MeteorologicalEvent        int `json:"meteorological_event"`
	GeologicalEvent            int `json:"geological_event"`
	HydrologicalEvent          int `json:"hydrological_event"`
	NaturalHazardEvent         int `json:"natural_hazard_event"`
	InfrastructureFailureEvent int `json:"infrastructure_failure_event"`
	AirborneParticlesEvent     int `json:"airborne_particles_event"`
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
	ID                     int     `json:"id"`
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
