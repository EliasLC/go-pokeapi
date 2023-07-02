package repositories

import (
	"fmt"
	"log"
	"poke-api-graphql/graph/model"
	"poke-api-graphql/pokeapi"
	"sync"
)

const apiResourseName = "pokemon/"
const apiGenerationResourseName = "generation/"

type PokemonRepository struct{}

func (repository PokemonRepository) FindByName(name string) (*model.Pokemon, error) {
	resourceUrl := fmt.Sprintf("%s%s", apiResourseName, name)

	return findPokemon(resourceUrl)
}

func (repository PokemonRepository) FindByNationalPokedexNumber(nationalPokedexNumber int) (*model.Pokemon, error) {
	requestUrl := fmt.Sprintf("%s%d", apiResourseName, nationalPokedexNumber)

	return findPokemon(requestUrl)
}

func FindPokemonByGeneration(generation int) ([]*model.Pokemon, error) {
	requestUrl := fmt.Sprintf("%s%d", apiGenerationResourseName, generation)

	return findPokemons(requestUrl)
}

func findPokemons(requestUrl string) ([]*model.Pokemon, error) {
	requestData, err := pokeapi.MakeApiRequest(requestUrl)

	if err != nil {
		return nil, err
	}

	pokemonsData := requestData["pokemon_species"].([]interface{})

	return createPokemonList(pokemonsData)
}

func createPokemonList(pokemonsData []interface{}) ([]*model.Pokemon, error) {
	var wg sync.WaitGroup
	result := make([]*model.Pokemon, len(pokemonsData))

	for index, pokemon := range pokemonsData {
		pokemonData := pokemon.(map[string]interface{})
		pokemonName := pokemonData["name"].(string)
		pokemonRequestUrl := fmt.Sprintf("%s%s", apiResourseName, pokemonName)

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
	abilities := createAbilities(requestData["abilities"])
	sprites := createSpriteField(requestData["sprites"])

	result := model.Pokemon{
		Name:                  fmt.Sprintf("%v", requestData["name"]),
		NationalPokedexNumber: int(requestData["id"].(float64)),
		Types:                 types,
		Abilities:             abilities,
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

func createAbilities(abilitiesRawData interface{}) []*model.Ability {
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
