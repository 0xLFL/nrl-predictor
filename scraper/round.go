package main

import (
	"sync"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

type Round struct {
	matches []*Match
	startDay string
	endDay string
	roundName string
	roundIndex int
}

type RoundMatch struct {
	homeTeam string
	awayTeam string
	url string
}

func (r *Round) String() string {
	matches := createListStr(r.matches)
	return fmt.Sprintf(`{
		"round": "%s",
		"weekIndex": "%d",
		"startDay": "%s",
		"endDay": "%s",
		"matches": %s
	}`, r.roundName, r.roundIndex, r.startDay, r.endDay, matches)
}

func ExtractAllMatches(html string) ([]RoundMatch, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	var matches []RoundMatch

	doc.Find(".match").Each(func(i int, s *goquery.Selection) {
		home := strings.TrimSpace(s.Find(".match-team__name--home").Text())
		away := strings.TrimSpace(s.Find(".match-team__name--away").Text())

		url := ""
		s.Find(`a.match--highlighted.u-flex-column.u-flex-align-items-center.u-width-100`).Each(func(_ int, a *goquery.Selection) {
			if href, exists := a.Attr("href"); exists {
				url = href
			}
		})

		if home != "" && away != "" && url != "" {
			matches = append(matches, RoundMatch{
				homeTeam: home,
				awayTeam: away,
				url: url,
			})
		}
	})

	return matches, nil
}

func scrapeRound(round *Round, season *Season, f Fetcher, wg *sync.WaitGroup, stats *StatsTracker) {
	defer wg.Done()
	stats.Start()
	defer stats.Finish()

	content, err := f.Fetch(
		fmt.Sprintf("https://www.nrl.com/draw/?competition=111&round=%v&season=%s", round.roundIndex, season.year),
		chromedp.Tasks{},
		true,
	)

	if err != nil {
		fmt.Println(round, season)
		panic(err)
	}

	writeToFile(content, fmt.Sprintf("%d.html", round.roundIndex))
	matches, err := ExtractAllMatches(content)
	if err != nil {
		panic(err)
	}

	start, end, err := parseRoundDates(content)
	if err != nil {
		panic(err)
	}

	round.startDay = start
	round.endDay = end
	fmt.Printf("days saved, %s, %s", start, end)
	fmt.Println(matches)

	for _, v := range matches {
		wg.Add(1)
		match := &Match{
			homeTeam: v.homeTeam,
			awayTeam: v.awayTeam,
		}
		round.matches = append(round.matches, match)
	
		scrapeMatch(match, v.url, f, wg, stats)
	}
}

func scrapeRounds(rounds []string, season *Season, f Fetcher, wg *sync.WaitGroup, stats *StatsTracker) {
	defer wg.Done()
	stats.Start()
	defer stats.Finish()

	for i, v := range rounds {
		if (i == 0) {
			wg.Add(1)

			round := &Round{roundName: v, roundIndex: i + 1 }
			season.rounds = append(season.rounds, round)
			go scrapeRound(round, season, f, wg, stats)
		}
	}
}

func parseRoundDates(html string) (string, string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", "", err
	}

	var dates []string

	doc.Find("p.match-header__title").Each(func(i int, s *goquery.Selection) {
		dateStr := strings.TrimSpace(s.Text())
		if dateStr == "" {
			return
		}

		dates = append(dates, dateStr)
	})

	if len(dates) == 0 {
		return "", "", fmt.Errorf("no dates found")
	}

	minDate := dates[0]
	maxDate := dates[len(dates) - 1]

	return minDate, maxDate, nil
}
