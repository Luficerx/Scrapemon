package scraper

type Pokemon struct {
	// Original url page.
	Url string `json:"url"`

	// Pokemon name.
	Name string `json:"name"`

	// Image url.
	Image string `json:"image"`

	// National number.
	National string `json:"national"`

	// -- GENERIC DATA

	// Pokemon type.
	Type []string `json:"type"`

	// Pokemon species.
	Species string `json:"species"`

	// Pokemon height.
	Height string `json:"height"`

	// Weight height.
	Weight string `json:"weight"`

	// Pokemon abilities, hidden included.
	Abilities []string `json:"abilities"`

	// Places pokemon appear at.
	Locales []string `json:"locales"`

	// -- TRAINING DATA

	// Pokemon ev yield.
	EvYield string `json:"evyield"`

	// Pokemon catch rate.
	CatchRate string `json:"catchrate"`

	// Pokemon base friendship
	BaseFriendship string `json:"basefriendship"`

	// Pokemon base experience.
	BaseExp int `json:"baseexp"`

	// Pokemon growth rate
	GrowthRate string `json:"growthrate"`

	// -- BREEDING DATA

	// Pokemon egg groups
	EggGroups []string `json:"egggroups"`

	// Pokemon gender percentage
	Gender []string `json:"gender"`

	// Pokemon cycles
	Cycles string `json:"cycles"`
}
