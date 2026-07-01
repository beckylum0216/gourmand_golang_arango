package entities

type Recipe struct {
	Id		  	string   `json:"_id,omitempty"`
	Title       string `json:"title"`
	Description string `json:"description"`
	PrepTime   	uint   `json:"preptime"`
	CookTime   	uint   `json:"cooktime"`
	Servings    uint   `json:"servings"`
	Url		 	string `json:"url"`
	Method 		string `json:"method"`
	RecipeText 	string `json:"recipe_text"`
}