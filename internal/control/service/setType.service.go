package control

import (
	"errors"
	"qira/internal/interfaces"
)

// setTypes sets the appropriate table structure based on ControlType
func setTypes(data interfaces.InputControlls) (interface{}, error) {
	var table interface{}
	switch data.ControlType {
	case "Relevance":
		table = &interfaces.Relevance{
			AuthenticationAttack:       data.AuthenticationAttack,
			AuthorisationAttack:        data.AuthorisationAttack,
			CommunicationAttack:        data.CommunicationAttack,
			DenialOfServiceAttack:      data.DenialOfServiceAttack,
			InformationLeakageAttack:   data.InformationLeakageAttack,
			MalwareAttack:              data.MalwareAttack,
			MisconfigurationAttack:     data.MisconfigurationAttack,
			MisuseAttack:               data.MisuseAttack,
			PhysicalAttack:             data.PhysicalAttack,
			ReconnaissanceActivities:   data.ReconnaissanceActivities,
			SocialEngineeringAttack:    data.SocialEngineeringAttack,
			SoftwareExploitationAttack: data.SoftwareExploitationAttack,
			SupplyChainAttack:          data.SupplyChainAttack,
			PeopleFailure:              data.PeopleFailure,
			ProcessFailure:             data.ProcessFailure,
			TechnologyFailure:          data.TechnologyFailure,
			BiologicalEvent:            data.BiologicalEvent,
			MeteorologicalEvent:        data.MeteorologicalEvent,
			GeologicalEvent:            data.GeologicalEvent,
			HydrologicalEvent:          data.HydrologicalEvent,
			NaturalHazardEvent:         data.NaturalHazardEvent,
			InfrastructureFailureEvent: data.InfrastructureFailureEvent,
			AirborneParticlesEvent:     data.AirborneParticlesEvent,
		}
	case "Strength":
		table = &interfaces.Strength{
			AuthenticationAttack:       data.AuthenticationAttack,
			AuthorisationAttack:        data.AuthorisationAttack,
			CommunicationAttack:        data.CommunicationAttack,
			DenialOfServiceAttack:      data.DenialOfServiceAttack,
			InformationLeakageAttack:   data.InformationLeakageAttack,
			MalwareAttack:              data.MalwareAttack,
			MisconfigurationAttack:     data.MisconfigurationAttack,
			MisuseAttack:               data.MisuseAttack,
			PhysicalAttack:             data.PhysicalAttack,
			ReconnaissanceActivities:   data.ReconnaissanceActivities,
			SocialEngineeringAttack:    data.SocialEngineeringAttack,
			SoftwareExploitationAttack: data.SoftwareExploitationAttack,
			SupplyChainAttack:          data.SupplyChainAttack,
			PeopleFailure:              data.PeopleFailure,
			ProcessFailure:             data.ProcessFailure,
			TechnologyFailure:          data.TechnologyFailure,
			BiologicalEvent:            data.BiologicalEvent,
			MeteorologicalEvent:        data.MeteorologicalEvent,
			GeologicalEvent:            data.GeologicalEvent,
			HydrologicalEvent:          data.HydrologicalEvent,
			NaturalHazardEvent:         data.NaturalHazardEvent,
			InfrastructureFailureEvent: data.InfrastructureFailureEvent,
			AirborneParticlesEvent:     data.AirborneParticlesEvent,
		}
	case "Propused":
		table = &interfaces.Propused{
			AuthenticationAttack:       data.AuthenticationAttack,
			AuthorisationAttack:        data.AuthorisationAttack,
			CommunicationAttack:        data.CommunicationAttack,
			DenialOfServiceAttack:      data.DenialOfServiceAttack,
			InformationLeakageAttack:   data.InformationLeakageAttack,
			MalwareAttack:              data.MalwareAttack,
			MisconfigurationAttack:     data.MisconfigurationAttack,
			MisuseAttack:               data.MisuseAttack,
			PhysicalAttack:             data.PhysicalAttack,
			ReconnaissanceActivities:   data.ReconnaissanceActivities,
			SocialEngineeringAttack:    data.SocialEngineeringAttack,
			SoftwareExploitationAttack: data.SoftwareExploitationAttack,
			SupplyChainAttack:          data.SupplyChainAttack,
			PeopleFailure:              data.PeopleFailure,
			ProcessFailure:             data.ProcessFailure,
			TechnologyFailure:          data.TechnologyFailure,
			BiologicalEvent:            data.BiologicalEvent,
			MeteorologicalEvent:        data.MeteorologicalEvent,
			GeologicalEvent:            data.GeologicalEvent,
			HydrologicalEvent:          data.HydrologicalEvent,
			NaturalHazardEvent:         data.NaturalHazardEvent,
			InfrastructureFailureEvent: data.InfrastructureFailureEvent,
			AirborneParticlesEvent:     data.AirborneParticlesEvent,
		}
	default:
		return nil, errors.New("invalid control type")
	}
	return table, nil
}

func setTypesGet(controlType string) (interface{}, error) {
	switch controlType {
	case "Relevance":
		return &interfaces.Relevance{}, nil
	case "Strength":
		return &interfaces.Strength{}, nil
	case "Propused":
		return &interfaces.Propused{}, nil
	case "Library":
		return &interfaces.ControlLibrary{}, nil
	case "Implementation":
		return &interfaces.ControlImplementation{}, nil
	default:
		return nil, errors.New("invalid control type")
	}
}
