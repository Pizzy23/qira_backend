package db

type AssetsInventory struct { //post get:all/id
	ID                int     `json:"id" xorm:"pk autoincr 'id' INT"`
	Name              string  `json:"name" xorm:"VARCHAR(255) notnull"`
	Description       string  `json:"description" xorm:"TEXT notnull"`
	Location          string  `json:"location" xorm:"VARCHAR(255) notnull"`
	Responsible       string  `json:"responsible" xorm:"VARCHAR(255) notnull"`
	BusinessValue     int     `json:"business_value" xorm:"INT notnull"`
	ReplacementCost   float64 `json:"replacement_cost" xorm:"DOUBLE notnull"`
	Criticality       string  `json:"criticality" xorm:"VARCHAR(50) notnull"`
	Users             string  `json:"users" xorm:"TEXT notnull"`
	TargetEnvironment string  `json:"target_environment" xorm:"VARCHAR(255) notnull"`
}

type ThreatEventCatalogue struct { //Ameaças
	ID          int    `json:"id" xorm:"pk autoincr 'id' INT"`
	ThreatGroup string `json:"threat_group" xorm:"VARCHAR(255) notnull"`
	ThreatEvent string `json:"threat_event" xorm:"VARCHAR(255) notnull"`
	Description string `json:"description" xorm:"TEXT notnull"`
	InScope     string `json:"in_scope" xorm:"VARCHAR(50) notnull"`
}

type Frequency struct { //Criar ameaça, ai cria a frequencia junto 0, Editar a frequencia]
	ID                  int    `json:"id" xorm:"pk autoincr 'id' INT"`
	ThreatEventID       int    `json:"threat_event_id" xorm:"INT notnull"`
	ThreatEvent         string `json:"threat_event" xorm:"VARCHAR(255) notnull"`
	MinFrequency        int    `json:"min_frequency" xorm:"INT notnull"`
	MaxFrequency        int    `json:"max_frequency" xorm:"INT notnull"`
	MostCommonFrequency int    `json:"most_common_frequency" xorm:"INT notnull"`
	SupportInformation  string `json:"support_information" xorm:"TEXT notnull"`
}

type ThreatEventAssets struct { // Vai selecionar Assets, Oque foi afetado Pode ser varios afetados
	ID            int      `json:"id" xorm:"pk autoincr 'id' INT"`
	ThreatID      int      `json:"threat_id" xorm:"INT notnull"`
	ThreatEvent   string   `json:"threat_event" xorm:"VARCHAR(255) notnull"`
	AffectedAsset []string `json:"affected_asset" xorm:"JSON notnull"`
}

type LossHigh struct { //Vinculado ao Events com maximo e min de perca.
	ID             int      `json:"id" xorm:"pk autoincr 'id' INT"`
	ThreatEventID  int      `json:"threat_event_id" xorm:"INT notnull"`
	ThreatEvent    string   `json:"threat_event" xorm:"VARCHAR(255) notnull"`
	Assets         []string `json:"assets" xorm:"JSON notnull"`
	LossType       string   `json:"loss_type" xorm:"VARCHAR(50) notnull"`
	MinimumLoss    float64  `json:"minimum_loss" xorm:"DOUBLE notnull"`
	MaximumLoss    float64  `json:"maximum_loss" xorm:"DOUBLE notnull"`
	MostLikelyLoss float64  `json:"most_likely_loss" xorm:"DOUBLE notnull"`
}

type RiskCalculator struct { //so da input :)
	ID                int     `json:"id" xorm:"pk autoincr 'id' INT"`
	ThreatEventID     int     `json:"threat_event_id" xorm:"INT notnull"`
	ThreatEvent       string  `json:"threat_event" xorm:"VARCHAR(255) notnull"`
	MinFrequency      int     `json:"min_frequency" xorm:"INT notnull"`
	MaxFrequency      int     `json:"max_frequency" xorm:"INT notnull"`
	ModeFrequency     int     `json:"mode_frequency" xorm:"INT notnull"`
	EstimateFrequency float64 `json:"estimate_frequency" xorm:"DOUBLE notnull"`
	MinLoss           int     `json:"min_loss" xorm:"INT notnull"`
	MaxLoss           int     `json:"max_loss" xorm:"INT notnull"`
	ModeLoss          int     `json:"mode_loss" xorm:"INT notnull"`
	EstimateLoss      float64 `json:"estimate_loss" xorm:"DOUBLE notnull"`
	MinRisk           int     `json:"min_risk" xorm:"INT notnull"`
	MaxRisk           int     `json:"max_risk" xorm:"INT notnull"`
	ModeRisk          int     `json:"mode_risk" xorm:"INT notnull"`
	EstimateRisk      float64 `json:"estimate_risk" xorm:"DOUBLE notnull"`
}

type Relevance struct { //Cria tudo RELEVANCE
	ID                         int `json:"id" xorm:"pk autoincr 'id' INT"`
	AuthenticationAttack       int `json:"authentication_attack" xorm:"INT notnull"`
	AuthorisationAttack        int `json:"authorisation_attack" xorm:"INT notnull"`
	CommunicationAttack        int `json:"communication_attack" xorm:"INT notnull"`
	DenialOfServiceAttack      int `json:"denial_of_service_attack" xorm:"INT notnull"`
	InformationLeakageAttack   int `json:"information_leakage_attack" xorm:"INT notnull"`
	MalwareAttack              int `json:"malware_attack" xorm:"INT notnull"`
	MisconfigurationAttack     int `json:"misconfiguration_attack" xorm:"INT notnull"`
	MisuseAttack               int `json:"misuse_attack" xorm:"INT notnull"`
	PhysicalAttack             int `json:"physical_attack" xorm:"INT notnull"`
	ReconnaissanceActivities   int `json:"reconnaissance_activities" xorm:"INT notnull"`
	SocialEngineeringAttack    int `json:"social_engineering_attack" xorm:"INT notnull"`
	SoftwareExploitationAttack int `json:"software_exploitation_attack" xorm:"INT notnull"`
	SupplyChainAttack          int `json:"supply_chain_attack" xorm:"INT notnull"`
	PeopleFailure              int `json:"people_failure" xorm:"INT notnull"`
	ProcessFailure             int `json:"process_failure" xorm:"INT notnull"`
	TechnologyFailure          int `json:"technology_failure" xorm:"INT notnull"`
	BiologicalEvent            int `json:"biological_event" xorm:"INT notnull"`
	MeteorologicalEvent        int `json:"meteorological_event" xorm:"INT notnull"`
	GeologicalEvent            int `json:"geological_event" xorm:"INT notnull"`
	HydrologicalEvent          int `json:"hydrological_event" xorm:"INT notnull"`
	NaturalHazardEvent         int `json:"natural_hazard_event" xorm:"INT notnull"`
	InfrastructureFailureEvent int `json:"infrastructure_failure_event" xorm:"INT notnull"`
	AirborneParticlesEvent     int `json:"airborne_particles_event" xorm:"INT notnull"`
}

type Strength struct {
	ID                         int `json:"id" xorm:"pk autoincr 'id' INT"`
	ControlID                  int `json:"controlId" xorm:"INT notnull"`
	AuthenticationAttack       int `json:"authentication_attack" xorm:"INT notnull"`
	AuthorisationAttack        int `json:"authorisation_attack" xorm:"INT notnull"`
	CommunicationAttack        int `json:"communication_attack" xorm:"INT notnull"`
	DenialOfServiceAttack      int `json:"denial_of_service_attack" xorm:"INT notnull"`
	InformationLeakageAttack   int `json:"information_leakage_attack" xorm:"INT notnull"`
	MalwareAttack              int `json:"malware_attack" xorm:"INT notnull"`
	MisconfigurationAttack     int `json:"misconfiguration_attack" xorm:"INT notnull"`
	MisuseAttack               int `json:"misuse_attack" xorm:"INT notnull"`
	PhysicalAttack             int `json:"physical_attack" xorm:"INT notnull"`
	ReconnaissanceActivities   int `json:"reconnaissance_activities" xorm:"INT notnull"`
	SocialEngineeringAttack    int `json:"social_engineering_attack" xorm:"INT notnull"`
	SoftwareExploitationAttack int `json:"software_exploitation_attack" xorm:"INT notnull"`
	SupplyChainAttack          int `json:"supply_chain_attack" xorm:"INT notnull"`
	PeopleFailure              int `json:"people_failure" xorm:"INT notnull"`
	ProcessFailure             int `json:"process_failure" xorm:"INT notnull"`
	TechnologyFailure          int `json:"technology_failure" xorm:"INT notnull"`
	BiologicalEvent            int `json:"biological_event" xorm:"INT notnull"`
	MeteorologicalEvent        int `json:"meteorological_event" xorm:"INT notnull"`
	GeologicalEvent            int `json:"geological_event" xorm:"INT notnull"`
	HydrologicalEvent          int `json:"hydrological_event" xorm:"INT notnull"`
	NaturalHazardEvent         int `json:"natural_hazard_event" xorm:"INT notnull"`
	InfrastructureFailureEvent int `json:"infrastructure_failure_event" xorm:"INT notnull"`
	AirborneParticlesEvent     int `json:"airborne_particles_event" xorm:"INT notnull"`
}

type Propused struct {
	ID                         int `json:"id" xorm:"pk autoincr 'id' INT"`
	ControlID                  int `json:"controlId" xorm:"INT notnull"`
	AuthenticationAttack       int `json:"authentication_attack" xorm:"INT notnull"`
	AuthorisationAttack        int `json:"authorisation_attack" xorm:"INT notnull"`
	CommunicationAttack        int `json:"communication_attack" xorm:"INT notnull"`
	DenialOfServiceAttack      int `json:"denial_of_service_attack" xorm:"INT notnull"`
	InformationLeakageAttack   int `json:"information_leakage_attack" xorm:"INT notnull"`
	MalwareAttack              int `json:"malware_attack" xorm:"INT notnull"`
	MisconfigurationAttack     int `json:"misconfiguration_attack" xorm:"INT notnull"`
	MisuseAttack               int `json:"misuse_attack" xorm:"INT notnull"`
	PhysicalAttack             int `json:"physical_attack" xorm:"INT notnull"`
	ReconnaissanceActivities   int `json:"reconnaissance_activities" xorm:"INT notnull"`
	SocialEngineeringAttack    int `json:"social_engineering_attack" xorm:"INT notnull"`
	SoftwareExploitationAttack int `json:"software_exploitation_attack" xorm:"INT notnull"`
	SupplyChainAttack          int `json:"supply_chain_attack" xorm:"INT notnull"`
	PeopleFailure              int `json:"people_failure" xorm:"INT notnull"`
	ProcessFailure             int `json:"process_failure" xorm:"INT notnull"`
	TechnologyFailure          int `json:"technology_failure" xorm:"INT notnull"`
	BiologicalEvent            int `json:"biological_event" xorm:"INT notnull"`
	MeteorologicalEvent        int `json:"meteorological_event" xorm:"INT notnull"`
	GeologicalEvent            int `json:"geological_event" xorm:"INT notnull"`
	HydrologicalEvent          int `json:"hydrological_event" xorm:"INT notnull"`
	NaturalHazardEvent         int `json:"natural_hazard_event" xorm:"INT notnull"`
	InfrastructureFailureEvent int `json:"infrastructure_failure_event" xorm:"INT notnull"`
	AirborneParticlesEvent     int `json:"airborne_particles_event" xorm:"INT notnull"`
}

type ControlLibrary struct {
	ID               int    `json:"id" xorm:"pk autoincr 'id' INT"`
	ControlID        int    `json:"controlId" xorm:"INT notnull"`
	ControlType      string `json:"control_type" xorm:"VARCHAR(255) notnull"`
	ControlReference string `json:"control_reference" xorm:"VARCHAR(255) notnull"`
	Information      string `json:"information" xorm:"TEXT notnull"`
	InScope          string `json:"in_scope" xorm:"VARCHAR(50) notnull"`
}

type ControlImplementation struct {
	ID                     int     `json:"id" xorm:"pk autoincr 'id' INT"`
	ControlID              int     `json:"controlId" xorm:"INT notnull"`
	CurrentImplementation  int     `json:"current_implementation" xorm:"INT notnull"`
	CurrentPercentValue    string  `json:"current_percent_value" xorm:"VARCHAR(50) notnull"`
	ProposedImplementation int     `json:"proposed_implementation" xorm:"INT notnull"`
	ProposedPercentValue   string  `json:"proposed_percent_value" xorm:"VARCHAR(50) notnull"`
	ProjectedCost          float64 `json:"projected_cost" xorm:"DOUBLE notnull"`
}

type AggregatedControlStrength struct { // COnversa com o event
	ID                      int    `json:"id" xorm:"pk autoincr 'id' INT"`
	ThreatEventID           int    `json:"threat_event_id" xorm:"INT notnull"`
	ThreatEvent             string `json:"threat_event" xorm:"VARCHAR(255) notnull"`
	CurrentControlStrength  string `json:"current_control_strength" xorm:"VARCHAR(50) notnull"`
	ProposedControlStrength string `json:"proposed_control_strength" xorm:"VARCHAR(50) notnull"`
}
