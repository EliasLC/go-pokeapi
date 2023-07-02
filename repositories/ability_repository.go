package repositories

import (
	"fmt"
	"poke-api-graphql/graph/model"
	"poke-api-graphql/pokeapi"
)

const abiltyApiResouseName = "ability/"

func FindAbilityData(ability *model.Ability) (*model.Ability, error) {
	requestUrl := fmt.Sprintf("%s%s", abiltyApiResouseName, ability.Name)

	requestData, err := pokeapi.MakeApiRequest(requestUrl)

	if err != nil {
		return nil, err
	}

	generationRawData := requestData["generation"].(map[string]interface{})
	generationName := generationRawData["name"].(string)
	effectRawData := requestData["effect_entries"].([]interface{})
	shortEffect, effect := getEffectData(effectRawData)

	newAbility := model.Ability{
		Name:        ability.Name,
		Generation:  generationName,
		ShortEffect: shortEffect,
		Effect:      effect,
	}

	return &newAbility, nil
}

func getEffectData(effectRawData []interface{}) (string, string) {
	var shortEffect string
	var effect string

	for _, effectLanguage := range effectRawData {
		effectLanguageData := effectLanguage.(map[string]interface{})
		languageData := effectLanguageData["language"].(map[string]interface{})
		languageName := languageData["name"].(string)

		if languageName == "en" {
			shortEffect = effectLanguageData["short_effect"].(string)
			effect = effectLanguageData["effect"].(string)
		}
	}

	return shortEffect, effect
}
