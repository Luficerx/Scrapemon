package scraper

import (
	"Scrapemon/dump"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

type batch [8]string

type Scrapemon struct {
	Cl  colly.Collector
	Url string

	pokemons []Pokemon
	urls     []batch
	mu       sync.Mutex

	found_urls int
}

func (sc *Scrapemon) ScrapePokemonWebsite() {
	log.Println("[SCRAPING]", sc.Url)
	var array batch
	var index int

	sc.Cl.OnHTML("main#main", func(elem *colly.HTMLElement) {
		elem.ForEach("div.infocard-list.infocard-list-pkmn-lg", func(_ int, elem *colly.HTMLElement) {
			elem.ForEach("div.infocard span.infocard-lg-img a", func(_ int, elem *colly.HTMLElement) {
				url := elem.Request.AbsoluteURL(elem.Attr("href"))
				sc.found_urls++

				if BatchInsert(&array, url, &index) {
					sc.mu.Lock()
					sc.urls = append(sc.urls, array)
					sc.mu.Unlock()
					array = batch{}
				}
			})

			// Ensure we have all elements
			for _, data := range array {
				if data != "" {
					sc.mu.Lock()
					sc.urls = append(sc.urls, array)
					sc.mu.Unlock()
					array = batch{}
					index = 0
					break
				}
			}
		})
	})

	sc.Cl.Visit(sc.Url)
	sc.Cl.Wait()

	log.Println("[URLS FOUND]", sc.found_urls)
}

func (sc *Scrapemon) ScrapePokemons() {
	var mu sync.Mutex

	for _, array := range sc.urls {
		var wg sync.WaitGroup

		for i, url := range array {
			if url == "" {
				continue
			}

			wg.Add(1)
			go func(_ int, url string) {
				defer wg.Done()

				var pokemon Pokemon
				c := sc.Cl.Clone()

				c.OnHTML("main#main", func(main *colly.HTMLElement) {
					pokemon.Name = main.ChildText("h1")
					pokemon.Url = url
					pokemon.Image = main.ChildAttr("a[rel='lightbox']", "href")

					main.ForEach("table.vitals-table tbody tr", func(_ int, elem *colly.HTMLElement) {
						key := strings.TrimSpace(elem.ChildText("th"))

						switch key {
						case "National №":
							pokemon.National = elem.ChildText("td")

						case "Type":
							elem.ForEach("td a", func(_ int, elem *colly.HTMLElement) {
								pokemon.Type = append(pokemon.Type, elem.Text)
							})

						case "Species":
							pokemon.Species = elem.ChildText("td")

						case "Height":
							pokemon.Height = elem.ChildText("td")

						case "Weight":
							pokemon.Weight = elem.ChildText("td")

						case "Abilities":
							pokemon.Abilities = append(pokemon.Abilities, elem.ChildText("span"))
							pokemon.Abilities = append(pokemon.Abilities, elem.ChildText("small"))

						case "Local №":
							elem.ForEach("td", func(_ int, elem *colly.HTMLElement) {
								location := ""

								elem.DOM.Contents().Each(func(_ int, node *goquery.Selection) {
									node_type := goquery.NodeName(node)

									switch node_type {

									case "#text":
										text := strings.TrimSpace(node.Text())
										if text != "" && text != "\n" {
											location = text
										}

									case "small":
										local := strings.TrimSpace(node.Text())
										location = location + " " + local

										pokemon.Locales = append(pokemon.Locales, location)

										location = ""
									case "br":
									}
								})
							})

						case "EV yield":
							pokemon.EvYield = elem.ChildText("td")

						case "Catch rate":
							pokemon.CatchRate = elem.ChildText("td")

						case "Base Friendship":
							pokemon.BaseFriendship = elem.ChildText("td")

						case "Base Exp.":
							value := elem.ChildText("td")
							number, err := strconv.Atoi(value)
							if err != nil {
								pokemon.BaseExp = 0
							} else {
								pokemon.BaseExp = number
							}

						case "Growth Rate":
							pokemon.GrowthRate = elem.ChildText("td")

						case "Egg Groups":
							for eggs := range strings.SplitSeq(elem.ChildText("td"), ",") {
								pokemon.EggGroups = append(pokemon.EggGroups, strings.TrimSpace(eggs))
							}

						case "Gender":
							for gender := range strings.SplitSeq(elem.ChildText("td"), ",") {
								pokemon.Gender = append(pokemon.Gender, strings.TrimSpace(gender))
							}

						case "Egg cycles":
							elem.ForEach("td", func(_ int, elem *colly.HTMLElement) {
								value := ""

								elem.DOM.Contents().Each(func(_ int, node *goquery.Selection) {
									node_type := goquery.NodeName(node)

									switch node_type {

									case "#text":
										text := strings.TrimSpace(node.Text())
										if text != "\n" {
											value = text
										}

									case "small":
										extras := strings.TrimSpace(node.Text())
										value = value + " " + extras

										pokemon.Cycles = value

										value = ""

									case "br":
									}
								})
							})

						}
					})
				})

				c.Visit(url)
				c.Wait()
				mu.Lock()
				sc.pokemons = append(sc.pokemons, pokemon)
				mu.Unlock()

			}(i, url)
		}
		wg.Wait()
	}

	if sc.found_urls != len(sc.pokemons) {
		log.Printf("[FAILURE] - Scrapemon found [%d] urls but only managed to create [%d] pokemons!\n", sc.found_urls, len(sc.pokemons))
	} else {
		log.Printf("[SUCCESS] - Scrapemon found [%d] pokemons! out of [%d] urls!\n", len(sc.pokemons), sc.found_urls)
	}
}

func (sc *Scrapemon) DumpPokemons() {
	bytes, err := json.MarshalIndent(sc.pokemons, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	dump.Dump("pokemon.json", bytes)
}

func BatchInsert(array *batch, value string, index *int) bool {
	array[*index] = value
	*index++

	if *index == len(array) {
		*index = 0
		return true
	}

	return false
}
