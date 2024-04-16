package compareDB

import(
	"src/reader"
	"fmt"
)

// compareDatabases compares two databases and prints the changes.
func CompareDatabases(original, stolen *reader.Recipe) {
	// Create maps for quick lookup
	originalCakes := make(map[string]reader.Cake)
	for _, cake := range original.Cakes {
		originalCakes[cake.Name] = cake
	}

	stolenCakes := make(map[string]reader.Cake)
	for _, cake := range stolen.Cakes {
		stolenCakes[cake.Name] = cake
	}

	// Check for added or removed cakes
	for name := range stolenCakes {
		if _, exists := originalCakes[name]; !exists {
			fmt.Printf("ADDED cake \"%s\"\n", name)
		}
	}

	for name := range originalCakes {
		if _, exists := stolenCakes[name]; !exists {
			fmt.Printf("REMOVED cake \"%s\"\n", name)
		}
	}

	// Check for changes in cooking time and ingredients
	for name, originalCake := range originalCakes {
		if stolenCake, exists := stolenCakes[name]; exists {
			if originalCake.StoveTime != stolenCake.StoveTime {
				fmt.Printf("CHANGED cooking time for cake \"%s\" - \"%s\" instead of \"%s\"\n", name, stolenCake.StoveTime, originalCake.StoveTime)
			}

			// Create maps for quick lookup of ingredients
			originalIngredients := make(map[string]reader.Item)
			for _, ingredient := range originalCake.Ingredients {
				originalIngredients[ingredient.Name] = ingredient
			}

			stolenIngredients := make(map[string]reader.Item)
			for _, ingredient := range stolenCake.Ingredients {
				stolenIngredients[ingredient.Name] = ingredient
			}

			// Check for added or removed ingredients
			for ingredientName := range stolenIngredients {
				if _, exists := originalIngredients[ingredientName]; !exists {
					fmt.Printf("ADDED ingredient \"%s\" for cake \"%s\"\n", ingredientName, name)
				}
			}

			for ingredientName := range originalIngredients {
				if _, exists := stolenIngredients[ingredientName]; !exists {
					fmt.Printf("REMOVED ingredient \"%s\" for cake \"%s\"\n", ingredientName, name)
				}
			}

			// Check for changes in ingredient count and unit
			for ingredientName, originalIngredient := range originalIngredients {
				if stolenIngredient, exists := stolenIngredients[ingredientName]; exists {
					if originalIngredient.Count != stolenIngredient.Count {
						fmt.Printf("CHANGED unit count for ingredient \"%s\" for cake \"%s\" - \"%.1f\" instead of \"%.1f\"\n", ingredientName, name, stolenIngredient.Count, originalIngredient.Count)
					}
					if originalIngredient.Unit != stolenIngredient.Unit {
						if originalIngredient.Unit != "" && stolenIngredient.Unit == "" {
							fmt.Printf("REMOVED unit \"%s\" for ingredient  \"%s\" for cake \"%s", originalIngredient.Unit, ingredientName, name)
							continue	
						} 
						fmt.Printf("CHANGED unit for ingredient \"%s\" for cake \"%s\" - \"%s\" instead of \"%s\"\n", ingredientName, name, stolenIngredient.Unit, originalIngredient.Unit)
					}
				}
			}
		}
	}
}