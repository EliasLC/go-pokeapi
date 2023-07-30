package repositories

import (
	"fmt"
	"log"
	"poke-api-graphql/graph/model"
	"poke-api-graphql/pokeapi"
	"poke-api-graphql/repositories/utils"
	"sync"
)

const pokemonResourceName = "pokemon/"

type PokemonRepository struct{}

func (repository PokemonRepository) FindByName(name string) (*model.Pokemon, error) {
	resourceUrl := fmt.Sprintf("%s%s%s", pokeapi.ApiUrl, pokemonResourceName, name)

	return findPokemon(resourceUrl)
}

func (repository PokemonRepository) FindByNationalPokedexNumber(nationalPokedexNumber int) (*model.Pokemon, error) {
	requestUrl := fmt.Sprintf("%s%s%d", pokeapi.ApiUrl, pokemonResourceName, nationalPokedexNumber)

	return findPokemon(requestUrl)
}

func FindPokemons(offset int) ([]*model.Pokemon, error) {
	requestUrl := fmt.Sprintf("%s%s?offset=%d", pokeapi.ApiUrl, pokemonResourceName, offset)

	return findPokemons(requestUrl)
}

func FindPokemosByGeneration(generation int, offset int) ([]*model.Pokemon, error) {
	pokeApiOffset := utils.GenerationInitialOffset(generation) + offset
	pokeApiOffsetLimit := getPokeApiLimit(generation, offset)

	log.Println("offset: ", pokeApiOffset, " limit: ", pokeApiOffsetLimit)

	requestUrl := fmt.Sprintf("%s%s?offset=%d&limit=%d", pokeapi.ApiUrl, pokemonResourceName, pokeApiOffset, pokeApiOffsetLimit)

	return findPokemons(requestUrl)
}

func getPokeApiLimit(generation int, offset int) int {
	generationPokemonCount := utils.GenerationPokemonCount(generation)

	var pokeApiLimit = 20

	if (offset + 20) > generationPokemonCount {
		return generationPokemonCount - offset
	}

	if pokeApiLimit <= 0 {
		return -1
	}

	return pokeApiLimit
}

func findPokemons(requestUrl string) ([]*model.Pokemon, error) {
	requestData, err := pokeapi.MakeApiRequest(requestUrl)

	if err != nil {
		return nil, err
	}

	pokemonsData := requestData["results"].([]interface{})

	return createPokemonList(pokemonsData)
}

func createPokemonList(pokemonsData []interface{}) ([]*model.Pokemon, error) {
	var wg sync.WaitGroup
	result := make([]*model.Pokemon, len(pokemonsData))

	for index, pokemon := range pokemonsData {
		pokemonData := pokemon.(map[string]interface{})
		pokemonRequestUrl := pokemonData["url"].(string)
		wg.Add(1)

		go func(index int) {
			defer wg.Done()

			pokemonData, err := findPokemon(pokemonRequestUrl)

			if err != nil {
				log.Println("pokemon request")
			}

			result[index] = pokemonData
		}(index)
	}

	wg.Wait()

	return result, nil
}

func findPokemon(requestUrl string) (*model.Pokemon, error) {
	requestData, err := pokeapi.MakeApiRequest(requestUrl)

	if err != nil {
		log.Println("request error")
		return nil, err
	}

	types := createTypesField(requestData["types"])
	stats := createStatsField(requestData["stats"])
	abilities := createAbilitiesField(requestData["abilities"])
	moves := createMovesField(requestData["moves"])
	sprites := createSpriteField(requestData["sprites"])

	result := model.Pokemon{
		Name:                  fmt.Sprintf("%v", requestData["name"]),
		NationalPokedexNumber: int(requestData["id"].(float64)),
		Types:                 types,
		Stats:                 stats,
		Abilities:             abilities,
		Moves:                 moves,
		Sprites:               sprites,
	}

	return &result, nil
}

func createTypesField(typesRawData interface{}) []string {
	types := make([]string, 0)

	typeList := typesRawData.([]interface{})

	for _, typeRawData := range typeList {
		typeMap := typeRawData.(map[string]interface{})
		typeData := typeMap["type"].(map[string]interface{})
		types = append(types, typeData["name"].(string))
	}

	return types
}

func createStatsField(statsRawData interface{}) []*model.Stat {
	statsData := statsRawData.([]interface{})
	result := make([]*model.Stat, len(statsData))

	for index, stat := range statsData {
		statData := stat.(map[string]interface{})
		statDetails := statData["stat"].(map[string]interface{})

		statName := statDetails["name"].(string)
		baseStat := int(statData["base_stat"].(float64))
		effort := int(statData["effort"].(float64))

		newStat := model.Stat{
			Name:     statName,
			BaseStat: baseStat,
			Effort:   effort,
		}

		result[index] = &newStat
	}

	return result
}

func createAbilitiesField(abilitiesRawData interface{}) []*model.Ability {
	rawData := abilitiesRawData.([]interface{})
	abilities := make([]*model.Ability, len(rawData))

	for index, ability := range rawData {
		abilityMap := ability.(map[string]interface{})
		abilityData := abilityMap["ability"].(map[string]interface{})

		newAbility := model.Ability{
			Name: abilityData["name"].(string),
		}

		abilities[index] = &newAbility
	}

	return abilities
}

func createMovesField(movesRawData interface{}) []*model.PokemonMoveDetails {
	moveListRawData := movesRawData.([]interface{})
	result := make([]*model.PokemonMoveDetails, len(moveListRawData))

	for index, moveRawData := range moveListRawData {
		moveMap := moveRawData.(map[string]interface{})
		moveDetails := moveMap["move"].(map[string]interface{})
		versionGroupDetails := moveMap["version_group_details"].([]interface{})

		moveName := moveDetails["name"].(string)
		versionDetails := make([]*model.MoveVersionGroupDetails, len(versionGroupDetails))

		for index, versionMoveDetails := range versionGroupDetails {
			details := createMoveVersionDetails(versionMoveDetails)

			versionDetails[index] = details
		}

		newMove := model.PokemonMoveDetails{
			Name:           moveName,
			VersionDetails: versionDetails,
		}

		result[index] = &newMove
	}

	return result
}

func createMoveVersionDetails(vgDetails interface{}) *model.MoveVersionGroupDetails {
	detailsMap := vgDetails.(map[string]interface{})
	learnMethodMap := detailsMap["move_learn_method"].(map[string]interface{})
	versionGroupMap := detailsMap["version_group"].(map[string]interface{})

	levelLearnedAt := int(detailsMap["level_learned_at"].(float64))
	learnMethod := learnMethodMap["name"].(string)
	version := versionGroupMap["name"].(string)

	return &model.MoveVersionGroupDetails{
		LevelLearnedAt: levelLearnedAt,
		LearnMethod:    learnMethod,
		Version:        version,
	}
}

func createSpriteField(spritesRawData interface{}) *model.Sprites {
	var sprites model.Sprites

	spritesMap, ok := spritesRawData.(map[string]interface{})

	backDefault := fmt.Sprintf("%v", spritesMap["back_default"])
	backFemale := fmt.Sprintf("%v", spritesMap["back_female"])
	frontDefault := fmt.Sprintf("%v", spritesMap["front_default"])
	frontFemale := fmt.Sprintf("%v", spritesMap["front_female"])
	frontShiny := fmt.Sprintf("%v", spritesMap["front_shiny"])
	frontShinyFemale := fmt.Sprintf("%v", spritesMap["front_shiny_female"])
	backShiny := fmt.Sprintf("%v", spritesMap["back_shiny"])
	backShinyFemale := fmt.Sprintf("%v", spritesMap["back_shiny_female"])

	if ok {
		sprites.BackDefault = &backDefault
		sprites.BackFemale = &backFemale
		sprites.FrontDefault = &frontDefault
		sprites.FrontFemale = &frontFemale
		sprites.FrontShiny = &frontShiny
		sprites.FrontShinyFemale = &frontShinyFemale
		sprites.BackShiny = &backShiny
		sprites.BackShinyFemale = &backShinyFemale
	}

	return &sprites
}
