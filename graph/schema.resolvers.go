package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.33

import (
	"context"
	"poke-api-graphql/graph/model"
)

// Abilities is the resolver for the abilities field.
func (r *pokemonResolver) Abilities(ctx context.Context, obj *model.Pokemon) ([]*model.Ability, error) {
	return r.abilityService.FindAbilitiesData(obj.Abilities)
}

// Pokemon is the resolver for the pokemon field.
func (r *queryResolver) Pokemon(ctx context.Context, filter model.PokemonFilter) (*model.Pokemon, error) {
	return r.pokemonService.FindPokemon(ctx, filter)
}

// Pokemons is the resolver for the pokemons field.
func (r *queryResolver) Pokemons(ctx context.Context, filter model.PokemonsFilter) ([]*model.Pokemon, error) {
	return r.pokemonService.FindPokemons(ctx, filter)
}

// Pokemon returns PokemonResolver implementation.
func (r *Resolver) Pokemon() PokemonResolver { return &pokemonResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type pokemonResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
