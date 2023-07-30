package services

import (
	"context"
	"poke-api-graphql/graph/model"
	"poke-api-graphql/repositories"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type PokemonService struct {
	pokemonRespository repositories.PokemonRepository
}

func NewPokemonService() *PokemonService {
	repository := repositories.PokemonRepository{}

	pokemonService := PokemonService{
		pokemonRespository: repository,
	}

	return &pokemonService
}

func (pokemonService PokemonService) FindPokemon(ctx context.Context, filter model.PokemonFilter) (*model.Pokemon, error) {
	var pokemon *model.Pokemon
	var err error

	if filter.Name != nil {
		pokemon, err = pokemonService.pokemonRespository.FindByName(*filter.Name)
	}

	if filter.NationalPokedexNumber != nil {
		pokemon, err = pokemonService.pokemonRespository.FindByNationalPokedexNumber(*filter.NationalPokedexNumber)
	}

	if pokemon == nil {
		graphql.AddError(ctx, gqlerror.Errorf(err.Error()))
	}

	return pokemon, err
}

func (pokemonService PokemonService) FindPokemons(ctx context.Context, filter model.PokemonsFilter) ([]*model.Pokemon, error) {
	var pokemons []*model.Pokemon
	var err error

	if filter.Generation != nil {
		pokemons, err = repositories.FindPokemosByGeneration(*filter.Generation, *filter.Offset)
	} else {
		pokemons, err = repositories.FindPokemons(*filter.Offset)
	}

	if err != nil {
		graphql.AddError(ctx, gqlerror.Errorf(err.Error()))
	}

	return pokemons, err
}
