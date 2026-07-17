package entities




type FoodDetail struct {
	FdcId          int             `json:"fdcId"`
	Description    string          `json:"description"`
	DataType       string          `json:"dataType"` // Foundation, SR Legacy, Branded, Survey (FNDDS)
	PublicationDate string         `json:"publicationDate"`
	FoodNutrients  []FoodNutrient  `json:"foodNutrients"`
	FoodPortions   []FoodPortion   `json:"foodPortions,omitempty"`

	// Branded-only fields
	BrandOwner     *string `json:"brandOwner,omitempty"`
	GtinUpc        *string `json:"gtinUpc,omitempty"`
	IngredientsRaw *string `json:"ingredients,omitempty"` // free text ingredient list, Branded only
	ServingSize    *float64 `json:"servingSize,omitempty"`
	ServingSizeUnit *string `json:"servingSizeUnit,omitempty"`

	// SR Legacy / Foundation
	NdbNumber      *string `json:"ndbNumber,omitempty"`
	FoodCategory   *FoodCategory `json:"foodCategory,omitempty"`
}

type FoodNutrient struct {
	Nutrient NutrientRef `json:"nutrient"`
	Amount   float64     `json:"amount"` // per 100g for Foundation/SR Legacy; check servingSize for Branded
}

type NutrientRef struct {
	Id       int    `json:"id"`
	Number   string `json:"number,omitempty"`
	Name     string `json:"name"`
	UnitName string `json:"unitName"`
}

type FoodPortion struct {
	Id          int      `json:"id"`
	Amount      float64  `json:"amount"`
	Modifier    string   `json:"modifier"` // e.g. "medium", "cup"
	GramWeight  float64  `json:"gramWeight"`
	Measure     *MeasureUnit `json:"measureUnit,omitempty"`
}

type MeasureUnit struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
}

type FoodCategory struct {
	Id          int    `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

type FoodsRequest struct {
	FdcIds    []int    `json:"fdcIds"`
	Format    string   `json:"format,omitempty"`
	Nutrients []int    `json:"nutrients,omitempty"`
}

