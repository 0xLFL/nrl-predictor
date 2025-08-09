package main

import (
	"sync"
	"fmt"
	"strings"
	"strconv"

	"github.com/chromedp/chromedp"
	"github.com/PuerkitoBio/goquery"
)

type Play struct {
	time string
	play string
	team string
	notes string
}

type Player struct {
	nameFirst string
	nameLast string
	position string
	number int
	playerStats PlayerStats
}

type MatchOffical struct {
	nameFirst string
	nameLast string
	role string
}

type Match struct {
	homeTeam string
	homeScore int
	homeTeamList []*Player

	awayTeam string
	awayScore int
	awayTeamList []*Player

	matchOfficals []*MatchOffical

	location string
	kickoffTime string
	datePlayed string
	weather	string

	playByPlay []*Play

	stats *MatchStats
}

func (p *Play) String() string {
	return fmt.Sprintf(`
	{ 
		"play": "%s",
		"team": "%s",
		"notes": "%s",
		"time": "%s"
	}`, p.play, p.team, p.notes, p.time)
}

func (p *Player) String() string {
	return fmt.Sprintf(`
	{ 
		"nameFirst": "%s",
		"nameLast": "%s",
		"position": "%s",
		"number": %d,
		"playerStats": ""
	}`, p.nameFirst, p.nameLast, p.position, p.number)
}

func (m *Match) String() string {
	homeTeamList := createListStr(m.homeTeamList)
	awayTeamList := createListStr(m.awayTeamList)
	playByPlay := createListStr(m.playByPlay)

	return fmt.Sprintf(`
		{ "homeTeam": "%s",
		"homeScore": %d,
		"homeTeamList": %s,
		"awayTeam": "%s", 
		"awayScore": %d,
		"awayTeamList": %s,
		"matchOfficals": "",
		"location": "%s",
		"datePlayed": "%s",
		"weather":	"%s",
		"playByPlay": %s,
		"stats": %s}`,
		m.homeTeam,
		m.homeScore,
		homeTeamList,
		m.awayTeam,
		m.awayScore,
		awayTeamList,
		m.location,
		m.datePlayed,
		m.weather,
		playByPlay,
		m.stats,
	)
}

func scrapeMatch(m *Match, url string, f Fetcher, wg *sync.WaitGroup, stats *StatsTracker) {
	defer wg.Done()
	stats.Start()
	defer stats.Finish()

	content, err := f.Fetch(
		fmt.Sprintf("https://www.nrl.com/%s", url),
		chromedp.Tasks{},
		true,
	)

	if err != nil {
		return;
	}

	wg.Add(3)
	go parseMatchStats(m.stats, content, wg)
	go parsePlaybyPlay(&m.playByPlay, content, wg)
	go parseTeamList(m, content, wg)
}

func parsePlaybyPlay(playByPlay *[]*Play, content string, wg *sync.WaitGroup) {
	defer wg.Done()

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return;
	}

	doc.Find("div.match-centre-event").Each(func(_ int, b *goquery.Selection) {
		play := &Play{}
		b.Find(".match-centre-event__team-name").Each(func(_ int, s *goquery.Selection) {
			play.team = strings.TrimSpace(s.Text())
		})

		b.Find(".match-centre-event__title").Each(func(_ int, s *goquery.Selection) {
			play.play = strings.TrimSpace(s.Text())
		})

		b.Find(".u-font-weight-500").Each(func(_ int, s *goquery.Selection) {
			play.notes = strings.TrimSpace(play.notes + " " + strings.Join(strings.Fields(s.Text()), " "))
		})

		b.Find("span.match-centre-event__timestamp").Each(func(_ int, s *goquery.Selection) {
			play.time = strings.TrimSpace(s.Text())
		})

		*playByPlay = append(*playByPlay, play)
	})
}

func parseTeamList(m *Match, content string, wg *sync.WaitGroup) {
	defer wg.Done()

	var doc *goquery.Document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	m.homeTeamList, m.awayTeamList, _ = ExtractTeamPlayers(doc)

	if err != nil {
		sel := doc.Find(".match-team__score.match-team__score--home").First()
		text := sel.Clone().Children().Remove().End().Text()
		trimmed := strings.TrimSpace(text)
		score, err := strconv.Atoi(trimmed)
		if err == nil {
			m.homeScore = score
		}

		sel = doc.Find(".match-team__score.match-team__score--away").First()
		text = sel.Clone().Children().Remove().End().Text()
		trimmed = strings.TrimSpace(text)
		score, err = strconv.Atoi(trimmed)
		if err == nil {
			m.awayScore = score
		}

		sel = doc.Find(".match-venue.o-text").First()
		text = sel.Clone().Children().Remove().End().Text()
		location := strings.TrimSpace(text)
		if err == nil {
			m.location = location
		}

		sel = doc.Find("p.match-header__title").First()
		dateStr := strings.TrimSpace(sel.Text())
		if dateStr != "" {
			m.datePlayed = dateStr
		}

		doc.Find("p.match-weather__text").Each(func(i int, s *goquery.Selection) {
			if strings.Contains(s.Text(), "Weather:") {
				m.weather = strings.TrimSpace(s.Find("span").Text())
			}
		})
	}
}

func ExtractTeamPlayers(doc *goquery.Document) ([]*Player, []*Player, error) {
	var hPlayers []*Player
	var aPlayers []*Player

	// Locate the home team block using heading text
	doc.Find("div.team-list__container > div.team-list").Each(func(_ int, b *goquery.Selection) {
		homeNumber := 0
		awayNumber := 0
		position := ""
		
		b.Find("div.team-list-position > span.team-list-position__text").Each(func(_ int, s *goquery.Selection) {
			position = strings.TrimSpace(s.Text())
		})

		b.Find("div.team-list-position > p > span.team-list-position__number:not(.u-text-align-left)").Each(func(_ int, s *goquery.Selection) {
			numText := strings.TrimSpace(s.Text())
			homeNumber, _ = strconv.Atoi(numText)
			awayNumber = homeNumber
		})

		b.Find("div.team-list-position > p > span.team-list-position__number.u-text-align-left").Each(func(_ int, s *goquery.Selection) {
			numText := strings.TrimSpace(s.Text())
			awayNumber, _ = strconv.Atoi(numText)
		})

		pHome := &Player{
			position: position,
			number: homeNumber,
		}
	
		pAway := &Player{
			position: position,
			number: awayNumber,
		}

		b.Find(".team-list-profile:not(.team-list-profile--away) > div.team-list-profile-content > div.team-list-profile__name").Each(func(_ int, s *goquery.Selection) {
			str := strings.TrimSpace(s.Text())
			name := strings.Fields(str)
			
			if len(name) >= 2 {
				pHome.nameFirst = name[0]
				pHome.nameLast = name[1]
			}
		})

		b.Find(".team-list-profile:not(.team-list-profile--home) > div.team-list-profile-content > div.team-list-profile__name").Each(func(_ int, s *goquery.Selection) {
			str := strings.TrimSpace(s.Text())
			name := strings.Fields(str)

			if len(name) >= 2 {
				pAway.nameFirst = name[0]
				pAway.nameLast = name[1]
			}
		})

		hPlayers = append(hPlayers, pHome)
		aPlayers = append(aPlayers, pAway)
	})

	return hPlayers, aPlayers, nil
}