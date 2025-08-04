package main

import (
	"sync"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

type Round struct {
	matches []Match
	startDay string
	endDay string
	roundName string
	roundIndex int
}

type RoundMatch struct {
	HomeTeam string
	AwayTeam string
}

func ExtractAllMatches(html string) ([]RoundMatch, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	var matches []RoundMatch

	// Loop through all .match elements
	doc.Find(".match").Each(func(i int, s *goquery.Selection) {
		home := strings.TrimSpace(s.Find(".match-team__name--home").Text())
		away := strings.TrimSpace(s.Find(".match-team__name--away").Text())

		if home != "" && away != "" {
			matches = append(matches, RoundMatch{HomeTeam: home, AwayTeam: away})
		}
	})

	return matches, nil
}

func scrapeRound(round string, roundIndex int, season string, f Fetcher, wg *sync.WaitGroup) {
	defer wg.Done()

	content, error_ := f.Fetch(
		fmt.Sprintf("https://www.nrl.com/draw/?competition=111&round=%v&season=%s", roundIndex, season),
		chromedp.Tasks{},
	)

	matches, err := ExtractAllMatches(content)
	fmt.Println(matches, round, season, err, error_)
}

func scrapeRounds(rounds []string, season string, f Fetcher, wg *sync.WaitGroup) {
	defer wg.Done()

	for i, v := range rounds {
		fmt.Printf("%s, %v\n", v, i + 1)
		wg.Add(1)
		go scrapeRound(v, i + 1, season, f, wg)
	}
}
