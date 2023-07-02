package services

import (
	"poke-api-graphql/graph/model"
	"poke-api-graphql/repositories"
)

type AbilityService struct{}

func (abilityService AbilityService) FindAbilitiesData(abilities []*model.Ability) ([]*model.Ability, error) {
	result := make([]*model.Ability, len(abilities))

	for index, ability := range abilities {
		abilityResult, err := repositories.FindAbilityData(ability)

		if err != nil {
			return nil, err
		}

		result[index] = abilityResult
	}

	return result, nil
}
