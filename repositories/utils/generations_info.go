package utils

func GenerationPokemonCount(generation int) int {
	return [9]int{151, 100, 135, 107, 156, 72, 88, 96, 105}[generation]
}

func GenerationInitialOffset(generation int) int {
	return [9]int{0, 151, 251, 386, 493, 649, 721, 809, 905}[generation]
}
