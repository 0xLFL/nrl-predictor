package main

import (
	"sync"
	"fmt"
	"strings"
	"context"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
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

func scrapeRound(roundIndex int, season string, roundID uuid.UUID, compID int, datesSet bool, f Fetcher, wg *sync.WaitGroup, stats *StatsTracker) {
	defer wg.Done()
	stats.Start()
	defer stats.Finish()

	fmt.Println(roundIndex)

	content, err := f.Fetch(
		fmt.Sprintf("https://www.nrl.com/draw/?competition=%d&round=%d&season=%s", compID, roundIndex, season),
		chromedp.Tasks{},
		true,
	)

	if err != nil {
		return
	}

	matches, err := ExtractAllMatches(content)
	if err != nil {
		return
	}

	if !datesSet {
		start, end, err := parseRoundDates(content)
		if err != nil {
			return
		}

		db, err := NewDB()
		if err != nil {
			return
		}
		defer db.Conn.Close() 

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		fmt.Println(roundID, start, end)
		db.SetRoundDates(ctx, roundID, start, end)
	}

  db, err := NewDB()
	if err != nil {
		return
	}
	defer db.Conn.Close() 

	for _, v := range matches {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		matchID, err := db.CreateMatch(ctx, roundID, v.homeTeam,v.awayTeam)
	
		if err == nil {
			wg.Add(1)
			go scrapeMatch(matchID, v.url, f, wg, stats)
		}
	}
}

func scrapeRounds(rounds []string, season string, seasonID uuid.UUID, compID int, f Fetcher, wg *sync.WaitGroup, stats *StatsTracker) {
	defer wg.Done()
	stats.Start()
	defer stats.Finish()

	db, err := NewDB()
	if err != nil {
		return
	}
	defer db.Conn.Close() 

	for i, v := range rounds {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		roundID, datesSet, err := db.CreateRound(ctx, i + 1, v, seasonID)

		if err == nil {
			wg.Add(1)
			fmt.Println(datesSet)
			go scrapeRound(i + 1, season, roundID, compID, datesSet, f, wg, stats)
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
