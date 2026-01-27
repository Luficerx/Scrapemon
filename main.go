package main

import (
	"Scrapemon/scraper"
	"github.com/gocolly/colly/v2"
)

func main() {
	sc := scraper.Scrapemon{Cl: *colly.NewCollector(), Url: "https://pokemondb.net/pokedex/national"}

	sc.ScrapePokemonWebsite()
	sc.ScrapePokemons()
	sc.DumpPokemons()
}
