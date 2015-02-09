package recipes

import (
	"fmt"
	proto "proto"
)

func DebugPrint(recipe proto.Recipe) {
	fmt.Println(*recipe.Name)
	fmt.Println(fmt.Sprintf("Prep: %dm, Cook: %dm, Ready: %dm",
		*recipe.Time.Prep, *recipe.Time.Cook, *recipe.Time.Ready))

	fmt.Println("\nIngredients:")
	for i, ingredient := range recipe.Ingredients {
		fmt.Println(fmt.Sprintf("  %d. %s", i+1, *ingredient.Name))
		fmt.Println(fmt.Sprintf("     %s", ingredient.Modifiers))
		fmt.Println(fmt.Sprintf("     %s", ingredient.Ingrids))
	}
}
