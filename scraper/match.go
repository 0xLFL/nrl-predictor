package main

import (
	"sync"
	"fmt"
	"time"
	"strings"
	"strconv"

	"github.com/chromedp/chromedp"
	"github.com/PuerkitoBio/goquery"
)

type Play struct {
	time int
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

type PosAndComp struct {
	homePosPer int
	awayPosPer int

	homePosTime int
	awayPosTime int

	homeSets int
	homeSetCompleated int

	awaySets int
	awaySetsCompleated int
}

type Attack struct {
	homeRuns int
	awayRuns int

	homeRunMeters int
	awayRunMeters int

	homePostContactMeters int
	awayPostContactMeters int

	homeLineBreaks int
	awayLineBreaks int

	homeAvgSetDistance float64
	awayAvgSetDistance float64

	homeKickReturnMeters int
	awayKickReturnMeters int

	homeAvgPlayTheBallSpeed int
	awayAvgPlayTheBallSpeed int
}

type Passing struct {
	homeOffloads int
	awayOffloads int

	homeReceipts int
	awayReceipts int

	homeTotalPasses int
	awayTotalPasses int

	homeDummyPasses int
	awayDummyPasses int
}

type Kicking struct {
	homeKicks int
	awayKicks int

	homeKickingMeters int
	awayKickingMeters int

	homeForcedDropOuts int
	awayForcedDropOuts int

	homeKickDefusal int
	awayKickDefusal int

	homeBombs int
	awayBombs int

	homeGrubbers int
	awayGrubbers int
}

type Defence struct {
	homeEffecTackle int
	awayEffecTackle int

	homeTacklesMade int
	awayTacklesMade int

	homeMissedTackles int
	awayMissedTackles int

	homeIntercepts int
	awayIntercepts int

	homeIneffecTackles int
	awayIneffecTackles int
}

type NegPlays struct {
	homeErrors int
	awayErrors int

	homePenCon int
	awayPenCon int

	homeRuckInf int
	awayRuckInf int

	homeInside10 int
	awayInside10 int

	homeOnReport int
	awayOnReport int
}

type MatchStats struct {
	posAndComp PosAndComp
	attack Attack
	passing Passing
	kicking Kicking
	defence Defence
	negPlays NegPlays
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

func createTeamListString(teamList []*Player) (str string) {
	str = "["
	for i, v := range teamList {
		str += v.String()

		if i < len(teamList) - 1 {
			str += ", "
		}
	}

	str += "]"
	return
}

func (m Match) String() string {
	homeTeamList := createTeamListString(m.homeTeamList)
	awayTeamList := createTeamListString(m.awayTeamList)

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
		"playByPlay": "",
		"stats": ""}`,
		m.homeTeam,
		m.homeScore,
		homeTeamList,
		m.awayTeam,
		m.awayScore,
		awayTeamList,
		m.location,
		m.datePlayed,
		m.weather,
	)
}

func scrapeMatch(m *Match, url string, f Fetcher, wg *sync.WaitGroup, stats *StatsTracker) {
	defer wg.Done()

	wg.Add(3)
	go scrapePlaybyPlay(m.playByPlay, url, f, wg,  stats)
	go scrapeTeamList(m, url, f, wg, stats)
	go scrapeMatchStats(m.stats, url, f, wg, stats)
}

func scrapePlaybyPlay(playByPlay []*Play, url string, f Fetcher, wg *sync.WaitGroup, stats *StatsTracker) {
	defer wg.Done()
	
	stats.Start()
	defer stats.Finish()
}

func scrapeTeamList(m *Match, url string, f Fetcher, wg *sync.WaitGroup, stats *StatsTracker) {
	defer wg.Done()

	stats.Start()
	defer stats.Finish()

	content, err := f.Fetch(
			fmt.Sprintf("https://www.nrl.com/%s", url),
			chromedp.Tasks{
				chromedp.Sleep(2 * time.Second),
				chromedp.Click(`//a[.//span[contains(text(), "Team Lists")]]`, chromedp.BySearch),
				chromedp.Sleep(1 * time.Second),
			},
			true,
	)

	var doc *goquery.Document
	if err == nil {
		doc, err = goquery.NewDocumentFromReader(strings.NewReader(content))
		m.homeTeamList, m.awayTeamList, _ = ExtractTeamPlayers(doc)
		
	} else {
		content, err = f.Fetch(
			fmt.Sprintf("https://www.nrl.com/%s", url),
			chromedp.Tasks{
				chromedp.Sleep(2 * time.Second),
			},
			true,
		)

		if err == nil {
			doc, err = goquery.NewDocumentFromReader(strings.NewReader(content))
		}
	}

	if err == nil {
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

func scrapeMatchStats(stats *MatchStats, url string, f Fetcher, wg *sync.WaitGroup, stats_ *StatsTracker) {
	defer wg.Done()

	stats_.Start()
	defer stats_.Finish()
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
			fmt.Println(str, name)
			
			if len(name) >= 2 {
				pHome.nameFirst = name[0]
				pHome.nameLast = name[1]
			}
		})

		b.Find(".team-list-profile:not(.team-list-profile--home) > div.team-list-profile-content > div.team-list-profile__name").Each(func(_ int, s *goquery.Selection) {
			str := strings.TrimSpace(s.Text())
			name := strings.Fields(str)
			fmt.Println(str, name)

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