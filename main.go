package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

const baseURL = "http://digitale.bnc.roma.sbn.it/tecadigitale/giornali/"
const cacheDir = "./bncr_cache"

func journals() {
	journals := colly.NewCollector(colly.CacheDir(cacheDir))
	journals.OnHTML("li a", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("href"), "giornali") {
			fmt.Println(e.Attr("href"))
		}
	})
	journals.Visit(baseURL)
}

func numbers(id string) {
	url := baseURL + id
	year := colly.NewCollector(colly.CacheDir(cacheDir))
	number := colly.NewCollector(colly.CacheDir(cacheDir))

	year.OnHTML("li a", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("href"), id) {
			number.Visit(e.Attr("href"))
		}
	})

	number.OnHTML("li a", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("href"), id) {
			fmt.Println(e.Attr("href"))
		}
	})

	year.Visit(url)
}
func metadata(id string) {
	var years []string
	var numbers []string
	metadata := make(map[string]string)

	year := colly.NewCollector(colly.CacheDir(cacheDir))
	year.OnHTML("li a", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("href"), id) {
			years = append(years, e.Text)
		}
	})

	url := baseURL + id
	year.Visit(url)
	firstYear := years[0]
	lastYear := years[len(years)-1]

	number := colly.NewCollector(colly.CacheDir(cacheDir))
	number.OnHTML("li a", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("href"), id) {
			numbers = append(numbers, e.Attr("href"))
		}
	})

	number.Visit(url + "/" + firstYear)
	firstYearFirstNumber := numbers[0]

	metadata["Annate"] = firstYear + ".." + lastYear
	metadata["Url"] = firstYearFirstNumber

	meta := colly.NewCollector(colly.CacheDir(cacheDir))
	// get journal thumbnail
	meta.OnHTML(".side-sx .image-box a img", func(e *colly.HTMLElement) {
		metadata["Thumbnail"] = e.Attr("src")
	})
	// get metadata
	meta.OnHTML(".data-detail dl dt", func(e *colly.HTMLElement) {
		metadata[e.Text] = strings.TrimSpace(e.DOM.Next().Text())
	})

	meta.Visit(firstYearFirstNumber)
	data, _ := json.MarshalIndent(metadata, "", "    ")
	fmt.Printf("%s", data)

}

func main() {
	getJournals := flag.Bool("get-journals", false, "get the url list of all journals")
	getNumbers := flag.String("get-numbers", "", "get the numbers list of a journal ID")
	getMetadata := flag.String("get-metadata", "", "get the metadata of a journal ID")

	flag.Parse()

	if *getJournals {
		journals()
	}
	if *getNumbers != "" {
		numbers(*getNumbers)
	}
	if *getMetadata != "" {
		metadata(*getMetadata)
	}
}
