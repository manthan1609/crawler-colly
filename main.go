package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

type star struct {
	Name      string
	Photo     string
	JobTitle  string
	BirthDate string
	Bio       string
	TopMovies []movie
}

type movie struct {
	Title string
	Year  string
}

func main() {
	month := flag.Int("month", 1, "Month to fetch birthdays for")
	day := flag.Int("day", 1, "Day to fetch birthdays for")

	flag.Parse()
	crawl(*month, *day)
}

func crawl(month int, day int) {

	c := colly.NewCollector(colly.AllowedDomains("imdb.com", "www.imdb.com"))

	infoCollector := c.Clone()

	c.OnHTML(".ipc-metadata-list", func(h *colly.HTMLElement) {
		profileUrl := h.ChildAttr("div.ipc-avatar > a", "href")
		profileUrl = h.Request.AbsoluteURL(profileUrl)
		infoCollector.Visit(profileUrl)
	})

	// c.OnHTML("button.ipc-btn", func(h *colly.HTMLElement) {

	// 	nextPage := h.Request.AbsoluteURL(h.Attr("href"))
	// 	c.Visit(nextPage)
	// })

	infoCollector.OnHTML(".ipc-page-section", func(h *colly.HTMLElement) {
		tmpProfile := star{}
		tmpProfile.Name = h.ChildText("h > span.hero__primary-text")
		tmpProfile.Photo = h.ChildAttr(".ipc-poster > a.ipc-lockup-overlay", "href")
		tmpProfile.JobTitle = h.ChildText("h3.ipc-title__text")
		tmpProfile.BirthDate = h.ChildText("li.ipc-inline-list__item > a")
		tmpProfile.Bio = strings.TrimSpace(h.ChildText("#name-bio-text > div.name-trivia-bio-text > div.inline"))

		h.ForEach("div.ipc-sub-grid", func(i int, kf *colly.HTMLElement) {
			tmpMovie := movie{}

			tmpMovie.Title = kf.ChildText("div.ipc-list-card--span > a.ipc-primary-image-list-card__title")
			tmpMovie.Year = kf.ChildText("div.ipc-list-card--span > span.ipc-primary-image-list-card__secondary-text")

			tmpProfile.TopMovies = append(tmpProfile.TopMovies, tmpMovie)
		})

		js, err := json.MarshalIndent(tmpProfile, "", "\t")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(js))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting :", r.URL.String())
	})

	infoCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting profile :", r.URL.String())
	})

	startUrl := fmt.Sprintf("https://www.imdb.com/search/name/?birth_monthday=%d-%d", month, day)
	c.Visit(startUrl)
}
